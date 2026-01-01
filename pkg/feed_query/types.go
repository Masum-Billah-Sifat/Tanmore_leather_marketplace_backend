package feedquery

import (
	"database/sql/driver"
	"fmt"
	"strings"
)

// Postgres TEXT[] scanner
type StringArray []string

func (a *StringArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("expected []byte for StringArray, got %T", value)
	}
	s := string(bytes)

	// Strip curly braces and split
	s = strings.Trim(s, "{}")
	if s == "" {
		*a = []string{}
		return nil
	}
	*a = strings.Split(s, ",")
	return nil
}

func (a StringArray) Value() (driver.Value, error) {
	return "{" + strings.Join(a, ",") + "}", nil
}
