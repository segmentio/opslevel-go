package rest

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
)

// PayloadRequest represents a structured request to the OpsLevel payload webhook endpoint
type PayloadRequest struct {
	Service string      `validate:"required" json:"service"`
	Check   string      `validate:"required" json:"check"`
	Data    interface{} `validate:"required" json:"data"`
}

// Payload sends a PayloadRequest to the OpsLevel payload check integration at integrationID
func (c *Client) Payload(req PayloadRequest, integrationID string) error {
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

	fullURL := fmt.Sprintf("/integrations/payload/%s", integrationID)
	return c.do("POST", fullURL, bytes.NewReader(b), &resp)
}
