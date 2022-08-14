package model

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	TIMEOUT = 5 * time.Second
)

type UserStore struct {
	DB *pgxpool.Pool
}

type User struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password []byte    `json:"-"`
}

func (u *User) MatchedPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
	return err == nil
}

func (u *UserStore) InsertUser(ctx context.Context, email, name, password string) error {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()

	query := `INSERT INTO users(email, name, password) 
			  VALUES($1, $2, $3)`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = u.DB.Query(ctx, query, email, name, string(hashedPassword))

	if err != nil {
		println("cant query:", err.Error())
		return err
	}

	return nil
}

func (u *UserStore) GetUserByEmail(ctx context.Context, email string) (User, error) {
	ctx, cancel := context.WithTimeout(ctx, TIMEOUT)
	defer cancel()

	query := `SELECT id, name, email, password
			  FROM users
			  WHERE email = $1;`

	var user User

	err := u.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}
