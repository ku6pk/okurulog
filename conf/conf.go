package conf

// Import needed packages
import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// Define okurulog internal constants
const (
	OKURULOG_NAME         string = "OkuruLog"
	OKURULOG_DESCRIPTION  string = "The smaller, simpler, faster log shipper."
	OKURULOG_VERSION      string = "0.0.1a"
	OKURULOG_VERSION_NAME string = "Donskoy"
	OKURULOG_AUTHOR       string = "Patrick Kuti <code@introspect.in>"
	OKURULOG_CONFIG_DIR   string = "/etc/okurulog/"
)

// JSON struct for okurulog server configuration
type ServerConfiguration struct {
	Hostname string
	Port     struct {
		Client int
		GUI    int
	}
	CacheDirectory string
	LogDirectory   string
}

// JSON struct for okurulog client configuration
type ClientConfiguration struct {
	ServerHostname             string
	ServerPort                 int
	NumberOfTimesToRetry       int
	RetryInterval              int
	MaximumSizeOfALogFileCache int
	CacheDirectory             string
	LogDirectory               string
	WatchFiles                 []string
}

// SetDefaults sets default configuration values for a server
func (self *ServerConfiguration) SetDefaults() {
	self.Hostname = "okurulog-master"
	self.Port.Client = 7070
	self.Port.GUI = 7080
	self.CacheDirectory = "/var/cache/okurulog/server/"
	self.LogDirectory = "/var/log/okurulog/server/"
}

// SetDefaults sets default configuration values for a client
func (self *ClientConfiguration) SetDefaults() {
	self.ServerHostname = "okurulog-master"
	self.ServerPort = 7070
	self.NumberOfTimesToRetry = 5
	self.RetryInterval = 60
	self.MaximumSizeOfALogFileCache = 5242880
	self.CacheDirectory = "/var/cache/okurulog/client/"
	self.LogDirectory = "/var/log/okurulog/client/"
	self.WatchFiles = []string{"/var/log/messages", "/var/log/secure"}
}

// LoadConfig loads values from a configuration file into a client configuration
func (self *ClientConfiguration) LoadConfig() {

	// Set the config path for the client config file
	ConfigFilePath := OKURULOG_CONFIG_DIR + "client.json"

	// Define variables ahead of time
	var err error

	// Check if configuration file exists
	_, err = os.Stat(ConfigFilePath)
	if err != nil {
		log.Fatal("Error looking for configuration file: " + ConfigFilePath)
		os.Exit(1)
	}

	// Try to open the okurulog configuration file and read it into global variable
	ConfigFile, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal("Error reading configuration file, error returned was: ", err)
		os.Exit(1)
	}

	// Decode JSON from configuration file and dump into global configuration to override default values
	err = json.Unmarshal(ConfigFile, self)
	if err != nil {
		log.Fatal("Error parsing configuration file, please ensure it is valid JSON. Error returned was: ", err)
		os.Exit(1)
	}
}

// LoadConfig loads values from a configuration file into a server configuration
func (self *ServerConfiguration) LoadConfig(Configuration ServerConfiguration) {

	// Set the config path for the server config file
	ConfigFilePath := OKURULOG_CONFIG_DIR + "server.json"

	// Define variables ahead of time
	var err error

	// Check if configuration file exists
	_, err = os.Stat(ConfigFilePath)
	if err != nil {
		log.Fatal("Error looking for configuration file: " + ConfigFilePath)
		os.Exit(1)
	}

	// Try to open the okurulog configuration file and read it into global variable
	ConfigFile, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		log.Fatal("Error reading configuration file, error returned was: ", err)
		os.Exit(1)
	}

	// Decode JSON from configuration file and dump into global configuration to override default values
	err = json.Unmarshal(ConfigFile, Configuration)
	if err != nil {
		log.Fatal("Error parsing configuration file, please ensure it is valid JSON. Error returned was: ", err)
		os.Exit(1)
	}
}
