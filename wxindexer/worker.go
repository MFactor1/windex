package main

import (
	"fmt"
	"wxindexer/containers"
	"wxindexer/cleaners"
)

func index(body string, cleaner cleaners.Cleaner) containers.WordFrequencies {
	data := cleaner.Clean(body)
	for _, link := range(data.Links) {
		fmt.Println(link)
	}
	fmt.Printf("Cleaned Data: %s\n", data.Body)
	return containers.WordFrequencies{Words: make(map[string]int)}
}
