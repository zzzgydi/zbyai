package scrape

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
	"github.com/zzzgydi/zbyai/common/utils"
	"golang.org/x/net/html/charset"
)

func ScrapeByURL(rawUrl string) (string, error) {
	if !filterUrl(rawUrl) {
		return "", fmt.Errorf("not support url: %s", rawUrl)
	}

	req, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return "", err
	}

	// 添加请求头
	headers := map[string]string{
		"Accept":        "text/html;q=0.9, application/xhtml+xml;q=0.8",
		"Cache-Control": "no-cache",
		"Connection":    "keep-alive",
		"Pragma":        "no-cache",
		"User-Agent":    utils.RandomUserAgent(),
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := scrapeClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request status: %d, url: %s", resp.StatusCode, rawUrl)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("invalid Content-Type: %s", contentType)
	}

	reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		return "", err
	}

	domain := ""
	u, err := url.Parse(rawUrl)
	if err == nil {
		domain = u.Scheme + "://" + u.Host
	}

	options := &md.Options{
		GetAbsoluteURL: func(selec *goquery.Selection, rawURL, _ string) string {
			// 如果是相对路径，拼接成绝对路径
			if strings.HasPrefix(rawURL, "/") {
				return domain + rawURL
			}
			return rawURL
		},
	}

	converter := md.NewConverter("", true, options)
	markdown, err := converter.ConvertReader(reader)
	if err != nil {
		return "", err
	}

	return utils.RemoveMarkdownImages(markdown.String()), nil
}

type Scrape struct {
	rewiseDomain bool
	pipeline     []func(string) string
}

func NewScrape(rewiseDomain bool) *Scrape {
	return &Scrape{
		rewiseDomain: rewiseDomain,
	}
}

func (s *Scrape) AddPipeline(fn func(string) string) {
	s.pipeline = append(s.pipeline, fn)
}

func (s *Scrape) request(ctx context.Context, rawUrl string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", rawUrl, nil)
	if err != nil {
		return "", err
	}

	headers := map[string]string{
		"Accept":        "text/html;q=0.9, application/xhtml+xml;q=0.8",
		"Cache-Control": "no-cache",
		"Connection":    "keep-alive",
		"Pragma":        "no-cache",
		"User-Agent":    utils.RandomUserAgent(),
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := scrapeClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request status: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "text/html") {
		return "", fmt.Errorf("invalid Content-Type: %s", contentType)
	}

	reader, err := charset.NewReader(resp.Body, contentType)
	if err != nil {
		return "", err
	}

	domain := ""
	u, err := url.Parse(rawUrl)
	if err == nil {
		domain = u.Scheme + "://" + u.Host
	}

	options := &md.Options{
		GetAbsoluteURL: func(selec *goquery.Selection, rawURL, _ string) string {
			if !s.rewiseDomain {
				return rawURL
			}
			// 如果是相对路径，拼接成绝对路径
			if strings.HasPrefix(rawURL, "/") {
				return domain + rawURL
			}
			return rawURL
		},
	}

	article, err := readability.FromReader(reader, u)
	if err != nil {
		return "", err
	}

	converter := md.NewConverter("", true, options)
	markdown, err := converter.ConvertString(article.Content)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func (s *Scrape) Run(ctx context.Context, rawUrl string) (string, error) {
	ret, err := s.request(ctx, rawUrl)
	if err != nil {
		return "", err
	}

	for _, fn := range s.pipeline {
		ret = fn(ret)
	}

	return ret, nil
}
