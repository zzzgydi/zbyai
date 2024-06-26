package scrape_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
	"unicode/utf8"

	"github.com/go-shiori/go-readability"
	"github.com/zzzgydi/zbyai/service/scrape"
)

func TestScrapeRequest(t *testing.T) {
	scrape.Init()

	// 写一个测试用例，并行检测每一个url是否正常
	cases := []string{
		"https://www.aliyun.com/benefit/price/price_reduction?utm_content=m_1000390637",
		"https://www-cdn.anthropic.com/files/4zrzovbb/website/bd2a28d2535bfb0494cc8e2a3bf135d2e7523226.pdf",
		"https://unix.stackexchange.com/questions/148985/how-to-get-a-response-from-any-url",
		"https://zhuanlan.zhihu.com/p/645073085",
		"https://baike.baidu.com/item/%E7%BB%A7%E7%BB%AD/6186",
		"https://dictionary.cambridge.org/zhs/%E8%AF%8D%E5%85%B8/%E6%B1%89%E8%AF%AD-%E7%AE%80%E4%BD%93-%E8%8B%B1%E8%AF%AD/%E7%BB%A7%E7%BB%AD",
	}

	for _, url := range cases {
		t.Run(url, func(t *testing.T) {
			data, err := scrape.ScrapeByURL(url)
			if err != nil {
				fmt.Printf("[ERROR] URL: %s, error: %v\n", url, err)
				return
			}

			fmt.Printf("[INFO] URL: %s, length: %v\n", url, utf8.RuneCountInString(data))
		})
	}
}

func TestReadability(t *testing.T) {
	// 写一个测试用例，并行检测每一个url是否正常
	urls := []string{
		"https://www.aliyun.com/benefit/price/price_reduction?utm_content=m_1000390637",
		// "https://www-cdn.anthropic.com/files/4zrzovbb/website/bd2a28d2535bfb0494cc8e2a3bf135d2e7523226.pdf",
		// "https://unix.stackexchange.com/questions/148985/how-to-get-a-response-from-any-url",
		// "https://zhuanlan.zhihu.com/p/645073085",
		// "https://baike.baidu.com/item/%E7%BB%A7%E7%BB%AD/6186",
		// "https://dictionary.cambridge.org/zhs/%E8%AF%8D%E5%85%B8/%E6%B1%89%E8%AF%AD-%E7%AE%80%E4%BD%93-%E8%8B%B1%E8%AF%AD/%E7%BB%A7%E7%BB%AD",
	}
	for _, u := range urls {
		resp, err := http.Get(u)
		if err != nil {
			log.Fatalf("failed to download %s: %v\n", u, err)
		}
		defer resp.Body.Close()

		parsedURL, err := url.Parse(u)
		if err != nil {
			log.Fatalf("error parsing url")
		}

		article, err := readability.FromReader(resp.Body, parsedURL)
		if err != nil {
			log.Fatalf("failed to parse %s: %v\n", u, err)
		}

		// fmt.Printf("URL     : %s\n", u)
		// fmt.Printf("Title   : %s\n", article.Title)
		// fmt.Printf("Author  : %s\n", article.Byline)
		// fmt.Printf("Length  : %d\n", article.Length)
		// fmt.Printf("Excerpt : %s\n", article.Excerpt)
		// fmt.Printf("SiteName: %s\n", article.SiteName)
		// fmt.Printf("Image   : %s\n", article.Image)
		// fmt.Printf("Favicon : %s\n", article.Favicon)
		fmt.Printf("Content     : %s\n", article.Content)
	}
}
