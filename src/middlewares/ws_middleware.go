package middlewares

import (
	"context"
	"docker-control-go/src/constant"
	"encoding/json"
	"errors"
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
		// Baca pesan pertama dari klien
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("❌ Error reading authentication message:", err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized"}`))
			c.Close()
			return
		}

		// Pastikan JSON valid sebelum parsing
		if !json.Valid(message) {
			log.Println("❌ JSON is not valid")
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid JSON format"}`))
			c.Close()
			return
		}

		// Parsing JSON token
		var authMsg AuthMessage
		if err := json.Unmarshal(message, &authMsg); err != nil {
			log.Println("❌ Invalid authentication message format:", err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Invalid authentication format"}`))
			c.Close()
			return
		}

		// Validasi JWT
		claims, err := ValidateJWT(authMsg.Token)
		if err != nil {
			log.Println("❌ JWT validation failed:", err)
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized: Invalid Token"}`))
			c.Close()
			return
		}

		// Pastikan userID dan userRole ada dalam claims sebelum konversi
		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			log.Println("❌ Error: userID not found in token")
			c.WriteMessage(websocket.TextMessage, []byte(`{"error": "Unauthorized: Missing userID"}`))
			c.Close()
			return
		}

		userRole, _ := claims["userRole"].(string) // Bisa kosong jika tidak ada role

		log.Println("✅ User authenticated:", userID, "Role:", userRole)

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
		return nil, errors.New("JWT secret key is missing")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(SecretKey), nil
	})

	// Cek validitas token
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Ambil claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Cek apakah token expired
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("token expiration missing")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}
