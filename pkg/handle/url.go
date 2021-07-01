package handle

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AleksandrMac/ushort/pkg/models/url"
	"github.com/AleksandrMac/ushort/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

// tmpUrl хранит ссылки котырые были сгенерированы, но еще не попали в БД
var tmpURL map[string]time.Time

func (env *Env) setURLHandlers(r *chi.Mux) {
	tmpURL = make(map[string]time.Time)
	r.Get("/{urlId}", env.redirectTo)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(env.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Route("/url", func(r chi.Router) {
			r.Get("/", urlList)
			r.Post("/", createURL)
			r.Get("/generate", env.generateURL)

			r.Route("/{urlId}", func(r chi.Router) {
				r.Get("/", getURL)
				r.Patch("/", updateURL)
				r.Delete("/", deleteURL)
			})
		})
	})
}

func (env *Env) redirectTo(w http.ResponseWriter, r *http.Request) {}
func urlList(w http.ResponseWriter, r *http.Request)               {}
func createURL(w http.ResponseWriter, r *http.Request)             {}
func getURL(w http.ResponseWriter, r *http.Request)                {}
func updateURL(w http.ResponseWriter, r *http.Request)             {}
func deleteURL(w http.ResponseWriter, r *http.Request)             {}
func (env *Env) generateURL(w http.ResponseWriter, r *http.Request) {
	for {
		newURL := utils.RandString(env.Config.LengthURL)
		urlFromDB, err := url.SelectWithID(newURL, env.DB)
		if err != nil {
			log.Default().Println(err)
			if err.Error() != "sql: no rows in result set" {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if urlFromDB != nil {
			log.Default().Printf("Генерация URL: %s уже существует в БД\n", urlFromDB)
			continue
		}
		if t, ok := tmpURL[newURL]; ok {
			if time.Now().After(t) {
				log.Default().Printf("Генерация URL: %s зарезервирован\n", newURL)
				continue
			}
		}
		tmpURL[newURL] = time.Now().Add(time.Duration(env.Config.TmpURLLifeTime))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("{\"shortURL\": \"/%s\"}", newURL))); err != nil {
			log.Default().Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
