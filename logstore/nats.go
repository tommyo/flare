package logstore

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rivian/delta-go/logstore"
	"github.com/rivian/delta-go/storage"
	"github.com/tommyo/flare"
)

var _ logstore.LogStore = &NatsStore{}

func timeToBytes(t time.Time) []byte {
	return []byte(t.Format(time.RFC3339))
}

func bytesToTime(b []byte) (time.Time, error) {
	return time.Parse(time.RFC3339, string(b))
}

func buildPath(paths ...storage.Path) string {
	path := "log"
	for _, p := range paths {
		path = path + "." + p.Raw
	}
	return path
}

func encodeEntry(entry *logstore.CommitEntry) ([]byte, error) {
	b := new(bytes.Buffer)
	enc := gob.NewEncoder(b)

	if err := enc.Encode(entry); err != nil {
		return nil, fmt.Errorf("failed to encode entry: %w", err)
	}
	return b.Bytes(), nil
}

func decodeEntry(b []byte) (*logstore.CommitEntry, error) {
	entry := new(logstore.CommitEntry)
	dec := gob.NewDecoder(bytes.NewReader(b))

	if err := dec.Decode(entry); err != nil {
		return nil, fmt.Errorf("failed to decode entry: %w", err)
	}
	return entry, nil
}

type NatsStore struct {
	conf *flare.Config
	nats *nats.Conn
	kv   jetstream.KeyValue
	// NATS JetStream key-value store currently only supports store level expiry
	trueExiry time.Duration

	latests map[string]*logstore.CommitEntry
	watch   jetstream.KeyWatcher
}

// Client implements logstore.LogStore.
func (n *NatsStore) Client() any {
	return n.kv
}

// ExpirationDelaySeconds implements logstore.LogStore.
func (n *NatsStore) ExpirationDelaySeconds() uint64 {
	return uint64(n.trueExiry.Seconds())
}

// Get implements logstore.LogStore. Doesn't actually return the value?
// TODO add context support upstream
func (n *NatsStore) Get(tablePath storage.Path, fileName storage.Path) (*logstore.CommitEntry, error) {
	key := buildPath(tablePath, fileName)
	raw, err := n.kv.Get(context.TODO(), key)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	return decodeEntry(raw.Value())
}

// Latest implements logstore.LogStore.
// TODO should missing entry return an error?
func (n *NatsStore) Latest(tablePath storage.Path) (*logstore.CommitEntry, error) {
	entry, ok := n.latests[tablePath.Raw]
	if !ok {
		return nil, nil
	}
	return entry, nil
}

// Put implements logstore.LogStore.
func (n *NatsStore) Put(entry *logstore.CommitEntry, overwrite bool) error {
	key := buildPath(entry.TablePath(), entry.FileName())

	val, err := encodeEntry(entry)
	if err != nil {
		return err
	}

	if overwrite {
		_, err = n.kv.Put(context.TODO(), key, val)
		return err
	}
	_, err = n.kv.Create(context.TODO(), key, val)
	return err
}

func (n *NatsStore) Connect() error {
	var err error
	n.nats, err = nats.Connect(n.conf.String("nats.url"))
	if err != nil {
		return err
	}
	js, err := jetstream.New(n.nats)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	kv := n.conf.String("nats.kvstore")
	expiry := n.conf.Duration("nats.expiry")

	n.kv, err = js.KeyValue(ctx, kv)
	if err != nil && errors.Is(err, nats.ErrBucketNotFound) {
		// TODO flesh out the config
		cfg := jetstream.KeyValueConfig{
			Bucket: kv,
			TTL:    expiry,
		}
		n.kv, err = js.CreateKeyValue(ctx, cfg)
	}
	if err != nil {
		return err
	}

	status, err := n.kv.Status(ctx)
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	n.trueExiry = status.TTL()

	// FIXME switch to context.WithCancel
	background := context.Background()
	n.watch, err = n.kv.Watch(background, "log.*")
	if err != nil {
		return err
	}

	go n.onUpdate(background)

	return nil
}

func (n *NatsStore) onUpdate(ctx context.Context) {
	for {
		select {
		case v := <-n.watch.Updates():
			if v == nil {
				continue
			}
			// FIXME handle errors
			val, _ := decodeEntry(v.Value())
			tablePath := val.TablePath()
			fileName := val.FileName()
			entry, ok := n.latests[tablePath.Raw]
			if !ok {
				n.latests[tablePath.Raw] = val
				continue
			}
			if fileName.Raw > entry.FileName().Raw {
				n.latests[tablePath.Raw] = val
			}
		case <-ctx.Done():
			return
		}
	}
}

func (n *NatsStore) Close() error {
	if err := n.watch.Stop(); err != nil {
		return err
	}
	n.nats.Close()
	return nil
}

func NewNatsStore(conf *flare.Config) *NatsStore {
	conf.RegisterDefault("nats.url", "", nats.DefaultURL, "NATS server URL")
	conf.RegisterDefault("nats.kvstore", "", "flare", "NATS JetStream key-value store name")
	conf.RegisterDefault("nats.expiry", "", "1h", "NATS JetStream key-value store expiry duration")

	return &NatsStore{
		conf:    conf,
		latests: make(map[string]*logstore.CommitEntry),
	}
}
