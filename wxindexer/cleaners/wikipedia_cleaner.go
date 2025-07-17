package cleaners

import (
	"regexp"
	"strings"
	"fmt"
	"net/http"
	"net/url"
	"encoding/json"
	"wxindexer/containers"
)

var (
	reLinkExtract        = regexp.MustCompile(`\[\[([^\|\]]+)`)
	reRefTag             = regexp.MustCompile(`(?s)<ref[^>]*?>.*?</ref>`)
	reSelfClosingRef     = regexp.MustCompile(`(?s)<ref[^>]*/>`)
	reTemplate           = regexp.MustCompile(`(?s)\{\{.*?\}\}`)
	reTable              = regexp.MustCompile(`(?s)\{\|.*?\|\}`)
	reFileLink           = regexp.MustCompile(`\[\[File:[^\]]*\]\]`)
	reImageLink          = regexp.MustCompile(`\[\[Image:[^\]]*\]\]`)
	reCategory           = regexp.MustCompile(`\[\[Category:[^\]]*\]\]`)
	reInternalLink       = regexp.MustCompile(`\[\[([^\|\]]*\|)?([^\]]+)\]\]`)
	reExternalLink       = regexp.MustCompile(`\[(https?://[^\s\]]+)(\s+[^\]]+)?\]`)
	reHTMLComment        = regexp.MustCompile(`(?s)<!--.*?-->`)
	reHTMLTag            = regexp.MustCompile(`</?[a-zA-Z]+.*?>`)
	reBoldItalic         = regexp.MustCompile(`'''''(.*?)'''''`)
	reBold               = regexp.MustCompile(`'''(.*?)'''`)
	reItalic             = regexp.MustCompile(`''(.*?)''`)
	reQuotes             = regexp.MustCompile(`"(.*?)"`)
	reNonAlphanumeric    = regexp.MustCompile(`[^a-zA-Z0-9\s]+`)
	reExtraWhitespace    = regexp.MustCompile(`[ \t]+`)
	reWhitespaceLines    = regexp.MustCompile(`(?m)^[ \t\r\f\v]+$`)
	reMultipleNewlines   = regexp.MustCompile(`\n{2,}`)
	invalidPrefixes      = get_invalid_namespaces()
	linkPrefix 			 = "https://en.wikipedia.org/wiki/"
)

type WikipediaCleaner struct {}

func NewWikipediaCleaner() Cleaner {
	var c WikipediaCleaner
	return &c
}

func (v *WikipediaCleaner) Clean(text string) containers.Doc {
	var links []string
	linkSet := make(map[string]bool)

	entity_replacements := map[string]string{
		"&nbsp;": " ", "&amp;": " ", "&lt;": " ", "&gt;": " ", "&quot;": "",
	}

	// find and save all links
	matches := reLinkExtract.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		link := strings.TrimSpace(match[1])
		parts := strings.Split(link, ":")
		if link != "" && !linkSet[link] && (len(parts) <= 1 || !invalidPrefixes.Contains(parts[0])){
			link = strings.ReplaceAll(link, " ", "_")
			linkSet[link] = true
			link = linkPrefix + url.PathEscape(link)
			links = append(links, link)
		}
	}

	// remove metadata junk
	text = reRefTag.ReplaceAllString(text, "")
	text = reSelfClosingRef.ReplaceAllString(text, "")
	text = reTemplate.ReplaceAllString(text, "")
	text = reTable.ReplaceAllString(text, "")
	text = reFileLink.ReplaceAllString(text, "")
	text = reImageLink.ReplaceAllString(text, "")
	text = reCategory.ReplaceAllString(text, "")
	text = reInternalLink.ReplaceAllString(text, "$2")
	text = reExternalLink.ReplaceAllString(text, "$2")
	text = reHTMLComment.ReplaceAllString(text, "")
	text = reHTMLTag.ReplaceAllString(text, "")

	// remove formatting
	text = reBoldItalic.ReplaceAllString(text, "$1")
	text = reBold.ReplaceAllString(text, "$1")
	text = reItalic.ReplaceAllString(text, "$1")
	text = reQuotes.ReplaceAllString(text, "$1")

	// Remove unwanted HTML entities
	for k, v := range entity_replacements {
		text = strings.ReplaceAll(text, k, v)
	}

	// Remove any remaining non-alphanumeric characters
	text = reNonAlphanumeric.ReplaceAllString(text, "")

	// Remove excessive whitespace
	text = reExtraWhitespace.ReplaceAllString(text, " ")
	text = reWhitespaceLines.ReplaceAllString(text, "")
	text = reMultipleNewlines.ReplaceAllString(text, "\n\n")

	// Lowercase everything
	text = strings.ToLower(text)

	return containers.Doc{Body: strings.TrimSpace(text), Links: links}
}

func get_invalid_namespaces() *containers.Set {
	url := "https://en.wikipedia.org/w/api.php?action=query&meta=siteinfo&siprop=namespaces&format=json"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		panic(err)
	}

	query, exists := data["query"].(map[string]any)
	if !exists {
		panic(fmt.Errorf("Expected 'query' key in namespaces JSON: %s", data))
	}

	namespaces, exists := query["namespaces"].(map[string]any)
	if !exists {
		panic(fmt.Errorf("Expected 'namespaces' key in namespaces JSON: %s", query))
	}

	invalid_namespaces := containers.NewSet()
	for _, namespace := range namespaces {
		if nsMap, ok := namespace.(map[string]any); ok {
			if name, exists := nsMap["*"]; exists && name != "" {
				invalid_namespaces.Add(name.(string))
				fmt.Printf("Found invalid namespace: %s\n", name)
			}
		}
	}

	return invalid_namespaces
}
