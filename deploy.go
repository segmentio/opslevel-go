package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"gopkg.in/go-playground/validator.v9"
)

// Deployer represents the entity taking the action
type Deployer struct {
	Email string `validate:"required" json:"email"`
	Name  string `json:"name"`
}

// Commit represents the commit being deployed
type Commit struct {
	SHA            string    `json:"sha"`
	Message        string    `json:"message"`
	Branch         string    `json:"branch"`
	Date           time.Time `json:"date"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	AuthorName     string    `json:"author_name"`
	AuthorEmail    string    `json:"author_email"`
	AuthoringDate  time.Time `json:"authoring_date"`
}

// DeployRequest represents a structured request to the OpsLevel deploys webhook endpoint
type DeployRequest struct {
	Service      string    `validate:"required" json:"service"`
	Deployer     Deployer  `validate:"required" json:"deployer"`
	DeployedAt   time.Time `validate:"required" json:"deployed_at"`
	Description  string    `validate:"required" json:"description"`
	Environment  string    `json:"environment"`
	DeployURL    string    `json:"deploy_url"`
	DeployNumber string    `json:"deploy_number"`
	Commit       Commit    `json:"commit"`
	DedupID      string    `json:"dedup_id"`
}

// Deploy sends a DeployRequest to the OpsLevel deploy integration at integrationID
func (c *Client) Deploy(req DeployRequest, integrationID string) error {
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

	fullURL := fmt.Sprintf("/integrations/deploy/%s", integrationID)
	return c.do("POST", fullURL, bytes.NewReader(b), &resp)
}
