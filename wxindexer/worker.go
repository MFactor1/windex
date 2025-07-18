package main

import (
	"strings"
	"wxindexer/containers"
	"wxindexer/cleaners"
	"context"
	"common"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func index(
	page common.PageData,
	cleaner cleaners.Cleaner,
	stopwords *containers.Set,
	rdb *redis.Client,
) containers.PageTF {

	// Clean raw text
	data := cleaner.Clean(page.Body)

	// Tokenize
	words := strings.Split(data.Body, " ")

	// Index
	frequencies := make(map[string]int)
	var word_count int64 = 0

	for _, word := range words {
		if !stopwords.Contains(word) {
			frequencies[word]++
			word_count++
		}
	}
	flushToRedis(rdb, frequencies)
	return containers.PageTF{Title: page.Title, URL: page.URL, Links: data.Links, Words: frequencies}
}

func flushToRedis(rdb *redis.Client, wordCounts map[string]int) error {
	_, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for word := range wordCounts {
			pipe.HIncrBy(ctx, "df_map", word, 1)
		}
		pipe.Incr(ctx, "total_pages")
		return nil
	})
	return err
}
