package flare

// Import necessary packages
import (
	"sync"
)

func in[T comparable](list []T, item T) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

// TODO SessionData and SessionEventsManager may be artificially abstract as a result of the direct translation from Scala.
// TODO

// SessionStatus represents the status of a session with predefined states.
type SessionStatus int

const (
	// Define possible states of SessionStatus
	Pending SessionStatus = iota
	Started
	Closed
)

func (p SessionStatus) String() string {
	// Define string representations for each state
	switch p {
	case Pending:
		return "Pending"
	case Started:
		return "Started"
	case Closed:
		return "Closed"
	default:
		return "Unknown"
	}
}

// SessionKey represents a unique key for a session.
type SessionKey struct {
	UserId    string
	SessionId string
}

type Session struct {
	// sessionID represents the unique identifier for the session.
	id        string
	userID    string
	eventTime int64
	extraTags map[string]string
}

// SessionID returns the session ID.
func (p *Session) ID() string {
	return p.id
}

// Key returns a unique identifier for the session.
func (p *Session) Key() SessionKey {
	return SessionKey{UserId: p.userID, SessionId: p.id}
}

// SessionEventsManager manages session events, posting them to a listener bus.
type SessionManager struct {
	session *Session      // Session for which the events are generated
	status  SessionStatus // Current status of the session, private to the package
	// TODO review mutex implementation
	mu sync.Mutex // Mutex to protect status changes
	// DataFrameCache                     sync.Map // Concurrent map for DataFrame equivalent in Go
	// ListenerCache                      sync.Map // Concurrent map for StreamingQueryListener equivalent in Go
	// StreamingForeachBatchRunnerCleaner *StreamingForeachBatchHelperCleanerCache
	// StreamingServersideListenerHolder  *ServerSideListenerHolder
}

// NewSessionEventsManager creates a new instance of SessionEventsManager.
func NewSessionManager(session *Session) *SessionManager {
	return &SessionManager{
		session: session,
		status:  Pending, // Initialize with Pending status
	}
}

// SessionID returns the session ID from the session holder.
func (sem *SessionManager) SessionID() string {
	return sem.session.ID()
}

// Status returns the current status of the session.
func (sem *SessionManager) Status() SessionStatus {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	return sem.status
}

// SetStatus safely sets the session's status.
func (sem *SessionManager) setStatus(status SessionStatus) {
	sem.mu.Lock()
	defer sem.mu.Unlock()
	sem.status = status
}

func (sem *SessionManager) PostStarted() {
	if !in([]SessionStatus{Pending}, Started) {
		// error
	}
	defer sem.setStatus(Started)
	// SessionData.session.sparkContext.listenerBus.post(SparkListenerConnectSessionStarted(sessionData.sessionId, sessionData.userId, clock.getTimeMillis()))
}

func (sem *SessionManager) PostClosed() {
	if !in([]SessionStatus{Started}, Closed) {
		// error
	}
	defer sem.setStatus(Closed)
	// SessionData.session.sparkContext.listenerBus.post(SparkListenerConnectSessionClosed(sessionData.sessionId, sessionData.userId, clock.getTimeMillis()))
}

// ServerSessionId returns the server side session ID.
func (sh *SessionManager) ServerSessionId() string {
	// Placeholder for logic to return server session ID
	return ""
}
