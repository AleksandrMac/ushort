package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (c *Controller) setMetricsControllers(r *chi.Mux) {
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
}
