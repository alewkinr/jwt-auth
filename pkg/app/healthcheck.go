package app

type healthCheckFunc func() error

// AddHealthCheck добавляет healthcheck
func (a *App) AddHealthCheck(f func() error) {
	a.debugServer.addHealthCheck(f)
}
