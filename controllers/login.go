package controllers

import (
	"fmt"
	"log"
	"net/http"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// we show standard output data requests
	log.Println(r.Method, "path", r.URL.Path)
	// not implemented
	fmt.Fprintf(w, "login", "hola")
}
