package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedirectTo(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		"GET",
		"/",
		bytes.NewBuffer([]byte("{\"email\":\""+email+"\",\"password\":\""+password+"\"}")))
	checkError(err, t)
	done := make(chan int)

	go func() {
		http.HandlerFunc(ctrl.RedirectTo).ServeHTTP(rec, req)
		assert.Equal(t, 303, rec.Code)
		assert.Equal(t, redirectTo, rec.HeaderMap.Get("Location"))
		done <- 1
	}()
	ListenChan(ctrl, done)
}

func TestURLList(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		"GET",
		"/url/",
		bytes.NewBuffer([]byte("{\"email\":\""+email+"\",\"password\":\""+password+"\"}")))
	checkError(err, t)
	done := make(chan int)

	go func() {
		http.HandlerFunc(ctrl.URLList).ServeHTTP(rec, req)
		assert.Equal(t, 200, rec.Code)
		assert.JSONEq(t, `[{"urlID":"qwqwq","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","redirectTo":"https://google.com","description":"","UserID":"ererqwerqerwe"}]`, rec.Body.String())
		done <- 1
	}()
	ListenChan(ctrl, done)
}

// func TestSignIn(t *testing.T) {
// 	rec := httptest.NewRecorder()
// 	req, err := http.NewRequest(
// 		"POST",
// 		"/auth/sign-in",
// 		bytes.NewBuffer([]byte("{\"email\":\""+email+"\",\"password\":\""+password+"\"}")))
// 	checkError(err, t)
// 	done := make(chan int)

// 	go func() {
// 		http.HandlerFunc(ctrl.SignUp).ServeHTTP(rec, req)

// 		//assert.NoError(t, err)

// 		assert.Equal(t, 200, rec.Code)
// 		//assert.Equal(t, respBody, rr.Body.Bytes())
// 		done <- 1
// 	}()
// 	ListenChan(ctrl, done)
// }
