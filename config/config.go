package config

import (
	"flag"
	"os"
	"strings"
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultAccrualAddress = "http://localhost:8888"
)

// Config contains application settings
type Config struct {
	ServerAddress  string
	DatabaseDSN    string
	AccrualAddress string
}

var defaultConfig = Config{
	ServerAddress:  defaultServerAddress,
	AccrualAddress: defaultAccrualAddress,
}

// New creates config by merging default settings with flags, then with env variables
// The last nonempty value takes precedence (default < flag < env) except for
// DATABASE_URI env variable which overrides -d flags even if empty
// New also handles validation and returns non-nil error if validation failed
func New() (Config, error) {
	c := defaultConfig
	c.parseFlags()
	c.parseEnvVars()
	err := c.Validate()
	return c, err
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", defaultServerAddress, "network address the server listens on")
	flag.StringVar(&c.DatabaseDSN, "d", "", `database dsn (default "")`)
	flag.StringVar(&c.AccrualAddress, "r", defaultAccrualAddress, "accrual system address")

	flag.Parse()
}

func (c *Config) parseEnvVars() {
	sa := os.Getenv("RUN_ADDRESS")
	if sa != "" {
		c.ServerAddress = sa
	}

	dd, ok := os.LookupEnv("DATABASE_URI")
	if ok {
		// empty string is valid here, overrides -d flag and returns the default ""
		c.DatabaseDSN = dd
	}

	aa := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if aa != "" {
		c.AccrualAddress = aa
	}
}

func (c *Config) Validate() error {
	// Remove leading and trailing spaces without complaining
	// Other mistakes and typos are to be considered as errors
	c.ServerAddress = strings.TrimSpace(c.ServerAddress)
	c.AccrualAddress = strings.TrimSpace(c.AccrualAddress)

	if err := ValidateServerAddress(c.ServerAddress); err != nil {
		return err
	}

	if err := ValidateURL(c.AccrualAddress); err != nil {
		return err
	}

	// No need to validate c.DatabaseDSN, storage itself will do the job
	return nil
}
