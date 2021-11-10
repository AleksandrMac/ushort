package controller_test

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/controller"
	"github.com/AleksandrMac/ushort/pkg/model"

	"github.com/go-chi/jwtauth/v5"
	"github.com/opentracing/opentracing-go/mocktracer"
)

const (
	userID        = "usertest"
	urlID         = "dfsdlic4fr"
	email         = "proto@example.com"
	password      = "12345"
	redirectTo    = "https://shop.com/items?param1=somevalue1&param2=somevalue2&param3=somevalue3"
	description   = "instagram promo"
	Authorization = "BEARER eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcnRlc3QifQ.FiaKOEqI-pGGn8OhzlgZPmBgSEvRg3kiML2EIIgZiFw"
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
func (u *mockDB) Read(table model.Table) error {
	switch table {
	case model.TableUser:
		u.mock.password = fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	}
	return nil
}
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

type mock struct {
	password string
}

func (u *mock) Fields() ([]string, error)                    { return nil, nil }
func (u *mock) Values() (map[model.Field]interface{}, error) { return nil, nil }
func (u *mock) Value(field model.Field) (interface{}, error) {
	switch field {
	case model.FieldID:
		return userID, nil
	case model.FieldEmail:
		return email, nil
	case model.FieldPassword:
		return u.password, nil
	case model.FieldRedirectTo:
		return redirectTo, nil
	}
	return nil, fmt.Errorf("not find field")
}
func (u *mock) SetValue(field model.Field, val interface{}) error {
	switch field {
	case model.FieldPassword:
		u.password = val.(string)
	}
	return nil
}
func (u *mock) JSON() ([]byte, error) {
	return []byte(`[{"urlID":"qwqwq","createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z","redirectTo":"` + redirectTo + `","description":"","UserID":"ererqwerqerwe"}]`), nil
}
func (u *mock) FromJSON([]byte) error {
	u.password = password
	return nil
}

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

var tracer = mocktracer.New()

var ctrl = &controller.Controller{
	DB:        &mockDB{mock: new(mock)},
	TokenAuth: jwtauth.New("HS256", []byte("secret"), nil),
	Config:    &config.Config{LengthURL: 10},
	Tracer:    tracer,
	Ctx:       context.Background(),
	Info:      make(chan string),
	Debug:     make(chan string),
	Err:       make(chan error),
	Warn:      make(chan error),
	Critical:  make(chan error),
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

func testRequest(t *testing.T, ts *httptest.Server, method, path string, header http.Header, body io.Reader) (int, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	if header != nil {
		for k, v := range header {
			req.Header.Set(k, v[0])
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
		return 0, ""
	}
	defer resp.Body.Close()

	return resp.StatusCode, string(respBody)
}
