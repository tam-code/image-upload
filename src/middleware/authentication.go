package middleware

import (
	"log"
	"net/http"
)

var tokenUsers map[string]string = map[string]string{
	"00000000": "user0",
	"aaaaaaaa": "userA",
	"05f717e5": "randomUser",
}

func ValidateSecretToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secretToken := r.Header.Get("X-Secret-Token")
		user, ok := tokenUsers[secretToken]
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		log.Printf("Authenticated user %s\n", user)

		next.ServeHTTP(w, r)
	})
}
