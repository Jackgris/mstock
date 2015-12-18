package middleware_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/jackgris/mstock/middleware"
	"github.com/jackgris/mstock/models"
	"github.com/modocache/gory"
	"github.com/unrolled/render"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var dbName string = "server_test"
var session *mgo.Session

var _ = ginkgo.Describe("Middleware", func() {

	Expect := gomega.Expect
	Describe := ginkgo.Describe
	Context := ginkgo.Context
	It := ginkgo.It
	BeforeEach := ginkgo.BeforeEach
	AfterEach := ginkgo.AfterEach

	Describe("Server", func() {

		var serve http.Server
		var request *http.Request
		var recorder *httptest.ResponseRecorder

		BeforeEach(func() {
			dburl := "localhost"
			session = newSession(dburl)
			recorder = httptest.NewRecorder()
			s := http.NewServeMux()
			s.Handle("/fake", middleware.AuthMiddleware(fakeHandler{}))
			serve = http.Server{
				Addr:    ":8080",
				Handler: s,
			}
			serve.ListenAndServe()
		})

		AfterEach(func() {
			// Clean the database
			session.DB(dbName).DropDatabase()
			session.Close()
		})

		Describe("GET /fake URL created to perform the test", func() {

			Context("Without authorization header", func() {

				BeforeEach(func() {
					body := ""
					request, _ = http.NewRequest("GET", "/fake",
						strings.NewReader(body))
				})

				It("Return status code 400", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(400))
				})
			})

			Context("With authorization header", func() {

				BeforeEach(func() {
					token, _ := models.GenerateToken("test")
					body := ""
					request, _ = http.NewRequest("GET", "/fake",
						strings.NewReader(body))
					request.Header.Add("Authorization", "Bearer "+token.Hash)
				})

				It("Return status code 200", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(200))
				})
			})

			Context("With a real user", func() {

				var user *models.User

				BeforeEach(func() {
					user = gory.Build("userOk").(*models.User)
					token, _ := models.GenerateToken("testOk")
					user.Token = token
					saveEntity("users", user)
				})

				It("Check if user with token is saved on the database", func() {
					Expect(1).To(gomega.Equal(2))
				})
			})
		})
	})
})

// This entity will use to simulate a URL that requires authentication to be accessed
type fakeHandler struct{}

func (f fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := render.New()
	rd.JSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

// With this function we will be able to save data for test
func saveEntity(table string, e interface{}) {
	c := session.DB(dbName).C(table)
	insert := bson.M{"$set": e}
	err := c.Insert(insert)
	if err != nil {
		log.Println("Middelware test: saveEntity", err)
	}
}

func newSession(dburl string) *mgo.Session {
	s, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	return s
}
