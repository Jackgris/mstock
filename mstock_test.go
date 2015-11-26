package main_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"

	mstock "github.com/jackgris/mstock"
	models "github.com/jackgris/mstock/models"
	"github.com/modocache/gory"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
)

var _ = ginkgo.Describe("Mstock", func() {

	Expect := gomega.Expect
	Describe := ginkgo.Describe
	Context := ginkgo.Context
	It := ginkgo.It
	BeforeEach := ginkgo.BeforeEach
	AfterEach := ginkgo.AfterEach

	Describe("Server", func() {

		var dbName string
		var session *mgo.Session
		var server *mstock.Server
		var request *http.Request
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			dbName = "server_test"
			dburl := "localhost"
			session = NewSession(dburl)
			server = mstock.NewServer()
			recorder = httptest.NewRecorder()
		})

		AfterEach(func() {
			session.DB(dbName).DropDatabase()
			session.Close()
		})

		Describe("POST /auth/signup", func() {

			Context("With invalid JSON", func() {
				BeforeEach(func() {
					body, _ := json.Marshal(gory.Build("userBad"))
					request, _ = http.NewRequest("POST", "/auth/signup",
						bytes.NewReader(body))
				})

				It("Return a status code of 400", func() {
					server.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(400))
				})

				It("User must not be saved in the database", func() {
					user := gory.Build("userBad").(*models.User)
					userTest := models.User{
						IdUser: user.IdUser,
					}
					userTest, err := userTest.Get()
					errTest := errors.New("not found")
					Expect(err).To(gomega.Equal(errTest))
					Expect(userTest).To(gomega.BeZero())
				})
			})

			Context("With valid JSON", func() {

				BeforeEach(func() {
					body, _ := json.Marshal(gory.Build("userOk"))
					request, _ = http.NewRequest("POST", "/auth/signup",
						bytes.NewReader(body))
				})

				It("Return a status code of 200", func() {
					server.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(200))
				})

				It("User must be saved in the database", func() {
					user := gory.Build("userOk").(*models.User)
					userTest := models.User{
						IdUser: user.IdUser,
					}
					userTest, err := userTest.Get()

					Expect(err).To(gomega.BeNil())
					Expect(user.Name).To(gomega.Equal(userTest.Name))
				})

				It("Response need have hash token", func() {
					server.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(200))
					Expect(recorder.HeaderMap["Content-Type"][0]).
						To(gomega.ContainSubstring("application/json; charset=UTF-8"))

					data := myCloser{bytes.NewBufferString(recorder.Body.String())}
					token, err := DecodeToken(data)
					Expect(err).To(gomega.BeNil())
					Expect(token).ShouldNot(gomega.BeZero())
				})
			})
		})

		Describe("POST /auth/login", func() {

			Context("With invalid JSON", func() {

			})

			Context("With valid JSON", func() {

			})

		})
	})
})

func DecodeToken(r io.ReadCloser) (*models.Token, error) {
	defer r.Close()
	var t models.Token
	err := json.NewDecoder(r).Decode(&t)
	return &t, err
}

func NewSession(dburl string) *mgo.Session {
	s, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	return s
}

type myCloser struct {
	io.Reader
}

func (myCloser) Close() error { return nil }
