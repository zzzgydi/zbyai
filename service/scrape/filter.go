package scrape

import "net/url"

var (
	NotSupportHostname = map[string]bool{
		"music.163.com": true,
	}
)

// 如果为true，表示url是可以爬取的
func filterUrl(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		// 不是标准的url，就过滤掉
		return false
	}

	// 过滤掉一些不需要的域名
	if _, ok := NotSupportHostname[u.Host]; ok {
		return false
	}

	return true
}
