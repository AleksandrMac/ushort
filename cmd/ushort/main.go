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

	"github.com/AleksandrMac/ushort/pkg/config/env"
	"github.com/AleksandrMac/ushort/pkg/connect"
	"github.com/AleksandrMac/ushort/pkg/controller"
	"github.com/rs/zerolog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
)

func main() {
	cfg, err := env.New()
	if err != nil {
		log.Panic(err)
	}

	var logger zerolog.Logger
	switch cfg.LogLevel {
	case "trace":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "trace",
			JSON:     true})
	case "debug":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "debug",
			JSON:     true})
	case "info":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "info",
			JSON:     true})
	case "warn":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "warn",
			JSON:     true})
	case "error":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "error",
			JSON:     true})
	case "critical":
		logger = httplog.NewLogger("ushort", httplog.Options{
			LogLevel: "critical",
			JSON:     true})
	}

	db, err := connect.NewDB(cfg.DB.URL)
	if err != nil {
		log.Panic(err)
	}
	ctxMain, cancelMain := context.WithCancel(context.Background())
	defer cancelMain()

	ctrl := &controller.Controller{
		DB:        db,
		Config:    cfg,
		TokenAuth: jwtauth.New("HS256", cfg.Auth.PrivateKey, nil),
		Logger:    &logger,
		Info:      make(chan string),
		Debug:     make(chan error),
		Err:       make(chan error),
		Critical:  make(chan error),
	}

	r := chi.NewRouter()
	r.Use(httplog.RequestLogger(logger))
	ctrl.SetControllers(r)

	srv := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: r,
	}

	go watchSignals(cancelMain, ctrl)

	go func() {
		ctrl.Err <- srv.ListenAndServe()
	}()

	for {
		select {
		case <-ctxMain.Done():
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ctrl.Config.ServerGraceFullTime)*time.Second)
			ctrl.Err <- srv.Shutdown(ctx)
			defer cancel()
			return
		case err := <-ctrl.Err:
			logger.Error().Msg(err.Error())
		}
	}
}

func watchSignals(cancel context.CancelFunc, ch *controller.Controller) {
	osSignalChan := make(chan os.Signal, 1)

	signal.Notify(osSignalChan,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM)
	sig := <-osSignalChan
	ch.Info <- fmt.Sprintf("got signal %q", sig.String())
	ch.Info <- "Server stoped"

	// если сигнал получен, отменяем контекст работы
	cancel()
}
