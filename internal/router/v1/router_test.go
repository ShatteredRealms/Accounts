package v1_test

import (
	"github.com/ShatteredRealms/Accounts/internal/log"
	"github.com/ShatteredRealms/Accounts/internal/option"
	v1 "github.com/ShatteredRealms/Accounts/internal/router/v1"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var logger log.LoggerService = log.NewLogger(log.Info, "")

var _ = Describe("Router", func() {
	var keyDir string
	var releaseMode string
	var config option.Config
	

	BeforeEach(func() {
		keyDir = "../../../test/auth"
		releaseMode = gin.TestMode
		config = option.DefaultConfig
		config.KeyDir.Value = &keyDir
		config.Mode.Value = &releaseMode
	})

	Context("release mode", func() {
		It("should panic if db is invalid", func() {
			Expect(startRoutingNil_Release).To(Panic())
		})

		It("should work with valid input", func() {
			db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())
			router, err := v1.InitRouter(db, config, logger)

			Expect(router).NotTo(BeNil())
			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})

func startRoutingNil_Release() {
	v1.InitRouter(nil, option.Config{}, logger)
}
