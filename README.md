# About

I'm a Golang newbie and this project is my take on [the Form3 take-home interview test](https://github.com/form3tech-oss/interview-accountapi).

Most commands in this README file assume that the working directory is the root project directory.

# Documentation

Apart from this README there is also some API documentation in [godoc](https://pkg.go.dev/golang.org/x/tools/cmd/godoc) format.

To read it install godoc:

    go install -v golang.org/x/tools/cmd/godoc@latest

And then run the HTTP documentation server:

    godoc

After the server starts the documentation can be found at [this link](http://localhost:6060/pkg/github.com/jannis-baratheon/form3-take-home-exercise/).

# Getting started

## Creating Form3 API client

```go
import (
    "net/http"

    "github.com/jannis-baratheon/form3-take-home-exercise/form3apiclient"
)

// ...

var apiURL string // API URL
var httpClient *http.Client // the HTTP client to be used for communication with Form3 API

// sample values
/*
    apiURL = "http://localhost:8080/v1"
    httpClient = &http.Client{}
*/

// ...

client := form3apiclient.NewForm3APIClient(apiURL, httpClient)

// ...
```

## Creating an account

```go
import (
    "github.com/jannis-baratheon/form3-take-home-exercise/form3apiclient"
)

// ...

var accountData form3apiclient.AccountData // the resource to be created

// sample value
/*
    accountData = form3apiclient.AccountData{
        ID:             uuid.NewString(),
        OrganisationID: uuid.NewString(),
        Type:           "accounts",
        Attributes: form3apiclient.AccountAttributes{
            AccountClassification: "Personal",
            Name:                  []string{"Jan Kowalski", "Jasiu Kowalski"},
            Country:               "PL",
        },
    }
*/

// ...

resp, err := client.Accounts().Create(context.Background(), accountData)

// ...
```

## Fetching an account

```go
import (
    "github.com/jannis-baratheon/form3-take-home-exercise/form3apiclient"
)

// ...

var resourceID string // the ID of the resource to fetch (e.g. "3e93cb04-7d07-11ec-90d6-0242ac120003")

// ...

accountData, err := client.Accounts().Get(context.Background(), resourceID)

// ...
```

## Deleting an account

```go
import (
    "github.com/jannis-baratheon/form3-take-home-exercise/form3apiclient"
)

// ...

var resourceID string // the ID of the resource to be deleted (e.g. "3e93cb04-7d07-11ec-90d6-0242ac120003")
var resourceVersion int64 // concurrency control (optimistic locking) - the version number of the resource to be deleted (e.g. 0)

// ...

err := client.Accounts().Delete(context.Background(), resourceID, resourceVersion)

// ...
```

# Static analysis

The project uses [golangci-lint](https://golangci-lint.run) for [static analysis](https://en.wikipedia.org/wiki/Static_program_analysis).

To run code analysis use this command after [installing golangci-lint](https://golangci-lint.run/usage/install/) on your machine:

    golang-ci run

You can find golangci-lint configuration for this project [here](.golangci.yml).

# Tests

The library uses [Ginkgo](https://onsi.github.io/ginkgo/) test framework and [Gomega](https://onsi.github.io/gomega/) assertion/matcher library. Please refer to [Ginkgo documentation](https://onsi.github.io/ginkgo/) on how to get started.

## Non-E2E Integration/Unit tests

To run non-E2E tests use this command line:

    ginkgo --label-filter="!e2e" -r

## E2E tests

Apart from the usual integration/unit tests there are also E2E tests that require a working Form3 API environment.

### Environment

You can use [docker-compose](https://docs.docker.com/compose/) and the provided [docker-compose YML file](docker/docker-compose.yml) to start a test enviroment locally:

    docker-compose -f docker/docker-compose.yml up -d

The enviroment will be accessible on port 8080.

### Running the tests

Command-line for executing E2E tests:

    FORM3_API_URL=<put environment URL here> ginkgo --label-filter="e2e" -r

The `FORM3_API_URL` environment variable has to point to a functional environment. The URL is `http://localhost:8080/v1` in case of the aforementioned local docker-compose environement.

# Continous integration

The project uses [Github Actions](https://github.com/features/actions) for CI. There are three pipelines run on every commit to the repository:

* [Static analysis](https://github.com/jannis-baratheon/form3-take-home-exercise/actions/workflows/static_analysis.yml) - self describing.
* [Tests](https://github.com/jannis-baratheon/form3-take-home-exercise/actions/workflows/test.yml) - runs the usual integration/unit tests.
* [E2E Tests](https://github.com/jannis-baratheon/form3-take-home-exercise/actions/workflows/e2e.yml) - runs E2E tests.

Code of the pipelines can be found in [.github/workflows](.github/workflows).

Please note that the pipelines can also be [run manually](https://docs.github.com/en/actions/managing-workflow-runs/manually-running-a-workflow).
