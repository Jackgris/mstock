package middleware_test

import (
	"log"
	"net/http"

	"github.com/jackgris/mstock/middleware"
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

		BeforeEach(func() {
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

				It("Return status code 400", func() {
					Expect(200).To(gomega.Equal(400))
				})
			})

			Context("With authorization header", func() {

				It("Return status code 200", func() {
					Expect(400).To(gomega.Equal(200))
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
