package middlewares_test

import (
	"github.com/ShatteredRealms/Accounts/pkg/helpers"
	"github.com/ShatteredRealms/Accounts/pkg/middlewares"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	path    string
	methods map[string]int
)
var _ = BeforeSuite(func() {
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

var _ = Describe("CORS", func() {
	var w *httptest.ResponseRecorder
	var r *gin.Engine
	var req *http.Request

	It("CONNECT should respond with status code 200 (OK)", func() {
		method := http.MethodConnect
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("DELETE should respond with status code 200 (OK)", func() {
		method := http.MethodDelete
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("GET should respond with status code 200 (OK)", func() {
		method := http.MethodGet
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("HEAD should respond with status code 200 (OK)", func() {
		method := http.MethodHead
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("PATCH should respond with status code 200 (OK)", func() {
		method := http.MethodPatch
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("POST should respond with status code 200 (OK)", func() {
		method := http.MethodPost
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("PUT should respond with status code 200 (OK)", func() {
		method := http.MethodPut
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("TRACE should respond with status code 200 (OK)", func() {
		method := http.MethodTrace
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("OPTIONS should respond with status code 204 (No Content)", func() {
		method := http.MethodOptions
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusNoContent))
	})

	It("CONNECT should respond with status code 200 (OK)", func() {
		method := http.MethodConnect
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		req.Header.Set("Origin", "localhost")
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})

	It("should not set headers if the Origin header is empty", func() {
		method := http.MethodGet
		w, _, r = helpers.SetupTestEnvironment(method)
		r.Use(middlewares.CORSMiddleWare())
		req, _ = http.NewRequest(method, "/"+path, nil)
		r.Handle(method, path, helpers.TestHandler)
		r.ServeHTTP(w, req)
		Expect(w.Code).To(Equal(http.StatusOK))
	})
})
