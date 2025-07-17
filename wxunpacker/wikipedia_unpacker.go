package main

import (
	//"compress/bzip2"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strings"
	"net"
	"net/url"
	"time"
	"github.com/vmihailenco/msgpack/v5"
	"common"
)

type Page struct {
	Title string `xml:"title"`
	Text string `xml:"revision>text"`
	Namespace string `xml:"ns"`
}

func main() {
	addr := "/tmp/windexIPC.sock"
	connection, err := net.Dial("unix", addr)
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	encoder := msgpack.NewEncoder(connection)

	//countPages("/run/media/matthewnesbitt/Linux 1TB SSD/WikiDump/enwiki-20250320-pages-articles-multistream.xml")
	file, err := os.Open("/run/media/matthewnesbitt/Linux 1TB SSD/WikiDump/enwiki-20250320-pages-articles-multistream.xml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//bz2_reader := bzip2.NewReader(file)
	decoder := xml.NewDecoder(file)

	var page Page
	var i = 0
	var diff = 0
	var wait int64 = 0

	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		switch t := tok.(type) {
		case xml.StartElement:
				if t.Name.Local == "page" {
					page = Page{}
					decoder.DecodeElement(&page, &t)
					if page.Namespace == "0" && !strings.HasPrefix(page.Text, "#REDIRECT") && page.Title != "" {
						if diff >= 1000 {
							fmt.Println("wxunpacker: Processed:", i)
							fmt.Println("wxunpacker: avg send wait time:", wait / 1000)
							diff = 0
							wait = 0
						}
						diff++
						i++

						url_title := strings.ReplaceAll(page.Title, " ", "_")
						url_data := "https://en.wikipedia.org/wiki/" + url.PathEscape(url_title)
						data := common.PageData{URL: url_data, Body: page.Text}
						start := time.Now()
						err := encoder.Encode(data)
						wait = time.Since(start).Microseconds() + wait
						if err != nil {
							panic(err)
						}
						fmt.Println(page.Text)
						return
					}
				}
		default:
		}
	}
}

func countPages(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	decoder := xml.NewDecoder(file)

	count := 0
	diff := 0
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		} else if err != nil {
			return count, err
		}

		if se, ok := tok.(xml.StartElement); ok && se.Name.Local == "page" {
			if diff >= 1000 {
				fmt.Println("Preprocessed:", count)
				diff = 0
			}
			diff++
			count++
		}
	}

	return count, nil
}

