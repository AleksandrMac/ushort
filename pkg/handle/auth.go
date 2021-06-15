package handle

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setAuthHandlers(r *chi.Mux) {
	r.Post("/login", login)
	r.Get("/logout", logout)
}

func login(w http.ResponseWriter, r *http.Request)  {}
func logout(w http.ResponseWriter, r *http.Request) {}
