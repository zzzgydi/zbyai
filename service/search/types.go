package search

type Searcher interface {
	Engine() string
	Search(option *SearchOption) (*SearchResult, error)
}

type SearchOption struct {
	Query    string
	Count    int
	Country  string
	Language string
}

// union
type SearchResult struct {
	Items   []*SearchItem    `json:"items,omitempty"`
	Related []*SearchRelated `json:"related,omitempty"`
}
type SearchItem struct {
	Title     string           `json:"title"`
	Link      string           `json:"link"`
	Snippet   string           `json:"snippet"`
	Sitelinks []serperSiteLink `json:"sitelinks,omitempty"`
}
type SearchRelated struct {
	Title   string `json:"title,omitempty"`
	Link    string `json:"link,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

// serper
type serperResponse struct {
	KnowledgeGraph  map[string]any        `json:"knowledgeGraph,omitempty"`
	Organic         []serperOrganicResult `json:"organic,omitempty"`
	PeopleAlsoAsk   []serperPeopleAlsoAsk `json:"peopleAlsoAsk,omitempty"`
	RelatedSearches []serperRelatedSearch `json:"relatedSearches,omitempty"`
}

type serperOrganicResult struct {
	Title      string           `json:"title"`
	Link       string           `json:"link"`
	Snippet    string           `json:"snippet,omitempty"`
	Sitelinks  []serperSiteLink `json:"sitelinks,omitempty"`
	Attributes map[string]any   `json:"attributes,omitempty"` // 可选字段
}

type serperSiteLink struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

type serperPeopleAlsoAsk struct {
	Question string `json:"question,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
	Title    string `json:"title,omitempty"`
	Link     string `json:"link,omitempty"`
}

type serperRelatedSearch struct {
	Query string `json:"query"`
}

// bing
type bingResponse struct {
	WebPages        *bingWebPages `json:"webPages,omitempty"`
	RelatedSearches struct {
		ID    string `json:"id"`
		Value []struct {
			Text         string `json:"text,omitempty"`
			DisplayText  string `json:"displayText,omitempty"`
			WebSearchURL string `json:"webSearchUrl,omitempty"`
		} `json:"value,omitempty"`
	} `json:"relatedSearches,omitempty"`
}

type bingWebPages struct {
	WebSearchUrl string             `json:"webSearchUrl,omitempty"`
	Value        []*bingSearchValue `json:"value,omitempty"`
}

type bingSearchValue struct {
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Url        string `json:"url,omitempty"`
	DisplayUrl string `json:"displayUrl,omitempty"`
	Snippet    string `json:"snippet,omitempty"`
	Language   string `json:"language,omitempty"`
}

// searxng

type searxResponse struct {
	Results []searxResult `json:"results,omitempty"`
}

type searxResult struct {
	Title   string `json:"title,omitempty"`
	Url     string `json:"url,omitempty"`
	Content string `json:"content,omitempty"`
}
