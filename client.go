package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	log "github.com/sirupsen/logrus"
)

// Client represents a rest http client and is used to send requests to OpsLevel integrations
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	Logger     *log.Logger
}

// NewClient returns a Client pointer
func NewClient() *Client {
	baseURL, _ := url.Parse("https://app.opslevel.com")
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		Logger:     log.StandardLogger(),
	}
}

func (c *Client) do(method string, path string, body io.Reader, recv interface{}) error {
	var err error
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	c.Logger.Debugf("Sending request to OpsLevel endpoint %s", url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.Logger.Debugf("Failed to send request to OpsLevel: %s", err.Error())
		return err
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	c.Logger.Debugf("Received status code %d", resp.StatusCode)
	if resp.StatusCode != http.StatusAccepted {
		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(resp.Body)
		s := buf.String()
		return fmt.Errorf("status %d; %s", resp.StatusCode, s)
	}

	err = decoder.Decode(&recv)
	if err != nil {
		c.Logger.Debugf("Failed to decode response from OpsLevel: %s", err.Error())
		return err
	}
	return nil
}
