package controllers

import (
	"log"
	"net/http"

	"github.com/jackgris/mstock/models"
	"github.com/unrolled/render"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {

	// We obtain user data, and create our user with them
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

	token, err := models.GenerateToken(user.Name, "#", user.Pass)
	if err != nil {
		log.Println("Can't create token!", err)
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}

	user.Token = token
	// keep user data and token we created in the database
	err = user.Save()
	if err != nil {
		log.Println("Can't save the user on the database!")
		rd.JSON(w, http.StatusBadRequest, map[string]string{"token": ""})
		return
	}

	// if everything went correctly, we send the token to the client
	rd.JSON(w, http.StatusOK, map[string]string{"token": user.Token.Hash})
}
