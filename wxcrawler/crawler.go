package main

import (
	"fmt"
	"strings"

	"wxcrawler/validators"
	"wxcrawler/containers"

	"github.com/gocolly/colly/v2"
)

func scrape(url string, validator validators.Validator) (Result, error) {
	var links = containers.NewSet()
	var textContent strings.Builder
	collector := colly.NewCollector()

	collector.OnHTML("h1", func(e *colly.HTMLElement) {
		textContent.WriteString(e.Text)
	})

	collector.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if validator.Validate(link) {
			links.Add(link)
		}
	})

	collector.OnHTML("p", func(e *colly.HTMLElement) {
		textContent.WriteString(e.Text)
	})

	fmt.Printf("Visiting: %s\n", url)
	err := collector.Visit(url)
	if err != nil {
		return Result{}, fmt.Errorf("Error vising the site: %s: %s\n", url, err)
	}

	return Result {
		FromURL: url,
		Links: links,
		Text: textContent.String(),
	}, nil
}
