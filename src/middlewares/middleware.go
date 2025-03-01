package middlewares

import (
	"docker-control-go/src/database/models"
	"docker-control-go/src/helpers"
	"docker-control-go/src/services"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken - Membuat JWT Token untuk user
func GenerateToken(user *models.User) (string, error) {
	SecretKey := os.Getenv("SECRET_KEY")

	// Cek jika SecretKey kosong
	if SecretKey == "" {
		return "", fmt.Errorf("SECRET_KEY is not set in environment variables")
	}

	claims := jwt.MapClaims{
		"userID":   user.ID,
		"username": user.Username,
		"userRole": user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Expired dalam 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// Simpan session ke Redis
	services.StoreSessionInRedis(user.ID, signedToken, time.Hour*24)

	return signedToken, nil
}

// JWTMiddleware - Middleware untuk validasi JWT
func JWTMiddleware(c *fiber.Ctx) error {
	SecretKey := os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		return helpers.ErrorResponse(c, 500, "Internal Server Error", "JWT secret key is missing")
	}

	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return helpers.ErrorResponse(c, 401, "Unauthorized", "Missing Authorization Header")
	}

	// Pastikan token dalam format "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return helpers.ErrorResponse(c, 401, "Unauthorized", "Invalid Authorization Format")
	}

	tokenString := parts[1]

	// Parse Token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Pastikan metode tanda tangan adalah HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	// Token tidak valid
	if err != nil || !token.Valid {
		return helpers.ErrorResponse(c, 401, "Unauthorized", "Invalid Token")
	}

	// Ambil claims dari token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return helpers.ErrorResponse(c, 401, "Unauthorized", "Invalid Token Claims")
	}

	// Cek apakah token sudah expired
	expirationTime := int64(claims["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		return helpers.ErrorResponse(c, 401, "Unauthorized", "Token Expired")
	}

	// Simpan user info ke Fiber context
	c.Locals("userID", claims["userID"])
	c.Locals("userRole", claims["userRole"])

	return c.Next()
}
