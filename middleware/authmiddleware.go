package middleware

import (
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	models "github.com/jackgris/mstock/models"
	"github.com/unrolled/render"
)

func AuthMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it has authorization header
		h := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(h) != 2 || h[0] != "Bearer" {
			message := "hasn't authorization header"
			log.Println("Middleware auth:", message)
			e := errorH{}
			e.message = message
			e.ServeHTTP(w, r)
			return
		}

		// Check if has an valid token
		token := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		payload, err := models.ParseToken(token[1])

		// Verify diferents errors type, than can made from the token
		switch err.(type) {
		case nil:
			if !payload.Valid {
				message := "Invalid payload " + err.Error()
				log.Println("Middleware auth: ", message)
				e := errorH{}
				e.message = message
				e.ServeHTTP(w, r)
				return
			}
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				message := "Token expired " + err.Error()
				log.Println("Middleware auth: ", message)
				e := errorH{}
				e.message = message
				e.ServeHTTP(w, r)
				return
			default:
				message := "Error validation " + err.Error()
				log.Println("Middleware auth: ", message)
				e := errorH{}
				e.message = message
				e.ServeHTTP(w, r)
				return
			}
		default:
			if err != nil {
				message := "Error payload " + err.Error()
				log.Println("Middleware auth: ", message)
				e := errorH{}
				e.message = message
				e.ServeHTTP(w, r)
				return
			}
		}

		// Check for email in the database
		if email, ok := payload.Claims["sub"].(string); !ok {
			message := "Error get email from token"
			log.Println("Middleware auth: ", message)
			e := errorH{}
			e.message = message
			e.ServeHTTP(w, r)
			return
		} else {
			// FIXME unimplented
			log.Println("Middleware auth email:", email)
		}

		// If everything is okay, we return response Handler
		handler.ServeHTTP(w, r)
	})
}

// Handler to handle errors in requests that need to authenticacion
type errorH struct {
	message string
}

func (e errorH) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := render.New()
	rd.JSON(w, http.StatusBadRequest, map[string]string{"token": "", "message": e.message})
}
