package trace 

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// Configuration contains configuration info
type Configuration struct {
	Username string
	Password string
	Protocol string
	Address  string
	Port     string
}

// Init reads database backend configuration with TOML format
func (conf *Configuration) Init(configfile string) error {
	// Get database configuration file
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Configuration file missing", err)
		return err
	}

	// Read configuration file
	if _, err := toml.DecodeFile(configfile, &conf); err != nil {
		log.Fatal(err)
	}

	// Return
	return nil
}

// GetDSN returns data source name from the config file
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func (conf *Configuration) GetDSN() string {
	dsn := conf.Username + ":" + conf.Password + "@"
	dsn += conf.Protocol + "(" + conf.Address + ":" + conf.Port + ")"
	dsn += "/"

	return dsn
}
