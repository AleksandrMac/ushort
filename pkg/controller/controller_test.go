package controller_test

import (
	"testing"

	"github.com/AleksandrMac/ushort/pkg/controller"
	"github.com/AleksandrMac/ushort/pkg/model"
)

const (
	email      = "proto@example.com"
	password   = "12345"
	redirectTo = "https://google.com"
)

var (
	url1 model.URL = model.URL{
		Base: model.Base{
			ID: "qwqwq",
		},
		RedirectTo:  redirectTo,
		Description: "",
		UserID:      "ererqwerqerwe",
	}
)

type mockDB struct {
	mock *mock
}

func (u *mockDB) Model(table model.Table) model.Model { return u.mock }
func (u *mockDB) Create(model.Table) error            { return nil }
func (u *mockDB) Read(model.Table) error              { return nil }
func (u *mockDB) ReadAll(table model.Table, userID string) ([]model.Model, error) {
	switch table {
	case model.TableUser:
		return []model.Model{}, nil
	case model.TableURL:
		out := make([]model.Model, 1)
		out[0] = &url1
		return out, nil
	}
	return nil, nil
}
func (u *mockDB) Update(model.Table) error { return nil }
func (u *mockDB) Delete(model.Table) error { return nil }

type mock struct{}

func (u *mock) Fields() ([]string, error)                      { return nil, nil }
func (u *mock) Values() (map[model.DBField]interface{}, error) { return nil, nil }
func (u *mock) Value(field model.DBField) interface{} {
	switch field {
	case model.DBFieldEmail:
		return email
	case model.DBFieldPassword:
		return password
	case model.DBFieldRedirectTo:
		return redirectTo
	}
	return nil
}
func (u *mock) SetValue(field model.DBField, val interface{}) error { return nil }
func (u *mock) JSON() ([]byte, error)                               { return nil, nil }
func (u *mock) FromJSON([]byte) error                               { return nil }

//type mockUser struct{}
// func (u *mockUser) Fields() ([]string, error)                    { return nil, nil }
// func (u *mockUser) Values() (map[string]interface{}, error)      { return nil, nil }
// func (u *mockUser) Value(field string) interface{}               { return nil }
// func (u *mockUser) SetValue(field string, val interface{}) error { return nil }
// func (u *mockUser) JSON() ([]byte, error)                        { return nil, nil }

// type mockURL struct{}

// func (u *mockURL) Fields() ([]string, error)                    { return nil, nil }
// func (u *mockURL) Values() (map[string]interface{}, error)      { return nil, nil }
// func (u *mockURL) Value(field string) interface{}               { return nil }
// func (u *mockURL) SetValue(field string, val interface{}) error { return nil }
// func (u *mockURL) JSON() ([]byte, error)                        { return nil, nil }

var ctrl = &controller.Controller{
	DB:       &mockDB{},
	Info:     make(chan string),
	Debug:    make(chan error),
	Err:      make(chan error),
	Warn:     make(chan error),
	Critical: make(chan error),
}

func testController() *controller.Controller {
	return &controller.Controller{DB: &mockDB{mock: &mock{}}}
}

func checkError(err error, t *testing.T) {
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}
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
