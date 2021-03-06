package models_test

import (
	"errors"
	"strings"

	"github.com/jackgris/mstock/models"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Authentication", func() {

	Expect := gomega.Expect
	Describe := ginkgo.Describe
	Context := ginkgo.Context
	It := ginkgo.It
	BeforeEach := ginkgo.BeforeEach
	// AfterEach := ginkgo.AfterEach

	Describe("Token", func() {

		var email string
		var pass string
		var data string
		var token models.Token
		var err error

		Describe("Generating token", func() {

			Context("With valid data", func() {

				BeforeEach(func() {
					email = "hola@gmail.com"
					pass = "123hola456"
					data = email + pass
					token, err = models.GenerateToken(data)
				})

				It("Generate token", func() {
					Expect(err).To(gomega.BeNil())
					Expect(token).ShouldNot(gomega.BeZero())
				})

				It("Create correct size", func() {
					n := strings.SplitN(token.Hash, ".", 3)
					Expect(len(n)).To(gomega.Equal(3))
				})
			})

			Context("With invalid data", func() {

				It("Invalid token and message error", func() {
					badData := 1234
					token, err := models.GenerateToken(badData)
					errTest := errors.New("invalid data, you can not create the token")
					Expect(token).To(gomega.BeZero())
					Expect(err).To(gomega.Equal(errTest))
				})
			})
		})
	})
})
