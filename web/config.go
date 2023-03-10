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
	// Base path of the server (default: "/")
	BasePath string `toml:"base_path,comment"`
	// Runs in release mode by default, unless this is true in which case gin is run in debug mode.
	DebugMode bool `toml:"debug_mode"`
	// Route for health check, set to empty to disable
	HealthRoute string `toml:"health_route,comment"`
	// Settings for get-run
	GetRun ConfigGetRun `toml:"getrun"`
	// Settings for stats
	Stats ConfigStats `toml:"stats"`
	// Settings for upload
	Upload ConfigUpload `toml:"upload"`
}

type ConfigGetRun struct {
	// Get-run route
	Route string `toml:"route,comment"`
	// Require authentication and "getrun" scope
	Auth bool `toml:"auth,comment"`
}

type ConfigUpload struct {
	// Upload route
	Route string `toml:"route,comment"`
	// Store uploads in the database, the main point of this whole thing
	StoreToDb bool `toml:"store_to_db,comment"`
	// If true, stores raw upload json to a file on the disk
	SaveRawToDisk bool `toml:"save_raw_to_disk,comment"`
	// If true, stores raw upload json to the database
	SaveRawToDb bool `toml:"save_raw_to_db,comment"`
	// Directory where runs are stored if save_raw_to_disk is true
	RunsDir string `toml:"runs_dir,comment"`
}

type ConfigStats struct {
	// Route to proxy stats from
	Route string `toml:"route,comment"`
	// Stats HTTP address
	Upstream string `toml:"upstream,comment"`
	// If true, require authentication AND the stats:view scope to access the stats page(s).
	Auth bool `toml:"auth,comment"`
}

func (c Config) Default() Config {
	return Config{
		BasePath:  "/",
		DebugMode: false,
		GetRun: ConfigGetRun{
			Route: "/getrun",
			Auth:  true,
		},
		Stats: ConfigStats{
			Route: "/stats",
			Auth:  true,
		},
		Upload: ConfigUpload{
			Route:       "/upload",
			StoreToDb:   true,
			SaveRawToDb: true,
			RunsDir:     "data/runs",
		},
	}
}

func NewConfig() *Config {
	c := Config{}.Default()
	return &c
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
