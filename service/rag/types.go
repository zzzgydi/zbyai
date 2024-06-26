package rag

import (
	"github.com/sashabaranov/go-openai"
	"github.com/zzzgydi/zbyai/model"
)

type RAGChunk struct {
	Text      string
	Embedding openai.Embedding
	Score     float32
}

type RAGExec interface {
	Run() (model.SearchList, error)
	WaitContent() string
	WaitResult() model.SearchList
	WaitResultUnlock()
}
