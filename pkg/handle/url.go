package handle

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (env *Env) setURLHandlers(r *chi.Mux) {
	r.Get("/{urlId}", env.redirectTo)

	r.Route("/url", func(r chi.Router) {
		r.Get("/", urlList)
		r.Post("/", createURL)
		r.Get("/generate", generateURL)

		r.Route("/{urlId}", func(r chi.Router) {
			r.Get("/", url)
			r.Patch("/", updateURL)
			r.Delete("/", deleteURL)
		})
	})
}

func (env *Env) redirectTo(w http.ResponseWriter, r *http.Request) {}
func urlList(w http.ResponseWriter, r *http.Request)               {}
func createURL(w http.ResponseWriter, r *http.Request)             {}
func url(w http.ResponseWriter, r *http.Request)                   {}
func updateURL(w http.ResponseWriter, r *http.Request)             {}
func deleteURL(w http.ResponseWriter, r *http.Request)             {}
func generateURL(w http.ResponseWriter, r *http.Request)           {}
