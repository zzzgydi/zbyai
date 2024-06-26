package model

import "time"

type SearchCacheStatus int

const (
	SCStatusInit  = 0
	SCStatusDone  = 1
	SCStatusError = 2
)

type ThreadSearch struct {
	Id string `json:"id" db:"id" gorm:"primary_key"` // xid, objectid
	// Status    int       `json:"status" db:"status" gorm:"index"`
	Title     string    `json:"title" db:"title"`
	Link      string    `json:"link" db:"link"`
	Snippet   string    `json:"snippet" db:"snippet"`
	Page      string    `json:"page,omitempty" db:"page"`
	Token     int       `json:"token" db:"token" gorm:"index"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (s ThreadSearch) TableName() string {
	return "thread_search"
}
