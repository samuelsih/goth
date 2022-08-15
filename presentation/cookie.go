package presentation

import (
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/securecookie"
	_ "github.com/joho/godotenv/autoload"
)

var (
	securer = securecookie.New([]byte(os.Getenv("HASH_KEY")), nil)
	mutex   sync.RWMutex
)

func SetCookie(w http.ResponseWriter, value string) error {
	mutex.Lock()
	defer mutex.Unlock()

	encoded, err := securer.Encode("app_session", value)

	if err != nil {
		JSON(w, http.StatusInternalServerError, Map{
			"code":    http.StatusInternalServerError,
			"message": http.StatusText(http.StatusInternalServerError),
		})

		fmt.Println("ENCODED FAILED:", err.Error())
		return err
	}

	cookie := http.Cookie{
		Name:     "app_session",
		Value:    encoded,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		MaxAge:   24 * 3600,
	}

	http.SetCookie(w, &cookie)
	return nil
}

func ReadCookie(w http.ResponseWriter, r *http.Request) (string, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	cookie, err := r.Cookie("app_session")
	if err != nil {
		JSON(w, http.StatusBadRequest, Map{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})

		return "", err
	}

	var sessionID string

	err = securer.Decode("app_session", cookie.Value, &sessionID)
	if err != nil {
		JSON(w, http.StatusBadRequest, Map{
			"code":    http.StatusBadRequest,
			"message": http.StatusText(http.StatusBadRequest),
		})

		return "", err
	}

	return sessionID, err
}
