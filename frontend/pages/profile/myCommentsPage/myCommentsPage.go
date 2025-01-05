package mycommentspage

import (
	"html/template"
	"net/http"

	"forum/backend/requests"
)

func MyCommentsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
		return
	}

	data, errReq := requests.GetCommentDataForServeWithReq("http://localhost:8080/api/mycomments", cookie.Value)
	if errReq != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("frontend/pages/profile/myCommentsPage/myCommentsPage.html")
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

func DeleteMyComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	commentId := r.FormValue("id")

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
		return
	}

	err := requests.DeleteCommentRequest("http://localhost:8080/api/deletecomment", commentId, cookie.Value)
	if err != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/mycomments", http.StatusSeeOther)
}
