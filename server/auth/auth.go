package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var (
	ErrExpiredToken     error = errors.New("Token is expired")
	ErrNotValidYetToken error = errors.New("Token not active yet")
	ErrMalformedToken   error = errors.New("That's not even a token")
	ErrInvalidToken           = errors.New("Invalid token")
)

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       time.Duration
	tokenType     string
}

// TokenInfo jwt info
type TokenInfo struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Expire      int64  `json:"expire"`
}

var (
	defaultKey = []byte("easybot")
	Options    = options{
		tokenType:     "Bearer",
		expired:       time.Duration(7200) * time.Second,
		signingMethod: jwt.SigningMethodHS512,
		signingKey:    defaultKey,
		keyfunc: func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}
			return defaultKey, nil
		},
	}
)

func SetTokenType(tokenType string) {
	Options.tokenType = tokenType
}

func SetExpired(duration time.Duration) {
	Options.expired = duration
}

func SetSigningMethod(method *jwt.SigningMethodHMAC) {
	Options.signingMethod = method
}

func SetSigningKey(secret string) {
	Options.signingKey = []byte(secret)
}

// GenerateToken create new jwt token
func GenerateToken(userID, audience string) (*TokenInfo, error) {
	now := time.Now()
	expiresAt := now.Add(Options.expired).Unix()

	token := jwt.NewWithClaims(Options.signingMethod, &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		Audience:  audience,
		ExpiresAt: expiresAt,
		NotBefore: now.Unix(),
		Subject:   userID,
	})

	tokenString, err := token.SignedString(Options.signingKey)
	if err != nil {
		return nil, err
	}
	return &TokenInfo{
		AccessToken: tokenString,
		TokenType:   Options.tokenType,
		Expire:      expiresAt,
	}, nil
}

// ParseToken resolve TokenString to Info
func ParseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, Options.keyfunc)
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, ErrMalformedToken
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				// Token is expired
				return nil, ErrExpiredToken
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, ErrNotValidYetToken
			} else {
				return nil, ErrInvalidToken
			}
		}
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func RefreshToken(tokenString string) (*TokenInfo, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, Options.keyfunc)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		now := time.Now()
		expiresAt := now.Add(Options.expired).Unix()
		claims.ExpiresAt = time.Now().Add(Options.expired).Unix()

		token := jwt.NewWithClaims(Options.signingMethod, claims)
		tokenString, err := token.SignedString(Options.signingKey)
		if err != nil {
			return nil, err
		}
		return &TokenInfo{
			AccessToken: tokenString,
			TokenType:   Options.tokenType,
			Expire:      expiresAt,
		}, nil
	}

	return nil, ErrInvalidToken
}

// GetTokenFromRequest get token info from http request
func GetTokenFromRequest(request *http.Request) (*jwt.StandardClaims, error) {
	tokenString := request.Header.Get("Authorization")
	prefix := fmt.Sprintf("%s ", Options.tokenType)
	if tokenString != "" && strings.HasPrefix(tokenString, prefix) {
		tokenString = tokenString[len(prefix):]
	} else {
		return nil, ErrMalformedToken
	}

	return ParseToken(tokenString)
}
