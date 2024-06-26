package thread

import (
	"github.com/zzzgydi/zbyai/model"
)

type ThreadDetail struct {
	Id      string             `json:"id"`
	Title   string             `json:"title"`
	History []*model.ThreadRun `json:"history"`
	Current uint64             `json:"current"`
	Status  model.RunStatus    `json:"status"`
}
