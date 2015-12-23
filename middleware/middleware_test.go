package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"gopkg.in/mgo.v2"

	"github.com/jackgris/mstock/middleware"
	"github.com/jackgris/mstock/models"
	"github.com/modocache/gory"
	"github.com/unrolled/render"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

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
			models.DB_NAME = "server_test"
		})

		AfterEach(func() {
			// Clean the database
			session.DB(models.DB_NAME).DropDatabase()
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

			Context("With authorization header but a user who doesn't exist in the database", func() {

				BeforeEach(func() {
					token, _ := models.GenerateToken("test")
					body := ""
					request, _ = http.NewRequest("GET", "/fake",
						strings.NewReader(body))
					request.Header.Add("Authorization", "Bearer "+token.Hash)
				})

				It("Return status code 400", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(400))
				})
			})

			Context("With a real user", func() {

				var user *models.User

				BeforeEach(func() {
					user = gory.Build("userOk").(*models.User)
					token, _ := models.GenerateToken(user.Email + "#" + user.Pass)
					user.Token = token
					user.Save()

					body := ""
					request, _ = http.NewRequest("GET", "/fake",
						strings.NewReader(body))
					request.Header.Add("Authorization", "Bearer "+token.Hash)
				})

				It("Check if user with token is saved on the database", func() {
					chkUser := models.User{}
					chkUser.IdUser = user.IdUser
					chkUser, err := chkUser.Get()

					Expect(err).To(gomega.BeNil())
					Expect(chkUser).ShouldNot(gomega.BeZero())
				})

				It("Return status code 200", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(200))
				})
			})

			Context("With a user who does not exist in the database", func() {

				var user *models.User

				BeforeEach(func() {
					user = gory.Build("userBad").(*models.User)
					token, _ := models.GenerateToken("testOk")
					user.Token = token

					body := ""
					request, _ = http.NewRequest("GET", "/fake",
						strings.NewReader(body))
					request.Header.Add("Authorization", "Bearer "+token.Hash)
				})

				It("Check if user with token is saved on the database", func() {
					chkUser := models.User{}
					chkUser.IdUser = user.IdUser
					chkUser, err := chkUser.Get()

					Expect(err).ShouldNot(gomega.BeNil())
					Expect(chkUser).Should(gomega.BeZero())
				})

				It("Return status code 400", func() {
					serve.Handler.ServeHTTP(recorder, request)
					Expect(recorder.Code).To(gomega.Equal(400))
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

// We create a session in the database to perform the test
func newSession(dburl string) *mgo.Session {
	s, err := mgo.Dial(dburl)
	if err != nil {
		panic(err)
	}
	return s
}
