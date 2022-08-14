package presentation

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/samuelsih/goth/business"
)

type Map map[string]any

func ReadInput[T any](w http.ResponseWriter, r *http.Request, in *T, cr *business.CommonRequest) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	err := decoder.Decode(&in)

	if err != nil {
		JSON(w, http.StatusBadRequest, Map{
			"error":   http.StatusBadRequest,
			"message": err,
		})
		return false
	}

	return true
}

func WriteOutput[T any](w http.ResponseWriter, out T, cr *business.CommonResponse) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(cr.Code)
	json.NewEncoder(w).Encode(out)
}

func WriteOutputWithCookie[T any](w http.ResponseWriter, sessionID string, out T, cr *business.CommonResponse) {
	if sessionID == "" {
		cookie := http.Cookie{
			Name:     "app_session",
			Value:    "",
			Secure:   true,
			MaxAge:   -1,
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}

		http.SetCookie(w, &cookie)
		WriteOutput(w, out, cr)
		return
	}

	if err := SetCookie(w, sessionID); err != nil {
		return
	}

	WriteOutput(w, out, cr)
}

func JSON(w http.ResponseWriter, status int, output Map) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(output)
}
