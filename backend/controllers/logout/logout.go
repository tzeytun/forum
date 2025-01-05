package logout

import (
	"fmt"
	"net/http"

	"forum/backend/auth"
	"forum/backend/database"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "ERROR: You are not logged in", http.StatusBadRequest)
		return
	}

	db, errDb := database.OpenDb(w)
	if errDb != nil {
		http.Error(w, "ERROR: Database cannot open", http.StatusBadRequest)
		return
	}
	defer db.Close()

	_, err = db.Exec("UPDATE USERS SET session_token = '' WHERE session_token = ?", cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	auth.RemoveCookie(w, r, cookie.Name)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User successfully logged out")
}
