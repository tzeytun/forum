package loginpage

import (
	"net/http"

	"forum/backend/requests"
)

const loginApiUrl = "http://localhost:8080/api/login"

func LoginPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "frontend/pages/loginPage/loginPage.html")
	case "POST":
		email := r.FormValue("email")
		password := r.FormValue("password")

		err := requests.LoginRequest(loginApiUrl, email, password, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
