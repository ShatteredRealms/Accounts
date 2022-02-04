package health_test

import (
	"encoding/json"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/productivestudy/auth/cmd/auth/controller/v1/health"
	"github.com/productivestudy/auth/tests/helpers"
)

var _ = Describe("Health", func() {

	w, c, _ := helpers.SetupTestEnvironment("GET")
	health.Health(c)

	resp := health.Response{}
	Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
	Expect(w.Code).To(Equal(http.StatusOK))
	Expect(resp.Health).To(Equal("ok"))
})
