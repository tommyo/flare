package flare

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type Config struct {
	*koanf.Koanf
	flags *pflag.FlagSet
}

func NewConfig() *Config {

	flags := pflag.NewFlagSet("config", pflag.ContinueOnError)
	flags.Usage = func() {
		fmt.Println(flags.FlagUsages())
		os.Exit(0)
	}

	flags.StringP("conf", "c", "", "path to the configuration file")

	flags.Bool("devel", false, "enable development mode")

	return &Config{
		koanf.NewWithConf(koanf.Conf{
			Delim:       ".",
			StrictMerge: true,
		}),
		flags,
	}
}

func (p *Config) RegisterDefault(param string, short string, value string, usage string) {
	param = strings.Replace(param, ".", "-", -1)
	if short != "" {
		p.flags.StringP(param, short, value, usage)
	} else {
		p.flags.String(param, value, usage)
	}
}

func (p *Config) Register(param string, short string, usage string) {
	param = strings.Replace(param, ".", "-", -1)
	if short != "" {
		p.flags.StringP(param, short, "", usage)
	} else {
		p.flags.String(param, "", usage)
	}
}

func (p *Config) RegisterOption(param string, short string, usage string) {
	param = strings.Replace(param, ".", "-", -1)
	if short != "" {
		p.flags.StringP(param, short, "", usage)
	} else {
		p.flags.BoolP(param, short, false, usage)
	}
}

func (p *Config) Parse() error {

	err := p.flags.Parse(os.Args[1:])
	if err != nil {
		return err
	}

	// if we've set a config file, load it first
	if f, _ := p.flags.GetString("conf"); f != "" {
		var parser koanf.Parser
		switch filepath.Ext(f) {
		case ".json":
			parser = json.Parser()
		case ".yaml", ".yml":
			parser = yaml.Parser()
		case ".toml":
			parser = toml.Parser()
		default:
			return fmt.Errorf("unsupported config file format: %s", f)
		}
		err := p.Load(file.Provider(f), parser)
		if err != nil {
			return err
		}
	}

	envParser := func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "FLARE_")), "_", ".", -1)
	}

	// Load environment variables next
	err = p.Load(env.Provider("FLARE_", ".", envParser), nil)
	if err != nil {
		return err
	}

	// Load the .env files
	err = p.Load(file.Provider(".env"), dotenv.ParserEnv("", ".", envParser))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if devel, _ := p.flags.GetBool("devel"); devel || strings.HasPrefix(p.String("env"), "dev") {
		err = p.Load(file.Provider(".env.devel"), dotenv.ParserEnv("", ".", envParser))
	} else {
		err = p.Load(file.Provider(".env.production"), dotenv.ParserEnv("", ".", envParser))
	}
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	err = p.Load(file.Provider(".env.local"), dotenv.ParserEnv("", ".", envParser))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	// Load the flags and defaults
	p.flags.VisitAll(func(f *pflag.Flag) {
		name := strings.ReplaceAll(f.Name, "-", ".")
		if f.Changed {
			p.Set(name, f.Value.String())
			return
		}
		if f.DefValue != "" && !p.Exists(name) {
			p.Set(name, f.DefValue)
		}
	})

	return nil
}

func (p *Config) Usage() {
	p.flags.Usage()
}
