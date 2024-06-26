package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

type RunStatus int

const (
	RunStatusInit  RunStatus = 0
	RunStatusWork  RunStatus = 1
	RunStatusDone  RunStatus = 2
	RunStatusError RunStatus = 3
)

type AnswerList []*ThreadAnswer
type SearchList []*ThreadSearch

type ThreadRun struct {
	Id        uint64            `json:"id" db:"id" gorm:"primary_key;autoIncrement"`
	ThreadId  string            `json:"thread_id" db:"thread_id" gorm:"index"`
	Status    RunStatus         `json:"status" db:"status" gorm:"index"`
	Query     string            `json:"query" db:"query"`
	AnswerIds IdList            `json:"-" db:"answer_ids" gorm:"type:json"`
	SearchIds IdList            `json:"-" db:"search_ids" gorm:"type:json"`
	Answer    AnswerList        `json:"answer,omitempty" db:"-" gorm:"-"`
	Search    SearchList        `json:"search,omitempty" db:"-" gorm:"-"`
	Setting   *ThreadRunSetting `json:"setting" db:"setting" gorm:"json"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`

	fetchAnswer bool `json:"-" db:"-" gorm:"-"`
	fetchSearch bool `json:"-" db:"-" gorm:"-"`
}

func (t ThreadRun) TableName() string {
	return "thread_run"
}

func NewThreadRun(threadId, query string) *ThreadRun {
	return &ThreadRun{
		ThreadId:  threadId,
		Status:    RunStatusInit,
		Query:     query,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *ThreadRun) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *ThreadRun) CopyTo(threadId string) *ThreadRun {
	return &ThreadRun{
		ThreadId:  threadId,
		Status:    t.Status,
		Query:     t.Query,
		AnswerIds: t.AnswerIds,
		SearchIds: t.SearchIds,
		Setting:   t.Setting,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *ThreadRun) PrefetchSearch(db *gorm.DB) error {
	if t.fetchSearch || t.SearchIds == nil || len(t.SearchIds) == 0 {
		return nil
	}
	t.fetchSearch = true
	searchList := make([]*ThreadSearch, 0)
	if err := db.Model(&ThreadSearch{}).Omit("page", "created_at").
		Where("id IN (?)", []string(t.SearchIds)).Order("id").Find(&searchList).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	t.Search = searchList
	return nil
}

func (t *ThreadRun) PrefetchAnswer(db *gorm.DB) error {
	if t.fetchAnswer || t.AnswerIds == nil || len(t.AnswerIds) == 0 {
		return nil
	}
	t.fetchAnswer = true
	answerList := make([]*ThreadAnswer, 0)
	if err := db.Model(&ThreadAnswer{}).Select("*").
		Where("id IN (?)", []string(t.AnswerIds)).Order("id").Find(&answerList).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	t.Answer = answerList
	return nil
}

type ThreadRunSetting struct {
	UseSearch     bool     `json:"use_search,omitempty"`
	Model         string   `json:"model,omitempty"`
	IsProgramming bool     `json:"is_programming,omitempty"`
	Language      string   `json:"language,omitempty"`
	QueryList     []string `json:"query_list,omitempty"`
	// Rephrased     string   `json:"rephrased,omitempty"`
}

func (trs ThreadRunSetting) GormDataType() string {
	return "json"
}

func (trs ThreadRunSetting) Value() (driver.Value, error) {
	return json.Marshal(trs)
}

func (trs *ThreadRunSetting) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), trs)
}
