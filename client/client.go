package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
)

// Define okurulog internal constants (TODO: move config_dir const to a command flag)
const (
	OKURULOG_NAME         string = "OkuruLog"
	OKURULOG_DESCRIPTION  string = "The small, simple, faster log shipper"
	OKURULOG_VERSION      string = "0.0.1a"
	OKURULOG_VERSION_NAME string = "Donskoy"
	OKURULOG_AUTHOR       string = "Patrick Kuti <code@introspect.in>"
	OKURULOG_CONFIG_DIR   string = "/etc/okurulog/"
)

// JSON struct for okurulog client configuration
type ClientConfiguration struct {
	ServerHostname            string
	ServerPort                int
	NumberOfTimesToRetry      int
	RetryInterval             int
	MaximumSizeOfLogFileCache int
	LogDirectory              string
}

// Struct for okurulog server
type LogServer struct {
	Hostname   string
	Port       int
	IsOnline   bool
	Connection net.Conn
}

// Struct for okurulog client
type LogClient struct {
}

// Struct for logfile
type LogLine struct {
	UUID string
}

// Struct for a line in a logfile
type LogFile struct {
	UUID string
}

// Variable for global client configuration
var OkuruLogClientConfiguration *ClientConfiguration = &ClientConfiguration{}

// setDefaults sets default configuration values
func (self *ClientConfiguration) setDefaults() {
	self.ServerHostname = "okurulog-master"
	self.ServerPort = 7070
	self.NumberOfTimesToRetry = 5
	self.RetryInterval = 60
	self.MaximumSizeOfLogFileCache = 5242880
	self.LogDirectory = "/var/log/okurulog/"
}

// Init parses configuration files, compiles global variables
// sets up logging and watchers
func init() {

	// define variables ahead of time
	var err error

	// Check if okurulog client configuration file exists
	OkuruLogClientConfigurationFilePath := OKURULOG_CONFIG_DIR + "config.json"
	_, err = os.Stat(OkuruLogClientConfigurationFilePath)
	if err != nil {
		log.Fatal("Error looking for configuration file: " + OkuruLogClientConfigurationFilePath)
		os.Exit(1)
	}

	// Load default configuration
	OkuruLogClientConfiguration.setDefaults()

	// Attempt to load and parse json configuration file into global configuration
	_loadConfig(OkuruLogClientConfigurationFilePath)

	// Check if log directory exists and is a directory
	LogDirectory, err := os.Stat(OkuruLogClientConfiguration.LogDirectory)
	if err != nil {
		log.Fatal("Error trying to access log directory ("+OkuruLogClientConfiguration.LogDirectory+"), error returned was: ", err)
		os.Exit(1)
	}
	if !LogDirectory.IsDir() {
		log.Fatal("Log directory is not a directory, error returned was: ", err)
		os.Exit(1)
	}

}

// _loadConfig loads values from a configuration file into global configuration
func _loadConfig(OkuruLogClientConfigurationFilePath string) {

	// Try to open the okurulog client configuration file and read it into global variable
	OkuruLogClientConfigurationFile, err := ioutil.ReadFile(OkuruLogClientConfigurationFilePath)
	if err != nil {
		log.Fatal("Error reading configuration file, error returned was: ", err)
		os.Exit(1)
	}

	// Decode JSON from configuration file and dump into global configuration to override default values
	err = json.Unmarshal(OkuruLogClientConfigurationFile, OkuruLogClientConfiguration)
	if err != nil {
		log.Fatal("Error parsing configuration file, please ensure it is valid JSON. Error returned was: ", err)
		os.Exit(1)
	}

}

// connect tries to connect to the logserver's hostname and port,
// Returns true or false depending on if it was able to connect
func (self *LogServer) connect() (net.Conn, error) {

	// Define variables ahead of time
	var err error

	// Try to connect to logserver and error out on failure
	logServerConnection, err := net.DialTimeout("tcp", self.Hostname+":"+strconv.Itoa(self.Port), 10)
	if err != nil {
		self.IsOnline = false
		log.Fatal("Error connecting to LogServer, error return was: ", err)
		os.Exit(1)
	}

	self.IsOnline = true

	return logServerConnection, err
}

// setHostname sets the hostname in the logserver struct
func (self *LogServer) setHostname(Hostname string) {
	self.Hostname = Hostname
}

// setPort sets the port in the logserver struct
func (self *LogServer) setPort(Port int) {
	self.Port = Port
}

func main() {

	/**
	 * setup, init, load configs,
	 * check connection to server
	 * check access to directories
	 * check access to logfiles ()
	 * check backlog on logfiles, if there is send from there at reduced speed/rate
	 * if not, start sending to server via msgpack
	 */

	// Variable for global client configuration
	var LogServer *LogServer = &LogServer{}
	LogServer.setHostname(OkuruLogClientConfiguration.ServerHostname)
	LogServer.setPort(OkuruLogClientConfiguration.ServerPort)

	LogServer.connect()

}
