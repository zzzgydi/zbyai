package model

import (
	"errors"
	"time"

	"github.com/lithammer/shortuuid/v4"
	"gorm.io/gorm"
)

// 一个thread包含多个threadRun

type Thread struct {
	Id        string        `json:"id" db:"id" gorm:"primary_key"`
	UserId    string        `json:"user_id" db:"user_id" gorm:"index"`
	Visible   bool          `json:"visible" db:"visible" gorm:"index"`
	Title     string        `json:"title" db:"title"`
	Options   *ThreadOption `json:"options,omitempty" db:"options" gorm:"type:json"`
	CreatedAt time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt time.Time     `json:"updated_at" db:"updated_at"`
}

func (t Thread) TableName() string {
	return "thread"
}

func NewThread(userId, title string, visible bool) *Thread {
	return &Thread{
		Id:        shortuuid.New(),
		UserId:    userId,
		Visible:   visible,
		Title:     title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ForkThread(thread *Thread, userId string, visible bool) *Thread {
	return &Thread{
		Id:        shortuuid.New(),
		UserId:    userId,
		Visible:   visible,
		Title:     thread.Title,
		Options:   thread.Options,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (t *Thread) Create(db *gorm.DB) error {
	return db.Create(t).Error
}

func (t *Thread) AllRuns(db *gorm.DB) ([]*ThreadRun, error) {
	threadRuns := make([]*ThreadRun, 0)
	if err := db.Model(&ThreadRun{}).Where("thread_id = ?", t.Id).
		Order("id").Find(&threadRuns).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}
	return threadRuns, nil
}

type ThreadOption struct {
	Model string `json:"model"`
}
