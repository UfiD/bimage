package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PostObjectHandlerRequest struct {
	Task string `json:"task"`
}

func GetObjectPostRequest(r *http.Request) (string, error) {
	var req PostObjectHandlerRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&req); err != nil {
		return "", fmt.Errorf("error while decode json: %v", err)
	}
	return req.Task, nil
}
