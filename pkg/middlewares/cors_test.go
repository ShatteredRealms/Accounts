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

	BeforeEach(func() {
		path = helpers.RandString(5)
	})

	methods := map[string]int{
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

	for method, expectedStatus := range methods {
		It(fmt.Sprintf("should respond to %s requests with %d status code", method, expectedStatus), func() {
			testMethod := method
			testExpectedStatus := expectedStatus
			w, _, r = helpers.SetupTestEnvironment(testMethod)
			r.Use(CORSMiddleWare())
			req, _ = http.NewRequest(testMethod, "/"+path, nil)
			r.Handle(testMethod, path, helpers.TestHandler)
			r.ServeHTTP(w, req)
			Expect(w.Code).To(Equal(testExpectedStatus))
		})
	}
})
