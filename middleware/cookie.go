package middleware

import (
	"net/http"

	"github.com/samuelsih/goth/presentation"
)

func CookieExists() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, err := r.Cookie("app_session")
			if err != nil {
				presentation.JSON(
					w,
					http.StatusUnauthorized,
					presentation.Map{
						"code":    http.StatusUnauthorized,
						"message": "you must login first",
					},
				)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func CookieNotExists() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, err := r.Cookie("app_session")
			if err == nil {
				presentation.JSON(
					w,
					http.StatusNotAcceptable,
					presentation.Map{
						"code":    http.StatusNotAcceptable,
						"message": "you're already login",
					},
				)
				return
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
