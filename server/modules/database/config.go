package database

import (
	"encoding/json"
	"os"

	"github.com/golang/glog"
)

// config contains connection info of the database
type config struct {
	Container string `json:"container"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Protocol  string `json:"protocol"`
	Address   string `json:"address"`
	Port      string `json:"port"`
}

// read config file
func (conf *config) read(path string) error {
	// Get configuration file
	configFile, err := os.Open(path)
	if err != nil {
		glog.Fatal("opening config file:", err.Error())
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&conf); err != nil {
		glog.Fatal("parsing config file: ", err.Error())
	}

	// Return
	return nil
}

// getDSN returns data source name from the config file
// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func (conf *config) getDSN() string {
	dsn := conf.Username + ":" + conf.Password + "@"
	dsn += conf.Protocol + "(" + conf.Address + ":" + conf.Port + ")"
	dsn += "/?parseTime=true"

	return dsn
}
