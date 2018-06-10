package main

// Import needed packages
import (
	//	"fmt"
	"fmt"
	"github.com/ActiveState/tail"
	olConf "github.com/patyx7/okurulog/conf"
	"log"
	"net"
	"os"
	"strconv"
	"time"
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
	logServerConnection, err := net.DialTimeout("tcp", self.Hostname+":"+strconv.Itoa(self.Port), time.Duration(10)*time.Second)
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

func tailFileandSend() {

}

// Worker struct
type Worker struct {
	ID          int
	Work        chan WorkRequest
	WorkerQueue chan chan WorkRequest
	QuitChan    chan bool
}

// WorkRequest struct
type WorkRequest struct {
	Id int
}

// WorkerQueue struct
var WorkerQueue chan chan WorkRequest

// A buffered channel that we can send work requests on.
var WorkQueue = make(chan WorkRequest, 500)

func StartDispatcher(nworkers int) {
	// First, initialize the channel we are going to but the workers' work channels into.
	WorkerQueue = make(chan chan WorkRequest, nworkers)
	// Now, create all of our workers.
	for i := 0; i < nworkers; i++ {
		fmt.Println("Starting worker", i+1)
		worker := NewWorker(i+1, WorkerQueue)
		worker.Start()
	}
	go func() {
		for {
			select {
			case work := <-WorkQueue:
				fmt.Println("Received work requeust")
				go func() {
					worker := <-WorkerQueue
					fmt.Println("Dispatching work request")
					worker <- work
				}()
			}
		}
	}()
}

// NewWorker creates, and returns a new Worker object. Its only argument
// is a channel that the worker can add itself to whenever it is done its
// work.
func NewWorker(id int, workerQueue chan chan WorkRequest) Worker {
	// Create, and return the worker.
	worker := Worker{
		ID:          id,
		Work:        make(chan WorkRequest),
		WorkerQueue: workerQueue,
		QuitChan:    make(chan bool)}
	return worker
}

// This function "starts" the worker by starting a goroutine, that is
// an infinite "for-select" loop.
func (w Worker) Start() {
	go func() {
		for {
			// Add ourselves into the worker queue.
			w.WorkerQueue <- w.Work
			select {
			case work := <-w.Work:
				ingestJob(work.Id)
				// Receive a work request.
				fmt.Printf("worker%d: Received work request, for ingest job %d, this is cool\n", w.ID, work.Id)
				time.Sleep(500)
			case <-w.QuitChan:
				// We have been asked to stop.
				fmt.Printf("worker%d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop tells the worker to stop listening for work requests.
//
// Note that the worker will only stop *after* it has finished its work.
func (w Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}

func main() {

	// Default generic log handler or fatal
	LogFileHandler := _startGenericLogFileHandler()
	defer LogFileHandler.Close()
	log.SetOutput(LogFileHandler)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// Start up and ping mysql datastore or fatal
	datastore.mysql = _initMySQLDatastore()
	datastore.mysql.SetMaxIdleConns(0)   // this is the root problem! set it to 0 to remove all idle connections
	datastore.mysql.SetMaxOpenConns(500) // or whatever is appropriate for your setup.
	defer datastore.mysql.Close()

	// Starting Dispatcher
	StartDispatcher(10)

	for {

		time.Sleep(500 * time.Millisecond)

		// Get pending ingest jobs record from mysql datastore
		ingestJob, err := datastore.GetPendingIngestJob()
		if err != nil {
			log.Printf("Error trying to get a pending ingest job from mysql datastore:  %v", err)
			continue
		}

		// Only work if there is a job
		if ingestJob.Id != 0 {
			// Now, we take the ingestJob and make a WorkRequest out of them.
			work := WorkRequest{Id: ingestJob.Id}

			// Push the work onto the queue.
			WorkQueue <- work

			_, err = datastore.QueuePendingIngestJob(ingestJob)
			if err != nil {
				log.Printf("Error trying to update queued status of pending ingest job from mysql datastore:  %v", err)
				continue

			}

		} else {
			time.Sleep(5000 * time.Millisecond)
		}

	}

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

	conn, _ := LogServer.connect()

	seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
	t, _ := tail.TailFile("/var/log/cron", tail.Config{Follow: true, ReOpen: true, Location: seekInfo})

	tailFileandSend

	for line := range t.Lines {
		_, err = conn.Write([]byte(line.Text))
		if err != nil {
			log.Fatal("Error snding to LogServer, error return was: ", err)
			os.Exit(1)
		}
		fmt.Println(line.Text)
	}

}
