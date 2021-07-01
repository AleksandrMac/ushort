package handle

import (
	"crypto/sha256"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	"github.com/AleksandrMac/ushort/pkg/models/user"
)

type Claims struct {
	jwt.StandardClaims
	UserEmail string
}

func (env *Env) setAuthHandlers(r *chi.Mux) {
	r.Post("/auth/sign-up", env.signUp)
	r.Post("/auth/sign-in", env.signIn)

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(env.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/auth/sign-out", env.signOut)
	})
}

// авторизация в данном экземпляре носить второстепенную роль
// поэтому токен выдается без синхронизации с БД, и без даты истечения срока действия
// вместо логаута заглушка
func (env *Env) signUp(w http.ResponseWriter, r *http.Request) {
	usr := user.User{
		Email:    r.Header.Get("email"),
		Password: r.Header.Get("password"),
	}

	if usr.Email == "" || usr.Password == "" {
		err := http.StatusText(http.StatusBadRequest)
		log.Default().Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}
	usr.Password = fmt.Sprintf("%x", sha256.Sum256([]byte(usr.Password)))

	id, err := usr.Insert(env.DB)
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

	n, err := w.Write([]byte(fmt.Sprintf("{\"id\": %s}", id)))
	if err != nil {
		log.Default().Println(err)
		return
	}
	log.Default().Printf("Записано байт: %d\n", n)
	http.Redirect(w, r, "/auth/sign-in", http.StatusSeeOther)
}
func (env *Env) signIn(w http.ResponseWriter, r *http.Request) {
	usr := user.User{
		Email:    r.Header.Get("email"),
		Password: r.Header.Get("password"),
	}

	usrFromDB, err := user.SelectWithEmail(usr.Email, env.DB)
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

	_, tokenString, err := env.TokenAuth.Encode(map[string]interface{}{"user_email": usr.Email})
	if err != nil {
		log.Default().Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", "BEARER "+tokenString)
	w.WriteHeader(http.StatusOK)
}
func (env *Env) signOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
