# opslevel-go [![CircleCI](https://circleci.com/gh/segmentio/opslevel-go.svg?style=shield)](https://circleci.com/gh/segmentio/opslevel-go) [![Go Report Card](https://goreportcard.com/badge/github.com/segmentio/opslevel-go)](https://goreportcard.com/report/github.com/segmentio/opslevel-go) [![GoDoc](https://godoc.org/github.com/segmentio/opslevel-go?status.svg)](https://godoc.org/github.com/segmentio/opslevel-go)

`opslevel-go` is a client library for the [OpsLevel](https://www.opslevel.com/) integrations API

To get started create a new client:

```
client := opslevel.NewClient()
```

## Deploys Integration

The Deploys Integration requires the following fields:

```
deployRequest := opslevel.DeployRequest{
    Service:     "my_service",
    Description: "my_service was deployed",
    Deployer: rest.Deployer{
        Email: "deployer@myapp.com",
    },
    Environment: "env",
    DeployedAt:  time.Now(),
}
err := client.Deploy(deployRequest, "my-integration-uuid")
```

For a full list fields, see the docs.

## Checks Integration

The Deploys Integration requires the following fields:

```
checkRequest := CheckRequest{
    Service: "my_service",
    Check:   "my_check",
    Message: "Checks passed",
    Status:  "passed",
}
err := client.Check(checkRequest, "my-integration-uuid")
```

The `Message` field is optional and `Status` should be one of `passed` or `failed`.
