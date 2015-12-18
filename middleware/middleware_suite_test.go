package middleware_test

import (
	"time"

	"github.com/jackgris/mstock/models"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMiddleware(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Middleware Suite")
}

// Define factories: for a valid user,
// and for an invalid one
func defineFactories() {
	date := time.Now()
	gory.Define("userOk", models.User{},
		func(factory gory.Factory) {
			factory["IdUser"] = "123456"
			factory["Name"] = "Existe"
			factory["Pass"] = "el1usuario1esta1enla1basededatos"
			factory["Email"] = "hola@gmail.com"
			factory["LastLogin"] = date
			factory["CreatedAt"] = date
			factory["UpdateAt"] = date
		})

	gory.Define("userBad", models.User{},
		func(factory gory.Factory) {
			factory["IdUser"] = "noexiste1234"
			factory["Name"] = "NoExisteElUsuario"
			factory["Pass"] = "asdasf1231@@24sdmk09i0342"
			factory["Email"] = "estonoesunmailcom"
			factory["LastLogin"] = date
			factory["CreatedAt"] = date
			factory["UpdateAt"] = date
		})
}
