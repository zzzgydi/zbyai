package search

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

var _ Searcher = (*Bing)(nil)

type Bing struct {
	endpoint   string
	apiKey     string
	httpClient *http.Client
}

func (b *Bing) Engine() string {
	return "bing"
}

func (b *Bing) Search(option *SearchOption) (*SearchResult, error) {
	parsedURL, err := url.Parse(b.endpoint)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("q", option.Query)
	params.Add("count", fmt.Sprintf("%d", option.Count))
	params.Add("textDecorations", "false")
	params.Add("mkt", option.Country)
	params.Add("setLang", option.Language)
	params.Add("SafeSearch", "Moderate")
	parsedURL.RawQuery = params.Encode()

	req, err := http.NewRequest("GET", parsedURL.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Ocp-Apim-Subscription-Key", b.apiKey)

	resp, err := b.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("request status: %d, body: %s", resp.StatusCode, string(body))
	}

	jsonRes := &bingResponse{}
	err = json.NewDecoder(resp.Body).Decode(jsonRes)
	if err != nil {
		return nil, err
	}

	if jsonRes.WebPages == nil || jsonRes.WebPages.Value == nil {
		return nil, fmt.Errorf("web pages is nil")
	}

	items := make([]*SearchItem, 0, len(jsonRes.WebPages.Value))
	for _, v := range jsonRes.WebPages.Value {
		items = append(items, &SearchItem{
			Title:   v.Name,
			Link:    v.Url,
			Snippet: v.Snippet,
		})
	}
	related := make([]*SearchRelated, 0)
	for _, v := range jsonRes.RelatedSearches.Value {
		related = append(related, &SearchRelated{
			Title:   v.Text,
			Link:    v.WebSearchURL,
			Snippet: v.DisplayText,
		})
	}

	return &SearchResult{
		Items:   items,
		Related: related,
	}, nil
}
