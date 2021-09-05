package service

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type jwtAuthKeyType string

type JwtCustomClaim struct {
	ID string `json:"id" example:"d270f4e3-02ea-483d-6525-d7e7a22021507"`
	jwt.StandardClaims
}

var (
	jwtSecret  = []byte(os.Getenv("SECRET"))
	AccExpired = 60 * 15
	RemExpired = 60 * 60 * 24 * 7
)

func jwtCreate(id string, duration time.Duration) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtCustomClaim{
		ID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})

	token, err := t.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func JwtCreateAccessToken(id string) (string, error) {
	return jwtCreate(id, time.Second*time.Duration(AccExpired))
}

func JwtCreateRememberToken(id string) (string, error) {
	return jwtCreate(id, time.Second*time.Duration(RemExpired))
}

func JwtNew(id string) (accessToken string, rememberToken string, err error) {
	accessToken, err = JwtCreateAccessToken(id)
	if err != nil {
		return
	}

	rememberToken, err = JwtCreateRememberToken(id)
	if err != nil {
		return
	}

	return
}

func JwtValidate(token string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(token, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token signing method")
		}
		return jwtSecret, nil
	})
}

func RefreshToken(refreshToken string) (accessToken string, rememberToken string, err error) {
	dataClaim, err := JwtValidate(refreshToken)
	if err != nil || !dataClaim.Valid {
		if err == nil {
			err = ErrInvalidJWT
		}
		return
	}

	claims, _ := dataClaim.Claims.(*JwtCustomClaim)

	getUser, err := UserGetByID(context.Background(), claims.ID)
	if err != nil {
		return
	}

	if getUser.RememberToken != nil && refreshToken == *getUser.RememberToken {
		accessToken, rememberToken, err = JwtNew(claims.ID)
		if err != nil {
			return
		}

		_, err = UserUpdateRememberToken(context.Background(), claims.ID, rememberToken)
		if err != nil {
			return "", "", err
		}
	} else {
		err = errors.New("bad request")
		return
	}

	return
}
