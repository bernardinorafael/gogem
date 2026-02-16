package dbutil

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JSONB map[string]string

func (j JSONB) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

func (j *JSONB) Scan(value any) error {
	if value == nil {
		*j = make(map[string]string)
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}
	return json.Unmarshal(bytes, j)
}
