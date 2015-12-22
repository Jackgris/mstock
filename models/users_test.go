package models_test

import (
	"errors"
	"time"

	"github.com/jackgris/mstock/models"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Users", func() {

	Expect := gomega.Expect
	Describe := ginkgo.Describe
	It := ginkgo.It
	Context := ginkgo.Context

	Describe("Interacting with database", func() {

		It("Saving or update data", func() {
			user := models.User{
				IdUser:    "test1234",
				Name:      "test",
				Pass:      "1234",
				LastLogin: time.Now(),
				CreatedAt: time.Now(),
				UpdateAt:  time.Now(),
			}

			err := user.Save()
			Expect(err).To(gomega.BeNil())
		})

		It("Get user with Id", func() {
			user := models.User{
				IdUser: "test1234",
			}

			user, err := user.Get()
			Expect(err).To(gomega.BeNil())
			Expect(user.Name).To(gomega.Equal("test"))
			Expect(user.Pass).To(gomega.Equal("1234"))
		})

		It("Delete user with Id", func() {
			user := models.User{
				IdUser: "test1234",
			}
			err := user.Delete()
			Expect(err).To(gomega.BeNil())
		})

		Context("If user deleted", func() {
			user := models.User{
				IdUser: "test1234",
			}

			It("We need received zero value for user and error", func() {
				user, err := user.Get()
				errTest := errors.New("not found")
				Expect(err).To(gomega.Equal(errTest))
				Expect(user.Name).To(gomega.BeZero())
				Expect(user.Pass).To(gomega.BeZero())
			})
		})
	})
})
