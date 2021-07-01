package handle

import (
	"github.com/go-chi/chi/v5"

	"github.com/AleksandrMac/ushort/pkg/config"
)

type Handler struct {
	Env *config.Env
}

func (h *Handler) SetHandlers(r *chi.Mux) {
	h.setAuthHandlers(r)
	h.setURLHandlers(r)
}
