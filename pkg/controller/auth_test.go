package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignUp(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		"/auth/sign-up",
		bytes.NewBuffer([]byte("{\"email\":\""+email+"\",\"password\":\""+password+"\"}")))
	checkError(err, t)
	done := make(chan int)

	go func() {
		http.HandlerFunc(ctrl.SignUp).ServeHTTP(rec, req)

		//assert.NoError(t, err)

		assert.Equal(t, 200, rec.Code)
		//assert.Equal(t, respBody, rr.Body.Bytes())
		done <- 1
	}()
	ListenChan(ctrl, done)
}

func TestSignIn(t *testing.T) {
	rec := httptest.NewRecorder()
	req, err := http.NewRequest(
		"POST",
		"/auth/sign-in",
		bytes.NewBuffer([]byte("{\"email\":\""+email+"\",\"password\":\""+password+"\"}")))
	checkError(err, t)
	done := make(chan int)

	go func() {
		http.HandlerFunc(ctrl.SignIn).ServeHTTP(rec, req)

		//assert.NoError(t, err)

		assert.Equal(t, 200, rec.Code)
		assert.Equal(t, `{"Authorization": "`+Authorization+`"}`, string(rec.Body.Bytes()))
		done <- 1
	}()
	ListenChan(ctrl, done)
}
