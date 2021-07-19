package controller

import "net/http"

type Responser interface {
	JSON() ([]byte, error)
}

func Response(w http.ResponseWriter, statusCode int, r Responser, ctrl *Controller) {
	w.WriteHeader(statusCode)
	json, err := r.JSON()
	if err != nil {
		ctrl.Err <- err
		return
	}

	if _, err := w.Write(json); err != nil {
		ctrl.Warn <- err
	}
}
