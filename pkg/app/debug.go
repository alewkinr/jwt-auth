package app

import (
	"context"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type debugServer struct {
	s   *http.Server
	hcs []healthCheckFunc
}

func (a *App) newDebugServer() *debugServer {
	d := &debugServer{}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/live", http.HandlerFunc(livenessHandler))
	mux.Handle("/ready", http.HandlerFunc(d.readinessHandler))
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))

	d.s = &http.Server{
		Addr:    defaultDebugPort,
		Handler: mux,
	}

	return d
}

// Run запускает сервер
func (ds *debugServer) Run() error {
	return ds.s.ListenAndServe()
}

// Stop останавливает сервер
func (ds *debugServer) Stop(ctx context.Context) error {
	log.Info("Stopping debug server")
	return ds.s.Shutdown(ctx)
}

func (ds *debugServer) addHealthCheck(f func() error) {
	ds.hcs = append(ds.hcs, f)
}

func livenessHandler(w http.ResponseWriter, r *http.Request) {
	// nolint:errcheck
	w.Write([]byte(`{"ok": true}`))
}

func (ds *debugServer) readinessHandler(w http.ResponseWriter, r *http.Request) {
	for _, hc := range ds.hcs {
		err := hc()
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
