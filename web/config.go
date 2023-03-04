package web

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	// Set of IP addresses/ports to listen on
	Listen string `toml:"listen"`
	// Set of IP addresses/ports for the admin interface to listen on
	AdminListen string `toml:"admin_listen"`
	// Runs in release mode by default, unless this is true in which case
	// gin is run in debug mode.
	DebugMode bool `toml:"debug_mode"`
	// Directory where runs are stored
	RunsDir string `toml:"runs_dir"`
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
