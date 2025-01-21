package middleware

import (
	"fmt"
	"net/http"
	"strings"
)

func IsAuthenticated(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			next.ServeHTTP(w, r)
			return
		}

		token = strings.TrimPrefix(token, "Bearer ")
		fmt.Println(token)
		next.ServeHTTP(w, r)
	})
}
