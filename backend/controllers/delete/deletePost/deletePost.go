package deletepost

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/backend/auth"
	"forum/backend/database"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postId := r.FormValue("id")
	postIdInt, atoiErr := strconv.Atoi(postId)
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

	var postUserID int
	var postUserName string

	err := db.QueryRow("SELECT UserID, UserName FROM POSTS WHERE ID = ?", postIdInt).Scan(&postUserID, &postUserName)
	if err != nil {
		http.Error(w, "ERROR: Cannot find post from database", http.StatusBadRequest)
		return
	}

	if AreTheCredentialsMatch(postUserID, userId, postUserName, userName) {
		_, errDel := db.Exec(`DELETE FROM POSTS WHERE ID = ?`, postIdInt)
		if errDel != nil {
			http.Error(w, "ERROR: Unable to delete post", http.StatusInternalServerError)
			return
		}
		_, errComDel := db.Exec(`DELETE FROM COMMENTS WHERE PostId = ?`, postIdInt)
		if errComDel != nil {
			http.Error(w, "ERROR: Unable to delete post comments", http.StatusInternalServerError)
			return
		}
		_, errUserLikeDel := db.Exec(`DELETE FROM USERLIKES WHERE DeleteID = ?`, postIdInt)
		if errUserLikeDel != nil {
			http.Error(w, "ERROR: Unable to delete user likes", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Post successfully deleted")
}

func AreTheCredentialsMatch(DbUserID int, userID int, DbUserName string, userName string) bool {
	return DbUserID == userID && DbUserName == userName
}
