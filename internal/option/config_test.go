package option_test

import (
	"fmt"
	"github.com/ShatteredRealms/Accounts/internal/option"
	"github.com/ShatteredRealms/Accounts/pkg/helpers"
	"os"
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config", func() {
	var config option.Config
	BeforeEach(func() {
		config = option.DefaultConfig
	})

	Context("validity", func() {
		It("does not have any conflicting flags", func() {
			v := reflect.ValueOf(config)
			found := make(map[string]bool)
			for i := 0; i < v.NumField(); i++ {
				c := v.Field(i).Interface().(option.ConfigOption)
				Expect(found[c.Flag]).Should(BeFalse())
				found[c.Flag] = true
			}
		})

		It("has all defaults set", func() {
			v := reflect.ValueOf(config)
			for i := 0; i < v.NumField(); i++ {
				Expect(v.Field(i).Interface().(option.ConfigOption).Default).ShouldNot(BeNil())
				Expect(v.Field(i).Interface().(option.ConfigOption).EnvVar).ShouldNot(BeEmpty())
				Expect(v.Field(i).Interface().(option.ConfigOption).Flag).ShouldNot(BeEmpty())
				Expect(v.Field(i).Interface().(option.ConfigOption).Usage).ShouldNot(BeEmpty())
			}
		})
	})

	It("Generates the address correctly", func() {
		host := helpers.RandString(10)
		port := helpers.RandString(4)

		Expect(config.Address()).To(Equal(fmt.Sprintf("%s:%s", config.Host.Default, config.Port.Default)))

		config.Host.Value = &host
		config.Port.Value = &port
		Expect(config.Address()).To(Equal(fmt.Sprintf("%s:%s", host, port)))
		host = helpers.RandString(10)
		Expect(config.Address()).To(Equal(fmt.Sprintf("%s:%s", host, port)))
		port = helpers.RandString(4)
		Expect(config.Address()).To(Equal(fmt.Sprintf("%s:%s", host, port)))
	})

	It("gets default value if value is nil", func() {
		val := helpers.RandString(10)

		Expect(config.Host.GetValue()).To(Equal(config.Host.Default))

		config.Host.Value = &val
		Expect(config.Host.GetValue()).To(Equal(val))
	})

	It("gives the correct output for IsRelease()", func() {
		Expect(config.IsRelease()).To(BeFalse())
		mode := option.ReleaseMode
		config.Mode.Value = &mode
		Expect(config.IsRelease()).To(BeTrue())
	})

	It("reads the environment variables correctly", func() {
		port := helpers.RandString(5)
		host := helpers.RandString(5)
		mode := helpers.RandString(5)
		keydir := helpers.RandString(5)
		dbfile := helpers.RandString(5)
		os.Setenv("PS_AUTH_PORT", port)
		os.Setenv("PS_AUTH_HOST", host)
		os.Setenv("PS_AUTH_MODE", mode)
		os.Setenv("PS_KEY_DIR", keydir)
		os.Setenv("PS_DB_FILE", dbfile)

		config = option.NewConfig()
		Expect(config.Port.GetValue()).To(Equal(port))
		Expect(config.Host.GetValue()).To(Equal(host))
		Expect(config.Mode.GetValue()).To(Equal(mode))
		Expect(config.KeyDir.GetValue()).To(Equal(keydir))
		Expect(config.DBFile.GetValue()).To(Equal(dbfile))
	})
})
