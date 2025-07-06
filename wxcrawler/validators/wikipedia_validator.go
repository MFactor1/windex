package validators

import (
	"strings"
	"fmt"
	"net/http"
	"encoding/json"

	"wxcrawler/containers"
)

type WikipediaValidator struct {
	Invalid_namespaces *containers.Set
	Valid_prefix string
}

func NewWikipediaValidator() (Validator, error) {
	var v WikipediaValidator
	var err error
	v.Invalid_namespaces, err = get_invalid_namespaces()
	if err != nil {
		return nil, fmt.Errorf("Error getting invalid namespaces: %v", err)
	}
	v.Valid_prefix = "https://en.wikipedia.org/wiki/"
	return &v, nil
}

func (v *WikipediaValidator) Validate(link string) bool {
	if !strings.HasPrefix(link, v.Valid_prefix) {
		return false
	}

	page_name := strings.TrimPrefix(link, v.Valid_prefix)
	parts := strings.Split(page_name, ":")
	if len(parts) > 1 {
		if v.Invalid_namespaces.Contains(parts[0]) {
			return false
		}
	}

	return true
}

func get_invalid_namespaces() (*containers.Set, error) {
	url := "https://en.wikipedia.org/w/api.php?action=query&meta=siteinfo&siprop=namespaces&format=json"

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error fetching wikipedia namespaces: %v", err)
	}
	defer resp.Body.Close()

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("Error unmarshaling namespaces JSON: %v", err)
	}

	query, exists := data["query"].(map[string]any)
	if !exists {
		return nil, fmt.Errorf("Expected 'query' key in namespaces JSON: %v", err)
	}

	namespaces, exists := query["namespaces"].(map[string]any)
	if !exists {
		return nil, fmt.Errorf("Expected 'namespaces' key in namespaces JSON: %v", err)
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

	return invalid_namespaces, nil
}
