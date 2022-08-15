package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
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

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

		<-quit

		var wg sync.WaitGroup
		wg.Add(3)

		go shutdownServer(&wg, ctx, server)
		go closePgx(&wg, pgDB)
		go closeRedis(&wg, redisDB)

		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
		default:
		}

		wg.Wait()
	}()

	server.ListenAndServeTLS("example.crt", "example.key")
}

func shutdownServer(wg *sync.WaitGroup, ctx context.Context, server *http.Server) {
	defer wg.Done()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown error:", err)
		return
	}

	log.Println("Shutdown server")
}

func closePgx(wg *sync.WaitGroup, conn *pgxpool.Pool) {
	defer wg.Done()

	conn.Close()

	log.Println("Postgres disconnected")
}

func closeRedis(wg *sync.WaitGroup, conn *redis.Client) {
	defer wg.Done()

	if err := conn.Close(); err != nil {
		log.Fatal(err)
	}

	log.Println("Redis disconnected")
}
