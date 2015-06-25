package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/tomsteele/veil-evasion-api/myrpc"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/nabeken/negroni-auth"
	"github.com/tomsteele/veil-evasion-api/handlers"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkEnv(v string) string {
	l := os.Getenv(v)
	if l == "" {
		log.Fatalf("Env variable %s missing\n", v)
	}
	return l
}

func main() {
	c, err := myrpc.Dial("tcp", checkEnv("VEIL_LISTENER"))
	check(err)
	h := handlers.New(&handlers.H{C: c})
	r := mux.NewRouter()
	r.HandleFunc("/api/version", h.Version).Methods("GET")
	r.HandleFunc("/api/payloads", h.Payloads).Methods("GET")
	r.HandleFunc("/api/options", h.PayloadOptions).Methods("GET")
	r.HandleFunc("/api/generate", h.Generate).Methods("POST")
	n := negroni.New()
	n.Use(negroni.NewLogger())
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewStatic(http.Dir(checkEnv("VEIL_OUTPUT_DIR"))))
	n.Use(auth.Basic(checkEnv("ADMIN_USER"), checkEnv("ADMIN_PASS")))
	n.Use(negroni.NewStatic(http.Dir("public")))
	n.Use(negroni.HandlerFunc(func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		if req.Method != "POST" {
			next(w, req)
			return
		}
		if !strings.Contains(req.Header.Get("content-type"), "application/json") {
			h.JSON400(w, errors.New("Content-Type must be application/json"))
			return
		}
		next(w, req)
	}))
	n.UseHandler(r)
	n.Run(checkEnv("SERVER_LISTENER"))
}
