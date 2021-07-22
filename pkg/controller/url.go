package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	r.Get("/{urlID}", c.RedirectTo)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Route("/url", func(r chi.Router) {
			r.Get("/", c.URLList)
			r.Post("/", c.CreateURL)
			r.Patch("/", c.UpdateURL)
			r.Get("/generate", c.GenerateURL)
			r.Route("/{urlID}", func(r chi.Router) {
				r.Get("/", c.GetURL)
				r.Delete("/", c.DeleteURL)
			})
		})
	})
}

func (c *Controller) RedirectTo(w http.ResponseWriter, r *http.Request) {
	url := c.DB.Model(model.TableURL)
	if url == nil {
		c.Err <- fmt.Errorf("redirectTo: nil")
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	err := url.SetValue(model.DBFieldID, chi.URLParam(r, "urlID"))
	if err != nil {
		c.Info <- fmt.Sprintf("redirectTo: %v\n", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	err = c.DB.Read(model.TableURL)
	if err != nil {
		if err == model.SQLResult[model.SQLNoResult] {
			Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
			return
		}
		c.Info <- fmt.Sprintf("redirectTo: %v\n", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	http.Redirect(w, r, url.Value(model.DBFieldRedirectTo).(string), http.StatusSeeOther)
}

func (c *Controller) URLList(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		c.Info <- fmt.Sprintf("urlList: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	userID := ""
	if claims["user_id"] != nil {
		userID = claims["user_id"].(string)
	}

	urls, err := c.DB.ReadAll(model.TableURL, userID)
	if err != nil {
		if err != nil {
			c.Info <- fmt.Sprintf("urlList: %v\n", err)
			Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
			return
		}
	}

	buf, err := json.Marshal(urls)
	if err != nil {
		c.Debug <- fmt.Errorf("URLList: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(buf)
	if err != nil {
		c.Err <- fmt.Errorf("urlList: %w", err)
	}
}

func (c *Controller) CreateURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		c.Info <- fmt.Sprintf("CreateURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	userID := ""
	if claims["user_id"] == nil {
		// не должно возникать, по идее пользователь должен был отвалиться на этапе аутентификации
		message := fmt.Sprintf("CreateURL: %v\n", "user_id не обнаружен")
		c.Info <- message
		Response(w, http.StatusUnauthorized, model.ErrorResponseMap[http.StatusUnauthorized], c)
		return
	}
	userID = claims["user_id"].(string)

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Debug <- fmt.Errorf("CreateURL: %w", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	url := c.DB.Model(model.TableURL)
	err = url.FromJSON(requestBody)
	if err != nil {
		c.Debug <- fmt.Errorf("CreateURL: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	err = url.SetValue(model.DBFieldUserID, userID)
	if err != nil {
		c.Debug <- fmt.Errorf("CreateURL: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	if url.Value(model.DBFieldID) == "" || url.Value(model.DBFieldRedirectTo) == "" {
		message := "CreateURL: поля urlID и redirectTo не могут быть пустыми"
		c.Info <- message
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	if len(url.Value(model.DBFieldID).(string)) > int(c.Config.LengthURL) {
		message := fmt.Sprintf("CreateURL: максимальная длина короткого urlID равна %d", c.Config.LengthURL)
		c.Info <- message
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	err = c.DB.Create(model.TableURL)
	if err != nil {
		c.Info <- fmt.Sprintf("CreateURL: %v", err)

		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) GetURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		c.Info <- fmt.Sprintf("getURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	err = c.DB.Model(model.TableURL).SetValue(model.DBFieldUserID, claims["user_id"])
	if err != nil {
		c.Err <- fmt.Errorf("getURL: %v", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	err = c.DB.Read(model.TableURL)
	if err != nil {
		if err != nil {
			c.Info <- fmt.Sprintf("getURL: %v\n", err)
			Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
			return
		}
	}

	response, err := c.DB.Model(model.TableURL).JSON()
	if err != nil {
		c.Err <- fmt.Errorf("getURL: %v", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		c.Info <- fmt.Sprintf("getURL: %v", err)
	}
}

func (c *Controller) UpdateURL(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Info <- fmt.Sprintf("UpdateURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		c.Info <- fmt.Sprintf("UpdateURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	url := c.DB.Model(model.TableURL)
	err = url.FromJSON(requestBody)
	if err != nil {
		c.Err <- fmt.Errorf("UpdateURL: %v", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	err = url.SetValue(model.DBFieldUserID, claims["user_id"])
	if err != nil {
		c.Err <- fmt.Errorf("UpdateURL: %v", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	err = c.DB.Update(model.TableURL)
	if err != nil {
		if err != nil {
			c.Info <- fmt.Sprintf("UpdateURL: %v\n", err)
			Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) DeleteURL(w http.ResponseWriter, r *http.Request) {
	_, claims, err := jwtauth.FromContext(r.Context())
	if err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	url := c.DB.Model(model.TableURL)

	if err = url.SetValue(model.DBFieldID, chi.URLParam(r, "urlID")); err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = url.SetValue(model.DBFieldUserID, claims["user_id"]); err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err = c.DB.Delete(model.TableURL); err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) GenerateURL(w http.ResponseWriter, r *http.Request) {
	for {
		url := c.DB.Model(model.TableURL)
		if err := url.SetValue(model.DBFieldID, utils.RandString(c.Config.LengthURL)); err != nil {
			c.Err <- fmt.Errorf("GenerateURL: %v", err)
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
			return
		}

		if err := c.DB.Read(model.TableURL); err != nil {
			if err != model.SQLResult[model.SQLNoResult] {
				c.Err <- fmt.Errorf("GenerateURL: %v", err)
				Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
				return
			}
		}
		if url.Value(model.DBFieldRedirectTo).(string) != "" {
			c.Info <- fmt.Sprintf("GenerateURL: %s уже существует в БД\n", url.Value(model.DBFieldID))
			continue
		}
		if t, ok := tmpURL[url.Value(model.DBFieldID).(string)]; ok {
			if time.Now().After(t) {
				c.Info <- fmt.Sprintf("GenerateURL: %s зарезервирован\n", url.Value(model.DBFieldID))
				continue
			}
		}
		tmpURL[url.Value(model.DBFieldID).(string)] = time.Now().Add(time.Duration(c.Config.TmpURLLifeTime))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("{\"shortURL\": \"%s\"}", url.Value(model.DBFieldID)))); err != nil {
			c.Err <- fmt.Errorf("GenerateURL: %v", err)
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
			return
		}
		return
	}
}
