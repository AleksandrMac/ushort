package controller_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AleksandrMac/ushort/pkg/controller"
	"github.com/AleksandrMac/ushort/pkg/model"
	"github.com/stretchr/testify/assert"
)

type mockDB struct {
	user *mockUser
	url  *mockURL
}

const (
	email    = "proto@example.com"
	password = "12345"
)

type mockUser struct{}

func (u *mockDB) Model(model.Table) model.Model { return &mockUser{} }
func (u *mockDB) Create(model.Table) error      { return nil }
func (u *mockDB) Read(model.Table) error        { return nil }
func (u *mockDB) Update(model.Table) error      { return nil }
func (u *mockDB) Delete(model.Table) error      { return nil }

func (u *mockUser) Fields() ([]string, error)                    { return nil, nil }
func (u *mockUser) Values() (map[string]interface{}, error)      { return nil, nil }
func (u *mockUser) Value(field string) interface{}               { return nil }
func (u *mockUser) SetValue(field string, val interface{}) error { return nil }
func (u *mockUser) JSON() ([]byte, error)                        { return nil, nil }

type mockURL struct{}

func (u *mockURL) Fields() ([]string, error)                    { return nil, nil }
func (u *mockURL) Values() (map[string]interface{}, error)      { return nil, nil }
func (u *mockURL) Value(field string) interface{}               { return nil }
func (u *mockURL) SetValue(field string, val interface{}) error { return nil }
func (u *mockURL) JSON() ([]byte, error)                        { return nil, nil }

var ctrl = &controller.Controller{
	DB:       &mockDB{},
	Info:     make(chan string),
	Debug:    make(chan error),
	Err:      make(chan error),
	Warn:     make(chan error),
	Critical: make(chan error),
}

func testController() *controller.Controller {
	return &controller.Controller{DB: &mockDB{user: &mockUser{}}}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
}

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

func ListenChan(ctrl *controller.Controller, done chan int) {
	for {
		select {
		case <-done:
			return
		case <-ctrl.Info:
		case <-ctrl.Debug:
		case <-ctrl.Err:
		case <-ctrl.Warn:
		case <-ctrl.Critical:
		}
	}
}
