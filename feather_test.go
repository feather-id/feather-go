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
		BasePath:   feather.String("/v1"),
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

type publicKeyResponse struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	PEM    string `json:"pem"`
}

var samplePublicKeyResponse = publicKeyResponse{
	ID:     "0",
	Object: "publicKey",
	PEM: `-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAwovomIOamL39k/Q7OfSxRf7ipn0kuQMLfEY0UWHcwq7ubKfjs368
wMcAa7vhlHamZnrONMTtZUNStbhrMBVlzGcSkhSrOENKg+g6KG29WD5VhupKmSGt
hDjQRlx2nvgZSSVdjx8S+BDArPpIWMviViswjRCucWdqFHR6av0v/bvMRYO3qRXK
pGn+LuiDlCi9sgiK72Ayt9unTjodyugchx6Y+RyboKOWZmiLFWRdkMZkvBaxxgaT
S/y1TneJR9eg5EPxh0YQYYEPT3/CYgaw34s/HtqbILWcr4VSG1lrKDQXOneYL+xj
svTcv2z81qX1WN+qhGasUw/dEwjYmbidNwIDAQAB
-----END RSA PUBLIC KEY-----
`,
}

var sampleSessionTokenValidButStale = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzc3Mzk0LCJleHAiOjE1ODkzNzc5OTQsImlhdCI6MTU4OTM3NzM5NCwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6IlVDc3pybnFTeGhzNkZmeDUwdFp6ZHVwTmpsWmZQM1VrZzI3VUZvVHhReVlwa0hZN2VtMFZndEtSTDNQempReXpia3JoNGVONDl0WlRTS1dJaVowN05UQlJoSHY1Y25JeXNPSzciLCJyYXQiOm51bGwsInNlcyI6IlNFU18xMDgzNmNiNi05OTRkLTQwZjYtOTUwYy0zNjE3YmUxN2I3YzMiLCJzdWIiOiJVU1JfYTY4NzViMzQtNDZlYS00MjlkLTg1OGMtM2FhNmEzNDNiNTM0IiwidHlwIjoiYXV0aGVudGljYXRlZCJ9.wgXjg4eY6ziujbGpOuwFyNB9hQrQSFd98Ey4gMhVarZK3OmdXbB0QqDahg1ON6Ebzr_oydjTyk1yD-eJ-5Rf5YwIvl7tC9fTPDkH2rYSIH6qfi0a5k-8-Km8E4x7TY4YPybdMmA4ycJUXvyPEl7N2awHb1YduCpDptUR9A2y_ASzyK4Lw01EdazEjho0OW2sJ7BjInirRbLuK1dKvrUicI8Sj-glr9WRlD1XF0zBeOTcwIa7sMieBrkCUtzPb1QWWTXbExmtGyDr0lyGX_dXdZGO_Q53PRI7m01HxNzrCf1GF_zXoKogg6iQnpbXopSTp51hwe5m-QYBPd0IT2YTsw`

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

func TestSessionsValidate(t *testing.T) {
	var requestCount = 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)

		switch requestCount {
		case 0:
			assert.Equal(t, r.Method, http.MethodGet)
			assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/publicKeys/0"))
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(samplePublicKeyResponse)

		case 1:
			assert.Equal(t, r.Method, http.MethodPost)
			assert.Equal(t, r.URL.String(), "/v1/sessions/SES_10836cb6-994d-40f6-950c-3617be17b7c3/validate")
			assert.Equal(t, r.FormValue("session_token"), sampleSessionTokenValidButStale)
			w.WriteHeader(200)
			json.NewEncoder(w).Encode(sampleSessionAuthenticated)
		default:
			break
		}
		requestCount += 1
	}))
	defer server.Close()
	client := createTestClient(server)
	_, err := client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenValidButStale),
	})
	// assert.Equal(t, sampleSessionAuthenticated, *session)
	assert.Nil(t, err)
	assert.Equal(t, 2, requestCount)
}

// * * * * * Users * * * * * //

var sampleUserEmpty = feather.User{
	ID:        "USR_foo",
	Object:    "user",
	Email:     nil,
	Username:  nil,
	Metadata:  map[string]string{},
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	UpdatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
}

var sampleUser = feather.User{
	ID:        "USR_bar",
	Object:    "user",
	Email:     nil,
	Username:  feather.String("foobar"),
	Metadata:  map[string]string{"highScore": "123"},
	CreatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
	UpdatedAt: time.Date(2020, 01, 01, 01, 01, 01, 0, time.UTC),
}

var sampleUserList = feather.UserList{
	ListMeta: feather.ListMeta{
		Objet:      "list",
		URL:        "/v1/users",
		HasMore:    false,
		TotalCount: 2,
	},
	Data: []*feather.User{
		&sampleUserEmpty,
		&sampleUser,
	},
}

func TestUsersList(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users?"))
		assert.Equal(t, r.URL.Query().Get("limit"), "42")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(sampleUserList)
	}))
	defer server.Close()
	client := createTestClient(server)
	userList, err := client.Users.List(feather.UsersListParams{
		ListParams: feather.ListParams{
			Limit: feather.UInt32(42),
		},
	})
	assert.Equal(t, sampleUserList, *userList)
	assert.Nil(t, err)
}

func TestUsersList_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users?"))
		assert.Equal(t, r.URL.Query().Get("limit"), "42")
		w.WriteHeader(429)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeRateLimit,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	userList, err := client.Users.List(feather.UsersListParams{
		ListParams: feather.ListParams{
			Limit: feather.UInt32(42),
		},
	})
	assert.Nil(t, userList)
	assert.Equal(t, "An error message", err.Error())
}

func TestUsersRetrieve(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users/USR_foo"))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(sampleUserEmpty)
	}))
	defer server.Close()
	client := createTestClient(server)
	user, err := client.Users.Retrieve("USR_foo")
	assert.Equal(t, sampleUserEmpty, *user)
	assert.Nil(t, err)
}

func TestUsersRetrieve_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users/"))
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
	user, err := client.Users.Retrieve("")
	assert.Nil(t, user)
	assert.Equal(t, "An error message", err.Error())
}

func TestUsersUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users/USR_bar"))
		assert.Equal(t, r.FormValue("username"), "foobar")
		assert.Equal(t, r.FormValue("metadata[highScore]"), "123")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(sampleUser)
	}))
	defer server.Close()
	client := createTestClient(server)
	user, err := client.Users.Update("USR_bar", feather.UsersUpdateParams{
		Username: feather.String("foobar"),
		Metadata: &map[string]string{
			"highScore": "123",
		},
	})
	assert.Equal(t, sampleUser, *user)
	assert.Nil(t, err)
}

func TestUsersUpdate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodPost)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users/"))
		assert.Equal(t, r.FormValue("username"), "foobar")
		assert.Equal(t, r.FormValue("metadata[highScore]"), "123")
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
	user, err := client.Users.Update("", feather.UsersUpdateParams{
		Username: feather.String("foobar"),
		Metadata: &map[string]string{
			"highScore": "123",
		},
	})
	assert.Nil(t, user)
	assert.Equal(t, "An error message", err.Error())
}
