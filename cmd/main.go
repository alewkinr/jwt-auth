package main

import (
	"net/http"
	"os"
	"time"

	"example.com/back/auth/internal/random"

	"example.com/back/auth/internal/permissions"

	val10 "github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pressly/goose"
	notification "gitlab.example.com/example.com/client/notification.git"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"example.com/back/auth/internal/config"
	"example.com/back/auth/internal/handlers"
	"example.com/back/auth/internal/session"
	"example.com/back/auth/internal/token"
	"example.com/back/auth/internal/user"
	"example.com/back/auth/pkg/app"
)

func main() {
	cfg := config.MustInitConfig()

	db := sqlx.MustConnect("pgx", cfg.DatabaseDSN)
	defer func() {
		dbErr := db.Close()
		if dbErr != nil {
			log.Errorf("db close error: %v ", dbErr)
		}
	}()

	if os.Getenv("DEPLOY_ENV") != "production" {
		err := goose.Up(db.DB, "./migrations")
		if err != nil {
			log.Fatal(err)
		}
	}
	validator := val10.New()
	t := token.NewManager(cfg.AccessTokenKey, cfg.RefreshTokenKey)
	um := user.NewManager(db)
	sm := session.NewManager(db)
	notificationHTTP := notification.NewClient(cfg.NotificationBaseURL, notification.WithTimeout(7*time.Second))
	rand := random.New()

	router := mux.NewRouter()
	router.Handle("/v1/phone_sessions",
		handlers.NewPhoneSessionsPOSTHandler(sm, sm, validator, notificationHTTP)).
		Methods("POST")
	router.Handle("/v2/sessions/phone_session",
		handlers.NewPostV2PhoneSessionsHandler(um, sm, sm, notificationHTTP, validator)).Methods("POST")

	router.Handle("/v1/phone_sessions/{sessionId}/sms_code",
		handlers.NewPhoneSessionsPUTHandler(sm, sm, um, t, validator)).
		Methods("PUT")

	router.Handle("/v1/login", handlers.NewLoginHandler(um, t))
	router.Handle("/v1/logout", handlers.NewLogoutHandler(um, t))
	router.Handle("/v1/refresh", handlers.NewRefreshHandler(um, t))
	router.Handle("/v1/validate", handlers.NewValidateHandler(t))
	router.Handle("/v1/passwords/password_reset", handlers.NewPasswordResetHandler(um, notificationHTTP, validator, rand)).Methods("POST")
	router.Handle("/v1/users/{userID:[0-9]+}/password", handlers.NewPasswordChangeHandler(um, validator)).Methods("PATCH")

	router.Handle("/v1/users", handlers.NewCreateHandler(um, t, validator, sm, sm)).
		Methods("POST")
	router.Handle("/v1/users", handlers.NewQueryHandler(um)).
		Methods("GET")

	router.Handle("/v1/users/{userID:[0-9]+}", handlers.NewGetUserHandler(um,
		permissions.NewRoleOrOwnerChecker([]permissions.Role{
			permissions.RoleAdmin,
			permissions.RoleManager,
		}),
	)).
		Methods("GET")

	a := app.New(corsMiddleware(router))
	a.AddHealthCheck(db.Ping)

	if err := a.Run(); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions || r.Header.Get("X-Forwarded-Method") == http.MethodOptions {
			return
		}
		next.ServeHTTP(rw, r)
	})
}
