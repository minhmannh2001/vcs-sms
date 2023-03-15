package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/minhmannh2001/sms/service"
)

// AuthorizeJWT validates the token from the http request, returning a 401 if it's not valid
func AuthorizeJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer "
		authHeader := ctx.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]
		// tokenString, err := ctx.Cookie("token")
		// if err != nil {
		// 	if err == http.ErrNoCookie {
		// 		// If the cookie is not set, return an unauthorized status
		// 		ctx.AbortWithStatus(http.StatusUnauthorized)
		// 		return
		// 	}
		// 	// For any other type of error, return a bad request status
		// 	ctx.AbortWithStatus(http.StatusUnauthorized)
		// 	return
		// }
		token, err := service.NewJWTService().ValidateToken(tokenString)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			exp := claims["exp"].(float64)
			if int64(exp) < time.Now().Local().Unix() {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token Expired"})
				return
			}
			log.Println("Claims[Name]: ", claims["name"])
			log.Println("Claims[Admin]: ", claims["admin"])
			log.Println("Claims[Issuer]: ", claims["iss"])
			log.Println("Claims[IssuedAt]: ", claims["iat"])
			log.Println("Claims[ExpiresAt]: ", claims["exp"])
		} else {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
