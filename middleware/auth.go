package middleware

import (
	"fmt"
	"mysqlapi/controllers"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

var AuthenticateJWT fiber.Handler = func(c *fiber.Ctx) error {
	const BEARER_SCHEMA = "Bearer"
	authHeader := c.Get("Authorization")
	fmt.Println("authHeader")
	fmt.Println(authHeader)
	if len(authHeader) > 0 {
		tokenString := authHeader[len(BEARER_SCHEMA):]
		token, err := controllers.ValidateToken((strings.Trim(tokenString, " ")))
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			fmt.Println(claims["email"])
			return c.Next()
		} else {
			fmt.Println(err)
			body := make(map[string]interface{})
			body["message"] = "Unauthorize"
			return c.Status(http.StatusUnauthorized).JSON(body)
		}
	} else {
		body := make(map[string]interface{})
		body["message"] = "Token required"
		return c.Status(http.StatusUnauthorized).JSON(body)
	}

}
