package main_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "cim-identity"
)

var _ = Describe("Microservice backend communications", func() {

	var (
		// server *Server
		opt Options
	)

	Describe("when authenticating", func() {

		Context("with valid auth URL, client ID and secret", func() {
			BeforeEach(func() {
				opt = Options{
					"seed": "ginkgo",
					// "log":          ioutil.Discard,
					"serviceURL":   "https://dataidentity.test-hafslundnett.io/alias/bulk",
					"authURL":      "https://hid.test-hafslundnett.io/connect/token",
					"tokenURL":     "https://hid.test-hafslundnett.io/connect/token",
					"clientID":     "sesam",
					"clientSecret": "eDvS3YaJwasn618ab924FE6mr7We9k",
				}
				if len(os.Getenv("LOGOUTPUT")) == 0 {
					opt["log"] = ioutil.Discard
				}

			})
			It("authenticates successfully", func() {
				_, err := NewServer(NewOptions(&opt))
				Expect(err).To(BeNil())
			})
		})
	})

})
