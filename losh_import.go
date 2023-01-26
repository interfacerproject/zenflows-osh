package main

import (
	"log"
	"net/http"
)

func main() {
	conf, err := loadConfig()
	if err != nil {
		log.Fatalf("bad error config: %s\n", err.Error())
	}
	log.Printf("Starting service on %s\n", conf.addr)

	mux := http.NewServeMux()
	mux.HandleFunc("/clone", cloneHandler)

	log.Fatal(http.ListenAndServe(conf.addr, mux))
}
