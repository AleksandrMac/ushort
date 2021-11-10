package controller

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	"github.com/go-chi/jwtauth/v5"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog"
	jgConfig "github.com/uber/jaeger-client-go/config"

	"github.com/AleksandrMac/ushort/pkg/config"
	"github.com/AleksandrMac/ushort/pkg/config/env"
	"github.com/AleksandrMac/ushort/pkg/model"
)

type Controller struct {
	DB           model.CRUD
	dbCloser     io.Closer
	Config       *config.Config
	TokenAuth    *jwtauth.JWTAuth
	Logger       *zerolog.Logger
	Tracer       opentracing.Tracer
	TracerCloser io.Closer
	Ctx          context.Context
	CtxCancel    context.CancelFunc
	Info         chan string
	Debug        chan string
	Err          chan error
	Warn         chan error
	Critical     chan error
}

func New() (*Controller, error) {
	cfg, err := env.New()
	if err != nil {
		return nil, err
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

	db, err := model.NewDB(cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	tracer, tCloser, err := initJaeger("ushort", &logger)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Controller{
		DB:           db,
		dbCloser:     db,
		Config:       cfg,
		TokenAuth:    jwtauth.New("HS256", cfg.Auth.PrivateKey, nil),
		Logger:       &logger,
		Tracer:       tracer,
		TracerCloser: tCloser,
		Ctx:          ctx,
		CtxCancel:    cancel,
		Info:         make(chan string),
		Debug:        make(chan string),
		Err:          make(chan error),
		Warn:         make(chan error),
		Critical:     make(chan error),
	}, nil
}

func (c *Controller) Close() {
	c.dbCloser.Close()
	c.TracerCloser.Close()
	c.CtxCancel()
}

type zeroLogWrapper struct {
	logger *zerolog.Logger
}

func (z *zeroLogWrapper) Error(msg string) {
	z.logger.Err(errors.New(msg))
}

func (z *zeroLogWrapper) Infof(msg string, args ...interface{}) {
	z.logger.Info().Msgf(msg, args...)
}

func initJaeger(service string, logger *zerolog.Logger) (opentracing.Tracer, io.Closer, error) {
	cfg, err := jgConfig.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return nil, nil, err
	}

	cfg.ServiceName = service
	cfg.Sampler = &jgConfig.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	cfg.Reporter = &jgConfig.ReporterConfig{
		LogSpans: true,
	}

	tracer, closer, err := cfg.NewTracer(jgConfig.Logger(&zeroLogWrapper{logger: logger}))
	if err != nil {
		return nil, nil, fmt.Errorf("cannot init jaeger: %w", err)
	}

	return tracer, closer, err
}

func (c *Controller) SetControllers(r *chi.Mux) {
	c.setAuthControllers(r)
	c.setURLControllers(r)
	c.setMetricsControllers(r)
}
