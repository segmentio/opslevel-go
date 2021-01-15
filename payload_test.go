package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type myData struct {
	High   int `validate:"required" json:"high"`
	Medium int `validate:"required" json:"medium"`
	Low    int `validate:"required" json:"low"`
}

func TestPayload(t *testing.T) {
	t.Run("Doesn't return an error for a valid request", func(t *testing.T) {

		payloadRequest := PayloadRequest{
			Service: "my_service",
			Check:   "my_check",
			Data: myData{
				High:   1,
				Medium: 2,
				Low:    3,
			},
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Payload(payloadRequest, "uuid")
		assert.NoError(t, err)
	})

	t.Run("Returns a Bad Request error on a 422", func(t *testing.T) {
		payloadRequest := PayloadRequest{
			Service: "my_service",
			Check:   "my_check",
			Data: myData{
				High:   1,
				Medium: 2,
				Low:    3,
			},
		}

		body := `
		{
			"errors":[
				{"status":400,"title":"Check Error","detail":"param is missing or the value is empty: data"}
			]
		}
		`
		client, testServer := setupTest(400, body)
		defer func() { testServer.Close() }()

		err := client.Payload(payloadRequest, "uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 400")
	})

	t.Run("Returns a Something Went Wrong error on all other status codes", func(t *testing.T) {
		payloadRequest := PayloadRequest{
			Service: "my_service",
			Check:   "my_check",
			Data: myData{
				High:   1,
				Medium: 2,
				Low:    3,
			},
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

		err := client.Payload(payloadRequest, "uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "status 503")
	})

	t.Run("Returns an error for an invalid request", func(t *testing.T) {
		payloadRequest := PayloadRequest{
			Service: "my_service",
			Check:   "my_check",
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Payload(payloadRequest, "uuid")
		assert.EqualError(t,
			err,
			"Key: 'PayloadRequest.Data' Error:Field validation for 'Data' failed on the 'required' tag",
		)
	})

	t.Run("Returns an error for an invalid request on nested Data interface", func(t *testing.T) {
		payloadRequest := PayloadRequest{
			Service: "my_service",
			Check:   "my_check",
			Data: myData{
				Medium: 2,
				Low:    3,
			},
		}

		body := `
		{
			"result": "ok"
		}
		`
		client, testServer := setupTest(202, body)
		defer func() { testServer.Close() }()

		err := client.Payload(payloadRequest, "uuid")
		assert.EqualError(t,
			err,
			"Key: 'PayloadRequest.Data.High' Error:Field validation for 'High' failed on the 'required' tag",
		)
	})
}
