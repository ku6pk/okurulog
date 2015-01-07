package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
)

type P struct {
	M, N int64
}

func handleConnection(conn net.Conn) {
	dec := gob.NewDecoder(conn)
	p := &P{}
	dec.Decode(p)
	fmt.Printf("Received : %+v", p)
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

	// define variables ahead of time
	var err error
	var LogFileHandler *os.File

	fmt.Println("start")
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept() // this blocks until connection or error
		if err != nil {
			// handle error
			continue
		}
		go handleConnection(conn) // a goroutine handles conn so that the loop can accept other connections
	}
}
