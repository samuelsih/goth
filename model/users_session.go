package model

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/rs/xid"
)

var (
	guid = xid.New()

	day = 24 * 60 * 3600
)

type UserSession struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

type UserSessionStore struct {
	Conn *redis.Client
}

func (u *UserSessionStore) Save(ctx context.Context, user UserSession) (string, error) {
	buf, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	value := string(buf)

	sessionID := guid.String()
	if err != nil {
		return "", err
	}

	cmd := u.Conn.Set(ctx, sessionID, value, time.Duration(day)*time.Second)

	return sessionID, cmd.Err()
}

func (u *UserSessionStore) Get(ctx context.Context, key string) (UserSession, error) {
	var user UserSession

	userFromDB, err := u.Conn.Get(ctx, key).Result()
	if err != nil {
		return user, err
	}

	r := bytes.NewReader([]byte(userFromDB))

	err = json.NewDecoder(r).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, err
}

func (u *UserSessionStore) Delete(ctx context.Context, key ...string) error {
	cmd := u.Conn.Del(ctx, key...)

	return cmd.Err()
}
