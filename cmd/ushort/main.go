package main

import (
	"log"
	"net/http"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/go-chi/chi/v5"
)

func main() {
	var (
		conf *config.Config
		err  error
	)

	conf, err = config.New("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()
	if err = http.ListenAndServe(conf.Server.Port, router); err != nil {
		log.Fatal(err)
	}
}
