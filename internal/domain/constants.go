package domain

import "fmt"

// Error types
type ErrorType string

const (
	ErrTypeValidation     ErrorType = "VALIDATION_ERROR"
	ErrTypeConnection     ErrorType = "CONNECTION_ERROR"
	ErrTypeAuthentication ErrorType = "AUTHENTICATION_ERROR"
	ErrTypeAuthorization  ErrorType = "AUTHORIZATION_ERROR"
	ErrTypeDatabase       ErrorType = "DATABASE_ERROR"
	ErrTypeSession        ErrorType = "SESSION_ERROR"
	ErrTypeTransaction    ErrorType = "TRANSACTION_ERROR"
	ErrTypeQuery          ErrorType = "QUERY_ERROR"
	ErrTypeSecurity       ErrorType = "SECURITY_ERROR"
	ErrTypeNotFound       ErrorType = "NOT_FOUND_ERROR"
	ErrTypeConflict       ErrorType = "CONFLICT_ERROR"
	ErrTypeInternal       ErrorType = "INTERNAL_ERROR"
)

// ApplicationError represents an application-level error
type ApplicationError struct {
	Type    ErrorType
	Message string
	Details string
	Code    int
}

func (e *ApplicationError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Common errors
var (
	// Validation errors
	ErrEmptyUsername        = &ApplicationError{Type: ErrTypeValidation, Message: "username cannot be empty", Code: 400}
	ErrEmptyPassword        = &ApplicationError{Type: ErrTypeValidation, Message: "password cannot be empty", Code: 400}
	ErrInvalidConnectionStr = &ApplicationError{Type: ErrTypeValidation, Message: "invalid connection string format", Code: 400}
	ErrInvalidQuery         = &ApplicationError{Type: ErrTypeValidation, Message: "invalid SQL query", Code: 400}
	ErrInvalidWhereClause   = &ApplicationError{Type: ErrTypeValidation, Message: "invalid WHERE clause", Code: 400}

	// Connection errors
	ErrConnectionFailed = &ApplicationError{Type: ErrTypeConnection, Message: "failed to connect to database", Code: 503}
	ErrNoAccessibleDB   = &ApplicationError{Type: ErrTypeConnection, Message: "user has no accessible databases", Code: 403}

	// Authentication errors
	ErrInvalidCredentials    = &ApplicationError{Type: ErrTypeAuthentication, Message: "invalid username or password", Code: 401}
	ErrAuthenticationFailed  = &ApplicationError{Type: ErrTypeAuthentication, Message: "authentication failed", Code: 401}
	ErrProbeConnectionFailed = &ApplicationError{Type: ErrTypeAuthentication, Message: "failed to probe user connection", Code: 401}

	// Authorization errors
	ErrUnauthorized            = &ApplicationError{Type: ErrTypeAuthorization, Message: "unauthorized access", Code: 403}
	ErrInsufficientPermissions = &ApplicationError{Type: ErrTypeAuthorization, Message: "insufficient permissions", Code: 403}
	ErrTableAccessDenied       = &ApplicationError{Type: ErrTypeAuthorization, Message: "access denied to table", Code: 403}

	// Session errors
	ErrInvalidSession   = &ApplicationError{Type: ErrTypeSession, Message: "invalid session", Code: 401}
	ErrSessionExpired   = &ApplicationError{Type: ErrTypeSession, Message: "session expired", Code: 401}
	ErrSessionNotFound  = &ApplicationError{Type: ErrTypeSession, Message: "session not found", Code: 404}
	ErrMultipleSessions = &ApplicationError{Type: ErrTypeSession, Message: "multiple sessions not supported", Code: 409}

	// Transaction errors
	ErrNoActiveTransaction     = &ApplicationError{Type: ErrTypeTransaction, Message: "no active transaction", Code: 400}
	ErrActiveTransactionExists = &ApplicationError{Type: ErrTypeTransaction, Message: "transaction already active", Code: 409}
	ErrTransactionExpired      = &ApplicationError{Type: ErrTypeTransaction, Message: "transaction expired", Code: 408}
	ErrCommitFailed            = &ApplicationError{Type: ErrTypeTransaction, Message: "failed to commit transaction", Code: 500}
	ErrRollbackFailed          = &ApplicationError{Type: ErrTypeTransaction, Message: "failed to rollback transaction", Code: 500}

	// Database/Query errors
	ErrQueryFailed      = &ApplicationError{Type: ErrTypeDatabase, Message: "query execution failed", Code: 500}
	ErrTableNotFound    = &ApplicationError{Type: ErrTypeDatabase, Message: "table not found", Code: 404}
	ErrSchemaNotFound   = &ApplicationError{Type: ErrTypeDatabase, Message: "schema not found", Code: 404}
	ErrDatabaseNotFound = &ApplicationError{Type: ErrTypeDatabase, Message: "database not found", Code: 404}

	// Security errors
	ErrCookieTampering      = &ApplicationError{Type: ErrTypeSecurity, Message: "cookie tampering detected", Code: 400}
	ErrSQLInjectionDetected = &ApplicationError{Type: ErrTypeSecurity, Message: "potential SQL injection detected", Code: 400}

	// Not found errors
	ErrNotFound = &ApplicationError{Type: ErrTypeNotFound, Message: "resource not found", Code: 404}

	// Conflict errors
	ErrConflict = &ApplicationError{Type: ErrTypeConflict, Message: "resource conflict", Code: 409}

	// Internal errors
	ErrInternal = &ApplicationError{Type: ErrTypeInternal, Message: "internal server error", Code: 500}
)

// Constants for configuration and limits
const (
	// Query limits
	QueryResultHardLimit    = 1000
	QueryResultPageSize     = 50
	QueryResultDisplayLimit = 1000

	// Pagination
	CursorPaginationDefaultLimit = 50
	CursorPaginationMaxLimit     = 50

	// Session
	SessionTokenLength       = 32
	SessionExpirationTime    = 24 * 60 * 60     // 24 hours in seconds
	PasswordCookieExpiration = 15 * 60          // 15 minutes in seconds
	IdentityCookieExpiration = 7 * 24 * 60 * 60 // 7 days in seconds

	// Transaction
	TransactionTimeout = 60 * 60 // 1 hour in seconds

	// Database
	DefaultPostgresPort = "5432"
	DefaultSchema       = "public"

	// Encryption
	EncryptionKeyLength = 32
	NounceLength        = 12
)

// Cookie names
const (
	CookieSessionID = "session_id"
	CookieUsername  = "username"
	CookiePassword  = "password"
	CookieNonce     = "nonce"
	CookieSignature = "signature"
)

// Session and query status
const (
	StatusActive   = "active"
	StatusExpired  = "expired"
	StatusInactive = "inactive"
	StatusError    = "error"
)

// Sort directions
const (
	SortDirectionASC  = "ASC"
	SortDirectionDESC = "DESC"
)

// Permission types
const (
	PermissionSelect  = "SELECT"
	PermissionInsert  = "INSERT"
	PermissionUpdate  = "UPDATE"
	PermissionDelete  = "DELETE"
	PermissionConnect = "CONNECT"
	PermissionUsage   = "USAGE"
)
