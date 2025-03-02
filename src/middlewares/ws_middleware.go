package middlewares

import (
	"context"
	"docker-control-go/src/constant"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Struktur untuk parsing pesan pertama (token)
type AuthMessage struct {
	Token string `json:"token"`
}

// Middleware WebSocket untuk autentikasi
func WebSocketAuthMiddleware(next func(*websocket.Conn, context.Context)) func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		defer c.Close()

		// Baca pesan pertama dari klien
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("❌ Error reading authentication message:", err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized"}`))
			return
		}

		// Pastikan JSON valid sebelum parsing
		if !json.Valid(message) {
			log.Println("❌ JSON is not valid")
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid JSON format"}`))
			return
		}

		// Parsing JSON token
		var authMsg AuthMessage
		if err := json.Unmarshal(message, &authMsg); err != nil {
			log.Println("❌ Invalid authentication message format:", err)
			log.Println("📌 Raw message:", string(message))
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid authentication format"}`))
			return
		}

		// Validasi JWT
		claims, err := ValidateJWT(authMsg.Token)
		if err != nil {
			log.Println("❌ JWT validation failed:", err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized: Invalid Token"}`))
			return
		}

		// Ambil informasi user dari token
		userID := claims["userID"].(string)
		userRole := claims["userRole"].(string)
		log.Println("User authenticated:", userID, "Role:", userRole)

		// Simpan userID ke dalam context
		ctx := context.WithValue(context.Background(), constant.UserIDKey, userID)
		ctx = context.WithValue(ctx, constant.UserRoleKey, userRole)

		// Kirim konfirmasi sukses ke klien
		successMessage := fmt.Sprintf(`{"message": "Authentication successful", "userID": "%s", "role": "%s"}`, userID, userRole)
		c.WriteMessage(websocket.TextMessage, []byte(successMessage))
		// Jika valid, lanjutkan ke handler utama dengan claims
		next(c, ctx)
	}
}

// Fungsi untuk validasi JWT
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	SecretKey := os.Getenv("SECRET_KEY")
	if SecretKey == "" {
		return nil, log.Output(1, "JWT secret key is missing")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, log.Output(1, "unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	// Cek validitas token
	if err != nil || !token.Valid {
		return nil, log.Output(1, "invalid token")
	}

	// Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, log.Output(1, "invalid token claims")
	}

	// Cek apakah token expired
	expirationTime := int64(claims["exp"].(float64))
	if time.Now().Unix() > expirationTime {
		return nil, log.Output(1, "token expired")
	}

	return claims, nil
}
