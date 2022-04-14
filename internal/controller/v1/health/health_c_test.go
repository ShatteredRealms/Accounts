package health_test

import (
	"encoding/json"
	"github.com/ShatteredRealms/Accounts/internal/controller/v1/health"
	"github.com/ShatteredRealms/Accounts/pkg/helpers"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Health", func() {

	w, c, _ := helpers.SetupTestEnvironment("GET")
	health.Health(c)

	resp := health.Response{}
	Expect(json.NewDecoder(w.Body).Decode(&resp)).ShouldNot(HaveOccurred())
	Expect(w.Code).To(Equal(http.StatusOK))
	Expect(resp.Health).To(Equal("ok"))
})
