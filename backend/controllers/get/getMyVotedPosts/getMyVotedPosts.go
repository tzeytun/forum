package getmyvotedposts

import (
	"encoding/json"
	"net/http"

	"forum/backend/auth"
	"forum/backend/controllers/structs"
	"forum/backend/database"
)

func GetMyVotedPosts(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query(`
        SELECT POSTS.ID, POSTS.Title, POSTS.Content, POSTS.UserName, POSTS.UserId, POSTS.LikeCount
        FROM POSTS
        INNER JOIN USERLIKES ON POSTS.ID = USERLIKES.PostID
        WHERE USERLIKES.UserID = ? AND USERLIKES.IsComment = 0 AND (USERLIKES.Liked = 1 OR USERLIKES.Disliked = 1)
        ORDER BY POSTS.PostDate DESC
    `, userId)
	if err != nil {
		http.Error(w, "ERROR: Query error", http.StatusBadRequest)
		return
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.UserName, &post.UserID, &post.LikeCount)
		if err != nil {
			http.Error(w, "ERROR: Database scan error", http.StatusBadRequest)
			return
		}
		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		http.Error(w, "ERROR: Row iteration error", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "ERROR: Failed to encode posts to JSON", http.StatusInternalServerError)
		return
	}
}
