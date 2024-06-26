package rag

import (
	"strings"

	"github.com/zzzgydi/zbyai/common/utils"
	"github.com/zzzgydi/zbyai/service/search"
)

var (
	mainSearch *utils.Chooser[string]
)

func init() {
	mainSearch = utils.NewChooser([]utils.Choice[string]{
		{Item: "searxng/searxng/ddg", Weight: 20},
		{Item: "bing/searxng/ddg", Weight: 20},
		{Item: "bing/ddg/searxng", Weight: 20},
		{Item: "serper/searxng", Weight: 2},
		{Item: "serper/bing", Weight: 2},
	})
}

func onceEngineList(size int) []search.Searcher {
	ret := []search.Searcher{}

	for {
		if len(ret) >= size {
			return ret[:size]
		}
		for _, eng := range strings.Split(*mainSearch.Pick(), "/") {
			switch eng {
			case "searxng":
				ret = append(ret, search.NewSearXNG())
			case "bing":
				ret = append(ret, search.NewBing())
			case "serper":
				ret = append(ret, search.NewSerper())
			case "ddg":
				ret = append(ret, search.NewDDG())
			default:
				ret = append(ret, search.NewSearXNG())
			}
		}
	}
}
