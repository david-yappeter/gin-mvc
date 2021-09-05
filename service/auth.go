package service

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

var jwtAuthKey = jwtAuthKeyType("jwt-auth")

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken, err := c.Cookie("access_token")
		if err == http.ErrNoCookie {
			rememberToken, err := c.Cookie("remember_token")
			if err == http.ErrNoCookie {
				c.Next()
				return
			}
			accessToken, rememberToken, err = RefreshToken(rememberToken)
			setToken(c, accessToken, rememberToken)
			if err != nil {
				deleteToken(c)
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		dataClaim, err := JwtValidate(accessToken)
		if err != nil {
			deleteToken(c)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, _ := dataClaim.Claims.(*JwtCustomClaim)

		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), jwtAuthKey, claims))

		c.Next()
	}
}

func CtxVal(ctx context.Context) *JwtCustomClaim {
	raw, _ := ctx.Value(jwtAuthKey).(*JwtCustomClaim)
	return raw
}

func setToken(c *gin.Context, accessToken string, rememberToken string) {
	c.SetCookie("access_token", accessToken, AccExpired, "", "", false, true)
	c.SetCookie("remember_token", rememberToken, RemExpired, "", "", false, true)
}

func deleteToken(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "", "", false, true)
	c.SetCookie("remember_token", "", -1, "", "", false, true)
}
