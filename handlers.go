package main

import (
	"net/http"

	"github.com/samuelsih/goth/business"
	"github.com/samuelsih/goth/presentation"
)

func Register(deps *business.GuestDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := business.RegisterIn{}

		if !presentation.ReadInput(w, r, &in, nil) {
			return
		}

		out := deps.GuestRegister(r.Context(), &in)

		presentation.WriteOutput(w, out, &out.CommonResponse)
	}
}

func Login(deps *business.GuestDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := business.LoginIn{}

		if !presentation.ReadInput(w, r, &in, nil) {
			return
		}

		cookieID, out := deps.GuestLogin(r.Context(), &in)

		presentation.WriteOutputWithCookie(w, cookieID, out, &out.CommonResponse)
	}
}

func Logout(deps *business.GuestDeps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		in := business.LogoutIn{}

		if !presentation.ReadInput(w, r, &in, nil) {
			return
		}

		cookie, err := r.Cookie("app_session")
		if err != nil {
			presentation.JSON(w, http.StatusBadRequest, presentation.Map{
				"code":    http.StatusBadRequest,
				"message": err,
			})

			return
		}

		out := deps.GuestLogout(r.Context(), cookie.Value)

		presentation.WriteOutputWithCookie(w, "", out, &out.CommonResponse)
	}
}

func Pong() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("You are authenticated"))
	}
}

func Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello!!!"))
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	presentation.JSON(w, http.StatusMethodNotAllowed, presentation.Map{
		"code":    http.StatusMethodNotAllowed,
		"message": http.StatusText(http.StatusMethodNotAllowed),
	})
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	presentation.JSON(w, http.StatusNotFound, presentation.Map{
		"code":    http.StatusNotFound,
		"message": http.StatusText(http.StatusNotFound),
	})
}
