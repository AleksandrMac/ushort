package main

import (
	"log"
	"net/http"

	"github.com/AleksandrMac/ushort/pkg/config"
	h "github.com/AleksandrMac/ushort/pkg/handle"

	"github.com/go-chi/chi/v5"
)

var (
	c   *config.Config
	err error
)

func main() {
	c, err = config.New("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	h.SetHandlers(r)

	if err = http.ListenAndServe(c.Server.Port, r); err != nil {
		log.Fatal(err)
	}
}
