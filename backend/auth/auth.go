package auth

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gofrs/uuid"
)

func CreateSessionToken() (uuid.UUID, error) {
	sessionToken, err := uuid.NewV4()
	if err != nil {
		return sessionToken, err
	}
	return sessionToken, nil
}

func SetCookie(w http.ResponseWriter, sessionToken uuid.UUID) {
	cookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.String(),
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func SetTokenInDatabase(w http.ResponseWriter, db *sql.DB, sessionToken uuid.UUID, userID int) error {
	_, err := db.Exec("UPDATE USERS SET session_token = ? WHERE ID = ?", sessionToken.String(), userID)
	if err != nil {
		return err
	}
	return nil
}

func RemoveCookie(w http.ResponseWriter, r *http.Request, cookieName string) {
	expiredCookie := &http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/", // Çerezin ayarlandığı yola dikkat edin
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
	}

	http.SetCookie(w, expiredCookie)
}

func IsAuthenticated(r *http.Request, db *sql.DB) (bool, int, string) {
	cookie, errCookie := r.Cookie("session_token")
	if errCookie != nil {
		return false, 0, ""
	}

	var userID int
	var userName string

	err := db.QueryRow("SELECT ID, UserName FROM USERS WHERE session_token = ?", cookie.Value).Scan(&userID, &userName)
	if err != nil {
		fmt.Println(err)
		return false, 0, ""
	}

	return true, userID, userName
}
