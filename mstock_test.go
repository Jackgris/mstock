package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	mstock "github.com/jackgris/mstock"
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

				It("return a status code of 400", func() {
					server.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(400))
				})
			})

			Context("With valid JSON", func() {

				BeforeEach(func() {
					body, _ := json.Marshal(gory.Build("userOk"))
					request, _ = http.NewRequest("POST", "/auth/signup",
						bytes.NewReader(body))
				})

				It("return a status code of 200", func() {
					server.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(200))
				})
			})
		})
	})
})

func NewSession(dburl string) *mgo.Session {
	s, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	return s
}
