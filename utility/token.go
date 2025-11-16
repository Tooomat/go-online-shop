package utility

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ACCESS_TOKEN_TIME_TO_LIFE  = time.Minute * 15
	REFRESH_TOKEN_TIME_TO_LIFE = time.Hour * 24 * 7
)

// token API
func CreateAccessToken(id string, role string, secret string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"id":       id,
		"role":     role,
		"expAt":    time.Now().Add(ACCESS_TOKEN_TIME_TO_LIFE).Unix(), //15 menit exp (short lived access token)
		"issuedAt": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func CreateRefreshToken(id string, role string, secret string) (tokenString string, err error) {
	claims := jwt.MapClaims{
		"id":       id,
		"role":     role,
		"expAt":    time.Now().Add(REFRESH_TOKEN_TIME_TO_LIFE).Unix(),
		"IssuedAt": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", nil
	}
	return tokenString, nil
}

func ParseAccessToken(tokenString string, secret string) (id, role string, expAt, issueAt time.Time, err error) {
	tokens, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		err = fmt.Errorf("failed to parse token: %v", err)
		return
	}

	claims, ok := tokens.Claims.(jwt.MapClaims)
	if ok && tokens.Valid {
		id = fmt.Sprintf("%v", claims["id"])
		role = fmt.Sprintf("%v", claims["role"])

		if exp, ok := claims["expAt"].(float64); ok {
			expAt = time.Unix(int64(exp), 0)
		}
		if issue, ok := claims["IssuedAt"].(float64); ok {
			issueAt = time.Unix(int64(issue), 0)
		}
		return
	}

	err = fmt.Errorf("unable to extract claims access token")
	return
}

func ParseRefreshToken(tokenString string, secret string) (id, role string, err error) {
	tokens, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}

		return []byte(secret), nil
	})

	if err != nil {
		err = fmt.Errorf("failed to parse token: %v", err)
		return
	}

	claims, ok := tokens.Claims.(jwt.MapClaims)
	if ok && tokens.Valid {
		id = fmt.Sprintf("%v", claims["id"])
		role = fmt.Sprintf("%v", claims["role"])
		return
	}

	err = fmt.Errorf("unable to extract claims refres token")
	return
}
