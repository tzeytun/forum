package deletecomment

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/backend/auth"
	deletepost "forum/backend/controllers/delete/deletePost"
	"forum/backend/database"
)

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	commentId := r.FormValue("id")
	commentIdInt, atoiErr := strconv.Atoi(commentId)
	if atoiErr != nil {
		http.Error(w, "ERROR: Invalid post ID format", http.StatusBadRequest)
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
		http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
		return
	}

	var comUserId int
	var comUserName string

	err := db.QueryRow("SELECT UserId, UserName FROM COMMENTS WHERE ID = ?", commentIdInt).Scan(&comUserId, &comUserName)
	if err != nil {
		http.Error(w, "ERROR: Comment cannot found in the database", http.StatusBadRequest)
		return
	}

	if deletepost.AreTheCredentialsMatch(comUserId, userId, comUserName, userName) {
		_, errComDel := db.Exec(`DELETE FROM COMMENTS WHERE ID = ?`, commentIdInt)
		if errComDel != nil {
			http.Error(w, "ERROR: Unable to delete comment", http.StatusInternalServerError)
			return
		}
		_, errUserLikeDel := db.Exec(`DELETE FROM USERLIKES WHERE PostID = ? AND IsComment = 1`, commentIdInt)
		if errUserLikeDel != nil {
			http.Error(w, "ERROR: Unable to delete user likes", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Comment successfully deleted")
}
