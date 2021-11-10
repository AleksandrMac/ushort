package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/AleksandrMac/ushort/pkg/controller"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
)

func main() {
	wd, _ := os.Getwd()
	log.Default().Println(wd)
	ctrl, err := controller.New()
	if err != nil {
		log.Fatal(err)
	}
	defer ctrl.Close()

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(*ctrl.Logger))
	ctrl.SetControllers(r)

	srv := &http.Server{
		Addr:    ":" + ctrl.Config.Server.Port,
		Handler: r,
	}

	go watchSignals(ctrl)

	go func() {
		ctrl.Info <- "server starting"
		ctrl.Critical <- srv.ListenAndServe()
	}()

	ListenChan(ctrl, srv)
}

func watchSignals(ch *controller.Controller) {
	osSignalChan := make(chan os.Signal, 1)

	signal.Notify(osSignalChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM)
	sig := <-osSignalChan
	ch.Info <- fmt.Sprintf("got signal %q", sig.String())

	// если сигнал получен, отменяем контекст работы
	ch.CtxCancel()
}

func ListenChan(ctrl *controller.Controller, srv *http.Server) {
	for {
		select {
		case <-ctrl.Ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ctrl.Config.ServerGraceFullTime)*time.Second)
			err := srv.Shutdown(ctx)
			if err != nil {
				ctrl.Logger.Err(err)
			}
			ctrl.Logger.Info().Msg("server stoped")
			defer cancel()
			return
		case info := <-ctrl.Info:
			ctrl.Logger.Info().Msg(info)
		case deb := <-ctrl.Debug:
			ctrl.Logger.Debug().Msg(deb)
		case err := <-ctrl.Err:
			ctrl.Logger.Error().Msg(err.Error())
		case err := <-ctrl.Warn:
			ctrl.Logger.Error().Msg(err.Error())
		case err := <-ctrl.Critical:
			ctrl.Logger.Fatal().Msg(err.Error())
		}
	}
}
