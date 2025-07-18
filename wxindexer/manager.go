package main

import (
	"log"
	"net"
	"os"
	"io"
	"time"
	"bufio"
	"sync"

	"common"
	"wxindexer/cleaners"
	"wxindexer/containers"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/redis/go-redis/v9"
)

var (
	workers = 7
	reader_group sync.WaitGroup
	indexer_group sync.WaitGroup
	writer_group sync.WaitGroup
)

func main() {
	log.Println("wxindexer: initalizing cleaner")
	cleaner := cleaners.NewWikipediaCleaner()

	log.Println("wxindexer: initializing redis client")
	rdb := newRedisClient()

	log.Println("wxindexer: loading stopwords")
	stopwords, err := loadStopWords()
	if err != nil {
		panic(err)
	}

	addr := "/tmp/windexIPC.sock"
	os.Remove(addr)

	listener, err := net.Listen("unix", addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	log.Println("wxindexer: waiting for connection...")
	connection, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	log.Println("wxindexer: connection established")

	decoder := msgpack.NewDecoder(connection)

	index_chan := make(chan common.PageData, 1000)
	write_chan := make(chan containers.PageTF, 1000)

	reader_group.Add(1)
	writer_group.Add(1)
	indexer_group.Add(workers)

	go socketReader(decoder, index_chan)
	go jsonWriter(write_chan)


	for i := range workers {
		go indexer(i, cleaner, stopwords, rdb, index_chan, write_chan)
	}

	reader_group.Wait()
	close(index_chan)
	indexer_group.Wait()
	close(write_chan)
	writer_group.Wait()
}

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options {
		Addr: "localhost:6380",
	})
}

func loadStopWords() (*containers.Set, error) {
	file, err := os.Open("./data/stopwords")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stopwords := containers.NewSet()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stopwords.Add(line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return stopwords, nil
}

func socketReader(decoder *msgpack.Decoder, out_chan chan <- common.PageData) {
	var diff = 0
	var wait int64 = 0

	for {
		var page common.PageData
		start := time.Now()
		err := decoder.Decode(&page)
		wait = time.Since(start).Microseconds() + wait
		if diff > 1000 {
			log.Println("wxindexer/reader: avg recieve wait time:", wait / 1000)
			diff = 0
			wait = 0
		}
		diff++
		if err != nil {
			if err == io.EOF {
				log.Println("wxindexer/reader: connection closed by sender. Exiting.")
				reader_group.Done()
				return
			}
			log.Printf("wxindexer/reader: decoder error: %v", err)
		}
		out_chan <- page
	}
}

func indexer(
	id int,
	cleaner cleaners.Cleaner,
	stopwords *containers.Set,
	rdb *redis.Client,
	in_chan <- chan common.PageData,
	out_chan chan <- containers.PageTF) {

	var tf containers.PageTF
	for {
		if page, ok := <- in_chan; ok {
			tf = index(page, cleaner, stopwords, rdb)
			out_chan <- tf
		} else {
			log.Printf("wxindexer/indexer@%d: exiting\n", id)
			break
		}
	}
	indexer_group.Done()
}
