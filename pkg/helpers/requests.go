package helpers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
)

const (
	// APIURL The base URI path
	APIURL = "/api/v1"
)

// URLFromRelative creates a full URL path including APIURL prefix. The relative path can start with a forward slash,
// but it is not required
func URLFromRelative(relative string) string {
	regex := regexp.MustCompile(`^/.+$`)
	if regex.MatchString(relative) {
		return APIURL + relative
	}
	return APIURL + "/" + relative
}

// SetupTestEnvironment creates a HTTPTest recorder, initializes a test gin engine, and initializes the context Request
// data for Method and allocates space for the Header.
func SetupTestEnvironment(method string) (*httptest.ResponseRecorder, *gin.Context, *gin.Engine) {
	w := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, r := gin.CreateTestContext(w)
	c.Request = &http.Request{
		Header: make(http.Header),
		Method: method,
		URL: &url.URL{
			Path:   "/testing/api/v1/",
			Host:   "localhost",
			Scheme: "http",
		},
	}
	return w, c, r
}

// ServeRequest creates a request from the given information, sets the relevant headers, and sends a request to the
// gin engine to serve the request. The response will be nil if everything is successful.
func ServeRequest(method string, path string, body interface{}, r *gin.Engine, w *httptest.ResponseRecorder) error {
	var req *http.Request
	var err error
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err = http.NewRequest(method, path, bytes.NewReader(data))
	} else {
		req, err = http.NewRequest(method, path, nil)
	}
	if err != nil {
		return err
	}

	if method != http.MethodGet {
		req.Header.Set("Content-Type", "application/json")
	}

	r.ServeHTTP(w, req)
	return nil
}

// TestHandler is an simple handler that responds with an ExamplePayload with a Status
// OK (200) every time.
func TestHandler(c *gin.Context) {
	c.JSON(http.StatusOK, ExamplePayload{
		Message: "Test",
	})
}
