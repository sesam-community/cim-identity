package main_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "sesam-shaid"
)

var _ = Describe("Microservice handlers", func() {

	var (
		opt           Options = Options{}
		seed          string  = "ginkgo"
		server        *Server
		request       *http.Request
		response      *httptest.ResponseRecorder
		url           string
		input         string
		output        string
		contentHeader string = "Content-Type"
		contentType   string = "application/json"
		contentFull   string = "application/json; charset=utf-8"
	)

	BeforeEach(func() {
		if len(os.Getenv("LOGOUTPUT")) == 0 {
			opt["log"] = ioutil.Discard
		}
		opt["level"] = "ALL"
		opt["seed"] = seed
		server, _ = NewServer(NewOptions(&opt))
		response = httptest.NewRecorder()
	})

	AfterEach(func() {
	})

	Describe("when POST to incorrectly", func() {

		Context("without body", func() {
			It("returns HTTP error 400", func() {
				request, _ = http.NewRequest("POST", "/", nil)
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with empty body", func() {
			It("returns HTTP error 400", func() {
				input = ``
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with garbage body", func() {
			It("returns HTTP error 400", func() {
				input = `garbage\\\s-d.,f-.,32423#¤%#¤%R:WEfecøw<slzlkcd<mesopip09`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with non-array JSON", func() {
			It("returns HTTP error 400", func() {
				input = `{"key":"val"}`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with non-closed array JSON", func() {
			It("returns HTTP error 400", func() {
				input = `[{"key":"val"}`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with erroneous JSON array closing", func() {
			It("returns HTTP error 400", func() {
				input = `[{"key":"val"}{`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with superflous trailing comma in JSON array", func() {
			It("returns HTTP error 400", func() {
				input = `[{"key":"val"},]`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

		Context("with unsupported JSON entity in array", func() {
			It("returns HTTP error 400", func() {
				input = `[{"key":"val"},"string"]`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
				Expect(response.Code).To(Equal(400))
			})
		})

	})

	Describe("when POST to /", func() {

		Context("with JSON empty array", func() {
			BeforeEach(func() {
				input = `[]`
				output = `[]`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of JSON empty array")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with entity without '_id' field", func() {
			BeforeEach(func() {
				input = `[{"key":"val","fields":2}]`
				output = input
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of unchanged input")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with entity containing '_id' field", func() {
			BeforeEach(func() {
				input = `[{"_id":"convert-to-sha1-UUID", "key":"val","fields":2}]`
				output = `[{"_id":"a60989a3-0af4-5d95-b632-72a604a96474", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed '_id' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with entity containing additional system fields using '_' prefixing", func() {
			BeforeEach(func() {
				input = `[{"_id":"convert-to-sha1-UUID", "_previous":null, "_deleted":false, "key":"val","fields":2}]`
				output = `[{"_id":"a60989a3-0af4-5d95-b632-72a604a96474", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", "/", strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response without extra system fields")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

	})

	Describe("POST to /<field>", func() {

		Context("with URL using 'shaid' as field", func() {
			BeforeEach(func() {
				url = `/shaid`
				input = `[{"shaid":"convert-to-sha1-UUID", "_previous":null, "_deleted":false, "key":"val","fields":2}]`
				output = `[{"shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field shortcut and entity is fully namespaced", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{"entity-namespace:shaid":"convert-to-sha1-UUID", "_previous":null, "_deleted":false, "entity-namespace:key":"val","entity-namespace:fields":2}]`
				output = `[{"entity-namespace:shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "entity-namespace:key":"val","entity-namespace:fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using 'shaid' as field shortcut and entity is fully namespaced", func() {
			BeforeEach(func() {
				url = `/shaid`
				input = `[{"entity-namespace:shaid":"convert-to-sha1-UUID", "_previous":null, "_deleted":false, "entity-namespace:key":"val","entity-namespace:fields":2}]`
				output = `[{"entity-namespace:shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "entity-namespace:key":"val","entity-namespace:fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using '.shaid' as suffixed field shortcut", func() {
			BeforeEach(func() {
				url = `/.shaid`
				input = `[{"entity-namespace:prefix.shaid":"convert-to-sha1-UUID", "_previous":null, "_deleted":false, "entity-namespace:key":"val","entity-namespace:fields":2}]`
				output = `[{"entity-namespace:prefix.shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "entity-namespace:key":"val","entity-namespace:fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field and has 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{"shaid":"convert-to-sha1-UUID", "rdf:type":"~:namespace:value", "key":"val","fields":2}]`
				output = `[{"shaid":"81ef0d83-320b-540f-9e42-5cb9a3676bdc", "rdf:type":"~:namespace:value", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':.shaid' as suffixed field shortcut and has 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:.shaid`
				input = `[{"entity-namespace:prefix.shaid":"convert-to-sha1-UUID", "rdf:type":"~:namespace:value", "_previous":null, "_deleted":false, "entity-namespace:key":"val","entity-namespace:fields":2}]`
				output = `[{"entity-namespace:prefix.shaid":"81ef0d83-320b-540f-9e42-5cb9a3676bdc", "rdf:type":"~:namespace:value", "entity-namespace:key":"val","entity-namespace:fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field, but without 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{"shaid":"convert-to-sha1-UUID", "key":"val","fields":2}]`
				output = `[{"shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field, but empty 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{"shaid":"convert-to-sha1-UUID", "rdf:type":"", "key":"val","fields":2}]`
				output = `[{"shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "rdf:type":"", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using 'shaid' as field shortcut, namespaced entity and has 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/shaid`
				input = `[{"entity-namespace:shaid":"convert-to-sha1-UUID", "rdf:type":"~:namespace:value", "_previous":null, "_deleted":false, "entity-namespace:key":"val","entity-namespace:fields":2}]`
				output = `[{"entity-namespace:shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "rdf:type":"~:namespace:value", "entity-namespace:key":"val","entity-namespace:fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field and multiple entities", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : "81ef0d83-320b-540f-9e42-5cb9a3676bdc",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"entity-namespace:shaid"  : "b2a2ff67-027e-5790-b046-10d9f044fd28",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

	})

	Describe("POST to /<fieldA>;<fieldB>;...", func() {

		Context("with URL using multiple fields ':shaid' and 'cimid'", func() {
			BeforeEach(func() {
				url = `/:shaid;cimid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "convert-to-sha1-UUID",
					"entity-namespace:cimid"  : "cim:Type:convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:cimid"  : "cim:Type:also-convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : "81ef0d83-320b-540f-9e42-5cb9a3676bdc",
					"entity-namespace:cimid"  : "cfeabbbd-7bc2-578d-a762-026bff4fb5cf",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"entity-namespace:shaid"  : "b2a2ff67-027e-5790-b046-10d9f044fd28",
					"entity-namespace:cimid"  : "7f5d34b0-ea5d-51c3-934a-8a1fa09ca627",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using multiple fields ':shaid' and ':oldid'", func() {
			BeforeEach(func() {
				url = `/:shaid;:oldid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "convert-to-sha1-UUID",
					"entity-namespace:oldid"  : "convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:oldid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : "81ef0d83-320b-540f-9e42-5cb9a3676bdc",
					"entity-namespace:oldid"  : "81ef0d83-320b-540f-9e42-5cb9a3676bdc",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"entity-namespace:shaid"  : "b2a2ff67-027e-5790-b046-10d9f044fd28",
					"entity-namespace:oldid"  : "b2a2ff67-027e-5790-b046-10d9f044fd28",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using multiple fields ':shaid' and ':oldid' when 'rdf:type' is missing", func() {
			BeforeEach(func() {
				url = `/:shaid;:oldid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "convert-to-sha1-UUID",
					"entity-namespace:oldid"  : "convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2
					},{
					"_id"                     : "entity-namespace:2",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:oldid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : "a60989a3-0af4-5d95-b632-72a604a96474",
					"entity-namespace:oldid"  : "a60989a3-0af4-5d95-b632-72a604a96474",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2
					},{
					"_id"                     : "entity-namespace:2",
					"entity-namespace:shaid"  : "0e374c4b-be1d-5eb3-8385-5f177fd9a432",
					"entity-namespace:oldid"  : "0e374c4b-be1d-5eb3-8385-5f177fd9a432",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using multiple fields 'shaid' and 'oldid'", func() {
			BeforeEach(func() {
				url = `/shaid;oldid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "convert-to-sha1-UUID",
					"entity-namespace:oldid"  : "convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:oldid"  : "also-convert-to-sha1-UUID",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : "a60989a3-0af4-5d95-b632-72a604a96474",
					"entity-namespace:oldid"  : "a60989a3-0af4-5d95-b632-72a604a96474",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					},{
					"_id"                     : "entity-namespace:2",
					"entity-namespace:shaid"  : "0e374c4b-be1d-5eb3-8385-5f177fd9a432",
					"entity-namespace:oldid"  : "0e374c4b-be1d-5eb3-8385-5f177fd9a432",
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:othervalue"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using multiple fields ':shaid' and '+oldid' for prefixes", func() {
			BeforeEach(func() {
				url = `/:shaid;+oldid`
				input = `[{
					"_id"                     : "entity:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity:shaid"            : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity:group.oldid"      : [ "ns:class:convert-to-sha1-UUID" , "ns:type:also-convert-to-sha1-UUID" , "ns:flavour:also-convert-to-sha1-UUID" ],
					"entity:key"              : "val",
					"entity:fields"           : 2,
					"rdf:type"                : "~:namespace:value"
					}]`
				output = `[{
					"_id"                     : "entity:1",
					"entity:shaid"            : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity:group.oldid"      : [ "~:class:ab4ffc42-b2ba-535d-911a-5ab2da2bc59e" , "~:type:cf526379-30af-563a-8789-710799ce33a7" , "~:flavour:0437a832-73a7-537c-8086-081cec7fb7f3" ],
					"entity:key"              : "val",
					"entity:fields"           : 2,
					"rdf:type"                : "~:namespace:value"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

	})

	Describe("POST to /<field> using composed values", func() {

		Context("with URL field ':shaid' as array", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' when 'rdf:type' is an array", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" ]
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" ]
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field, but empty array 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{"shaid":"convert-to-sha1-UUID", "rdf:type":[], "key":"val","fields":2}]`
				output = `[{"shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "rdf:type":[], "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using ':shaid' as field, but array 'rdf:type' has empty string", func() {
			BeforeEach(func() {
				url = `/:shaid`
				input = `[{"shaid":"convert-to-sha1-UUID", "rdf:type":[""], "key":"val","fields":2}]`
				output = `[{"shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "rdf:type":[""], "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

	})

	Describe("POST to /<field>/<namespace>", func() {

		Context("with URL field ':shaid' and namespace prefix 'namespace:' when 'rdf:type' is an array", func() {
			BeforeEach(func() {
				url = `/:shaid/namespace:`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' and namespace prefix 'namespace:' when 'rdf:type' is missing", func() {
			BeforeEach(func() {
				url = `/:shaid/namespace:`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "a60989a3-0af4-5d95-b632-72a604a96474" , "0e374c4b-be1d-5eb3-8385-5f177fd9a432" , "0e374c4b-be1d-5eb3-8385-5f177fd9a432" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' and namespace prefix 'namespace:' when 'rdf:type' has multiple matches", func() {
			BeforeEach(func() {
				url = `/:shaid/namespace:`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" ]
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" ]
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' and namespace prefix 'namespace:' when 'rdf:type' has no match", func() {
			BeforeEach(func() {
				url = `/:shaid/names-r-us:`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" ]
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "a60989a3-0af4-5d95-b632-72a604a96474" , "0e374c4b-be1d-5eb3-8385-5f177fd9a432" , "0e374c4b-be1d-5eb3-8385-5f177fd9a432" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" ]
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' and namespace prefix 'namespace:' doesn't match the 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid/namespace:`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:names-r-us:value"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "a60989a3-0af4-5d95-b632-72a604a96474" , "0e374c4b-be1d-5eb3-8385-5f177fd9a432" , "0e374c4b-be1d-5eb3-8385-5f177fd9a432" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:names-r-us:value"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' and given namespace 'namespace:value' with unused 'rdf:type' array", func() {
			BeforeEach(func() {
				url = `/:shaid/namespace:value`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:names-r-us:value" , "~:namespace:othervalue" ]
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : [ "~:names-r-us:value" , "~:namespace:othervalue" ]
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL field ':shaid' and given already default indirect namespace 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid/rdf:type`
				input = `[{
					"_id"                     : "entity-namespace:1",
					"_previous"               : null,
					"_deleted"                : false,
					"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					}]`
				output = `[{
					"_id"                     : "entity-namespace:1",
					"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
					"entity-namespace:key"    : "val",
					"entity-namespace:fields" : 2,
					"rdf:type"                : "~:namespace:value"
					}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		/*
			// Belayed for now...
			// FIXME: this unit-test is a bit complex - derive into other separate ones as well to precede this
			Context("with URL fields ':shaid', '::.oldid' and other inlined namespaced shortcut fields given explicit default namespace when 'rdf:type' is array", func() {
				BeforeEach(func() {
					url = `/:shaid;::.oldid;.someid:~:namespace:othervalue;.otherid::~:namespace:othervalue;.localid:_:namespace:othervalue;:.shortid:_;.hardid:_;.moreid:refname:~:namespace:othervalue;.lastid:~:~:namespace:othervalue/namespace:value`
					input = `[{
						"_id"                     : "entity:1",
						"_previous"               : null,
						"_deleted"                : false,
						"entity:shaid"            : "convert-to-sha1-UUID",
						"entity:Type.oldid"       : "convert-to-sha1-UUID",
						"entity:Type.someid"      : "also-convert-to-sha1-UUID",
						"entity:Type.otherid"     : "also-convert-to-sha1-UUID",
						"entity:Type.localid"     : "also-convert-to-sha1-UUID",
						"entity:Type.shortid"     : "also-convert-to-sha1-UUID",
						"entity:Type.hardid"      : "also-convert-to-sha1-UUID",
						"entity:Type.moreid"      : "also-convert-to-sha1-UUID",
						"entity:Type.lastid"      : "also-convert-to-sha1-UUID",
						"entity:key"              : "val",
						"entity:fields"           : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" , "~:names-r-us:value" ]
						}]`
					output = `[{
						"_id"                     : "entity:1"     ,
						"entity:shaid"            : "81ef0d83-320b-540f-9e42-5cb9a3676bdc",
						"entity:Type.oldid"       : "urn:uuid:81ef0d83-320b-540f-9e42-5cb9a3676bdc",
						"entity:Type.someid"      : "~:0e374c4b-be1d-5eb3-8385-5f177fd9a432",
						"entity:Type.otherid"     : "~:othervalue:0e374c4b-be1d-5eb3-8385-5f177fd9a432",
						"entity:Type.localid"     : "#_0e374c4b-be1d-5eb3-8385-5f177fd9a432",
						"entity:Type.shortid"     : "_052261c2-da4e-5d62-84e9-8f404c2babb0",
						"entity:Type.shortid"     : "_0e374c4b-be1d-5eb3-8385-5f177fd9a432",
						"entity:Type.moreid"      : "~:refname:0e374c4b-be1d-5eb3-8385-5f177fd9a432",
						"entity:Type.lastid"      : "~:namespace:0e374c4b-be1d-5eb3-8385-5f177fd9a432",
						"entity:key"              : "val",
						"entity:fields"           : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:namespace:othervalue" , "~:names-r-us:value" ]
						}]`
					request, _ = http.NewRequest("POST", url, strings.NewReader(input))
					request.Header.Add(contentHeader, contentType)
					server.ServeHTTP(response, request)
				})
				It("replies with", func() {
					By("HTTP status OK 200")
					Expect(response.Code).To(Equal(200))
					By("JSON Content-Type HTTP header")
					Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
					By("response of transformed 'shaid' field")
					Expect(response.Body.String()).To(MatchJSON(output))
				})
			})

			Context("with URL field ':shaid' and namespace shortcut '~' when 'rdf:type' is an array", func() {
				BeforeEach(func() {
					url = `/:shaid/~`
					input = `[{
						"_id"                     : "entity-namespace:1",
						"_previous"               : null,
						"_deleted"                : false,
						"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
						"entity-namespace:key"    : "val",
						"entity-namespace:fields" : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
						}]`
					output = `[{
						"_id"                     : "entity-namespace:1",
						"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
						"entity-namespace:key"    : "val",
						"entity-namespace:fields" : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
						}]`
					request, _ = http.NewRequest("POST", url, strings.NewReader(input))
					request.Header.Add(contentHeader, contentType)
					server.ServeHTTP(response, request)
				})
				It("replies with", func() {
					By("HTTP status OK 200")
					Expect(response.Code).To(Equal(200))
					By("JSON Content-Type HTTP header")
					Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
					By("response of transformed 'shaid' field")
					Expect(response.Body.String()).To(MatchJSON(output))
				})
			})

			Context("with URL field ':shaid' and namespace shortcut '_' when 'rdf:type' is an array", func() {
				BeforeEach(func() {
					url = `/:shaid/_`
					input = `[{
						"_id"                     : "entity-namespace:1",
						"_previous"               : null,
						"_deleted"                : false,
						"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
						"entity-namespace:key"    : "val",
						"entity-namespace:fields" : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
						}]`
					output = `[{
						"_id"                     : "entity-namespace:1",
						"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
						"entity-namespace:key"    : "val",
						"entity-namespace:fields" : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
						}]`
					request, _ = http.NewRequest("POST", url, strings.NewReader(input))
					request.Header.Add(contentHeader, contentType)
					server.ServeHTTP(response, request)
				})
				It("replies with", func() {
					By("HTTP status OK 200")
					Expect(response.Code).To(Equal(200))
					By("JSON Content-Type HTTP header")
					Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
					By("response of transformed 'shaid' field")
					Expect(response.Body.String()).To(MatchJSON(output))
				})
			})

			Context("with URL field ':shaid' and namespace shortcut ':' when 'rdf:type' is an array", func() {
				BeforeEach(func() {
					url = `/:shaid/:`
					input = `[{
						"_id"                     : "entity-namespace:1",
						"_previous"               : null,
						"_deleted"                : false,
						"entity-namespace:shaid"  : [ "convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" , "also-convert-to-sha1-UUID" ],
						"entity-namespace:key"    : "val",
						"entity-namespace:fields" : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
						}]`
					output = `[{
						"_id"                     : "entity-namespace:1",
						"entity-namespace:shaid"  : [ "81ef0d83-320b-540f-9e42-5cb9a3676bdc" , "052261c2-da4e-5d62-84e9-8f404c2babb0" , "052261c2-da4e-5d62-84e9-8f404c2babb0" ],
						"entity-namespace:key"    : "val",
						"entity-namespace:fields" : 2,
						"rdf:type"                : [ "~:namespace:value" , "~:names-r-us:othervalue" ]
						}]`
					request, _ = http.NewRequest("POST", url, strings.NewReader(input))
					request.Header.Add(contentHeader, contentType)
					server.ServeHTTP(response, request)
				})
				It("replies with", func() {
					By("HTTP status OK 200")
					Expect(response.Code).To(Equal(200))
					By("JSON Content-Type HTTP header")
					Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
					By("response of transformed 'shaid' field")
					Expect(response.Body.String()).To(MatchJSON(output))
				})
			})


		*/

	})

	Describe("POST to /<field>/", func() {

		Context("with URL using ':shaid' and empty namespace using trailing slash '/' and has 'rdf:type'", func() {
			BeforeEach(func() {
				url = `/:shaid/`
				input = `[{"shaid":"convert-to-sha1-UUID", "rdf:type":"~:namespace:value", "key":"val","fields":2}]`
				output = `[{"shaid":"81ef0d83-320b-540f-9e42-5cb9a3676bdc", "rdf:type":"~:namespace:value", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

		Context("with URL using 'shaid' and empty namespace using trailing slash '/'", func() {
			BeforeEach(func() {
				url = `/shaid/`
				input = `[{"shaid":"convert-to-sha1-UUID", "rdf:type":"~:namespace:value", "key":"val","fields":2}]`
				output = `[{"shaid":"a60989a3-0af4-5d95-b632-72a604a96474", "rdf:type":"~:namespace:value", "key":"val","fields":2}]`
				request, _ = http.NewRequest("POST", url, strings.NewReader(input))
				request.Header.Add(contentHeader, contentType)
				server.ServeHTTP(response, request)
			})
			It("replies with", func() {
				By("HTTP status OK 200")
				Expect(response.Code).To(Equal(200))
				By("JSON Content-Type HTTP header")
				Expect(response.Header().Get(contentHeader)).To(Equal(contentFull))
				By("response of transformed 'shaid' field")
				Expect(response.Body.String()).To(MatchJSON(output))
			})
		})

	})

})
