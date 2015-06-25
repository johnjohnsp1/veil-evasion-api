package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"os"

	"github.com/tomsteele/veil-evasion-api/myrpc"

	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"

	"testing"
)

func connectToVeil(t *testing.T) *rpc.Client {
	c, err := myrpc.Dial("tcp", os.Getenv("VEIL_LISTENER"))
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestVersion(t *testing.T) {
	Convey("Given a request to get version from Veil", t, func() {
		c := connectToVeil(t)
		h := New(&H{C: c})
		r := mux.NewRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/version", nil)
		r.HandleFunc("/api/version", h.Version)
		r.ServeHTTP(w, req)
		Convey("The return code should be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The response should contain a version number", func() {
			response := map[string]string{}
			data, _ := ioutil.ReadAll(w.Body)
			err := json.Unmarshal(data, &response)
			So(err, ShouldBeNil)
			So(response["version"], ShouldNotBeEmpty)
		})
	})
}

func TestPayloads(t *testing.T) {
	Convey("Given a request to payloads from Veil", t, func() {
		c := connectToVeil(t)
		h := New(&H{C: c})
		r := mux.NewRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/payloads", nil)
		r.HandleFunc("/api/payloads", h.Payloads)
		r.ServeHTTP(w, req)
		Convey("The return code should be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The response should contain a an array of strings", func() {
			response := []string{}
			data, _ := ioutil.ReadAll(w.Body)
			err := json.Unmarshal(data, &response)
			So(err, ShouldBeNil)
			So(len(response), ShouldNotEqual, 0)
		})
	})
}

func TestPayloadOptions(t *testing.T) {
	Convey("Given a request to get options for a payload", t, func() {
		c := connectToVeil(t)
		h := New(&H{C: c})
		r := mux.NewRouter()
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/options?payload=c/meterpreter/rev_http", nil)
		r.HandleFunc("/api/options", h.PayloadOptions)
		r.ServeHTTP(w, req)
		Convey("The return code should be 200", func() {
			So(w.Code, ShouldEqual, 200)
		})
		Convey("The response should an array of key value pairs", func() {
			response := []PayloadOption{}
			data, _ := ioutil.ReadAll(w.Body)
			err := json.Unmarshal(data, &response)
			So(err, ShouldBeNil)
			So(len(response), ShouldNotEqual, 0)
			So(response[0].Key, ShouldNotBeEmpty)
			So(response[0].Value, ShouldNotBeEmpty)
		})
	})
}
