package controller

import (
	"crypto/sha256"
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

// SignUp регистрация нового пользователя
// nolint: funlen	//длина функции высока из за большого количества проверок
func (c *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		requestBody []byte
		usr         model.Model
		pswd        interface{}
		email       interface{}
	)

	c.Debug <- "SignUp: Проверка Requust:Body == nil"
	if r.Body == nil {
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "SignUp: Чтение Requust:Body"
	if requestBody, err = ioutil.ReadAll(r.Body); err != nil {
		c.Info <- fmt.Sprintf("signUP: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	c.Debug <- fmt.Sprintf("SignUp: Request:Body('%s')", requestBody)

	c.Debug <- "SignUp: Получаем модель User"
	usr = c.DB.Model(model.TableUser)

	c.Debug <- "SignUp: Заполням структуру из Request:Body"
	if err = usr.FromJSON(requestBody); err != nil {
		c.Debug <- fmt.Sprintf("signUP: %v", err.Error())
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	// можно добавить проверку идентификатора на существование в БД

	c.Debug <- "SignUp: Получаем пароль из структуры"
	if pswd, err = usr.Value(model.FieldPassword); err != nil {
		c.Debug <- fmt.Sprintf("signUP: %v", err.Error())
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "SignUp: Получаем email из структуры"
	if email, err = usr.Value(model.FieldEmail); err != nil {
		c.Debug <- fmt.Sprintf("signUP: %v", err.Error())
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "SignUp: Проверяем email и пароль на пустоту"
	if email.(string) == "" || pswd.(string) == "" {
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	c.Debug <- "SignUp: генерируем и устанавливаем userID"
	if err = usr.SetValue(model.FieldID, uuid.New().String()); err != nil {
		c.Err <- fmt.Errorf("signUp-SetValue(password): %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignUp: хэшируем пароль"
	if err = usr.SetValue(model.FieldPassword, fmt.Sprintf("%x", sha256.Sum256([]byte(pswd.(string))))); err != nil {
		c.Err <- fmt.Errorf("signUp-SetValue(password): %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignUp: Устанавливаем значение id"
	if err = usr.SetValue(model.FieldID, uuid.New().String()); err != nil {
		c.Err <- fmt.Errorf("signUp-SetValue(id): %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignUp: создаем пользователя в БД"
	if err = c.DB.Create(model.TableUser); err != nil {
		switch err.Error() {
		case `pq: повторяющееся значение ключа нарушает ограничение уникальности "email"`:
			response := model.ErrorResponse{
				Code:    "200",
				Message: fmt.Sprintf("Пользователь с email: %s,  уже зарегистрирован.", email.(string)),
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

// nolint: funlen
func (c *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var (
		err         error
		requestBody []byte
		usr         model.Model
		usrID       interface{}
		inpswd      interface{}
		bdpswd      interface{}
		email       interface{}
		tokenString string
	)

	c.Debug <- "SignIn: Проверка Requust:Body == nil"
	if r.Body == nil {
		c.Debug <- "получено пустое Request:Body"
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "SignIn: Чтение Requust:Body"
	if requestBody, err = ioutil.ReadAll(r.Body); err != nil {
		c.Debug <- fmt.Sprintf("signIn: %v", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}
	c.Debug <- fmt.Sprintf("SignUp: Request:Body('%s')", requestBody)

	c.Debug <- "SignIn: Получаем модель User"
	usr = c.DB.Model(model.TableUser)

	c.Debug <- "SignUp: Заполням структуру из Request:Body"
	if err = usr.FromJSON(requestBody); err != nil {
		c.Err <- fmt.Errorf("signIn: %w", err)
		Response(w, http.StatusBadRequest, model.ErrorResponseMap[http.StatusBadRequest], c)
		return
	}

	c.Debug <- "SignIn: Получаем пароль из структуры"
	if inpswd, err = usr.Value(model.FieldPassword); err != nil {
		c.Err <- err
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignIn: Получаем user данные из БД"

	if err = c.DB.Read(model.TableUser); err != nil {
		c.Err <- fmt.Errorf("signIn: %w", err)
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignIn: Получаем ИД из структуры"
	if usrID, err = usr.Value(model.FieldID); err != nil {
		c.Err <- err
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignIn: Получаем Email из структуры"
	if email, err = usr.Value(model.FieldEmail); err != nil {
		c.Err <- err
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	c.Debug <- "SignIn: Проверяем email и пароль на пустоту"
	if usrID.(string) == "" || email.(string) == "" {
		Response(w,
			http.StatusOK,
			&model.ErrorResponse{
				Code:    "200",
				Message: fmt.Sprintf("Пользователь '%s' не найден\n", email.(string)),
			}, c)
		return
	}

	c.Debug <- "SignIn: хэшируем полученный от пользователя пароль"
	inpswd = fmt.Sprintf("%x", sha256.Sum256([]byte(inpswd.(string))))

	c.Debug <- "SignIn: Получаем пароль из структуры"
	if bdpswd, err = usr.Value(model.FieldPassword); err != nil {
		c.Err <- err
		Response(w, http.StatusInternalServerError, model.ErrorResponseMap[http.StatusInternalServerError], c)
		return
	}

	if bdpswd.(string) != inpswd {
		Response(w,
			http.StatusOK,
			&model.ErrorResponse{
				Code:    "200",
				Message: "Неверная пара Логин и Пароль",
			}, c)
		return
	}

	if _, tokenString, err = c.TokenAuth.Encode(map[string]interface{}{"user_id": usrID.(string)}); err != nil {
		log.Default().Println(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if _, err = w.Write([]byte(`{"Authorization": "BEARER ` + tokenString + `"}`)); err != nil {
		c.Err <- err
	}
	w.WriteHeader(http.StatusOK)
}
func (c *Controller) SignOut(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
