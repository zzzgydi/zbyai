package common

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	supa "github.com/nedpals/supabase-go"
	"github.com/zzzgydi/zbyai/common/config"
)

var Supabase *supa.Client

func InitSupabase() error {
	conf := &config.AppConf.Supabase
	if conf.Url == "" || conf.Key == "" {
		return fmt.Errorf("supabase conf error")
	}

	Supabase = supa.CreateClient(conf.Url, conf.Key)

	// 设置代理，暂时用gpt的代理设置
	gptConf := &config.AppConf.Gpt
	if gptConf.ProxyURL != "" {
		proxyUrl, err := url.Parse(gptConf.ProxyURL)
		if err != nil {
			return err
		}
		Supabase.HTTPClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		}
	}
	return nil
}

func init() {
	// initializer.Register("supabase", InitSupabase)
}

func VerifyToken(token string) (*supa.User, error) {
	if Supabase == nil {
		if err := InitSupabase(); err != nil {
			return nil, err
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	user, err := Supabase.Auth.User(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("verify error, %v", err)
	}
	return user, nil
}
