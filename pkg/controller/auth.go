package controller

import (
	"crypto/sha256"
	"fmt"
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
	r.Post("/auth/sign-up", c.signUp)
	r.Post("/auth/sign-in", c.signIn)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/auth/sign-out", c.signOut)
	})
}

// авторизация в данном экземпляре носить второстепенную роль
// поэтому токен выдается без синхронизации с БД, и без даты истечения срока действия
// вместо логаута заглушка
func (c *Controller) signUp(w http.ResponseWriter, r *http.Request) {
	usr := model.User{
		Model: model.Model{
			ID: uuid.New().String(),
		},
		Email:    r.Header.Get("email"),
		Password: r.Header.Get("password"),
	}
	// можно добавить проверку идентификатора на существование в БД

	if usr.Email == "" || usr.Password == "" {
		err := http.StatusText(http.StatusBadRequest)
		log.Default().Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}
	usr.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(usr.Password)))

	err := usr.Insert(c.DB)
	if err != nil {
		switch err.Error() {
		case `pq: повторяющееся значение ключа нарушает ограничение уникальности "email"`:
			errString := fmt.Sprintf("Пользователь с email: %s,  уже зарегистрирован.", usr.Email)
			http.Error(w, errString, http.StatusBadRequest)
		default:
			log.Default().Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	n, err := w.Write([]byte(fmt.Sprintf("{\"id\": %s}", usr.ID)))
	if err != nil {
		log.Default().Println(err)
		return
	}
	log.Default().Printf("Записано байт: %d\n", n)
	http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
}
func (c *Controller) signIn(w http.ResponseWriter, r *http.Request) {
	usr := model.User{
		Email:    r.Header.Get("email"),
		Password: r.Header.Get("password"),
	}

	usrFromDB, err := model.SelectWithEmail(usr.Email, c.DB)
	if err != nil {
		log.Default().Println(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if usrFromDB == nil {
		err := fmt.Sprintf("Пользователь '%s' не найден\n", usr.Email)
		log.Default().Printf(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}
	usr.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(usr.Password)))

	if usr.Password != usrFromDB.Password {
		err := "Неверная пара Логин и Пароль"
		log.Default().Printf(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	_, tokenString, err := c.TokenAuth.Encode(map[string]interface{}{"user_id": usrFromDB.ID})
	if err != nil {
		log.Default().Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", "BEARER "+tokenString)
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) signOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
