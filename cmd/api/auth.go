package main

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"personal-projects/webapp/pkg/data"
	"strings"
	"time"
)

var jwtTokenExpiry = time.Minute * 15
var refreshTokenExpiry = time.Hour * 24 // will not contain much information

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

func (app *application) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	// add a header
	w.Header().Add("Vary", "Authorization")

	// get the authorization header
	authHeader := r.Header.Get("Authorization")

	// sanity check
	if authHeader == "" {
		return "", nil, errors.New("Missing Authorization header")
	}
	// split the header on spaces
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("Invalid Authorization header")
	}
	// check to see if we have the word Bearer
	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("Unauthorized header: no Bearer found")
	}

	token := headerParts[1]

	// declare an empty claims variable
	claims := &Claims{}

	// Parse the token with our claims (we read into claims), using our secret (from the receiver)
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		// validate the signing algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.JWSecret), nil
	})
	// check for an erro; note that this catches expired token as well.
	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expired token")
		}
		return "", nil, err
	}
	// make sure we issued this token
	if claims.Issuer != app.Domain {
		return "", nil, errors.New("invalid issuer")
	}
	//valid token
	return token, claims, nil
}

func (app *application) generateTokenPair(user *data.User) (TokenPairs, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = app.Domain
	claims["iss"] = app.Domain
	if user.IsAdmin == 1 {
		claims["admin"] = true
	} else {
		claims["admin"] = false
	}
	// set the expiry
	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()

	// create signed token
	signedAccessToken, err := token.SignedString([]byte(app.JWSecret))
	if err != nil {
		return TokenPairs{}, err
	}
	// create refresh token
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	// set the expiry; must be longer than jwt expiry
	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	//create signed refresh token
	signedRefreshToken, err := refreshToken.SignedString([]byte(app.JWSecret))
	if err != nil {
		return TokenPairs{}, err
	}
	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	return tokenPairs, nil
}
