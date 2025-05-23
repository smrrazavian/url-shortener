package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/smrrazavian/url-shortener/internal/config"
)

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(auth, "Bearer ")
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}
		ok := validateJWT(tokenString, cfg.HmacSecret)
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func validateJWT(token string, secret []byte) bool {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return false
	}

	headerAndPayload := parts[0] + "." + parts[1]
	signature := parts[2]

	expectedSig := signHS256(headerAndPayload, secret)
	if !hmac.Equal([]byte(signature), []byte(expectedSig)) {
		return false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return false
	}

	var payload struct {
		Exp int64 `json:"exp"`
	}

	if err := json.Unmarshal(payloadBytes, &payload); err != nil {
		return false
	}

	if payload.Exp != 0 && time.Now().Unix() > payload.Exp {
		return false
	}

	return true

}

func signHS256(data string, secret []byte) string {
	h := hmac.New(sha256.New, secret)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}
