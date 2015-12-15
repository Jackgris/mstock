package middleware

import (
	"log"
	"net/http"
	"strings"
)

func AuthMiddleware(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if it has authorization header
		h := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(h) != 2 || h[0] != "Bearer" {
			log.Println("Middleware auth: hasn't authorization header")
			// redirect?
			return
		}
	})
}
