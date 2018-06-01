package httplib

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/KulinaID/kulina-api-core/model/httpmodel"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

type httpfunc func(http.ResponseWriter, *http.Request)

var jwtMiddleware *jwtmiddleware.JWTMiddleware
var signingKey []byte

func InitJWTMiddleware(secret []byte) {
	signingKey = secret
	jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: jwt.SigningMethodHS256,
	})
}

func ExtractClaim(r *http.Request, key string) (interface{}, error) {
	tokenStr, err := jwtMiddleware.Options.Extractor(r)
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return signingKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims[key], nil
	} else {
		return "", nil
	}

}

func JWTMid(h httpfunc) httpfunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := checkJWT(w, r)
		if err != nil {
			return
		}
		h(w, r)
	}
}

func checkJWT(w http.ResponseWriter, r *http.Request) error {

	if !jwtMiddleware.Options.EnableAuthOnOptions {
		if r.Method == "OPTIONS" {
			return nil
		}
	}

	token, err := jwtMiddleware.Options.Extractor(r)
	if err != nil {
		errorRes := &httpmodel.ErrorResponse{
			"400",
			"invalid request: " + err.Error(),
			"invalid request: " + err.Error(),
			"400",
			"https://kulina.id",
		}
		ResponseError(w, errorRes)
		return err
	}

	if token == "" {

		if jwtMiddleware.Options.CredentialsOptional {
			return nil
		}

		errorMsg := "Required authorization token not found"
		errorRes := &httpmodel.ErrorResponse{
			"400",
			"invalid request: " + errorMsg,
			"invalid request: " + errorMsg,
			"400",
			"https://kulina.id",
		}
		ResponseError(w, errorRes)
		return fmt.Errorf(errorMsg)
	}

	parsedToken, err := jwt.Parse(token, jwtMiddleware.Options.ValidationKeyGetter)
	if err != nil {

		errorMsg := "Error parsing token: " + err.Error()
		errorRes := &httpmodel.ErrorResponse{
			"400",
			"invalid request: " + errorMsg,
			"invalid request: " + errorMsg,
			"400",
			"https://kulina.id",
		}
		ResponseError(w, errorRes)
		return fmt.Errorf("Error parsing token: %v", err)
	}

	if jwtMiddleware.Options.SigningMethod != nil && jwtMiddleware.Options.SigningMethod.Alg() != parsedToken.Header["alg"] {
		errorMsg := fmt.Sprintf("Expected %s signing method but token specified %s",
			jwtMiddleware.Options.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		errorRes := &httpmodel.ErrorResponse{
			"400",
			"invalid request: " + errorMsg,
			"invalid request: " + errorMsg,
			"400",
			"https://kulina.id",
		}
		ResponseError(w, errorRes)
		return fmt.Errorf("Error validating token algorithm: %s", errorMsg)
	}

	if !parsedToken.Valid {
		errorMsg := "Token is invalid"
		errorRes := &httpmodel.ErrorResponse{
			"400",
			"invalid request: " + errorMsg,
			"invalid request: " + errorMsg,
			"400",
			"https://kulina.id",
		}
		ResponseError(w, errorRes)
		return errors.New("Token is invalid")
	}

	newRequest := r.WithContext(context.WithValue(r.Context(), jwtMiddleware.Options.UserProperty, parsedToken))
	*r = *newRequest
	return nil
}
