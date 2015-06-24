package main

import (
	"log"
	"net/http"
	"os"

	"github.com/tomsteele/veil-evasion-api/myrpc"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
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
	n.Use(negroni.NewStatic(http.Dir(checkEnv("VEIL_OUTPUT_DIR"))))
	n.UseHandler(r)
	n.Run(checkEnv("SERVER_LISTENER"))
}
