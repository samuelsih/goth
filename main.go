package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/samuelsih/goth/business"
	"github.com/samuelsih/goth/db"
	midd "github.com/samuelsih/goth/middleware"
	"github.com/samuelsih/goth/model"
)

func main() {
	r := chi.NewRouter()

	r.Use(
		httprate.LimitByIP(60, 1*time.Minute),
		middleware.Logger,
		middleware.Recoverer,
		middleware.GetHead,
		midd.CORS(),
	)

	pgDB := db.NewPostgres()
	redisDB := db.NewSessionRedis()

	userRepo := model.UserStore{DB: pgDB}
	userSessionRepo := model.UserSessionStore{Conn: redisDB}
	guest := business.GuestDeps{
		Conn: userRepo,
		Sess: userSessionRepo,
	}

	//Public
	r.Group(func(r chi.Router) {
		r.Get("/", Root)
	})

	// Auth
	r.Group(func(r chi.Router) {
		// Make sure they didnt have cookie
		r.Group(func(r chi.Router) {
			r.Use(midd.CookieNotExists())
			r.Post("/register", Register(&guest))
			r.Post("/login", Login(&guest))
		})

		r.Post("/logout", Logout(&guest))
	})

	//Private
	r.Group(func(r chi.Router) {
		r.Use(midd.CookieExists())
		r.Get("/pong", Pong())
	})

	// Custom Error On Some Method
	r.MethodNotAllowed(MethodNotAllowed)
	r.NotFound(NotFound)

	http.ListenAndServeTLS(":80", "example.crt", "example.key", r)
}
