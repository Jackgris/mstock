package middleware_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/jackgris/mstock/middleware"
	"github.com/jackgris/mstock/models"
	"github.com/unrolled/render"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

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
			log.Println("finish test")
		})

		Describe("GET /fake", func() {

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
		})
	})
})

type fakeHandler struct{}

func (f fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rd := render.New()
	rd.JSON(w, http.StatusOK, map[string]string{"ok": "true"})
}
