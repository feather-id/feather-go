# Feather Go Library

[![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/feather-id/feather-go) [![Build Status](https://travis-ci.org/feather-id/feather-go.svg?branch=master)](https://travis-ci.org/feather-id/feather-go) [![Coverage Status](https://coveralls.io/repos/github/feather-id/feather-go/badge.svg?branch=master)](https://coveralls.io/github/feather-id/feather-go?branch=master)

This library provides a convenient interface to the Feather API for applications running in a Go server environment.

## Installation

```sh
$ go get -u github.com/feather-id/feather-go
```

## Usage

Initialize a Feather client with your project's API key, available on the [Feather Dashboard](https://feather.id/dashboard).

```go
import "github.com/feather-id/feather-go"

client := feather.New("live_...")
```

The example below will walk through a simple and common authentication flow:

- Create a credential given a username and password from your user.
- Create an authenticated session with the credential token.
- Add custom metadata to the user to save their current state.

Note this example ignores errors for brevity. You should not ignore errors in production code!

```go
// 1) Create a credential
credential, _ := client.Credentials.Create(feather.CredentialsCreateParams{
	Type:     feather.CredentialTypeUsernamePassword,
	Username: feather.String("jdoe"),
	Password: feather.String("pa$$w0rd"),
})

// Inform the user of their credential status
switch credential.Status {
case feather.CredentialStatusRequiresOneTimeCode:
	log.Printf("Please check your email for a link to sign in")
	return

case feather.CredentialStatusInvalid:
	log.Printf("Your username and password did not match")
	return

case feather.CredentialStatusValid:
	// The username and password were valid!
	break
}

// 2) Create an authenticated session
session, _ := client.Sessions.Create(feather.SessionsCreateParams{
  CredentialToken: credential.Token,
})

// 3) Add custom metadata to the user
user, _ := client.Users.Update(session.UserID, feather.UsersUpdateParams{
	Metadata: &map[string]string{
		"highScore": "123",
	},
})

log.Printf("Your high score is: %v", user.Metadata["highScore"])
```

## Development

To run tests, simply call:

```sh
$ make test
```

## More Information

- [Feather Docs](https://feather.id/docs)
- [API Reference](https://feather.id/docs/reference/api)
- [Error Handling](https://feather.id/docs/reference/api#errors)
