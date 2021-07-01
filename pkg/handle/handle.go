package handle

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/models"
)

type Env struct {
	DB        *models.DB
	Config    *config.Config
	TokenAuth *jwtauth.JWTAuth
}

func (env *Env) SetHandlers(r *chi.Mux) {
	env.setAuthHandlers(r)
	env.setURLHandlers(r)
}
