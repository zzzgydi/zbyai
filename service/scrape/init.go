package scrape

import (
	"net/http"
	"net/url"
	"time"

	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var (
	scrapeClient *http.Client
)

func Init() error {
	scrapeClient = &http.Client{
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
		scrapeClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}

	return nil
}

func init() {
	initializer.Register("scrape", Init)
}
