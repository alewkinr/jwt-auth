package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	statsd "github.com/etsy/statsd/examples/go"
	"github.com/evalphobia/logrus_sentry"
	"github.com/getsentry/raven-go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	// TODO move into config
	defaultDebugPort       = ":8010"
	defaultPublicPort      = ":8080"
	defaultShutdownTimeout = 30 * time.Second
)

// App описывает структуру приложения
type App struct {
	name         string
	sentryClient *raven.Client
	statsdClient *statsd.StatsdClient
	publicServer *publicServer
	debugServer  *debugServer
}

// New создает новое приложение
func New(handler http.Handler) *App {
	log.SetFormatter(&log.JSONFormatter{})

	a := &App{
		name:         os.Getenv("APP_NAME"),
		sentryClient: initSentryClient(),
		statsdClient: initStatsD(),
	}

	a.debugServer = a.newDebugServer()
	a.publicServer = a.newPublicServer(handler)

	return a
}

// Run запускает приложение
func (a *App) Run() error {
	sigs := make(chan os.Signal, 1)
	errCh := make(chan error)
	go func() {
		<-sigs
		errCh <- nil
	}()

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Info("Starting debug server")
		err := a.debugServer.Run()
		if err != http.ErrServerClosed {
			errCh <- errors.Wrap(err, "debug server error")
		}
	}()
	go func() {
		log.Info("Starting public server")
		err := a.publicServer.Run()
		if err != http.ErrServerClosed {
			errCh <- errors.Wrap(err, "public server error")
		}
	}()

	err := <-errCh

	ctx, cancel := context.WithTimeout(context.Background(), defaultShutdownTimeout)
	defer cancel()

	stopErr := a.publicServer.Stop(ctx)
	if stopErr != nil {
		log.Printf("error stopping public server: %v", stopErr)
	}

	stopErr = a.debugServer.Stop(ctx)
	if stopErr != nil {
		log.Printf("error stopping debug server: %v", stopErr)
	}

	return err
}

func initStatsD() *statsd.StatsdClient {
	host := os.Getenv("STATSD_HOST")
	// nolint:errcheck
	port, _ := strconv.Atoi(os.Getenv("STATSD_PORT"))
	if host == "" || port == 0 {
		panic("set proper STATSD_HOST and STATSD_PORT values")
	}
	return statsd.New(host, port)
}

func initSentryClient() *raven.Client {
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		panic("set proper SENTRY_DSN value")
	}
	sentryClient, err := raven.New(dsn)
	if err != nil {
		panic("failed to initialize Sentry client")
	}

	sentryClient.SetEnvironment(os.Getenv("DEPLOY_ENV"))
	hook, err := logrus_sentry.NewWithClientSentryHook(sentryClient, []log.Level{
		log.PanicLevel,
		log.FatalLevel,
		log.ErrorLevel,
		log.WarnLevel,
	})
	if err == nil {
		hook.Timeout = 10 * time.Second
		hook.StacktraceConfiguration.Enable = true
		log.AddHook(hook)
	}

	return sentryClient
}
