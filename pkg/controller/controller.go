package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/config/env"
	"github.com/AleksandrMac/ushort/pkg/model"
)

type Controller struct {
	DB        model.CRUD
	Config    *config.Config
	TokenAuth *jwtauth.JWTAuth
	Logger    *zerolog.Logger
	Info      chan string
	Debug     chan string
	Err       chan error
	Warn      chan error
	Critical  chan error
}

func New() (*Controller, error) {
	cfg, err := env.New()
	if err != nil {
		return nil, err
	}

	var logger zerolog.Logger
	switch cfg.LogLevel {
	case "trace":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "trace",
			JSON:     true})
	case "debug":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "debug",
			JSON:     true})
	case "info":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "info",
			JSON:     true})
	case "warn":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "warn",
			JSON:     true})
	case "error":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "error",
			JSON:     true})
	case "critical":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "critical",
			JSON:     true})
	}

	db, err := model.NewDB(cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	return &Controller{
		DB:        db,
		Config:    cfg,
		TokenAuth: jwtauth.New("HS256", cfg.Auth.PrivateKey, nil),
		Logger:    &logger,
		Info:      make(chan string),
		Debug:     make(chan string),
		Err:       make(chan error),
		Warn:      make(chan error),
		Critical:  make(chan error),
	}, nil
}

func (c *Controller) SetControllers(r *chi.Mux) {
	c.setAuthControllers(r)
	c.setURLControllers(r)
}
