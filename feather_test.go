package feather

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	sampleAPIKey = "fooKey"
)

func createTestClient(server *httptest.Server) Client {
	comps := strings.SplitN(strings.TrimPrefix(server.URL, "http://"), ":", 2)
	return New(sampleAPIKey, &Config{
		Protocol:   String("http"),
		Host:       String(comps[0]),
		Port:       String(comps[1]),
		HTTPClient: server.Client(),
	})
}

// * * * * * Credentials * * * * * //

var sampleCredentialEmailRequiresOneTimeCode = Credential{
	ID:        "CRD_foo",
	Object:    "credential",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	ExpiresAt: time.Date(2020, 01, 01, 01, 11, 01, 0, time.UTC),
	Status:    "requires_one_time_code",
	Token:     String("qwerty"),
	Type:      "email",
}

var sampleCredentialEmailValid = Credential{
	ID:        "CRD_foo",
	Object:    "credential",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	ExpiresAt: time.Date(2020, 01, 01, 01, 11, 01, 0, time.UTC),
	Status:    "bar",
	Token:     String("qwerty"),
	Type:      "email",
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
	credential, err := client.Credentials.Create(CredentialsCreateParams{
		Type:  "email",
		Email: String("foo@bar.com"),
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
		json.NewEncoder(w).Encode(Error{
			Object:  "error",
			Type:    "foo_err_type", // TODO enum
			Code:    "foo_err_code", // TODO enum
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	credential, err := client.Credentials.Create(CredentialsCreateParams{
		Type:     "email",
		Email:    String("foo@bar.com"),
		Username: String("foobar"),
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
	credential, err := client.Credentials.Update("CRD_foo", CredentialsUpdateParams{
		OneTimeCode: String("foobar"),
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
		json.NewEncoder(w).Encode(Error{
			Object:  "error",
			Type:    "foo_err_type", // TODO enum
			Code:    "foo_err_code", // TODO enum
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	credential, err := client.Credentials.Update("CRD_foo", CredentialsUpdateParams{})
	assert.Nil(t, credential)
	assert.Equal(t, "An error message", err.Error())
}

// * * * * * Sessions * * * * * //

var sampleSessionAnonymous = Session{
	ID:        "SES_foo",
	Object:    "session",
	Type:      "anonymous", // TODO enum
	Status:    "active",    // TODO enum
	Token:     String("qwerty"),
	UserID:    "USR_foo",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	RevokedAt: nil,
}

var sampleSessionAuthenticated = Session{
	ID:        "SES_bar",
	Object:    "session",
	Type:      "authenticated", // TODO enum
	Status:    "active",        // TODO enum
	Token:     String("qwerty"),
	UserID:    "USR_foo",
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	RevokedAt: Time(time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC)),
}

var sampleSessionList = SessionList{
	ListMeta: ListMeta{
		Objet:      "list",
		URL:        "/v1/sessions",
		HasMore:    false,
		TotalCount: 2,
	},
	Data: []*Session{
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
	session, err := client.Sessions.Create(SessionsCreateParams{
		CredentialToken: String("bar"),
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
		json.NewEncoder(w).Encode(Error{
			Object:  "error",
			Type:    "foo_err_type", // TODO enum
			Code:    "foo_err_code", // TODO enum
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Create(SessionsCreateParams{
		CredentialToken: String("-1"),
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
	sessionList, err := client.Sessions.List(SessionsListParams{
		UserID: String("USR_foo"),
		ListParams: ListParams{
			Limit:         UInt32(42),
			StartingAfter: String("SES_foo"),
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
		json.NewEncoder(w).Encode(Error{
			Object:  "error",
			Type:    "foo_err_type", // TODO enum
			Code:    "foo_err_code", // TODO enum
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	sessionList, err := client.Sessions.List(SessionsListParams{
		ListParams: ListParams{
			Limit: UInt32(42),
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
		json.NewEncoder(w).Encode(Error{
			Object:  "error",
			Type:    "foo_err_type", // TODO enum
			Code:    "foo_err_code", // TODO enum
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
	session, err := client.Sessions.Upgrade("SES_bar", SessionsUpgradeParams{
		CredentialToken: String("qwerty"),
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
		json.NewEncoder(w).Encode(Error{
			Object:  "error",
			Type:    "foo_err_type", // TODO enum
			Code:    "foo_err_code", // TODO enum
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Upgrade("SES_bar", SessionsUpgradeParams{
		CredentialToken: String("-1"),
	})
	assert.Nil(t, session)
	assert.Equal(t, "An error message", err.Error())
}

// * * * * * Users * * * * * //
