package controllers

import (
	"log"
	"net/http"

	"github.com/jackgris/mstock/models"
	"github.com/unrolled/render"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	// We obtain user data, and create our user with them
	rd := render.New()
	user, err := DecodeUserData(r.Body)
	if err != nil {
		log.Println("Unmarshal json register error", err)
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}

	if !user.Valid(false) {
		log.Println("User data in login are invalid")
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}

	// if the user is registered, must exist in the database
	userSaved, err := user.Get()
	if err != nil && len(userSaved.Name) < 0 {
		log.Println("User attempting to login does not exist in the database")
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}

	userSaved.Token, err = models.GenerateToken(userSaved.Email, userSaved.Pass)
	if err != nil {
		log.Println("User can't create token on login")
		rd.JSON(w, http.StatusNotFound, map[string]string{"token": ""})
		return
	}

	// keep data token we created in the database
	if err = userSaved.Save(); err != nil {
		log.Println("User can't save the new token on login")
		rd.JSON(w, http.StatusNotFound, map[string]string{"token": ""})
		return
	}

	// if everything went correctly, we send the token to the client
	rd.JSON(w, http.StatusOK, map[string]string{"token": userSaved.Token.Hash})
}
