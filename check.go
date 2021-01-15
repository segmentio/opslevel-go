package rest

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// CheckRequest represents a structured request to the OpsLevel checks webhook endpoint
type CheckRequest struct {
	Service string `validate:"required" json:"service"`
	Check   string `validate:"required" json:"check"`
	Status  string `validate:"required,oneof=passed failed" json:"status"`
	Message string `json:"message"`
}

// Check sends a CheckRequest to the OpsLevel check integration at integrationID
func (c *Client) Check(req CheckRequest, integrationID string) error {
	v := validator.New()
	if err := v.Struct(req); err != nil {
		return err
	}

	b, err := json.Marshal(req)
	if err != nil {
		return err
	}
	var resp struct {
		Result string `json:"result"`
	}

	fullURL := fmt.Sprintf("/integrations/check/%s", integrationID)
	return c.do("POST", fullURL, bytes.NewReader(b), &resp)
}
