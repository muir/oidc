package op

import (
	"context"
	"time"

	"gopkg.in/square/go-jose.v2"

	"github.com/zitadel/oidc/pkg/oidc"
)

type AuthStorage interface {
	CreateAuthRequest(context.Context, *oidc.AuthRequest, string) (AuthRequest, error)
	AuthRequestByID(context.Context, string) (AuthRequest, error)
	AuthRequestByCode(context.Context, string) (AuthRequest, error)
	SaveAuthCode(context.Context, string, string) error
	DeleteAuthRequest(context.Context, string) error

	// The TokenRequest parameter of CreateAccessToken can be any of:
	// - TokenRequest as returned by ClientCredentialsStorage.ClientCredentialsTokenRequest
	// - RefreshTokenRequest as returned by AuthStorage.TokenRequestByRefreshToken
	// - AuthRequest as returned one of the AuthStorage methods above
	// - oidc.JWTTokenRequest created by decoding a JWT
	CreateAccessToken(context.Context, TokenRequest) (string, time.Time, error)

	// The TokenRequest parameter of CreateAccessAndRefreshTokens can be any of:
	// - TokenRequest as returned by ClientCredentialsStorage.ClientCredentialsTokenRequest
	// - RefreshTokenRequest as returned by AuthStorage.TokenRequestByRefreshToken
	// - AuthRequest as returned one of the AuthStorage methods above
	CreateAccessAndRefreshTokens(ctx context.Context, request TokenRequest, currentRefreshToken string) (accessTokenID string, newRefreshToken string, expiration time.Time, err error)
	TokenRequestByRefreshToken(ctx context.Context, refreshToken string) (RefreshTokenRequest, error)

	TerminateSession(ctx context.Context, userID string, clientID string) error
	RevokeToken(ctx context.Context, token string, userID string, clientID string) *oidc.Error

	GetSigningKey(context.Context, chan<- jose.SigningKey)
	GetKeySet(context.Context) (*jose.JSONWebKeySet, error)
}

type ClientCredentialsStorage interface {
	ClientCredentialsTokenRequest(ctx context.Context, clientID string, scopes []string) (TokenRequest, error)
}

type OPStorage interface {
	GetClientByClientID(ctx context.Context, clientID string) (Client, error)
	AuthorizeClientIDSecret(ctx context.Context, clientID, clientSecret string) error
	SetUserinfoFromScopes(ctx context.Context, userinfo oidc.UserInfoSetter, userID, clientID string, scopes []string) error
	SetUserinfoFromToken(ctx context.Context, userinfo oidc.UserInfoSetter, tokenID, subject, origin string) error
	SetIntrospectionFromToken(ctx context.Context, userinfo oidc.IntrospectionResponse, tokenID, subject, clientID string) error
	GetPrivateClaimsFromScopes(ctx context.Context, userID, clientID string, scopes []string) (map[string]interface{}, error)
	GetKeyByIDAndUserID(ctx context.Context, keyID, userID string) (*jose.JSONWebKey, error)
	ValidateJWTProfileScopes(ctx context.Context, userID string, scopes []string) ([]string, error)
}

// Storage is a required parameter for NewOpenIDProvider(). In addition to the
// embedded interfaces below, if the passed Storage implements ClientCredentialsStorage
// then the grant type "client_credentials" will be supported. In that case, the access
// token returned by CreateAccessToken should be a JWT.
// See https://datatracker.ietf.org/doc/html/rfc6749#section-1.3.4 for context.
type Storage interface {
	AuthStorage
	OPStorage
	Health(context.Context) error
}

type StorageNotFoundError interface {
	IsNotFound()
}

type EndSessionRequest struct {
	UserID      string
	ClientID    string
	RedirectURI string
}
