package utils

import (
	"encoding/json"
	"fmt"
	"strings"
)

func TryParseJson(str string, value any) error {
	l1 := strings.Index(str, "```json")
	if l1 == -1 {
		l1 = strings.Index(str, "```")
		if l1 == -1 {
			l1 = strings.Index(str, "{")
			if l1 == -1 {
				return fmt.Errorf("json prefix not found")
			}
		} else {
			l1 += 3
		}
	} else {
		l1 += 7
	}

	r1 := strings.LastIndex(str, "```")
	if r1 == -1 {
		r1 = strings.LastIndex(str, "}")
		if r1 == -1 {
			return fmt.Errorf("json suffix not found")
		} else {
			r1++
		}
	}
	if r1 > l1 {
		temp := str[l1:r1]
		if err := json.Unmarshal([]byte(temp), value); err == nil {
			return nil
		}
	}

	temp := str[l1:]
	r1 = strings.Index(temp, "```")
	if r1 == -1 {
		r1 = strings.Index(temp, "}")
		if r1 == -1 {
			return fmt.Errorf("json suffix not found")
		}
	}
	temp = temp[:r1]

	return json.Unmarshal([]byte(temp), value)
}
