package middlewares

import (
	"book/helper"
	"book/services"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func AuthorizeJWT(jwtService services.JWTService) gin.HandlerFunc {

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := helper.BuildErrorResponse("Failed To Process Request", "No token found", nil)
			c.AbortWithStatusJSON(http.StatusBadRequest, response)
			return
		}


		splitedAuthHeader := strings.Fields(authHeader)
		tokenAuth := strings.Join(splitedAuthHeader[1:], "")
		token, err := jwtService.ValidateToken(tokenAuth)

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			log.Println("Claim user id ", claims["user_id"])
			log.Println("Claim issuer ", claims["issuer"])
		} else {
			log.Println(err)
			response := helper.BuildErrorResponse("token is not valid", err.Error(), nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
		}

	}
}
