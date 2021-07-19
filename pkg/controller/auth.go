package controller

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"

	"github.com/AleksandrMac/ushort/pkg/model"
)

type Claims struct {
	jwt.StandardClaims
	UserEmail string
}

func (c *Controller) setAuthControllers(r *chi.Mux) {
	r.Post("/auth/sign-up", c.SignUp)
	r.Post("/auth/sign-in", c.SignIn)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/auth/sign-out", c.SignOut)
	})
}

// авторизация в данном экземпляре носит второстепенную роль
// поэтому токен выдается без синхронизации с БД, и без даты истечения срока действия
// вместо логаута заглушка
func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Debug <- fmt.Errorf("signUP: %w", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	usr, ok := c.DB.Model(model.UserTable).(*model.User)
	if !ok {
		c.Err <- fmt.Errorf("signUp-(*User): error convertation")
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	err = json.Unmarshal(requestBody, usr)
	if err != nil {
		c.Debug <- fmt.Errorf("signUP: %w", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	// можно добавить проверку идентификатора на существование в БД

	if usr.Email == "" || usr.Password == "" {
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	usr.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(usr.Password)))
	usr.ID = uuid.New().String()

	err = c.DB.Create(model.UserTable)
	if err != nil {
		switch err.Error() {
		case `pq: повторяющееся значение ключа нарушает ограничение уникальности "email"`:
			response := model.ErrorResponse{
				Code:    "200",
				Message: fmt.Sprintf("Пользователь с email: %s,  уже зарегистрирован.", usr.Email),
			}
			Response(w, http.StatusOK, &response, c)
		default:
			c.Err <- fmt.Errorf("signUp-insert: %w", err)
			Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
}
func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		c.Debug <- fmt.Errorf("signIn: %w", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	usr, ok := c.DB.Model(model.UserTable).(*model.User)
	if !ok {
		c.Err <- fmt.Errorf("signIn-(*User): error convertation")
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	err = json.Unmarshal(requestBody, usr)
	if err != nil {
		c.Debug <- fmt.Errorf("signIn: %w", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	inputPassword := usr.Password

	err = c.DB.Read(model.UserTable)
	if err != nil {
		c.Debug <- fmt.Errorf("signIn: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	if usr.ID == "" || usr.Email == "" {
		Response(w,
			http.StatusOK,
			&model.ErrorResponse{
				Code:    "200",
				Message: fmt.Sprintf("Пользователь '%s' не найден\n", usr.Email),
			}, c)
		return
	}

	inputPassword = fmt.Sprintf("%x", sha256.Sum256([]byte(inputPassword)))

	if usr.Password != inputPassword {
		Response(w,
			http.StatusOK,
			&model.ErrorResponse{
				Code:    "200",
				Message: "Неверная пара Логин и Пароль",
			}, c)
		return
	}

	_, tokenString, err := c.TokenAuth.Encode(map[string]interface{}{"user_id": usr.ID})
	if err != nil {
		log.Default().Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", "BEARER "+tokenString)
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) SignOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
