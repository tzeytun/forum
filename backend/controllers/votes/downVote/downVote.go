package downvote

import (
	"fmt"
	"net/http"
	"strconv"

	"forum/backend/auth"
	"forum/backend/database"
)

func DownVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	postOrCommentId := r.FormValue("id")
	isComment := r.FormValue("isComment") == "true"
	postID := r.FormValue("post_id")
	idInt, err := strconv.Atoi(postOrCommentId)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	postIdInt, err := strconv.Atoi(postID)
	if err != nil {
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	db, errDb := database.OpenDb(w)
	if errDb != nil {
		http.Error(w, "ERROR: Database cannot open", http.StatusBadRequest)
		return
	}
	defer db.Close()

	authenticated, userId, _ := auth.IsAuthenticated(r, db)
	if !authenticated {
		http.Error(w, "ERROR: You are not authorized to up vote", http.StatusUnauthorized)
		return
	}

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM UserLikes WHERE UserID = ? AND PostID = ? AND IsComment = ?)", userId, idInt, isComment).Scan(&exists)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	if !exists {
		_, err = db.Exec(`INSERT INTO UserLikes (UserID, PostID, IsComment, DeleteID, Liked, Disliked) VALUES (?, ?, ?, ?, 0, 0)`, userId, idInt, isComment, postIdInt)
		if err != nil {
			http.Error(w, "Failed to create user like record", http.StatusInternalServerError)
			return
		}
	}

	var dislikeExists bool
	var likeExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM UserLikes WHERE UserID = ? AND PostID = ? AND IsComment = ? AND Disliked = 1)", userId, idInt, isComment).Scan(&dislikeExists)
	if err != nil {
		http.Error(w, "ERROR: Invalid query", http.StatusBadRequest)
		return
	}
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM UserLikes WHERE UserID = ? AND PostID = ? AND IsComment = ? AND Liked = 1)", userId, idInt, isComment).Scan(&likeExists)
	if err != nil {
		http.Error(w, "ERROR: Invalid query", http.StatusBadRequest)
		return
	}

	if likeExists {
		if isComment {
			_, err = db.Exec(`UPDATE COMMENTS SET LikeCount = LikeCount - 2 WHERE ID = ?`, idInt)
		} else {
			_, err = db.Exec(`UPDATE POSTS SET LikeCount = LikeCount - 2 WHERE ID = ?`, idInt)
		}

		if err == nil {
			_, err = db.Exec(`UPDATE UserLikes SET Disliked = 1 WHERE UserID = ? AND PostID = ? AND IsComment = ?`, userId, idInt, isComment)
			if err != nil {
				http.Error(w, "ERROR: Database update error", http.StatusBadRequest)
			}
			_, err = db.Exec(`UPDATE UserLikes SET Liked = 0 WHERE UserID = ? AND PostID = ? AND IsComment = ?`, userId, idInt, isComment)
			if err != nil {
				http.Error(w, "ERROR: Database update error", http.StatusBadRequest)
			}
		}
	} else if dislikeExists {
		if isComment {
			_, err = db.Exec(`UPDATE COMMENTS SET LikeCount = LikeCount + 1 WHERE ID = ?`, idInt)
		} else {
			_, err = db.Exec(`UPDATE POSTS SET LikeCount = LikeCount + 1 WHERE ID = ?`, idInt)
		}

		if err == nil {
			_, err = db.Exec(`UPDATE UserLikes SET Disliked = 0 WHERE UserID = ? AND PostID = ? AND IsComment = ?`, userId, idInt, isComment)
			if err != nil {
				http.Error(w, "ERROR: Database update error", http.StatusBadRequest)
			}
		}

	} else {

		if isComment {
			_, err = db.Exec(`UPDATE COMMENTS SET LikeCount = LikeCount - 1 WHERE ID = ?`, idInt)
		} else {
			_, err = db.Exec(`UPDATE POSTS SET LikeCount = LikeCount - 1 WHERE ID = ?`, idInt)
		}

		if err == nil {
			_, err = db.Exec(`UPDATE UserLikes SET Disliked = 1 WHERE UserID = ? AND PostID = ? AND IsComment = ?`, userId, idInt, isComment)
			if err != nil {
				http.Error(w, "ERROR: Database update error", http.StatusBadRequest)
			}
		}
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User successfully down vote")
}
