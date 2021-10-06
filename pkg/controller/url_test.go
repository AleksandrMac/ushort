package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
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
		assert.JSONEq(t, `[{"urlID":"qwqwq","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","redirectTo":"`+redirectTo+`","description":"","UserID":"ererqwerqerwe"}]`, rec.Body.String())
		done <- 1
	}()
	ListenChan(ctrl, done)
}

func TestCreateURL(t *testing.T) {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(ctrl.TokenAuth))
	r.Get("/url/", ctrl.CreateURL)

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := http.Header{}
	h.Set("Authorization", Authorization)
	done := make(chan int)
	go func() {
		status, resp := testRequest(t, ts, "GET", "/url/", h, nil)
		assert.Equal(t, 200, status, resp)
		done <- 1
	}()
	ListenChan(ctrl, done)
}

func TestGetURL(t *testing.T) {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(ctrl.TokenAuth))
	r.Post("/url/{urlID}", ctrl.GetURL)

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := http.Header{}
	h.Set("Authorization", Authorization)
	done := make(chan int)
	go func() {
		status, resp := testRequest(t, ts, "POST", "/url/qwqwq", h, nil)
		assert.Equal(t, 200, status, resp)
		assert.JSONEq(t, `[{"urlID":"qwqwq","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","redirectTo":"`+redirectTo+`","description":"","UserID":"ererqwerqerwe"}]`, resp)
		done <- 1
	}()
	ListenChan(ctrl, done)
}

func TestUpdateURL(t *testing.T) {
	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(ctrl.TokenAuth))
	r.Patch("/url/", ctrl.UpdateURL)

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := http.Header{}
	h.Set("Authorization", Authorization)
	done := make(chan int)
	go func() {
		status, resp := testRequest(t, ts, "PATCH", "/url/", h,
			bytes.NewBuffer([]byte(`{
				"urlID": "dfsdlic4fr",
				"redirectTo": "https://shop.com/items?param1=somevalue1&param2=somevalue2&param3=somevalue3",
				"description": "instagram promo"
			}`)))
		assert.Equal(t, 200, status, resp)
		done <- 1
	}()
	ListenChan(ctrl, done)
}

func TestDeleteURL(t *testing.T) {

	r := chi.NewRouter()
	r.Use(jwtauth.Verifier(ctrl.TokenAuth))
	r.Delete("/url/{urlID}", ctrl.DeleteURL)

	ts := httptest.NewServer(r)
	defer ts.Close()

	h := http.Header{}
	h.Set("Authorization", Authorization)
	done := make(chan int)
	go func() {
		status, resp := testRequest(t, ts, "DELETE", "/url/qwqwq", h, nil)
		assert.Equal(t, 200, status, resp)
		done <- 1
	}()
	ListenChan(ctrl, done)
}
