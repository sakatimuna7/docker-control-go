package middleware

import (
	"docker-control-go/src/helpers"
	"docker-control-go/src/services"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func ActivityLogger(action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			return c.Next() // Lewati middleware jika tidak ada token
		}

		// Pastikan token dalam format "Bearer <token>"
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		tokenString = strings.TrimSpace(tokenString) // Hilangkan spasi ekstra

		if tokenString == "" {
			return helpers.ErrorResponse(c, 401, "Unauthorized", "Invalid token format")
		}

		// Parsing JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})

		// Cek apakah token valid
		if err != nil || !token.Valid {
			return c.Next()
		}

		// Ambil claims dari token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Next()
		}

		// Ambil user_id dan validasi
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Next()
		}

		userID := int64(userIDFloat)

		// Simpan aktivitas
		services.LogActivity(userID, action)

		return c.Next()
	}
}
