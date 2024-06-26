package search

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

var _ Searcher = (*SearXNG)(nil)

type SearXNG struct {
	endpoint   string
	httpClient *http.Client
}

func (s *SearXNG) Engine() string {
	return "searXNG"
}

func (s *SearXNG) Search(option *SearchOption) (*SearchResult, error) {
	parsedURL, err := url.Parse(s.endpoint + "/search")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("q", option.Query)
	params.Add("format", "json")
	parsedURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status: %d", resp.StatusCode)
	}

	jsonRes := &searxResponse{}
	if err := json.NewDecoder(resp.Body).Decode(jsonRes); err != nil {
		return nil, err
	}

	if jsonRes.Results == nil || len(jsonRes.Results) == 0 {
		return nil, fmt.Errorf("web pages is nil")
	}

	if option.Count > 0 && len(jsonRes.Results) > option.Count {
		jsonRes.Results = jsonRes.Results[:option.Count]
	}

	items := make([]*SearchItem, 0, len(jsonRes.Results))
	for _, v := range jsonRes.Results {
		items = append(items, &SearchItem{
			Title:   v.Title,
			Link:    v.Url,
			Snippet: v.Content,
		})
	}

	return &SearchResult{
		Items: items,
	}, nil
}
