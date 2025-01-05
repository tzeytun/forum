package deleteaccountpage

import (
	"fmt"
	"net/http"

	"forum/backend/requests"
)

func DeleteAccountPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "frontend/pages/deleteAccountPage/deleteAccountPage.html")
	case "POST":
		password := r.FormValue("password")

		cookie, cookieErr := r.Cookie("session_token")
		if cookieErr != nil {
			http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
			return
		}

		err := requests.DeleteAccountRequest("http://localhost:8080/api/deleteaccount", password, cookie.Value)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			fmt.Println(err)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
