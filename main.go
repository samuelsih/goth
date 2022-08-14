package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
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

	server := &http.Server{
		Addr:         ":80",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func ()  {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

		<- quit

		closedServerChan := make(chan bool, 1)
		closedPgxChan := make(chan bool, 1)
		closedRedisChan := make(chan bool, 1)

		go shutdownServer(closedServerChan, ctx, server)
		go closePgx(closedPgxChan, pgDB)
		go closeRedis(closedRedisChan, redisDB)

		<- closedPgxChan
		<- closedRedisChan
		<- closedServerChan

		if d := ctx.Done(); d != nil {
			log.Println("Context Done Hit!")
		}

		close(closedPgxChan)
		close(closedRedisChan)
		close(closedServerChan)
	}()

	server.ListenAndServeTLS("example.crt", "example.key")
}

func shutdownServer(done chan bool, ctx context.Context, server *http.Server) {
	log.Println("Shutdown server")

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
	}

	done <- true
}

func closePgx(done chan bool, conn *pgxpool.Pool) {
	log.Println("Closing postgres")
	
	conn.Close()

	done <- true
}

func closeRedis(done chan bool, conn *redis.Client) {
	log.Println("Closing redis")

	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}

	done <- true
}