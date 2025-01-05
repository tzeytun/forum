package getsearchedposts

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"forum/backend/controllers/structs"
	"forum/backend/database"
)

func GetSearchedPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	categorySelection := r.FormValue("category")
	search := r.FormValue("search")
	filter := r.FormValue("filter")

	db, errDb := database.OpenDb(w)
	if errDb != nil {
		http.Error(w, "ERROR: Database cannot open", http.StatusBadRequest)
		return
	}
	defer db.Close()

	var posts []structs.Post
	var err error

	if categorySelection == "" {
		// Kategori seçilmemişse tüm postları çek
		posts, err = fetchAllPosts(db)
	} else {
		// Kategoriye göre postları çek
		posts, err = fetchPostsByCategory(db, categorySelection)
	}

	if search != "" {
		posts = filterPostsBySearch(posts, search)
	}

	if filter != "" {
		posts = FilterThePosts(posts, filter)
	}

	if err != nil {
		http.Error(w, "ERROR: Posts cannot use", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		http.Error(w, "ERROR: Failed to encode posts to JSON", http.StatusInternalServerError)
		return
	}
}

func FilterThePosts(posts []structs.Post, order string) []structs.Post {
	if order == "top" {
		for i := 0; i < len(posts)-1; i++ {
			for j := i + 1; j < len(posts); j++ {
				if posts[i].LikeCount < posts[j].LikeCount {
					posts[i], posts[j] = posts[j], posts[i]
				}
			}
		}
	}
	return posts
}

func filterPostsBySearch(posts []structs.Post, search string) []structs.Post {
	var filteredPosts []structs.Post
	for _, post := range posts {
		if strings.Contains(strings.ToLower(post.Title), strings.ToLower(search)) || strings.Contains(strings.ToLower(post.Content), strings.ToLower(search)) {
			filteredPosts = append(filteredPosts, post)
		}
	}
	return filteredPosts
}

func fetchAllPosts(db *sql.DB) ([]structs.Post, error) {
	rows, err := db.Query("SELECT ID, Title, UserId, Content, UserName, LikeCount FROM POSTS ORDER BY PostDate DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Title, &post.UserID, &post.Content, &post.UserName, &post.LikeCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func fetchPostsByCategory(db *sql.DB, category string) ([]structs.Post, error) {
	queryCategory := fmt.Sprintf("SELECT PostID FROM CATEGORIES WHERE %s = 1", category)
	rows, err := db.Query(queryCategory)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var postIDs []int
	for rows.Next() {
		var postID int
		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		postIDs = append(postIDs, postID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(postIDs) == 0 {
		return nil, nil
	}

	query := fmt.Sprintf("SELECT ID, Title, UserId, Content, UserName, LikeCount FROM POSTS WHERE ID IN (%s) ORDER BY PostDate DESC", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(postIDs)), ","), "[]"))

	rows, err = db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Title, &post.UserID, &post.Content, &post.UserName, &post.LikeCount)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
