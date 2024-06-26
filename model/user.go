package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lithammer/shortuuid/v4"
	"github.com/zzzgydi/zbyai/common"
)

type AuthType int

const (
	AUTH_NONE     AuthType = 0
	AUTH_SUPABASE AuthType = 1
)

type User struct {
	Id        string    `json:"id" db:"id" gorm:"primary_key"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	IP        string    `json:"ip" db:"ip"`
	AuthType  AuthType  `json:"auth_type" db:"auth_type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u User) TableName() string {
	return "user"
}

func NewTouristUser(ip string) *User {
	return &User{
		Id:        uuid.New().String(),
		Name:      "User-" + shortuuid.New()[0:6],
		Email:     "",
		IP:        ip,
		AuthType:  AUTH_NONE,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func NewSupabaseUser(id, name, email, ip string) *User {
	return &User{
		Id:        id,
		Name:      name,
		Email:     email,
		IP:        ip,
		AuthType:  AUTH_SUPABASE,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func GetUserById(id string) (*User, error) {
	user := &User{}
	err := common.MDB.Where("id = ?", id).First(user).Error
	return user, err
}

func (u *User) Create() error {
	if u.Id == "" {
		u.Id = uuid.New().String()
	}
	return common.MDB.Create(u).Error
}

func (u *User) Update() error {
	if u.Id == "" {
		return errors.New("user id is empty")
	}
	u.UpdatedAt = time.Now()
	return common.MDB.Save(u).Error
}
