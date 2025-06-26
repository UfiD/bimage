package types

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type PostObjectHandlerRequest struct {
	Task     string `json:"task"`
	Language string `json:"language"`
}

func GetObjectPostRequest(r *http.Request) (*PostObjectHandlerRequest, error) {
	var req PostObjectHandlerRequest
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decode json: %v", err)
	}
	return &req, nil
}
