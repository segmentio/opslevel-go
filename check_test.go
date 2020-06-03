package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheck(t *testing.T) {
	t.Run("Doesn't return an error for a valid request", func(t *testing.T) {
		checkRequest := CheckRequest{
			Service: "my_service",
			Check:   "my_check",
			Message: "Deployed service",
			Status:  "passed",
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Check(checkRequest, "uuid")
		assert.NoError(t, err)
	})

	t.Run("Returns a Bad Request error on a 422", func(t *testing.T) {
		checkRequest := CheckRequest{
			Service: "my_service",
			Check:   "my_check",
			Message: "Deployed service",
			Status:  "passed",
		}

		body := `
		{
			"errors":[
				{"status":400,"title":"Check Error","detail":"param is missing or the value is empty: status"}
			]
		}
		`
		client, testServer := setupTest(400, body)
		defer func() { testServer.Close() }()

		err := client.Check(checkRequest, "uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 400")
	})

	t.Run("Returns a Something Went Wrong error on all other status codes", func(t *testing.T) {
		checkRequest := CheckRequest{
			Service: "my_service",
			Check:   "my_check",
			Message: "Deployed service",
			Status:  "passed",
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

		err := client.Check(checkRequest, "uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 503")
	})

	t.Run("Returns an error for an invalid request", func(t *testing.T) {
		checkRequest := CheckRequest{
			Service: "my_service",
			Check:   "my_check",
			Message: "Deployed service",
			Status:  "bad status",
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Check(checkRequest, "uuid")
		assert.EqualError(t,
			err,
			"Key: 'CheckRequest.Status' Error:Field validation for 'Status' failed on the 'oneof' tag",
		)
	})
}
