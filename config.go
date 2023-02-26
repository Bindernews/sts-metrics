package stms

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/bindernews/sts-msr/chart"
	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	// Set of IP addresses/ports to listen on
	Listen []string `toml:"listen"`
	// SQL section of the config
	Sql ConfigSql `toml:"sql"`
	// Directory where runs are stored
	RunsDir string `toml:"runs_dir"`
	// List of defined charts
	Charts []chart.ChartToml `toml:"chart"`
	// List of extra files to include
}

type ConfigSql struct {
	// File(s) that will be run on startup to create temporary functions,
	// views, etc.
	SetupFiles []string `toml:"setup_files"`
	// Inline startup script
	SetupScript string `toml:"setup_script"`
}

func NewConfig() *Config {
	c := new(Config)
	c.RunsDir = "data/runs"
	return c
}

// Loads the configuration from the named file. If 'mayNotExist' is true,
// the file not existing will NOT be an error, and the default configuration
// will be loaded instead.
func (c *Config) LoadFile(fpath string, mayNotExist bool) error {
	fpath, err := filepath.Abs(fpath)
	if err != nil {
		return err
	}
	read, err := os.Open(fpath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) && mayNotExist {
			return nil
		} else {
			return err
		}
	}
	defer read.Close()
	if err := toml.NewDecoder(read).Decode(c); err != nil {
		return err
	}
	log.Default().Println("loaded config file", fpath)
	return nil
}
