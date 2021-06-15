package handle

import "github.com/go-chi/chi/v5"

func SetHandlers(r *chi.Mux) {
	setAuthHandlers(r)
	setUserHandlers(r)
	setURLHandlers(r)
}
