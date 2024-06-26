package search

// dockdockgo

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/zzzgydi/zbyai/common/utils"
)

var _ Searcher = (*DDG)(nil)

var (
	ErrNoGoodResult = errors.New("no good search results found")
	ErrAPIResponse  = errors.New("duckduckgo api responded with error")
)

type DDG struct {
	httpClient *http.Client
}

func (d *DDG) Engine() string {
	return "ddg"
}

func (d *DDG) Search(option *SearchOption) (*SearchResult, error) {
	queryURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(option.Query))

	req, err := http.NewRequest("GET", queryURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", utils.RandomUserAgent())

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status: %d", resp.StatusCode)
	}

	return parseDuckDuckGoResponse(resp, option.Count)
}

func parseDuckDuckGoResponse(resp *http.Response, maxCount int) (*SearchResult, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	results := []*SearchItem{}

	sel := doc.Find(".web-result")
	for i := range sel.Nodes {
		// Break loop once required amount of results are add
		if maxCount == len(results) {
			break
		}
		node := sel.Eq(i)
		titleNode := node.Find(".result__a")

		info := node.Find(".result__snippet").Text()
		title := titleNode.Text()
		link := ""

		// TODO 优化
		if len(titleNode.Nodes) > 0 && len(titleNode.Nodes[0].Attr) > 2 {
			link, _ = url.QueryUnescape(
				strings.TrimPrefix(
					titleNode.Nodes[0].Attr[2].Val,
					"/l/?kh=-1&uddg=",
				),
			)
		}
		if link == "" {
			continue
		}

		results = append(results, &SearchItem{title, link, info, []serperSiteLink{}})
	}

	if len(results) == 0 {
		return nil, ErrNoGoodResult
	}

	return &SearchResult{Items: results}, nil
}
