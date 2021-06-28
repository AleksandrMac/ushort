package handle

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func SetHandlers(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello")); err != nil {
			log.Fatal(err)
		}
	})
	setAuthHandlers(r)
	setUserHandlers(r)
	setURLHandlers(r)
}
