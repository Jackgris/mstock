package models_test

import (
	"time"

	. "github.com/jackgris/mstock/models"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Users", func() {

	Expect := gomega.Expect
	Describe := ginkgo.Describe
	It := ginkgo.It

	Describe("Interacting with database", func() {

		It("Saving or update data", func() {
			user := User{
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
			user := User{
				IdUser: "test1234",
			}

			user, err := user.Get()
			Expect(err).To(gomega.BeNil())
			Expect(user.Name).To(gomega.Equal("test"))
			Expect(user.Pass).To(gomega.Equal("1234"))
		})

		It("Delete user with Id", func() {
			user := User{
				IdUser: "test1234",
			}
			err := user.Delete()
			Expect(err).To(gomega.BeNil())
		})
	})
})
