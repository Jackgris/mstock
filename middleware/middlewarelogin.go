package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/unrolled/render"
)

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if it has authorization header
		h := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(h) != 2 || h[0] != "Bearer" {
			log.Println("Middleware auth: hasn't authorization header")
			// redirect?
			e := errorH{}
			e.ServeHTTP(w, r)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

func returnError() http.Handler {
	return errorH{}
}

type errorH struct{}

func (e errorH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := render.New()
	rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
}
