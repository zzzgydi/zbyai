package llm

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/sashabaranov/go-openai"
	"github.com/zzzgydi/zbyai/common/config"
	"github.com/zzzgydi/zbyai/common/initializer"
)

var (
	gptConfig *openai.ClientConfig
)

func InitLLM() error {
	conf := &config.AppConf.Gpt
	if conf.ApiKey == "" || conf.Endpoint == "" {
		return fmt.Errorf("openai token is empty")
	}
	config := openai.DefaultConfig(conf.ApiKey)
	gptConfig = &config
	gptConfig.BaseURL = conf.Endpoint
	// 设置代理
	if conf.ProxyURL != "" {
		proxyUrl, err := url.Parse(conf.ProxyURL)
		if err != nil {
			return err
		}
		gptConfig.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	return nil
}

func init() {
	initializer.Register("llm", InitLLM)
}
