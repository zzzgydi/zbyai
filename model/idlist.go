package model

import (
	"database/sql/driver"
	"encoding/json"
)

type IdList []string

func (il IdList) GormDataType() string {
	return "json"
}

func (il IdList) Value() (driver.Value, error) {
	return json.Marshal(il)
}

func (il *IdList) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	return json.Unmarshal(value.([]byte), il)
}
