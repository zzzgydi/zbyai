package model

import "time"

type AnswerStatus int

const (
	AnswerInit  AnswerStatus = 0
	AnswerWork  AnswerStatus = 1
	AnswerDone  AnswerStatus = 2
	AnswerError AnswerStatus = 3
)

const (
	AnswerKeyMain string = "main"
)

type ThreadAnswer struct {
	Id        string       `json:"id" db:"id" gorm:"primary_key"`
	Key       string       `json:"key" db:"key" gorm:"index"`
	Status    AnswerStatus `json:"status" db:"status" gorm:"index"`
	Model     string       `json:"model" db:"model"`
	Content   string       `json:"content" db:"content"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
}

func (t ThreadAnswer) TableName() string {
	return "thread_answer"
}
