package main

import (
	"log"
	"net/http"

	"github.com/AleksandrMac/ushort/pkg/config/env"
	"github.com/AleksandrMac/ushort/pkg/handle"
	"github.com/AleksandrMac/ushort/pkg/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
)

// var (
// 	cfg *config.Config
// 	e   *config.Env
// )

func main() {
	cfg, err := env.New()
	if err != nil {
		log.Panic(err)
	}

	db, err := models.NewDB(cfg.DB.URL)
	if err != nil {
		log.Panic(err)
	}
	h := &handle.Env{
		DB:        db,
		Config:    cfg,
		TokenAuth: jwtauth.New("HS256", cfg.Auth.PrivateKey, nil),
	}

	logger := httplog.NewLogger("httplog-example", httplog.Options{
		JSON: true,
	})

	// Service
	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	h.SetHandlers(r)

	if err = http.ListenAndServe(":"+cfg.Server.Port, r); err != nil {
		log.Fatal(err)
	}
}
