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
	var (
		url   model.Model
		err   error
		urlTo interface{}
	)

	c.Debug <- "RedirectTo: получаем модель URL"
	if url = c.DB.Model(model.TableURL); url == nil {
		c.Err <- fmt.Errorf("redirectTo: nil")
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "RedirectTo: устанавливаем значение ИД для URL"
	if err = url.SetValue(model.FieldID, chi.URLParam(r, "urlID")); err != nil {
		c.Err <- fmt.Errorf("redirectTo: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "RedirectTo: получаем значение из БД"
	if err = c.DB.Read(model.TableURL); err != nil {
		if err == model.SQLResult[model.SQLNoResult] {
			c.Debug <- err.Error()
			Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
			return
		}
		c.Err <- fmt.Errorf("redirectTo: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "RedirectTo: получаем redirect to"
	if urlTo, err = url.Value(model.FieldRedirectTo); err != nil {
		c.Err <- fmt.Errorf("redirectTo: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	http.Redirect(w, r, urlTo.(string), http.StatusSeeOther)
}

func (c *Controller) URLList(w http.ResponseWriter, r *http.Request) {
	var (
		urls   []model.Model
		err    error
		usrID  interface{}
		claims map[string]interface{}
		jsn    []byte
	)
	c.Debug <- "URLList: получаем claims из токена"
	if _, claims, err = jwtauth.FromContext(r.Context()); err != nil {
		c.Debug <- fmt.Sprintf("urlList: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	if id, ok := claims["user_id"]; !ok {
		usrID = ""
	} else {
		usrID = id
	}

	c.Debug <- "URLList: получаем []Model"
	if urls, err = c.DB.ReadAll(model.TableURL, usrID.(string)); err != nil {
		if err != nil {
			c.Debug <- fmt.Sprintf("urlList: %v\n", err)
			Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
			return
		}
	}

	c.Debug <- "URLList: получаем JSON из []Model"
	if jsn, err = json.Marshal(urls); err != nil {
		c.Err <- fmt.Errorf("URLList: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "URLList: отправляем ответ"
	if _, err = w.Write(jsn); err != nil {
		c.Err <- fmt.Errorf("urlList: %w", err)
	}
	w.WriteHeader(http.StatusOK)
}

// nolint: funlen
func (c *Controller) CreateURL(w http.ResponseWriter, r *http.Request) {
	var (
		url         model.Model
		err         error
		usrID       interface{}
		urlID       interface{}
		urlTo       interface{}
		claims      map[string]interface{}
		requestBody []byte
	)
	c.Debug <- "CreateURL: получаем claims из токена"
	if _, claims, err = jwtauth.FromContext(r.Context()); err != nil {
		c.Debug <- fmt.Sprintf("CreateURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	usrID, ok := claims["user_id"]
	if !ok {
		// не должно возникать, по идее пользователь должен был отвалиться на этапе аутентификации
		c.Debug <- fmt.Sprintf("CreateURL: %v\n", "user_id не обнаружен")
		Response(w, http.StatusUnauthorized, model.ErrorResponseMap[http.StatusUnauthorized], c)
		return
	}

	c.Debug <- "CreateURL: Чтение Requust:Body"
	if requestBody, err = ioutil.ReadAll(r.Body); err != nil {
		c.Debug <- fmt.Sprintf("CreateURL: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	c.Debug <- fmt.Sprintf("CreateURL: Request:Body('%s')", requestBody)

	c.Debug <- "CreateURL: Получаем модель URL"
	url = c.DB.Model(model.TableURL)

	c.Debug <- "CreateURL: Заполням структуру из Request:Body"
	if err = url.FromJSON(requestBody); err != nil {
		c.Debug <- fmt.Sprintf("CreateURL: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "CreateURL: Устанавливаем значение user_id полченое из AuthToken"
	if err = url.SetValue(model.FieldUserID, usrID.(string)); err != nil {
		c.Err <- fmt.Errorf("CreateURL: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "CreateURL: Получаем ИД из структуры"
	if urlID, err = url.Value(model.FieldID); err != nil {
		c.Debug <- fmt.Sprintf("signUP: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "CreateURL: Получаем  redirecrTo из структуры"
	if urlTo, err = url.Value(model.FieldID); err != nil {
		c.Debug <- fmt.Sprintf("signUP: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "Проверям поля urlID и redirectTo"
	if urlID.(string) == "" || urlTo.(string) == "" {
		c.Debug <- "CreateURL: поля urlID и redirectTo не могут быть пустыми"
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	if len(urlID.(string)) > int(c.Config.LengthURL) {
		c.Debug <- fmt.Sprintf("CreateURL: максимальная длина короткого urlID равна %d", c.Config.LengthURL)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "CreateURL: создаем url в БД"
	if err = c.DB.Create(model.TableURL); err != nil {
		c.Debug <- fmt.Sprintf("CreateURL: %v", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) GetURL(w http.ResponseWriter, r *http.Request) {
	var (
		url    model.Model
		err    error
		usrID  interface{}
		urlID  interface{}
		claims map[string]interface{}
		jsn    []byte
	)
	c.Debug <- "GetURL: получаем claims из токена"
	if _, claims, err = jwtauth.FromContext(r.Context()); err != nil {
		c.Debug <- fmt.Sprintf("GetURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	usrID, ok := claims["user_id"]
	if !ok {
		// не должно возникать, по идее пользователь должен был отвалиться на этапе аутентификации
		c.Debug <- fmt.Sprintf("CreateURL: %v\n", "user_id не обнаружен")
		Response(w, http.StatusUnauthorized, model.ErrorResponseMap[http.StatusUnauthorized], c)
		return
	}
	urlID = chi.URLParam(r, "urlID")

	c.Debug <- "GetURL: Получаем модель URL"
	url = c.DB.Model(model.TableURL)

	c.Debug <- "GetURL: Устанавливаем значение user_id полченое из AuthToken"
	if err = url.SetValue(model.FieldUserID, usrID); err != nil {
		c.Err <- fmt.Errorf("GetURL: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "GetURL: Устанавливаем значение url_id полченое из URLParam"
	if err = url.SetValue(model.FieldID, urlID); err != nil {
		c.Err <- fmt.Errorf("GetURL: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "GetURL: Читаем данные из БД"
	if err = c.DB.Read(model.TableURL); err != nil {
		if err != nil {
			switch err.Error() {
			case model.SQLResult[model.SQLNoResult].Error():
				Response(w, http.StatusNotFound, model.ErrorResponseMap[http.StatusNotFound], c)
				return
			default:
				c.Debug <- fmt.Sprintf("GetURL: %v\n", err)
				Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
				return
			}
		}
	}

	c.Debug <- "GetURL: Формируем JSON"
	if jsn, err = url.JSON(); err != nil {
		c.Err <- fmt.Errorf("GetURL: %v", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}
	if _, err = w.Write(jsn); err != nil {
		c.Debug <- fmt.Sprintf("GetURL: %v", err)
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) UpdateURL(w http.ResponseWriter, r *http.Request) {
	var (
		url         model.Model
		err         error
		usrID       interface{}
		claims      map[string]interface{}
		requestBody []byte
	)
	c.Debug <- "UpdateURL: получаем claims из токена"
	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c.Debug <- "UpdateURL: Чтение Requust:Body"
	if requestBody, err = ioutil.ReadAll(r.Body); err != nil {
		c.Debug <- fmt.Sprintf("UpdateURL: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "UpdateURL: получаем claims из токена"
	if _, claims, err = jwtauth.FromContext(r.Context()); err != nil {
		c.Debug <- fmt.Sprintf("UpdateURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	usrID, ok := claims["user_id"]
	if !ok {
		// не должно возникать, по идее пользователь должен был отвалиться на этапе аутентификации
		c.Debug <- fmt.Sprintf("UpdateURL: %v\n", "user_id не обнаружен")
		Response(w, http.StatusUnauthorized, model.ErrorResponseMap[http.StatusUnauthorized], c)
		return
	}

	c.Debug <- "GetURL: Получаем модель URL"
	url = c.DB.Model(model.TableURL)

	c.Debug <- "UpdateURL: Заполням структуру из Request:Body"
	if err = url.FromJSON(requestBody); err != nil {
		c.Debug <- fmt.Sprintf("UpdateURL: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "UpdateURL: Устанавливаем значение user_id полченое из AuthToken"
	if err = url.SetValue(model.FieldUserID, usrID.(string)); err != nil {
		c.Err <- fmt.Errorf("UpdateURL: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	if err = c.DB.Update(model.TableURL); err != nil {
		if err != nil {
			c.Debug <- fmt.Sprintf("UpdateURL: %v\n", err)
			if err.Error()[:2] == "pq" {
				Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
				return
			}
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) DeleteURL(w http.ResponseWriter, r *http.Request) {
	var (
		url    model.Model
		err    error
		usrID  interface{}
		urlID  interface{}
		claims map[string]interface{}
	)
	c.Debug <- "DeleteURL: получаем claims из токена"
	if _, claims, err = jwtauth.FromContext(r.Context()); err != nil {
		c.Debug <- fmt.Sprintf("GetURL: %v\n", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	usrID, ok := claims["user_id"]
	if !ok {
		// не должно возникать, по идее пользователь должен был отвалиться на этапе аутентификации
		c.Debug <- fmt.Sprintf("DeleteURL: %v\n", "user_id не обнаружен")
		Response(w, http.StatusUnauthorized, model.ErrorResponseMap[http.StatusUnauthorized], c)
		return
	}
	urlID = chi.URLParam(r, "urlID")

	c.Debug <- "GetURL: Получаем модель URL"
	url = c.DB.Model(model.TableURL)

	c.Debug <- "DeleteURL: Устанавливаем значение url_id полченое из URLParam"
	if err = url.SetValue(model.FieldID, urlID); err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	c.Debug <- "DeleteURL: Устанавливаем значение user_id полченое из AuthToken"
	if err = url.SetValue(model.FieldUserID, usrID); err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	c.Debug <- "DeleteURL: Производим удаление из БД"
	if err = c.DB.Delete(model.TableURL); err != nil {
		c.Info <- fmt.Sprintf("DeleteURL: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) GenerateURL(w http.ResponseWriter, r *http.Request) {
	var (
		url   model.Model
		err   error
		urlID interface{}
		urlTo interface{}
	)
	c.Debug <- "GenerateURL: генерируем URL"
	for {
		c.Debug <- "GetURL: Получаем модель URL"
		url = c.DB.Model(model.TableURL)

		c.Debug <- "GenerateURL: устанавливае urlID"
		urlID = utils.RandString(c.Config.LengthURL)
		if err = url.SetValue(model.FieldID, urlID); err != nil {
			c.Err <- fmt.Errorf("GenerateURL: %v", err)
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
			return
		}

		c.Debug <- "GenerateURL: получаем данные по URL из БД"
		if err = c.DB.Read(model.TableURL); err != nil {
			if err.Error() != model.SQLResult[model.SQLNoResult].Error() {
				c.Err <- fmt.Errorf("GenerateURL: %v", err)
				Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
				return
			}
		}

		c.Debug <- "GenerateURL: получаем данные по urlTo"
		if urlTo, err = url.Value(model.FieldRedirectTo); err != nil {
			c.Err <- fmt.Errorf("GenerateURL: %v", err)
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
			return
		}

		c.Debug <- "GenerateURL: проверям redirectTo на корректность"
		if urlTo != "" {
			c.Debug <- fmt.Sprintf("GenerateURL: %s уже существует в БД", urlID)
			continue
		}

		c.Debug <- "GenerateURL: проверям не забронирована ли сгенерированная ссылка"
		if t, ok := tmpURL[urlID.(string)]; ok {
			if time.Now().After(t) {
				c.Info <- fmt.Sprintf("GenerateURL: %s зарезервирован\n", urlID.(string))
				continue
			}
		}
		tmpURL[urlID.(string)] = time.Now().Add(time.Duration(c.Config.TmpURLLifeTime))
		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(fmt.Sprintf("{\"shortURL\": \"%s\"}", urlID.(string)))); err != nil {
			c.Err <- fmt.Errorf("GenerateURL: %v", err)
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
			return
		}
		return
	}
}
