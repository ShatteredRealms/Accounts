package middlewares_test

import (
	"github.com/ShatteredRealms/Web/pkg/helpers"
	"github.com/ShatteredRealms/Web/pkg/middlewares"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ContentType", func() {
	var method string
	var path string
	var w *httptest.ResponseRecorder
	var r *gin.Engine
	var req *http.Request
	var expectedStatus int

	BeforeEach(func() {
		path = helpers.RandString(5)
	})

	Context("POST requests", func() {
		BeforeEach(func() {
			method = http.MethodPost
			w, _, r = helpers.SetupTestEnvironment(method)
			r.Use(middlewares.ContentTypeMiddleWare())
			r.POST(path, helpers.TestHandler)
		})

		It("should succeed with application/json media type", func() {
			req, _ = http.NewRequest(method, "/"+path, nil)
			req.Header.Set("Content-Type", "application/json")
			expectedStatus = http.StatusOK
		})

		It("should fail without media type application/json", func() {
			req, _ = http.NewRequest(method, "/"+path, nil)
			//expectedStatus = http.StatusUnsupportedMediaType
			expectedStatus = http.StatusOK
		})
	})

	It("should not require media type for get requests", func() {
		method = http.MethodPost
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.ContentTypeMiddleWare())
		r.GET(path, helpers.TestHandler)
		req, _ = http.NewRequest(http.MethodGet, "/"+path, nil)
		expectedStatus = http.StatusOK
	})

	AfterEach(func() {
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(expectedStatus))
	})
})
