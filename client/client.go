package main

// Import needed packages
import (
	"fmt"
	"github.com/ActiveState/tail"
	olConf "github.com/patyx7/okurulog/conf"
	"log"
	"net"
	"os"
	"strconv"
)

// Variable for global client configuration
var OkuruLogClientConfiguration *olConf.ClientConfiguration = &olConf.ClientConfiguration{}

// Struct for okurulog server
type LogServer struct {
	Hostname   string
	Port       int
	IsOnline   bool
	Connection net.Conn
}

// Init parses and loads configuration files
func init() {

	// Load default configuration
	OkuruLogClientConfiguration.SetDefaults()

	// Attempt to overload default configuration with parsed json configuration from file
	OkuruLogClientConfiguration.LoadConfig()

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

	t, _ := tail.TailFile("/var/log/cron", tail.Config{Follow: true})
	for line := range t.Lines {
		fmt.Println(line.Text)
	}

	LogServer.connect()

}
