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

var sampleSessionTokenInvalidAlg = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5MzgzMzUwLCJleHAiOjE1ODkzODM5NTAsImlhdCI6MTU4OTM4MzM1MCwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6IndqWEl0UDVrT29RR1EwekppUEp5NU4zenhuRjM1MjBBTk1CT0E4eEM2dklwMWN2S2o0TGpYOXdBN1VnckZiWkNHcWJzRzJLV21hZVNWaGs1aVlMTlRKTUw3YnQ4eGxTSWJ2Y3MiLCJyYXQiOm51bGwsInNlcyI6IlNFU184ZGZkYzJmOC05MWY1LTRlOTUtYTI2OC03ODU0MDk4YzI2ZmUiLCJzdWIiOiJVU1JfYTY4NzViMzQtNDZlYS00MjlkLTg1OGMtM2FhNmEzNDNiNTM0IiwidHlwIjoiYXV0aGVudGljYXRlZCJ9.xwBv15XVY6TZWJDNoMkyJG7OwCRRwpPF1Yv3p1_bEpk`

var sampleSessionTokenInvalidSignature = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzc3Mzk0LCJleHAiOjE1ODkzNzc5OTQsImlhdCI6MTU4OTM3NzM5NCwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6IlVDc3pybnFTeGhzNkZmeDUwdFp6ZHVwTmpsWmZQM1VrZzI3VUZvVHhReVlwa0hZN2VtMFZndEtSTDNQempReXpia3JoNGVONDl0WlRTS1dJaVowN05UQlJoSHY1Y25JeXNPSzciLCJyYXQiOm51bGwsInNlcyI6IlNFU18xMDgzNmNiNi05OTRkLTQwZjYtOTUwYy0zNjE3YmUxN2I3YzMiLCJzdWIiOiJVU1JfYTY4NzViMzQtNDZlYS00MjlkLTg1OGMtM2FhNmEzNDNiNTM0IiwidHlwIjoiYXV0aGVudGljYXRlZCJ9.wgXjg4eY6ziujbGpOuwFyNB9hQrQSFd98Ey4gMhVarZK3OmdXbB0QqDahg1ON6Ebzr_oydjTyk1yD-eJ-5Rf5YwIvl7tC9fTPDkH2rYSIH6qfi0a5k-8-Km8E4x7TY4YPybdMmA4ycJUXvyPEl7N2awHb1YduCpDptUR9A2y_ASzyK4Lw01EdazEjho0OW2sJ7BjInirRbLuK1dKvrUicI8Sj-glr9WRlD1XF0zBeOTcwIa7sMieBrkCUtzPb1QWWTXbExmtGyDr0lyGX_dXdZGO_Q53PRI7m01HxNzrCf1GF_zXoKogg6iQnpbXopSTp51hwe5m-QYBPd0IT2YTs`

var sampleSessionTokenModified = `fyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzc3Mzk0LCJleHAiOjE1ODkzNzc5OTQsImlhdCI6MTU4OTM3NzM5NCwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6IlVDc3pybnFTeGhzNkZmeDUwdFp6ZHVwTmpsWmZQM1VrZzI3VUZvVHhReVlwa0hZN2VtMFZndEtSTDNQempReXpia3JoNGVONDl0WlRTS1dJaVowN05UQlJoSHY1Y25JeXNPSzciLCJyYXQiOm51bGwsInNlcyI6IlNFU18xMDgzNmNiNi05OTRkLTQwZjYtOTUwYy0zNjE3YmUxN2I3YzMiLCJzdWIiOiJVU1JfYTY4NzViMzQtNDZlYS00MjlkLTg1OGMtM2FhNmEzNDNiNTM0IiwidHlwIjoiYXV0aGVudGljYXRlZCJ9.wgXjg4eY6ziujbGpOuwFyNB9hQrQSFd98Ey4gMhVarZK3OmdXbB0QqDahg1ON6Ebzr_oydjTyk1yD-eJ-5Rf5YwIvl7tC9fTPDkH2rYSIH6qfi0a5k-8-Km8E4x7TY4YPybdMmA4ycJUXvyPEl7N2awHb1YduCpDptUR9A2y_ASzyK4Lw01EdazEjho0OW2sJ7BjInirRbLuK1dKvrUicI8Sj-glr9WRlD1XF0zBeOTcwIa7sMieBrkCUtzPb1QWWTXbExmtGyDr0lyGX_dXdZGO_Q53PRI7m01HxNzrCf1GF_zXoKogg6iQnpbXopSTp51hwe5m-QYBPd0IT2YTsw`

var sampleSessionTokenMissingKeyId = `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzg0ODk0LCJleHAiOjE1ODkzODU0OTQsImlhdCI6MTU4OTM4NDg5NCwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6InBPUk1WQ29kcXZjVTVzbENkdUNiRUpsZ2o3SDhuVkJXQ0JoUWxDd0dvU2p1UXB5WHdneHJDUVVGNjdFZ0dNR0wzM20zQlpISmJ3eVZwemVtN1VjQVcyakUzdW9iU1V6eXZBRlMiLCJyYXQiOm51bGwsInNlcyI6IlNFU180OGExMjBlMS01YjlmLTQ1M2QtOGY4NC03OTNiYmE2M2YyZTgiLCJzdWIiOiJVU1JfYTY4NzViMzQtNDZlYS00MjlkLTg1OGMtM2FhNmEzNDNiNTM0IiwidHlwIjoiYXV0aGVudGljYXRlZCJ9.F4Qnba4SaCGWKFL_3ZPX1o6j-durl4wNi22l5SxjDlXJ7yHDZfus6I0AgA_ZGGTA82wki_sOroL8oA4UGlWJhFnWPAoNiEqw7gjkUaS82lXPJ0DlPhAW5lcABz5W0AL7DzT1YB6kgyCt-cPzbLZECE2l6rLvh4yhvgZ_s1RIMsKZcqZYHkZWGMhJPwgx22T3RzoLymA41JQVr0wrsz8LQfy5nKeDHg51RLk1rZB63n_hXgOFbdjN7a6Dv4Oskhev4dZdSTdCYDLj1fBBHoNT08m02qwKNRuRJ8ePXndf3Lib52m4OQAlvhWf2wKhHGc91PaTpKbRY6GqgK9VXYjWEQ`

var sampleSessionTokenInvalidIssuer = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzg0Mzc3LCJleHAiOjE1ODkzODQ5NzcsImlhdCI6MTU4OTM4NDM3NywiaXNzIjoiZm9vIiwianRpIjoiZmhyMDJYdWZaMW9sV2gxZWZValNTam9Jd2lyVnNKNTU0WmRUWVBVdUtmSmpIUG9zTHJzdVhsaUtTYU1sUUxlN2N6REFzOTIzaWVIT0M4RmhqcnhzODdIRXdQNXVoWUgwRG1IOSIsInJhdCI6bnVsbCwic2VzIjoiU0VTX2YzZTFlZjZlLTI1OWUtNDNmYy1iNzRhLTEyMjlmM2ZmMzk3MCIsInN1YiI6IlVTUl9hNjg3NWIzNC00NmVhLTQyOWQtODU4Yy0zYWE2YTM0M2I1MzQiLCJ0eXAiOiJhdXRoZW50aWNhdGVkIn0.aaWgoDEkWFXLmK2cr-MCacosr9jhqWycumPyCYUV0AUNe_Vw4Ll-eB980yIHlaeCNhOlWmxMClFrsQ1foLcOmFeWT_SQk66OyvQD5b3JElVF6BkmspuO856BaAegTNsyWoAKPdjOodHp5ibQ7C8vGagYnck7nkGSIAzjFOIW5shVYjlgb80If78WIuD5dmeCm1LP2CKKFagEtHsP24Eib3-P35XzBScTOWp828erNEn_XkZSa6RZNE_2FBklGAxrxq4EfW-NglbvDWbgEk4cntoXDUwt_Rqg_k8mYOJPn6YnbFGfVnaB2qzZAOB7WFsJ_-Hl-wm_88_zLvIOXbUxDA`

var sampleSessionTokenInvalidSubject = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzg0NzIxLCJleHAiOjE1ODkzODUzMjEsImlhdCI6MTU4OTM4NDcyMSwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6IjBrcFZBb0pNTDVSSkRPUzZEMXJxVVZUOEI1b2k4OXhON214R21TZERsYUxCWUc3VmxvOU1uYnZpcmVaV1BJSGNrM1RSeFFJUDF4MFlHVmhhTENncFJqSWsxS1JNbTVRRUp1WWEiLCJyYXQiOm51bGwsInNlcyI6IlNFU183NGQ2YWQ0NC04MzAxLTQxZTQtOGRjNC05MzNkOTcwM2NlNzkiLCJzdWIiOiJmb28iLCJ0eXAiOiJhdXRoZW50aWNhdGVkIn0.MWB7wxza-QFwGAdGYW-FSRrgDgQ8gDoGluCUNfKTNeAn92tontQ-AFjvw_9CiptORVSeAYgLJ8pgonYqi7-16PV6BN6opDkHxllCzMRgbtxvJrGLJjURPqoUU-9YIeCJxB1nOXiildqQwSU9mV72lYolursmym8VEnLp1joaAbAZdbof2QkoHu8mi5aEnLErnSN_2eOx6R8HUdkpfIDa2_1vBQ0s5SDAqmTP2ILjuwRcrSeQe7pDKrZJpgTFb1WSO4Tx7ZMTg1OO-75Olb3mAs5t-N7YY78LzXKCRXuhIp2J0PhtGjq8Y2nQhL4vqBQ5FOXeWurnAV3r0cmMqFpizQ`

var sampleSessionTokenInvalidAudience = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJmb28iLCJjYXQiOjE1ODkzODUxMTUsImV4cCI6MTU4OTM4NTcxNSwiaWF0IjoxNTg5Mzg1MTE1LCJpc3MiOiJmZWF0aGVyLmlkIiwianRpIjoiZ0g2RUQ2OVhjUFFpc3ZVTnd6Qm1wRmMwSWFSMW9neU13S083NHpQN1hyRHdaa0Fac3pZOFI1QmZ5bllmNGNhWG8zbU1OSnpCRWJ5Mk1CSUlyOE9mV2ZSUkNOek1yWVRPZ0dTWiIsInJhdCI6bnVsbCwic2VzIjoiU0VTXzc3NjM4YmI2LTY1YjYtNDEwNC1hNTQyLTc2ZDRjNWE4MmMxOSIsInN1YiI6IlVTUl9hNjg3NWIzNC00NmVhLTQyOWQtODU4Yy0zYWE2YTM0M2I1MzQiLCJ0eXAiOiJhdXRoZW50aWNhdGVkIn0.ZzueiC6fTJF5zjrmlQh4t3AHylS8DARxeGHZEOpKPIIy6RdT0oD3WmINb1Eu5DEX7PkQs2jpS3cAXQBZxqutUjkaf6ZI0pfzU72Tp1Yu5rtCw8cucWpCt9sRvFrCxSuIwq2yM04plRtB8D040PTGJpSf2Gojf4O6EhsxTktjYCjqotdrpmLj2rfvl7Jj-MAv343mNa-C5id03V5Cb0GOE7EIkSeuDn3_QRX6e11BacMmKChxMzQlivpIZ7htOGOM4IZjWojurF4aFZuUI07YzXhjVWKiaHED91KSX1fY4M0qsvr2f-6f-JuFYzuv5zku2CLO7AY6AgZ7VGiYB3pOlw`

var sampleSessionTokenInvalidSessionId = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzg1MjE1LCJleHAiOjE1ODkzODU4MTUsImlhdCI6MTU4OTM4NTIxNSwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6Im9JaE9oaXJMbVZoRzdQOFhJRHZMa2NVSExxMUhsOVNTOEtLQXpCdzIyQmtDZEdrZmNLTHNodmpMTkVPMWxZY2JmRWtQakl0QURtN2FxTWhxTThXT0JHYmhGNnFkSTVkcTkwdnMiLCJyYXQiOm51bGwsInNlcyI6ImZvbyIsInN1YiI6IlVTUl9hNjg3NWIzNC00NmVhLTQyOWQtODU4Yy0zYWE2YTM0M2I1MzQiLCJ0eXAiOiJhdXRoZW50aWNhdGVkIn0.lxTKv8ty3OdupmqTwSKLYDH3bTwUcDSbHmzKh25jDs9qNNzZzOb27ekY1JEzrOZc1zcQjUh1I-lyAfnwNQeL7JmM2L-2dCGArivGoeB9lS-272uovy0eFZMrFsHjbiaQFHDxoiUDiQg5_kShaes2C8nNjHMOK0qOSPkhK9JTeEkgWYcItt7RK1TFHDUoJtoQELLJmOJhk47rgUnQn-kArR4ITy6rPGbJWUmsko2L1Z4kN1BuNCUpHf6D7KJQ51u7BuqR-_X64a0nF5fp1Q5GItr5YgE2Tr4pm2KumdCUOnHkBbzZ6_bR31DgLoI5AIM4wByxVHVovIlGsFqKtnaC5A`

var sampleSessionTokenInvalidType = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzg1NDU5LCJleHAiOjE1ODkzODYwNTksImlhdCI6MTU4OTM4NTQ1OSwiaXNzIjoiZmVhdGhlci5pZCIsImp0aSI6Imc2QUV0Tnd5NmRyY3BKSDZKZ0NZNTdyQ0IzZ2dmM3FZZXBsZndFcmNBSkROakFnYUk2VG5DZ2ZkZWlFYW1Ha0dseEFJdk90OGlTd2V2dXJWN0lxWFR5ZHVqSHVaem84ejBldFciLCJyYXQiOm51bGwsInNlcyI6IlNFU19jYjA5NmYyZS0wNTE5LTRjMjItYmFhNC1iNTU4NTM1YjUwNWUiLCJzdWIiOiJVU1JfYTY4NzViMzQtNDZlYS00MjlkLTg1OGMtM2FhNmEzNDNiNTM0IiwidHlwIjoiZm9vIn0.MyIIkcL1Cvug5_-B8jLOAMhFuk0kAuReQ4ueeolYWZNcn_hW8kocCn7GCNdGDZIPFv272-ioKNP4wVKjkN55D9dLN3ySV8B4j-dX9b3KfNZ0jaQk6HaHUbRZm-YUF0NXNe8YaikbYVLGouLHakQxd6hKvY4QBy0lfOuIe4vJ3YwNOfFD-vwv0eDp_CLfHUyFIScZ2GSdoOZpepnZ8qLznQS81wpR_d1lLWxmcgYnqeZzCV7GCAxV34j3F7_8IgEt8LkcR-PDMTTN-M6JOKRMIwFIiDPEeOrTmOLXJr3IYOt6KTeGtPqjOE8A7nN8yOEgXxUcI17CgNFvTdhc6tIglQ`

var sampleSessionTokenInvalidCreatedAt = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoiZm9vIiwiZXhwIjoxNTg5Mzg2MzQzLCJpYXQiOjE1ODkzODU3NDMsImlzcyI6ImZlYXRoZXIuaWQiLCJqdGkiOiJtOWxBYWs4UXFTVmwxaHpjaTBQWnExWjduaFVFS2UxMElQUWlBSEpWaG1MYVB3cTVqcmg4NDlJb3JFb09jZnN6eFM1WU9XdVlkOXFnQzV4ZE9TaE5hcUtWb0locm1kR3pKZmMxIiwicmF0IjpudWxsLCJzZXMiOiJTRVNfNTY0MzlmODgtZTExOS00MDYzLWJlMTQtYzk0NTE3OTZhM2NhIiwic3ViIjoiVVNSX2E2ODc1YjM0LTQ2ZWEtNDI5ZC04NThjLTNhYTZhMzQzYjUzNCIsInR5cCI6ImF1dGhlbnRpY2F0ZWQifQ.J93P2kFKvnUx6QFvONyvZqsiAOs4cZSE38may0ssQ_zyZXgc2yBBgocZOOymXjQFJu5DmXODh_zR24wZOR8hVSq3ROsZZiUPeFfd5pUKzigbx0x52lvCaHWThCWj8OIxwd9GUSwT6Re_IJsEmYXAU25dBdQ5Dgvt9dZR6UBS8KVTUGvsQUDTYJoVfgC4iQvGPRUbKqaja9IkkDaT8qhOYl0-2s7LR0NKkc5iDMnKFi5IauZFJfiiC__Sb_KrHOF4_ZvnL4Z9a_x8Dvt1VX5qX16Hv3nB7ZNGPoBNEsM0I7VyLUvvl6n9S4m5Ynvfb9q5Z-CzFXO5qmtTBqNvRaLMhA`

var sampleSessionTokenInvalidExpiresAt = `eyJhbGciOiJSUzI1NiIsImtpZCI6IjAiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJQUkpfY2RiY2M5ODYtYWU2Ni00NjY2LWI5NDYtMTgyNmNmMGIyYjU3IiwiY2F0IjoxNTg5Mzg1ODQ4LCJleHAiOiJmb28iLCJpYXQiOjE1ODkzODU4NDgsImlzcyI6ImZlYXRoZXIuaWQiLCJqdGkiOiJTeUFmczBlM0lnWWhFVnB1eUdHZWlCMlJGNHI4bkIzemludjlOUW5JUVdPeVpGTUpBTVBjTmhVbU41N2J6dFU3SUdzS2xWdElEZHRmRHBCV2RKTExzUXg4YmtNQmZNZmRSYUNDIiwicmF0IjpudWxsLCJzZXMiOiJTRVNfOTQ2NmVhYWQtYTY0YS00MGE0LTk2YzktZjZhMTlhZGYxYzcyIiwic3ViIjoiVVNSX2E2ODc1YjM0LTQ2ZWEtNDI5ZC04NThjLTNhYTZhMzQzYjUzNCIsInR5cCI6ImF1dGhlbnRpY2F0ZWQifQ.Cekvy-Liu0xYJfirQNMY4upHzOybqmkmWrisXofiIeNbgp2y-_EvrlybsJctroOz9MGhldqPJgoddl3CZseHcvWc3u_vVJDdCoYFsVIMzbJr3hCBmPnHOncLmQy-_Mr-Z5vg4x6r3hgbwn_pWM2iBIPKAccIQzS85vJF4QZnOCgyqn9KfOliyCQOtRIUl4UFZg4gcT7tX4ehhDLPc0uJZTuM0voqznZ1K_txsZRxlXBm5iAaL7bVPgI1dX6qgi5AoGew2B2ghZjxly3RjfaxGctqlXRiU42U_2U9oJDyyWSvCtF0Vjv2xisNwBINcsYDe3aBBH0VZrdh6IRYMwvHlA`

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
	session, err := client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenValidButStale),
	})
	assert.Equal(t, sampleSessionAuthenticated, *session)
	assert.Nil(t, err)
	assert.Equal(t, 2, requestCount)
}

func TestSessionsValidate_Nil(t *testing.T) {
	client := feather.New(sampleAPIKey)
	session, err := client.Sessions.Validate(feather.SessionsValidateParams{})
	assert.Nil(t, session)
	assert.Equal(t, "No session tokens were not provided for validation", err.Error())
}

func TestSessionsValidate_Invalid(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/publicKeys/0"))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(samplePublicKeyResponse)
	}))
	defer server.Close()
	client := createTestClient(server)

	// Invalid signature algorithm
	session, err := client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidAlg),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid signature
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidSignature),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Modified token
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenModified),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Missing key ID
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenMissingKeyId),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid issuer
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidIssuer),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid subject
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidSubject),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid audience
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidAudience),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid session ID
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidSessionId),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid session type
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidType),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid created at date
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidCreatedAt),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())

	// Invalid expires at date
	session, err = client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenInvalidExpiresAt),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())
}

func TestSessionsValidate_PublicKeyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/publicKeys/0"))
		w.WriteHeader(404)
		json.NewEncoder(w).Encode(feather.Error{
			Object:  "error",
			Type:    feather.ErrorTypeValidation,
			Message: "An error message",
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenValidButStale),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())
}

func TestSessionsValidate_PublicKeyParsingError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/publicKeys/0"))
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(publicKeyResponse{
			ID:     "0",
			Object: "publicKey",
			PEM: `-----BEGIN RSA PUBLIC KEY-----
		foo
		-----END RSA PUBLIC KEY-----
		`,
		})
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenValidButStale),
	})
	assert.Nil(t, session)
	assert.Equal(t, "The session token is invalid", err.Error())
}

func TestSessionsValidate_GatewayError(t *testing.T) {
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
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(feather.Error{
				Object:  "error",
				Type:    feather.ErrorTypeValidation,
				Message: "An error message",
			})
		default:
			break
		}
		requestCount += 1
	}))
	defer server.Close()
	client := createTestClient(server)
	session, err := client.Sessions.Validate(feather.SessionsValidateParams{
		SessionToken: feather.String(sampleSessionTokenValidButStale),
	})
	assert.Nil(t, session)
	assert.Equal(t, "An error message", err.Error())
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

// * * * * * Gateway * * * * * //

func TestGateway_UnparsableResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, _, _ := r.BasicAuth()
		assert.Equal(t, username, sampleAPIKey)
		assert.Equal(t, r.Method, http.MethodGet)
		assert.True(t, strings.HasPrefix(r.URL.String(), "/v1/users/USR_foo"))
		w.WriteHeader(404)
		w.Write([]byte("foo"))
	}))
	defer server.Close()
	client := createTestClient(server)
	user, err := client.Users.Retrieve("USR_foo")
	assert.Nil(t, user)
	assert.Equal(t, "The gateway received an unparsable response with status code 404", err.Error())
}
