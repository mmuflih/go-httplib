package httplib

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	jwtMid "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

type httpfunc func(http.ResponseWriter, *http.Request)

var jwtMiddleware *jwtMid.JWTMiddleware
var signingKey []byte
var myrole map[string][]string

func InitJWTMiddleware(secret []byte) {
	signingKey = secret
	jwtMiddleware = jwtMid.New(jwtMid.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: jwt.SigningMethodHS512,
	})
}

func InitJWTMiddlewareCustomSigningKey(secret []byte, signingMethod jwt.SigningMethod) {
	signingKey = secret
	jwtMiddleware = jwtMid.New(jwtMid.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: signingMethod,
	})
}

func InitJWTMiddlewareWithRole(secret []byte, signingMethod jwt.SigningMethod, role map[string][]string) {
	signingKey = secret
	myrole = role
	jwtMiddleware = jwtMid.New(jwtMid.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		},
		SigningMethod: signingMethod,
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
		err := checkJWT(w, r, "")
		if err != nil {
			return
		}
		h(w, r)
	}
}

func JWTMidWithRole(h httpfunc, role string) httpfunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := checkJWT(w, r, role)
		if err != nil {
			return
		}
		h(w, r)
	}
}

func checkJWT(w http.ResponseWriter, r *http.Request, role string) error {

	if !jwtMiddleware.Options.EnableAuthOnOptions {
		if r.Method == "OPTIONS" {
			return nil
		}
	}

	token, err := jwtMiddleware.Options.Extractor(r)
	if err != nil {
		eExtractor := errors.New("400")
		ResponseException(w, eExtractor, 400)
		return eExtractor
	}

	if token == "" {

		if jwtMiddleware.Options.CredentialsOptional {
			return nil
		}

		eReqiredToken := errors.New("Required authorization token not found")
		ResponseException(w, eReqiredToken, 401)
		return eReqiredToken
	}

	parsedToken, err := jwt.Parse(token, jwtMiddleware.Options.ValidationKeyGetter)
	if err != nil {
		ePassingToken := errors.New("Error parsing token: " + err.Error())
		ResponseException(w, ePassingToken, 401)
		return ePassingToken
	}
	fmt.Println(parsedToken.Claims, "token")

	if jwtMiddleware.Options.SigningMethod != nil && jwtMiddleware.Options.SigningMethod.Alg() != parsedToken.Header["alg"] {
		errorMsg := fmt.Sprintf("Expected %s signing method but token specified %s",
			jwtMiddleware.Options.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		eTokenSpecified := errors.New(errorMsg)
		ResponseException(w, eTokenSpecified, 401)
		return eTokenSpecified
	}

	if !parsedToken.Valid {
		eInvalidToken := errors.New("Token is invalid")
		ResponseException(w, eInvalidToken, 401)
		return eInvalidToken
	}

	newRequest := r.WithContext(context.WithValue(r.Context(), jwtMiddleware.Options.UserProperty, parsedToken))
	*r = *newRequest

	/** check role */
	if role == "" {
		return nil
	}
	tokenRole, _ := ExtractClaim(r, "role")
	fmt.Println(tokenRole, "role")
	for k, r := range myrole {
		if k == role {
			for _, c := range r {
				if c == tokenRole {
					return nil
				}
			}
			break
		}
	}
	e := errors.New("Access is not permitted")
	ResponseException(w, e, 401)
	return e
}
