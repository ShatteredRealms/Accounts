package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/productivestudy/auth/tests/helpers"
)

var _ = Describe("CORS", func() {
	var method string
	var path string
	var w *httptest.ResponseRecorder
	var r *gin.Engine
	var req *http.Request
	var expectedStatus int

	BeforeEach(func() {
		path = helpers.RandString(5)
	})

	methods := []string{
		http.MethodConnect,
		http.MethodDelete,
		http.MethodGet,
		http.MethodHead,
		http.MethodPatch,
		http.MethodPost,
		http.MethodPut,
		http.MethodTrace,
	}
	for _, v := range methods {
		method = v
		It(fmt.Sprintf("should allow %s requests", method), func() {
			w, _, r = helpers.SetupTestEnvironment(method)
			r.Use(CORSMiddleWare())
			req, _ = http.NewRequest(method, "/"+path, nil)
			expectedStatus = http.StatusOK
		})
	}

	It("should abort for option requests", func() {
		method = http.MethodOptions
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		expectedStatus = http.StatusNoContent
	})

	AfterEach(func() {
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(expectedStatus))
	})
})
