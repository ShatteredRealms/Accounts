package option

import (
	"flag"
	"fmt"
	"os"
)

// ConfigOption Stores the value for the configuration
type ConfigOption struct {
	Value   *string
	Default string
	EnvVar  string
	Flag    string
	Usage   string
}

// Config Represent all the configuration options there are for the application
type Config struct {
	// Port The port for the application
	Port ConfigOption

	// Host The host address for the application
	Host ConfigOption

	// Mode The running mode for the application
	Mode ConfigOption

	// KeyDir The path to the keys for JWT auth
	KeyDir ConfigOption

	// DBFile The path to the connection info for the DB
	DBFile ConfigOption

	// TLSCert The path to the TLS certificate for HTTPS
	TLSCert ConfigOption

	// TLSKey The path to the TLS key for HTTPS
	TLSKey ConfigOption
}

const (
	// DebugMode A config mode allowing for the most verbose logging. Used heavily for development.
	DebugMode = "debug"

	// ReleaseMode A config mode allowing for minimal logging. Intended to be used for production.
	ReleaseMode = "release"
)

var (
	DefaultConfig = Config{
		Port: ConfigOption{
			Default: "8080",
			EnvVar:  "SRO_AUTH_PORT",
			Flag:    "port",
			Usage:   "The port for the application",
		},
		Host: ConfigOption{
			Default: "",
			EnvVar:  "SRO_AUTH_HOST",
			Flag:    "host",
			Usage:   "The host address for the application",
		},
		Mode: ConfigOption{
			Default: DebugMode,
			EnvVar:  "SRO_AUTH_MODE",
			Flag:    "mode",
			Usage:   "The running mode for the application",
		},
		KeyDir: ConfigOption{
			Default: "/etc/sro/auth",
			EnvVar:  "SRO_KEY_DIR",
			Flag:    "keys",
			Usage:   "The path to the keys for JWT auth",
		},
		DBFile: ConfigOption{
			Default: "/etc/sro/db.yaml",
			EnvVar:  "SRO_DB_FILE",
			Flag:    "db",
			Usage:   "The path to the connection info for the DB",
		},
		TLSCert: ConfigOption{
			Default: "/etc/sro/tls/tls.crt",
			EnvVar:  "SRO_TLS_CRT_PATH",
			Flag:    "tlscrt",
			Usage:   "The path to the TLS certificate for HTTPS",
		},
		TLSKey: ConfigOption{
			Default: "/etc/sro/tls/tls.key",
			EnvVar:  "SRO_TLS_CRT_PATH",
			Flag:    "tlskey",
			Usage:   "The path to the TLS key for HTTPS",
		},
	}
)

func NewConfig() Config {
	config := DefaultConfig
	config.readFlags()
	config.readEnvs()
	return config
}

// Address Gets the full address for the HTTP server
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.Host.GetValue(), c.Port.GetValue())
}

// GetValue Gets the value for the Config option if it's set, otherwise it gets the default value
func (co *ConfigOption) GetValue() string {
	v := co.Value
	if v == nil {
		return co.Default
	}
	return *co.Value
}

// IsRelease Returns true if the application mode is release
func (c *Config) IsRelease() bool {
	return c.Mode.GetValue() == ReleaseMode
}

func (c *Config) readFlags() {
	c.Port.readFlag()
	c.Host.readFlag()
	c.Mode.readFlag()
	c.KeyDir.readFlag()
	c.DBFile.readFlag()
	c.TLSCert.readFlag()
	c.TLSKey.readFlag()
	flag.Parse()
}

func (co *ConfigOption) readFlag() {
	co.Value = flag.String(co.Flag, co.Default, co.Usage)
}

func (c *Config) readEnvs() {
	c.Port.readEnv()
	c.Host.readEnv()
	c.Mode.readEnv()
	c.KeyDir.readEnv()
	c.DBFile.readEnv()
	c.TLSCert.readEnv()
	c.TLSKey.readEnv()
}

func (co *ConfigOption) readEnv() {
	env, found := os.LookupEnv(co.EnvVar)
	if found && !isFlagPassed(co.Flag) {
		co.Value = &env
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
