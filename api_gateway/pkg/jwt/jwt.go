package jwt

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var SignedKey = []byte("MGJd@Ro]yKoCc)mVY1^c:upz~4rn9Pt!hYd]>c8dt#+%")

func GenJWT(m map[interface{}]interface{}) (string, string, error) {
	var (
		accessToken, refreshToken *jwt.Token
		claims                    jwt.MapClaims
	)

	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	rClaims := refreshToken.Claims.(jwt.MapClaims)

	for k, v := range m {
		claims[k.(string)] = v
		rClaims[k.(string)] = v
	}

	claims["iss"] = "user"
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().AddDate(0, 0, 1).Unix()

	rClaims["iss"] = "user"
	rClaims["iat"] = time.Now().Unix()
	rClaims["exp"] = time.Now().AddDate(0, 0, 10).Unix()

	accessTokenString, err := accessToken.SignedString(SignedKey)
	if err != nil {
		err = fmt.Errorf("access_token generating error: %s", err)
		return "", "", err
	}

	refreshTokenString, err := refreshToken.SignedString(SignedKey)
	if err != nil {
		err = fmt.Errorf("refresh_token generating error: %s", err)
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil

}

func ExtractClaims(tokenStr string) (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err   error
	)
	token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// check token signing method etc
		return SignedKey, nil
	})
	if err != nil {
		token, err = jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			// check token signing method etc
			return SignedKey, nil
		})
		if err != nil {
			return nil, err
		}
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		err = fmt.Errorf("Invalid JWT Token")
		return nil, err
	}
	return claims, nil
}
