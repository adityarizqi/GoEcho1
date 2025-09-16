package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func extractClaims(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func RoleMiddleware(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("token")
			if err != nil {
				if strings.Contains(err.Error(), "named cookie not present") {
					return c.Redirect(http.StatusSeeOther, "/login")
				}
				return c.String(http.StatusBadRequest, "bad request")
			}

			tokenString := cookie.Value
			claims, err := extractClaims(tokenString)
			if err != nil {
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			userRole := claims["role"].(string)
			for _, role := range roles {
				if role == userRole {
					c.Set("user_id", uint(claims["user_id"].(float64)))
					c.Set("role", claims["role"])
					return next(c)
				}
			}

			return c.String(http.StatusForbidden, "You are not authorized!")
		}
	}
}
