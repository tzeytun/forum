package postpage

import (
	"html/template"
	"net/http"

	"forum/backend/requests"
)

func PostPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	postId := r.FormValue("id")

	data, err := requests.GetPostWithComments("http://localhost:8080/api/postandcomments", postId)
	if err != nil {
		http.Error(w, "ERROR: Cannot get post and comments", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("frontend/pages/postPage/postPage.html")
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

func PostPageCreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postId := r.FormValue("id")
	comment := r.FormValue("comment")

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to create comment", http.StatusUnauthorized)
		return
	}

	err := requests.CreateCommentRequest("http://localhost:8080/api/createcomment", postId, comment, cookie.Value)
	if err != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/post?id="+postId, http.StatusSeeOther)
}

func PostPageUpVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	isComment := r.FormValue("isComment")
	postId := r.FormValue("post_id")

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to up vote", http.StatusUnauthorized)
		return
	}

	err := requests.VoteRequest("http://localhost:8080/api/upvote", id, isComment, postId, cookie.Value)
	if err != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/post?id="+postId, http.StatusSeeOther)
}

func PostPageDownVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.FormValue("id")
	isComment := r.FormValue("isComment")
	postId := r.FormValue("post_id")

	cookie, cookieErr := r.Cookie("session_token")
	if cookieErr != nil {
		http.Error(w, "ERROR: You are not authorized to down vote", http.StatusUnauthorized)
		return
	}

	err := requests.VoteRequest("http://localhost:8080/api/downvote", id, isComment, postId, cookie.Value)
	if err != nil {
		http.Error(w, "ERROR: Bad request", http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/post?id="+postId, http.StatusSeeOther)
}
