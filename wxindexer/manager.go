package main

import (
	"fmt"
	"net"
	"os"
	"log"
	"io"
	"time"
	"github.com/vmihailenco/msgpack/v5"
	"common"
	"wxindexer/cleaners"
)

func main() {
	addr := "/tmp/windexIPC.sock"
	os.Remove(addr)

	listener, err := net.Listen("unix", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("wxindexer: waiting for connection...")
	connection, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	fmt.Println("wxindexer: connection established")

	decoder := msgpack.NewDecoder(connection)

	var diff = 0
	var wait int64 = 0

	for {
		var page common.PageData
		start := time.Now()
		err := decoder.Decode(&page)
		wait = time.Since(start).Microseconds() + wait
		if diff > 1000 {
			fmt.Println("wxindexer: avg recieve wait time:", wait / 1000)
			diff = 0
			wait = 0
		}
		diff++
		if err != nil {
			if err == io.EOF {
				fmt.Println("wxindexer: connection closed by sender. Exiting.")
				break
			}
			log.Printf("wxindexer: decoder error: %v", err)
		}
		cleaner := cleaners.NewWikipediaCleaner()
		index(page.Body, cleaner)
		//fmt.Printf("Processing page: %s, len=%d\n", page.URL, len(page.Body))
	}
}
