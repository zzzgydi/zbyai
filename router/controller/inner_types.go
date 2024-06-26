package controller

type PostSearchRequest struct {
	Type     string `json:"type"`
	Query    string `json:"query"`
	Language string `json:"language"`
}

type PostThreadCreateRequest struct {
	Query string `json:"query"`
}
type PostThreadAppendRequest struct {
	Id    string `json:"id"`
	Query string `json:"query"`
}
type PostThreadRewriteRequest struct {
	Id    string `json:"id"`
	RunId uint64 `json:"runId"`
}

type PostThreadStreamRequest struct {
	Id    string `json:"id"`
	RunId uint64 `json:"runId"`
}
type PostThreadDetailRequest struct {
	Id string `json:"id"`
}
type PostThreadDeleteRequest struct {
	Id string `json:"id"`
}
