package feather_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/feather-id/feather-go"
	"github.com/stretchr/testify/assert"
)

const (
	sampleAPIKey = "fooKey"
)

func createTestClient(server *httptest.Server) feather.Client {
	comps := strings.SplitN(strings.TrimPrefix(server.URL, "http://"), ":", 2)
	return feather.New(sampleAPIKey, &feather.Config{
		Protocol:   feather.String("http"),
		Host:       feather.String(comps[0]),
		Port:       feather.String(comps[1]),
		HTTPClient: server.Client(),
	})
}

// * * * * * Credentials * * * * * //

var sampleCredentialEmailRequiresOneTimeCode = feather.Credential{
	ID:        "CRD_foo",
	Object:    "credential",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	ExpiresAt: time.Date(2020, 01, 01, 01, 11, 01, 0, time.UTC),
	Status:    feather.CredentialStatusRequiresOneTimeCode,
	Token:     feather.String("qwerty"),
	Type:      feather.CredentialTypeEmail,
}

var sampleCredentialEmailValid = feather.Credential{
	ID:        "CRD_foo",
	Object:    "credential",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	ExpiresAt: time.Date(2020, 01, 01, 01, 11, 01, 0, time.UTC),
	Status:    feather.CredentialStatusValid,
	Token:     feather.String("qwerty"),
	Type:      feather.CredentialTypeEmail,
}

func TestCredentialsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.String(), "/v1/credentials")
		assert.Equal(t, r.FormValue("type"), "email")
		assert.Equal(t, r.FormValue("email"), "foo@bar.com")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleCredentialEmailRequiresOneTimeCode)
	}))
	defer server.Close()
	client := createTestClient(server)
	credential, err := client.Credentials.Create(feather.CredentialsCreateParams{
		Type:  feather.CredentialTypeEmail,
		Email: feather.String("foo@bar.com"),
	})
	assert.Equal(t, sampleCredentialEmailRequiresOneTimeCode, *credential)
	assert.Nil(t, err)
}

func TestCredentialsCreate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.String(), "/v1/credentials")
		assert.Equal(t, r.FormValue("type"), "email")
		assert.Equal(t, r.FormValue("email"), "foo@bar.com")
		assert.Equal(t, r.FormValue("username"), "foobar")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Code:    feather.ErrorCodeParameterInvalid,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	credential, err := client.Credentials.Create(feather.CredentialsCreateParams{
		Type:     feather.CredentialTypeEmail,
		Email:    feather.String("foo@bar.com"),
		Username: feather.String("foobar"),
	})
	assert.Nil(t, credential)
	assert.Equal(t, "An error message", err.Error())
}

func TestCredentialsUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.String(), "/v1/credentials/CRD_foo")
		assert.Equal(t, r.FormValue("one_time_code"), "foobar")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleCredentialEmailValid)
	}))
	defer server.Close()
	client := createTestClient(server)
	credential, err := client.Credentials.Update("CRD_foo", feather.CredentialsUpdateParams{
		OneTimeCode: feather.String("foobar"),
	})
	assert.Equal(t, sampleCredentialEmailValid, *credential)
	assert.Nil(t, err)
}

func TestCredentialsUpdate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.String(), "/v1/credentials/CRD_foo")
		assert.Equal(t, r.FormValue("one_time_code"), "")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Code:    feather.ErrorCodeParameterInvalid,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	credential, err := client.Credentials.Update("CRD_foo", feather.CredentialsUpdateParams{})
	assert.Nil(t, credential)
	assert.Equal(t, "An error message", err.Error())
}

// * * * * * Sessions * * * * * //

var sampleSessionAnonymous = feather.Session{
	ID:        "SES_foo",
	Object:    "session",
	Type:      feather.SessionTypeAnonymous,
	Status:    feather.SessionStatusActive,
	Token:     feather.String("qwerty"),
	UserID:    "USR_foo",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	RevokedAt: nil,
}

var sampleSessionAuthenticated = feather.Session{
	ID:        "SES_bar",
	Object:    "session",
	Type:      feather.SessionTypeAuthenticated,
	Status:    feather.SessionStatusRevoked,
	Token:     feather.String("qwerty"),
	UserID:    "USR_foo",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	RevokedAt: feather.Time(time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC)),
}

var sampleSessionList = feather.SessionList{
	ListMeta: feather.ListMeta{
		Objet:      "list",
		URL:        "/v1/sessions",
		HasMore:    false,
		TotalCount: 2,
	},
	Data: []*feather.Session{
		&sampleSessionAnonymous,
		&sampleSessionAuthenticated,
	},
}

func TestSessionsCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.String(), "/v1/sessions")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleSessionAnonymous)
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Create(feather.SessionsCreateParams{
		CredentialToken: feather.String("bar"),
	})
	assert.Equal(t, sampleSessionAnonymous, *session)
	assert.Nil(t, err)
}

func TestSessionsCreate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.String(), "/v1/sessions")
		assert.Equal(t, r.FormValue("credential_token"), "-1")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Code:    feather.ErrorCodeParameterInvalid,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Create(feather.SessionsCreateParams{
		CredentialToken: feather.String("-1"),
	})
	assert.Nil(t, session)
	assert.Equal(t, err.Error(), "An error message")
}

func TestSessionsList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/sessions?"))
		assert.Equal(t, r.URL.Query().Get("user_id"), "USR_foo")
		assert.Equal(t, r.URL.Query().Get("limit"), "42")
		assert.Equal(t, r.URL.Query().Get("starting_after"), "SES_foo")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(sampleSessionList)
	}))
	defer server.Close()
	client := createTestClient(server)
	sessionList, err := client.Sessions.List(feather.SessionsListParams{
		UserID: feather.String("USR_foo"),
		ListParams: feather.ListParams{
			Limit:         feather.UInt32(42),
			StartingAfter: feather.String("SES_foo"),
		},
	})
	assert.Equal(t, sampleSessionList, *sessionList)
	assert.Nil(t, err)
}

func TestSessionsList_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/sessions?"))
		assert.Equal(t, r.URL.Query().Get("user_id"), "")
		assert.Equal(t, r.URL.Query().Get("limit"), "42")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Code:    feather.ErrorCodeParameterInvalid,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	sessionList, err := client.Sessions.List(feather.SessionsListParams{
		ListParams: feather.ListParams{
			Limit: feather.UInt32(42),
		},
	})
	assert.Nil(t, sessionList)
	assert.Equal(t, "An error message", err.Error())
}

func TestSessionsRetrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/sessions/SES_foo"))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(sampleSessionAnonymous)
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Retrieve("SES_foo")
	assert.Equal(t, sampleSessionAnonymous, *session)
	assert.Nil(t, err)
}

func TestSessionsRetrieve_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/sessions/"))
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Code:    feather.ErrorCodeParameterInvalid,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Retrieve("")
	assert.Nil(t, session)
	assert.Equal(t, "An error message", err.Error())
}

func TestSessionsUpgrade(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/sessions/SES_bar/upgrade"))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(sampleSessionAuthenticated)
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Upgrade("SES_bar", feather.SessionsUpgradeParams{
		CredentialToken: feather.String("qwerty"),
	})
	assert.Equal(t, sampleSessionAuthenticated, *session)
	assert.Nil(t, err)
}

func TestSessionsUpgrade_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/sessions/SES_bar/upgrade"))
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Code:    feather.ErrorCodeParameterInvalid,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Upgrade("SES_bar", feather.SessionsUpgradeParams{
		CredentialToken: feather.String("-1"),
	})
	assert.Nil(t, session)
	assert.Equal(t, "An error message", err.Error())
}

// * * * * * Users * * * * * //
