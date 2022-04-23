package middlewares

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/pkg/helpers"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CORS", func() {
	var path string
	var w *httptest.ResponseRecorder
	var r *gin.Engine
	var req *http.Request

	var methods map[string]int

	BeforeEach(func() {
		path = helpers.RandString(5)
		methods = map[string]int{
			http.MethodConnect: http.StatusOK,
			http.MethodDelete:  http.StatusOK,
			http.MethodGet:     http.StatusOK,
			http.MethodHead:    http.StatusOK,
			http.MethodPatch:   http.StatusOK,
			http.MethodPost:    http.StatusOK,
			http.MethodPut:     http.StatusOK,
			http.MethodTrace:   http.StatusOK,
			http.MethodOptions: http.StatusNoContent,
		}
	})

	for method, expectedStatus := range methods {
		When("input is "+method, func() {
			It(fmt.Sprintf("should respond with status code %d", expectedStatus), func() {
				w, _, r = helpers.SetupTestEnvironment(method)
				r.Use(CORSMiddleWare())
				req, _ = http.NewRequest(method, "/"+path, nil)
				r.Handle(method, path, helpers.TestHandler)
				r.ServeHTTP(w, req)
				Expect(w.Code).To(Equal(expectedStatus))
			})
		})
	}
})
