package feather

// ErrorType represents the type of error produced.
type ErrorType string

const (
	// Failed to connect to the Feather API.
	ErrorTypeAPIConnection ErrorType = "api_connection_error"

	// The API request did not have proper authentication.
	ErrorTypeAPIAuthentication ErrorType = "api_authentication_error"

	// Too many requests have been sent to quickly with this API key.
	ErrorTypeRateLimit ErrorType = "rate_limit_error"

	// The API request did not pass validation checks.
	ErrorTypeValidation ErrorType = "validation_error"

	// Any other type of error (eg temporary problem with the server).
	ErrorTypeAPI ErrorType = "api_error"
)

// ErrorCode provides a value which can be used to handle the error programmatically.
type ErrorCode string

const (
	ErrorCodeAPIKeyExpired                 ErrorCode = "api_key_expired"
	ErrorCodeAPIKeyInsufficientPermissions ErrorCode = "api_key_insufficient_permissions"
	ErrorCodeAPIKeyMissing                 ErrorCode = "api_key_missing"
	ErrorCodeAPIKeyInvalid                 ErrorCode = "api_key_invalid"
	ErrorCodeAPIKeyRevoked                 ErrorCode = "api_key_revoked"
	ErrorCodeBearerTokenInvalid            ErrorCode = "bearer_token_invalid"
	ErrorCodeCredentialAlreadyUsed         ErrorCode = "credential_already_used"
	ErrorCodeCredentialExpired             ErrorCode = "credential_expired"
	ErrorCodeCredentialInvalid             ErrorCode = "credential_invalid"
	ErrorCodeCredentialStatusNotValid      ErrorCode = "credential_status_not_valid"
	ErrorCodeCredentialStatusImmutable     ErrorCode = "credential_status_immutable"
	ErrorCodeCredentialTokenInvalid        ErrorCode = "credential_token_invalid"
	ErrorCodeCredentialTokenExpired        ErrorCode = "credential_token_expired"
	ErrorCodeHeaderEmpty                   ErrorCode = "header_empty"
	ErrorCodeHeaderMissing                 ErrorCode = "header_missing"
	ErrorCodeNotFound                      ErrorCode = "not_found"
	ErrorCodeOneTimeCodeInvalid            ErrorCode = "one_time_code_invalid"
	ErrorCodeOneTimeCodeUsed               ErrorCode = "one_time_code_used"
	ErrorCodeParsingFailed                 ErrorCode = "parsing_failed"
	ErrorCodeParameterEmpty                ErrorCode = "parameter_empty"
	ErrorCodeParameterInvalid              ErrorCode = "parameter_invalid"
	ErrorCodeParameterMissing              ErrorCode = "parameter_missing"
	ErrorCodeParameterUnknown              ErrorCode = "parameter_unknown"
	ErrorCodeParametersExclusive           ErrorCode = "parameters_exclusive"
	ErrorCodePasswordInvalid               ErrorCode = "password_invalid"
	ErrorCodePublicKeyNotFound             ErrorCode = "public_key_not_found"
	ErrorCodeSessionExpired                ErrorCode = "session_expired"
	ErrorCodeSessionInactive               ErrorCode = "session_inactive"
	ErrorCodeSessionRevoked                ErrorCode = "session_revoked"
	ErrorCodeSessionTokenInvalid           ErrorCode = "session_token_invalid"
	ErrorCodeSessionTokenExpired           ErrorCode = "session_token_expired"
	ErrorCodeUserBlocked                   ErrorCode = "user_blocked"
)

// Error is the Feather error object.
// https://feather.id/docs/reference/api#errors
type Error struct {
	Object  string    `json:"object"`
	Type    ErrorType `json:"type"`
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}
