package main_test

import (
	"time"

	"github.com/jackgris/mstock/models"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestMstock(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mstock Suite")
}

// Define factories: for a valid user,
// and for an invalid one
func defineFactories() {
	date := time.Now()
	gory.Define("userOk", models.User{},
		func(factory gory.Factory) {
			factory["IdUser"] = "esto1234es12una44prueba"
			factory["Name"] = "Juan"
			factory["Pass"] = "asdasf123124sdmk09i0342"
			factory["Email"] = "hola@gmail.com"
			factory["LastLogin"] = date.Add(time.Minute * (-5))
			factory["CreatedAt"] = date.AddDate(-1, 0, 0)
			factory["UpdateAt"] = date.Add(time.Hour * (-5))
		})
	gory.Define("userBad", models.User{},
		func(factory gory.Factory) {
			factory["IdUser"] = "esto1234es12una44prueba"
			factory["Name"] = "Juan123"
			factory["Pass"] = "asdasf1231@@24sdmk09i0342"
			factory["Email"] = "estonoesunmailcom"
			factory["LastLogin"] = date.Add(time.Minute * 5)
			factory["CreatedAt"] = date.AddDate(2, 0, 0)
			factory["UpdateAt"] = date.Add(time.Hour * 5)
		})
}
