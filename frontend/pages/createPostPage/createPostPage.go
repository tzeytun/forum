package createpostpage

import (
	"net/http"
	"path/filepath"
	"strings"

	"forum/backend/requests"
)

const createPostApiUrl = "http://localhost:8080/api/createpost"

func CreatePostPage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		http.ServeFile(w, r, "frontend/pages/createPostPage/createPostPage.html")
	case "POST":
		cookie, cookieErr := r.Cookie("session_token")
		if cookieErr != nil {
			http.Error(w, "ERROR: You are not authorized to create post", http.StatusUnauthorized)
			return
		}
		title := r.FormValue("title")
		content := r.FormValue("content")

		file, handler, err := r.FormFile("photo")
		if err != nil {
			http.Error(w, "ERROR: Failed to upload photo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		fileExt := strings.ToLower(filepath.Ext(handler.Filename))
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" {
			http.Error(w, "ERROR: Unsupported file type. Only JPEG, PNG, and GIF are allowed.", http.StatusBadRequest)
			return
		}

		categoryDatas := GetCategoryDatas(r)

		err = requests.CreatePostRequest(createPostApiUrl, title, content, categoryDatas, cookie.Value, file, handler.Filename)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/myposts", http.StatusSeeOther)
	}
}

func GetCategoryDatas(r *http.Request) map[string]string {
	categories := []string{"go", "html", "css", "php", "python", "c", "cpp", "csharp", "js", "assembly", "react", "flutter", "rust"}
	categoryDatas := make(map[string]string)

	for _, category := range categories {
		if r.FormValue(category) == "true" {
			categoryDatas[category] = "true"
		} else {
			categoryDatas[category] = "false"
		}
	}
	return categoryDatas
}
