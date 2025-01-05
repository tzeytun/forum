package registerpage

import (
	"net/http"

	"forum/backend/requests"
)

const registerApiUrl = "http://localhost:8080/api/register"

func RegisterPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "frontend/pages/registerPage/registerPage.html")
	case "POST":
		email := r.FormValue("email")
		userName := r.FormValue("username")
		password := r.FormValue("password")

		err := requests.RegisterRequest(registerApiUrl, email, userName, password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "http://localhost:8080/login", http.StatusSeeOther)
	}
}
