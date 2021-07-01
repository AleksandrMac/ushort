package handle

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AleksandrMac/ushort/pkg/models"
	"github.com/AleksandrMac/ushort/pkg/models/url"
	"github.com/AleksandrMac/ushort/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

// tmpUrl хранит ссылки котырые были сгенерированы, но еще не попали в БД
var tmpURL map[string]time.Time

func (h *Handler) setURLHandlers(r *chi.Mux) {
	tmpURL = make(map[string]time.Time)
	r.Get("/{urlID}", h.redirectTo)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(h.Env.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Route("/url", func(r chi.Router) {
			r.Get("/", urlList)
			r.Post("/", h.createURL)
			r.Get("/generate", h.generateURL)

			r.Route("/{urlId}", func(r chi.Router) {
				r.Get("/", getURL)
				r.Patch("/", updateURL)
				r.Delete("/", deleteURL)
			})
		})
	})
}

func (h *Handler) redirectTo(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "urlID")
	urlFormDB, err := url.SelectWithID(urlID, h.Env.DB)
	if err != nil {
		if err == models.SQLResult[models.NoResult] {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Default().Printf("redirectTo: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, urlFormDB.RedirectTo, http.StatusSeeOther)
}
func urlList(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) createURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Default().Printf("CreateURL: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newURL := &url.URL{
		ID:          r.Header.Get("urlID"),
		RedirectTo:  r.Header.Get("redirectTo"),
		Description: r.Header.Get("description"),
		UserID:      claims["user_id"].(string),
	}

	if newURL.ID == "" || newURL.RedirectTo == "" {
		message := "CreateURL: поля urlID и redirectTo не могут быть пустыми"
		log.Default().Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	if len(newURL.ID) > int(h.Env.Config.LengthURL) {
		message := fmt.Sprintf("CreateURL: максимальная длина короткого urlID равна %d", h.Env.Config.LengthURL)
		log.Default().Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	err = newURL.Insert(h.Env.DB)
	if err != nil {
		log.Default().Printf("CreateURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func getURL(w http.ResponseWriter, r *http.Request)    {}
func updateURL(w http.ResponseWriter, r *http.Request) {}
func deleteURL(w http.ResponseWriter, r *http.Request) {}
func (h *Handler) generateURL(w http.ResponseWriter, r *http.Request) {
	for {
		newURL := utils.RandString(h.Env.Config.LengthURL)
		urlFromDB, err := url.SelectWithID(newURL, h.Env.DB)
		if err != nil {
			log.Default().Println(err)
			if err.Error() != "sql: no rows in result set" {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if urlFromDB != nil {
			log.Default().Printf("generateURL: %s уже существует в БД\n", urlFromDB)
			continue
		}
		if t, ok := tmpURL[newURL]; ok {
			if time.Now().After(t) {
				log.Default().Printf("generateURL: %s зарезервирован\n", newURL)
				continue
			}
		}
		tmpURL[newURL] = time.Now().Add(time.Duration(h.Env.Config.TmpURLLifeTime))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("{\"shortURL\": \"/%s\"}", newURL))); err != nil {
			log.Default().Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
