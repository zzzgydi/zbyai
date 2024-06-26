package search

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var (
	httpClient   *http.Client
	searchConfig config.SearchConfig
)

func initSearch() error {
	searchConfig = config.AppConf.Search

	// check
	if searchConfig.Github == nil {
		return fmt.Errorf("search github config is not set")
	}
	if searchConfig.Bing == nil || searchConfig.Bing.ApiKey == "" {
		return fmt.Errorf("bing config is not set")
	}
	if searchConfig.Serper == nil || searchConfig.Serper.ApiKey == "" {
		return fmt.Errorf("serper config is not set")
	}
	if searchConfig.Searxng == nil || searchConfig.Searxng.Endpoint == "" {
		return fmt.Errorf("searxng config is not set")
	}

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          30,
			IdleConnTimeout:       3600 * time.Second,
			ResponseHeaderTimeout: 360 * time.Second,
			ExpectContinueTimeout: 360 * time.Second,
			DisableCompression:    false,
		},
	}

	// 设置代理
	if config.AppConf.Gpt.ProxyURL != "" {
		proxyUrl, err := url.Parse(config.AppConf.Gpt.ProxyURL)
		if err != nil {
			return err
		}
		httpClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	return nil
}

func init() {
	initializer.Register("search", initSearch)
}

func NewSearXNG() *SearXNG {
	return &SearXNG{
		endpoint:   searchConfig.Searxng.Endpoint,
		httpClient: httpClient,
	}
}

func NewSerper() *Serper {
	return &Serper{
		apiKey:     searchConfig.Serper.ApiKey,
		httpClient: httpClient,
	}
}

func NewBing() *Bing {
	return &Bing{
		endpoint:   searchConfig.Bing.Endpoint,
		apiKey:     searchConfig.Bing.ApiKey,
		httpClient: httpClient,
	}
}

func NewDDG() *DDG {
	return &DDG{
		httpClient: httpClient,
	}
}

func NewGithubCode() *Github {
	githubClient := github.NewClient(nil).WithAuthToken(searchConfig.Github.Tokens[0])
	return &Github{
		client: githubClient,
	}
}
