package hoverfly_test

import (
	"encoding/json"
	"io/ioutil"

	"github.com/SpectoLabs/hoverfly/functional-tests"
	"github.com/dghubble/sling"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hoverfly cache", func() {

	var (
		hoverfly *functional_tests.Hoverfly
	)

	BeforeEach(func() {
		hoverfly = functional_tests.NewHoverfly()
		hoverfly.Start()
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		hoverfly.Proxy(sling.New().Get("http://destination-server.com"))
	})

	AfterEach(func() {
		hoverfly.Stop()
	})

	It("should be flushed when changing to capture mode", func() {
		hoverfly.SetMode("capture")

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should be flushed when importing via the API", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)
		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should preload the cache with exact match request matcher when put into simulate mode", func() {
		hoverfly.ImportSimulation(functional_tests.JsonPayload)

		hoverfly.SetMode("simulate")

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(1))
	})

	It("should not preload the cache if simulation does not contain exact match request matcher", func() {
		hoverfly.ImportSimulation(functional_tests.JsonMatchSimulation)

		hoverfly.SetMode("simulate")

		cacheView := hoverfly.GetCache()

		Expect(cacheView.Cache).To(HaveLen(0))
	})

	It("should not stop matching on headers by caching the same request twice with different headers", func() {
		hoverfly.ImportSimulation(functional_tests.ExactMatchPayload)

		hoverfly.SetMode("simulate")

		req := sling.New().Get("http://localhost:" + hoverfly.GetAdminPort() + "/api/v2/cache")

		res := functional_tests.DoRequest(req)
		Expect(res.StatusCode).To(Equal(200))

		responseBytes, err := ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		var responseJson map[string]interface{}
		json.Unmarshal(responseBytes, &responseJson)

		Expect(responseJson["cache"]).To(HaveLen(1))

		req = sling.New().Get("http://test-server.com/path1").Add("Header", "value1")

		res = hoverfly.Proxy(req)
		Expect(res.StatusCode).To(Equal(200))

		responseBytes, err = ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		Expect(responseBytes).To(Equal([]byte("exact match 1")))

		req = sling.New().Get("http://test-server.com/path1").Add("Header", "value2")

		res = hoverfly.Proxy(req)
		Expect(res.StatusCode).To(Equal(200))

		responseBytes, err = ioutil.ReadAll(res.Body)
		Expect(err).To(BeNil())

		Expect(responseBytes).To(Equal([]byte("exact match 2")))
	})
})
