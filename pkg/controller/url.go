package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/AleksandrMac/ushort/pkg/model"
	"github.com/AleksandrMac/ushort/pkg/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

// tmpUrl хранит ссылки котырые были сгенерированы, но еще не попали в БД
var tmpURL map[string]time.Time

func (c *Controller) setURLControllers(r *chi.Mux) {
	tmpURL = make(map[string]time.Time)
	r.Get("/{urlID}", c.redirectTo)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Route("/url", func(r chi.Router) {
			r.Get("/", c.urlList)
			r.Post("/", c.createURL)
			r.Patch("/", c.updateURL)
			r.Get("/generate", c.generateURL)

			r.Route("/{urlID}", func(r chi.Router) {
				r.Get("/", c.getURL)
				r.Delete("/", c.deleteURL)
			})
		})
	})
}

func (c *Controller) redirectTo(w http.ResponseWriter, r *http.Request) {
	url := &model.URL{
		Model: model.Model{
			ID: chi.URLParam(r, "urlID"),
		},
		// UserID устанавливается "*" чтобы выборка прошла без учета user_id
		UserID: "*",
	}
	urlFormDB, err := url.SelectWithID(c.DB)
	if err != nil {
		if err == model.SQLResult[model.SQLNoResult] {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		log.Default().Printf("redirectTo: %v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, urlFormDB.RedirectTo, http.StatusSeeOther)
}
func (c *Controller) urlList(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Default().Printf("urlList: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	url := &model.URL{
		UserID: claims["user_id"].(string),
	}

	urls, err := url.Select(c.DB)
	if err != nil {
		if err != nil {
			log.Default().Printf("urlList: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	buf, err := json.Marshal(urls)
	if err != nil {
		log.Default().Printf("urlList: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		log.Default().Printf("urlList: %v", err)
	}
}
func (c *Controller) createURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Default().Printf("CreateURL: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := &model.URL{
		Model: model.Model{
			ID: r.Header.Get("urlID"),
		},
		RedirectTo:  r.Header.Get("redirectTo"),
		Description: r.Header.Get("description"),
		UserID:      claims["user_id"].(string),
	}

	if url.ID == "" || url.RedirectTo == "" {
		message := "CreateURL: поля urlID и redirectTo не могут быть пустыми"
		log.Default().Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	if len(url.ID) > int(c.Config.LengthURL) {
		message := fmt.Sprintf("CreateURL: максимальная длина короткого urlID равна %d", c.Config.LengthURL)
		log.Default().Println(message)
		http.Error(w, message, http.StatusBadRequest)
		return
	}
	err = url.Insert(c.DB)
	if err != nil {
		log.Default().Printf("CreateURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) getURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Default().Printf("getURL: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := &model.URL{
		Model: model.Model{
			ID: chi.URLParam(r, "urlID"),
		},
		RedirectTo:  r.Header.Get("redirectTo"),
		Description: r.Header.Get("description"),
		UserID:      claims["user_id"].(string),
	}

	urlFromDB, err := url.SelectWithID(c.DB)
	if err != nil {
		if err != nil {
			log.Default().Printf("getURL: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	response, err := json.Marshal(urlFromDB)
	if err != nil {
		log.Default().Printf("getURL: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Default().Printf("getURL: %v", err)
	}
}
func (c *Controller) updateURL(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Default().Printf("updateURL: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Default().Printf("updateURL: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	urlFromBody := &model.URL{
		Model:  model.Model{},
		UserID: claims["user_id"].(string),
	}
	err = json.Unmarshal(requestBody, urlFromBody)
	if err != nil {
		log.Default().Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = urlFromBody.Update(c.DB)
	if err != nil {
		if err != nil {
			log.Default().Printf("updateURL: %v\n", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) deleteURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		log.Default().Printf("deleteURL: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url := &model.URL{
		Model: model.Model{
			ID: chi.URLParam(r, "urlID"),
		},
		UserID: claims["user_id"].(string),
	}

	err = url.Delete(c.DB)
	if err != nil {
		log.Default().Printf("deleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) generateURL(w http.ResponseWriter, r *http.Request) {
	for {
		url := &model.URL{
			Model: model.Model{
				ID: utils.RandString(c.Config.LengthURL),
			},
			UserID: "*",
		}

		urlFromDB, err := url.SelectWithID(c.DB)
		if err != nil {
			log.Default().Println(err)
			if err != model.SQLResult[model.SQLNoResult] {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		if urlFromDB != nil {
			log.Default().Printf("generateURL: %s уже существует в БД\n", urlFromDB)
			continue
		}
		if t, ok := tmpURL[url.ID]; ok {
			if time.Now().After(t) {
				log.Default().Printf("generateURL: %s зарезервирован\n", url.ID)
				continue
			}
		}
		tmpURL[url.ID] = time.Now().Add(time.Duration(c.Config.TmpURLLifeTime))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("{\"shortURL\": \"%s\"}", url.ID))); err != nil {
			log.Default().Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
}
