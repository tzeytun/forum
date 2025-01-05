package createcomment

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/backend/auth"
	"forum/backend/database"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postId := r.FormValue("id")
	comment := r.FormValue("comment")

	if !IsForumCommentValid(postId, comment) {
		http.Error(w, "ERROR: Id or comment cannot empty", http.StatusBadRequest)
		return
	}

	postIdInt, atoiErr := strconv.Atoi(postId)
	if atoiErr != nil {
		http.Error(w, "ERROR: Invalid ID format", http.StatusBadRequest)
		return
	}

	db, errDb := database.OpenDb(w)
	if errDb != nil {
		http.Error(w, "ERROR: Database cannot open", http.StatusBadRequest)
		return
	}
	defer db.Close()

	authenticated, userId, userName := auth.IsAuthenticated(r, db)
	if !authenticated {
		http.Error(w, "ERROR: You are not authorized to create comment", http.StatusUnauthorized)
		return
	}
	var secondUserId int
	err := db.QueryRow("SELECT UserId FROM POSTS WHERE ID = ?", postIdInt).Scan(&secondUserId)
	if err != nil {
		http.Error(w, "ERROR: Invalid post ID", http.StatusBadRequest)
		return
	}

	_, errEx := db.Exec(`INSERT INTO COMMENTS (PostID, UserId, UserName, Comment) VALUES (?, ?, ?, ?)`, postIdInt, userId, userName, comment)
	if errEx != nil {
		http.Error(w, "ERROR: Post did not add to the database", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Comment successfully created")
}

func IsForumCommentValid(id, comment string) bool {
	return id != "" && comment != ""
}
