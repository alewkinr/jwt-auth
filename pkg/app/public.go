package app

import (
	"context"
	"net/http"

	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"

	"example.com/back/auth/pkg/app/middleware"
)

type publicServer struct {
	s *http.Server
}

func (a *App) newPublicServer(handler http.Handler) *publicServer {
	handler = middleware.Log(handler, a.name, a.statsdClient)
	handler = raven.Recoverer(handler)

	s := &http.Server{
		Addr:    defaultPublicPort,
		Handler: handler,
	}

	return &publicServer{
		s: s,
	}
}

// Run запускает сервер
func (ds *publicServer) Run() error {
	return ds.s.ListenAndServe()
}

// Stop останавливает сервер
func (ds *publicServer) Stop(ctx context.Context) error {
	log.Info("Stopping public server")
	return ds.s.Shutdown(ctx)
}
