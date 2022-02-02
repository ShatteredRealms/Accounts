package helpers_test

import (
	"github.com/ShatteredRealms/Web/pkg/helpers"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Requests", func() {
	It("should process relative urls corretly", func() {
		path := helpers.RandString(10)
		Expect(helpers.URLFromRelative(path)).To(Equal("/api/v1/" + path))
		Expect(helpers.URLFromRelative("/" + path)).To(Equal("/api/v1/" + path))
	})

	It("should setup gin text environment correctly", func() {
		method := http.MethodGet
		w, c, r := helpers.SetupTestEnvironment(method)
		Expect(w).NotTo(BeNil())
		Expect(c).NotTo(BeNil())
		Expect(r).NotTo(BeNil())
		Expect(gin.Mode()).To(Equal(gin.TestMode))
		Expect(c.Request.Method).To(Equal(method))
		Expect(c.Request.Header).NotTo(BeNil())
	})

	Context("Serving", func() {
		var method string
		var path string
		var w *httptest.ResponseRecorder
		var r *gin.Engine
		var body interface{}

		BeforeEach(func() {
			body = nil
			method = http.MethodPost
			path = "/" + helpers.RandString(5)
			w, _, r = helpers.SetupTestEnvironment(method)
			r.Handle(method, path, helpers.TestHandler)
		})

		It("should work with a body", func() {
			body = helpers.ExamplePayload{Message: helpers.RandString(10)}
		})

		It("should work without a body", func() {})

		AfterEach(func() {
			err := helpers.ServeRequest(method, path, body, r, w)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(w.Code).To(Equal(http.StatusOK))
		})
	})
})
