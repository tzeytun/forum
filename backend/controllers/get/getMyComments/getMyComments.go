package getmycomments

import (
	"encoding/json"
	"net/http"

	"forum/backend/auth"
	"forum/backend/controllers/structs"
	"forum/backend/database"
)

func GetMyComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
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
		http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
		return
	}

	rows, err := db.Query("SELECT ID, PostId, UserId, Comment, UserName, LikeCount FROM COMMENTS WHERE UserId = ? ORDER BY created_at desc", userId)
	if err != nil {
		http.Error(w, "ERROR: Query error", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.UserId, &comment.Comment, &comment.UserName, &comment.LikeCount)
		if err != nil {
			http.Error(w, "ERROR: Database scan error", http.StatusBadRequest)
			return
		}
		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "ERROR: Row iteration error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, "ERROR: Failed to encode posts to JSON", http.StatusInternalServerError)
		return
	}
}
