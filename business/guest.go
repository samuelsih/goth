package business

import (
	"context"
	"net/http"

	"github.com/samuelsih/goth/model"
)

const (
	statusErrInternal = http.StatusInternalServerError
	statusBadReq      = http.StatusBadRequest
	statusNotFound    = http.StatusNotFound
)

var (
	errInternal = "internal server error, please try again in a minutes"
	errPassword = "username or password doesn't match"
)

type GuestDeps struct {
	Conn model.UserStore
	Sess model.UserSessionStore
}

type LoginIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginOut struct {
	User model.UserSession `json:"user"`
	CommonResponse
}

func (g *GuestDeps) GuestLogin(ctx context.Context, in *LoginIn) (string, LoginOut) {
	var out LoginOut

	user, err := g.Conn.GetUserByEmail(ctx, in.Email)
	if err != nil {
		out.SetError(statusErrInternal, errInternal)
		return "", out
	}

	if !user.MatchedPassword(in.Password) {
		out.SetError(statusBadReq, errPassword)
		return "", out
	}

	sessionID, err := g.Sess.Save(ctx, model.UserSession{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	})

	if err != nil {
		out.SetError(statusErrInternal, err.Error())
		return "", out
	}

	out.User.ID = user.ID
	out.User.Email = user.Email
	out.User.Name = user.Name

	out.SetOK()

	return sessionID, out
}

type RegisterIn struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RegisterOut struct {
	CommonResponse
}

func (g *GuestDeps) GuestRegister(ctx context.Context, in *RegisterIn) RegisterOut {
	var out RegisterOut

	err := ValidateSignIn(*in)
	if err != nil {
		out.SetError(statusBadReq, err.Error())
		return out
	}

	err = g.Conn.InsertUser(ctx, in.Email, in.Name, in.Password)
	if err != nil {
		out.SetError(statusErrInternal, err.Error())
		return out
	}

	out.SetOK()

	return out
}

type LogoutIn struct {
	CommonRequest
}

type LogoutOut struct {
	CommonResponse
}

func (g *GuestDeps) GuestLogout(ctx context.Context, sessionID string) LogoutOut {
	var out LogoutOut

	err := g.Sess.Delete(ctx, sessionID)
	if err != nil {
		out.SetError(statusBadReq, err.Error())
		return out
	}

	out.SetOK()
	return out
}
