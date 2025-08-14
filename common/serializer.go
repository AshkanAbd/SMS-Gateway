package common

import (
	"encoding/json"
)

func JSONToValue[T any](str string, t *T) error {
	if err := json.Unmarshal([]byte(str), t); err != nil {
		return err
	}

	return nil
}

func ValueToJSON(t any) string {
	bytes, err := json.Marshal(t)
	if err != nil {
		return ""
	}

	return string(bytes)
}
