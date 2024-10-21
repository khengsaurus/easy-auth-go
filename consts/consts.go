package consts

type ContextKey string

var (
	UrlBase = "https://ice-milo.com/ea-api/api/"

	HeaderApiKey    = "Ea-Api-Key"
	HeaderUserToken = "Ea-User-Token"

	ErrorInvalidToken    = "invalid easy-auth user token"
	ErrorUserInfoMissing = "failed to retrieve user info"

	ContextKeyUser = ContextKey("easy_auth_user")
)
