package register

import (
	"fmt"
	"net/http"

	"forum/backend/database"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	password := r.FormValue("password")

	if !AreRegisterCredentialsCorrect(email, username, password) {
		http.Error(w, "ERROR: Please provide the register credentials correct", http.StatusBadRequest)
		return
	}

	hashedPasswd, errHash := HashThePasswd(password)
	if errHash != nil {
		http.Error(w, "ERROR: Internal server error", http.StatusInternalServerError)
		return
	}

	db, errDb := database.OpenDb(w)
	if errDb != nil {
		http.Error(w, "ERROR: Database cannot open", http.StatusBadRequest)
		return
	}
	defer db.Close()

	var userID int
	errMail := db.QueryRow("SELECT ID FROM USERS WHERE Email = ?", email).Scan(&userID)
	if errMail == nil {
		http.Error(w, "ERROR: Email already taken", http.StatusBadRequest)
		return
	}

	errUsername := db.QueryRow("SELECT ID FROM USERS WHERE UserName = ?", username).Scan(&userID)
	if errUsername == nil {
		http.Error(w, "ERROR: Username already taken", http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT INTO USERS (Email, UserName, Password, Role) VALUES (?, ?, ?, ?)", email, username, hashedPasswd, "User")
	if err != nil {
		http.Error(w, "ERROR: Bad Request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User successfully created")
}

func AreRegisterCredentialsCorrect(email, username, password string) bool {
	return email != "" && username != "" && password != ""
}

func HashThePasswd(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
