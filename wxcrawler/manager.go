package main

import (
	"log"
	"fmt"
	"sync"
	"time"
	"container/heap"
	"context"
	"strings"

	"wxcrawler/containers"
	"wxcrawler/validators"
)

type Result struct {
	FromURL string
	Links *containers.Set
	Text string
}

func main() {
	startURLs := []string {
		"https://en.wikipedia.org/wiki/Web_crawler",
		"https://en.wikipedia.org/wiki/Presidency_of_John_Tyler",
		"https://en.wikipedia.org/wiki/Wintjiya_Napaltjarri",
		"https://en.wikipedia.org/wiki/Asiana_Airlines_Flight_214",
		"https://en.wikipedia.org/wiki/Ludwig_Ahgren",
		"https://en.wikipedia.org/wiki/14th_Dalai_Lama",
		"https://en.wikipedia.org/wiki/Butts_for_Tour_Buses",
	}

	var wp_vldr validators.Validator
	var err error
	wp_vldr, err = validators.NewWikipediaValidator()
	if err != nil {
		log.Fatalf("Failed to get wikipedia validator: %v\n", err)
	}

	var urlQueue = make(containers.PriorityQueue, len(startURLs))
	for i, url := range startURLs {
		urlQueue[i] = &containers.Item {
			Value: url,
			Priority: 1,
			Index: i,
		}
	}
	var seenLock sync.Mutex
	var queuePopLock sync.Mutex
	var items = make(map[string]*containers.Item)
	var seen = containers.NewSet()
	ctx, stop := context.WithCancel(context.Background())
	heap.Init(&urlQueue)

	const numWorkers = 7
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := range numWorkers {
		go func(ctx context.Context, workerID int) {
			defer wg.Done()
			for {
				select {
				case <- ctx.Done():
					return
				default:
					fmt.Printf("Num Done: %d\nQueue size: %d : %d\n", len(*seen), urlQueue.Len(), len(items))
					queuePopLock.Lock()
					if urlQueue.Len() == 0 {
						queuePopLock.Unlock()
						time.Sleep(500 * time.Millisecond)
						continue
					}
					seenLock.Lock()
					item := heap.Pop(&urlQueue).(*containers.Item)
					queuePopLock.Unlock()
					url := item.Value
					seen.Add(url)
					delete(items, url)
					seenLock.Unlock()

					result, err := scrape(url, wp_vldr)
					if err != nil {
						fmt.Print(fmt.Errorf("Unable to scrape site: %s: %s", url, err))
						if strings.HasSuffix(err.Error(), "Too Many Requests\n") {
							fmt.Println("!!Being rate limited!!")
							heap.Push(&urlQueue, item)
							items[url] = item
							delete(*seen, url)
						}
						time.Sleep(10 * time.Second)
						continue
					}

					for link := range *result.Links {
						seenLock.Lock()
						if seen.Contains(link) {
							seenLock.Unlock()
							continue
						} else if item, exists := items[link]; exists {
							urlQueue.Update(item, item.Value, item.Priority + 1)
						} else {
							item := &containers.Item {
								Value: link,
								Priority: 1,
							}
							heap.Push(&urlQueue, item)
							items[link] = item
						}
						seenLock.Unlock()
					}
				}

			}
		}(ctx, i)
	}

	go func() {
		time.Sleep(300 * time.Second)
		stop()
	}()

	wg.Wait()
}
