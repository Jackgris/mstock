package controllers

import (
	"encoding/json"
	"io"

	"github.com/jackgris/mstock/models"
)

// This decode input data on json format and return the User with that data
func DecodeUserData(r io.ReadCloser) (*models.User, error) {
	defer r.Close()
	var u models.User
	err := json.NewDecoder(r).Decode(&u)
	return &u, err
}
