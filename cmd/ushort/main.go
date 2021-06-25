package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/config/env"
	h "github.com/AleksandrMac/ushort/pkg/handle"

	"github.com/go-chi/chi/v5"
)

var (
	c   *config.Config
	err error
)

func main() {
	c, err = env.New()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c)

	r := chi.NewRouter()
	h.SetHandlers(r)

	if err = http.ListenAndServe(":"+c.Server.Port, r); err != nil {
		log.Fatal(err)
	}
}
