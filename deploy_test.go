package rest

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestDeploy(t *testing.T) {
	t.Run("Doesn't return an error for a valid request", func(t *testing.T) {
		deployRequest := DeployRequest{
			Service:     "my_service",
			Description: "Deployed service",
			Deployer: Deployer{
				Email: "deployer@xyz.com",
			},
			Environment: "stage",
			DeployedAt:  time.Now(),
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Deploy(deployRequest, "uuid")
		assert.NoError(t, err)
	})

	t.Run("Returns a Bad Request error on a 422", func(t *testing.T) {
		deployRequest := DeployRequest{
			Service:     "my_service",
			Description: "Deployed service",
			Deployer: Deployer{
				Email: "deployer@xyz.com",
			},
			Environment: "stage",
			DeployedAt:  time.Now(),
		}

		body := `
		{
			"errors":[
				{"status":400,"title":"Deployed at Error","detail":"param is missing or the value is empty: deployed_at"}
			]
		}
		`
		client, testServer := setupTest(400, body)
		defer func() { testServer.Close() }()

		err := client.Deploy(deployRequest, "uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 400")
	})

	t.Run("Returns a Something Went Wrong error on all other status codes", func(t *testing.T) {
		deployRequest := DeployRequest{
			Service:     "my_service",
			Description: "Deployed service",
			Deployer: Deployer{
				Email: "deployer@xyz.com",
			},
			Environment: "stage",
			DeployedAt:  time.Now(),
		}

		body := `
		{
			"errors":[
				{"status":503,"title":"Service Unavailable"}
			]
		}
		`
		client, testServer := setupTest(503, body)
		defer func() { testServer.Close() }()

		err := client.Deploy(deployRequest, "uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 503")
	})

	t.Run("Returns an error for an invalid request", func(t *testing.T) {
		deployRequest := DeployRequest{
			Service:     "my_service",
			Description: "Deployed service",
			Deployer: Deployer{
				Email: "deployer@xyz.com",
			},
			DeployedAt: time.Now(),
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Deploy(deployRequest, "uuid")
		assert.EqualError(t,
			err,
			"Key: 'DeployRequest.Environment' Error:Field validation for 'Environment' failed on the 'required' tag",
		)
	})
}

func setupTest(statusCode int, body string) (*Client, *httptest.Server) {
	client := NewClient()
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(statusCode)
		res.Write([]byte(body))
	}))

	testServerURL, _ := url.Parse(testServer.URL)
	client.baseURL = testServerURL
	log.SetLevel(log.DebugLevel)

	return client, testServer
}
