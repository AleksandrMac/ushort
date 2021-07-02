package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/rs/zerolog"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/connect"
)

type Controller struct {
	DB        *connect.DB
	Config    *config.Config
	TokenAuth *jwtauth.JWTAuth
	Logger    *zerolog.Logger
	Info      chan string
	Debug     chan error
	Err       chan error
	Critical  chan error
}

func (c *Controller) SetControllers(r *chi.Mux) {
	c.setAuthControllers(r)
	c.setURLControllers(r)
}
