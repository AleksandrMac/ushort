package handle

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func setUserHandlers(r *chi.Mux) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/", userList)
		r.Post("/", createUser)

		r.Route("/{userId}", func(r chi.Router) {
			r.Get("/", user)
			r.Patch("/", updateUser)
			r.Delete("/", deleteUser)
			r.Post("/", createUserToken)
			r.Delete("/{api_key}", deleteUserToken)
		})
	})
}

func userList(w http.ResponseWriter, r *http.Request)        {}
func createUser(w http.ResponseWriter, r *http.Request)      {}
func user(w http.ResponseWriter, r *http.Request)            {}
func deleteUser(w http.ResponseWriter, r *http.Request)      {}
func updateUser(w http.ResponseWriter, r *http.Request)      {}
func createUserToken(w http.ResponseWriter, r *http.Request) {}
func deleteUserToken(w http.ResponseWriter, r *http.Request) {}
