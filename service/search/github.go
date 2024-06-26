package search

import (
	"context"

	"github.com/google/go-github/v60/github"
)

type Github struct {
	client *github.Client
}

func (g *Github) CodeSearch(query, language string) (*github.CodeSearchResult, error) {
	query = query + " language:" + language

	result, _, err := g.client.Search.Code(context.Background(), query, &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    1,
			PerPage: 5,
		},
	})

	return result, err
}
