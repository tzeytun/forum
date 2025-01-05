package mypostspage

import (
	"html/template"
	"net/http"

	"forum/backend/requests"
)

func MyPostsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
		return
	}

	data, errReq := requests.GetDataForServeWithReq("http://localhost:8080/api/myposts", cookie.Value)
	if errReq != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("frontend/pages/profile/myPostsPage/myPostsPage.html")
	if err != nil {
		http.Error(w, "ERROR: Unable to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "ERROR: Unable to execute template", http.StatusInternalServerError)
		return
	}
}

func DeleteMyPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postId := r.FormValue("id")

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
		return
	}

	err := requests.DeletePostRequest("http://localhost:8080/api/deletepost", postId, cookie.Value)
	if err != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/myposts", http.StatusSeeOther)
}
