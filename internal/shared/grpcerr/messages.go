package grpcerr

// Stable, client-safe error messages. Use with ClientError; map at the gateway to HTTP status.
const (
	MsgUsernameTaken       = "username already taken"
	MsgInvalidCredentials  = "invalid credentials"
	MsgInvalidRefreshToken = "invalid refresh token"
	MsgUserNotFound        = "user not found"
	MsgInvalidRequest      = "invalid request"
	MsgInternal            = "internal server error"
	MsgUpstreamUnavailable = "service temporarily unavailable"
	MsgUpstreamTimeout     = "service temporarily unavailable"
)
