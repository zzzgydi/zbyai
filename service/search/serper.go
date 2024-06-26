package search

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var _ Searcher = (*Serper)(nil)

type Serper struct {
	apiKey     string
	httpClient *http.Client
}

func (s *Serper) Engine() string {
	return "serper"
}

func (s *Serper) Search(option *SearchOption) (*SearchResult, error) {
	url := "https://google.serper.dev/search"
	reqData := map[string]any{
		"q":    option.Query,
		"page": 1,
	}
	if option.Country != "" {
		reqData["gl"] = option.Country
	}
	if option.Language != "" {
		reqData["hl"] = option.Language
	}

	jsonData, err := json.Marshal(reqData)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-API-KEY", s.apiKey)
	req.Header.Add("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request status: %d, body: %s", resp.StatusCode, string(body))
	}

	jsonRes := &serperResponse{}
	err = json.NewDecoder(resp.Body).Decode(jsonRes)
	if err != nil {
		return nil, err
	}

	if jsonRes.Organic == nil || len(jsonRes.Organic) == 0 {
		return nil, fmt.Errorf("web pages is nil")
	}

	items := make([]*SearchItem, 0, len(jsonRes.Organic))
	for _, v := range jsonRes.Organic {
		items = append(items, &SearchItem{
			Title:     v.Title,
			Link:      v.Link,
			Snippet:   v.Snippet,
			Sitelinks: v.Sitelinks,
		})
	}

	related := make([]*SearchRelated, 0)
	for _, v := range jsonRes.PeopleAlsoAsk {
		related = append(related, &SearchRelated{
			Title:   v.Title,
			Link:    v.Link,
			Snippet: v.Snippet,
		})
	}

	return &SearchResult{
		Items:   items,
		Related: related,
	}, nil
}
