package middlewares_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/productivestudy/auth/internal/requests"
	"github.com/productivestudy/auth/pkg/middlewares"
	"github.com/productivestudy/auth/tests/helpers"
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

		It("should accept application/json media type", func() {
			req, _ = http.NewRequest(method, "/"+path, nil)
			req.Header.Set(requests.ContentType, requests.JSONContent)
			expectedStatus = http.StatusOK
		})

		It("should fail if media type is not application/json", func() {
			req, _ = http.NewRequest(method, "/"+path, nil)
			expectedStatus = http.StatusUnsupportedMediaType
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
