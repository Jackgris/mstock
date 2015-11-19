package controllers

import (
	"log"
	"net/http"

	"github.com/unrolled/render"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	rd := render.New()
	user, err := DecodeUserData(r.Body)
	if err != nil {
		log.Println("Unmarshal json register error", err)
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}

	if !user.Valid() {
		log.Println("User data are invalid!")
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}
	// we show standard output data requests
	log.Println("RegisterHandler", r.Method, "path", r.URL.Path, "body", r.Body)
	log.Println("RegisterHandler", "unmarshal json", user)

	rd.JSON(w, http.StatusOK, map[string]string{"token": ""})
}
