package main_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/jackgris/mstock/models"
	"github.com/jackgris/mstock/server"
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
		var serve *server.Server
		var request *http.Request
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			dbName = "server_test"
			dburl := "localhost"
			session = newSession(dburl)
			serve = server.NewServer()
			recorder = httptest.NewRecorder()
		})

		AfterEach(func() {
			// Clean the database
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
					serve.Handler.ServeHTTP(recorder, request)
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
					serve.Handler.ServeHTTP(recorder, request)
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
					serve.Handler.ServeHTTP(recorder, request)
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

				BeforeEach(func() {
					user := models.User{}
					body, _ := json.Marshal(user)
					request, _ = http.NewRequest("POST", "/auth/login",
						bytes.NewReader(body))
				})

				It("Response should be 400 with empty input", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(400))

				})

				It("Response should be 400 with wrong data", func() {
					user := models.User{}
					putCorrectDataUser(&user)
					user.Name = "wrong"
					user.IdUser = "doesn't exist"
					user.Pass = "bad"

					body, _ := json.Marshal(user)
					request, _ = http.NewRequest("POST", "/auth/login",
						bytes.NewReader(body))
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(404))
				})
			})

			Context("With valid JSON", func() {

				var user models.User
				user = models.User{}
				putCorrectDataUser(&user)
				BeforeEach(func() {

					body, _ := json.Marshal(user)
					request, _ = http.NewRequest("POST", "/auth/login",
						bytes.NewReader(body))
				})

				It("Response should be 200", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(200))

				})
			})

		})
	})
})

func putCorrectDataUser(user *models.User) {

	date := time.Now()
	user.IdUser = "esto4es4una4buena4prueba4login"
	user.Name = "Juan Login"
	user.Pass = "asdasf123124sdmk09i0342"
	user.Email = "juan@gmail.com"
	user.LastLogin = date.Add(time.Minute * (-5))
	user.CreatedAt = date.AddDate(-1, 0, 0)
	user.UpdateAt = date.Add(time.Hour * (-5))
	token, err := models.GenerateToken(user.Name, user.Pass)

	user.Token = token
	err = user.Save()
	if err != nil {
		log.Println("Can't save user for test login", err)
	}
}

func DecodeToken(r io.ReadCloser) (*models.Token, error) {
	defer r.Close()
	var t models.Token
	err := json.NewDecoder(r).Decode(&t)
	return &t, err
}

func newSession(dburl string) *mgo.Session {
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
