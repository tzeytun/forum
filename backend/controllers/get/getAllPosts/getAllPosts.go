package getallposts

import (
	"encoding/json"
	"net/http"

	"forum/backend/controllers/structs"
	"forum/backend/database"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
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

	rows, err := db.Query("SELECT ID, UserID, UserName, Title, Content, LikeCount FROM POSTS ORDER BY PostDate desc")
	if err != nil {
		http.Error(w, "ERROR: Query execution failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []structs.Post

	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.UserID, &post.UserName, &post.Title, &post.Content, &post.LikeCount)
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
