package deleteaccount

import (
	"fmt"
	"net/http"

	"forum/backend/auth"
	"forum/backend/controllers/login"
	"forum/backend/database"
)

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	db, errDb := database.OpenDb(w)
	if errDb != nil {
		http.Error(w, "ERROR: Database cannot open", http.StatusBadRequest)
		return
	}
	defer db.Close()

	authenticated, userId, _ := auth.IsAuthenticated(r, db)
	if !authenticated {
		http.Error(w, "ERROR: You are not authorized to delete account", http.StatusUnauthorized)
		return
	}

	password := r.FormValue("password")

	var storedPassword string
	err := db.QueryRow("SELECT Password FROM USERS WHERE ID = ?", userId).Scan(&storedPassword)
	if err != nil {
		http.Error(w, "ERROR: Invalid query", http.StatusBadRequest)
		return
	}

	errComparePasswd := login.ArePasswordsMatching(storedPassword, password)
	if errComparePasswd != nil {
		http.Error(w, "ERROR: Invalid password", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM USERLIKES WHERE UserID = ?", userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM COMMENTS WHERE UserID = ?", userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Kullanıcının oluşturduğu tüm gönderileri sil
	_, err = db.Exec("DELETE FROM POSTS WHERE UserID = ?", userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Kullanıcıyı veritabanından sil
	_, err = db.Exec("DELETE FROM USERS WHERE ID = ?", userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM CATEGORIES WHERE USERID = ?", userId)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	auth.RemoveCookie(w, r, "session_token")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User successfully deleted")
}
