package util

import (
	"encoding/json"
	"net/http"
)

func ParseRequest(r *http.Request, v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		return err
	}

	return nil
}
