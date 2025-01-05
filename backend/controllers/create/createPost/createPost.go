package createpost

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"forum/backend/auth"
	"forum/backend/database"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(20 << 20)
	if err != nil {
		http.Error(w, "ERROR: File size exceeds the 20MB limit", http.StatusRequestEntityTooLarge)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if !IsForumPostValid(title, content) {
		http.Error(w, "ERROR: Content or title cannot be empty", http.StatusBadRequest)
		return
	}

	var PhotoPath string
	file, handler, err := r.FormFile("photo")
	if err == nil {
		defer file.Close()

		// Dosya türünü kontrol et
		fileExt := strings.ToLower(filepath.Ext(handler.Filename))
		if fileExt != ".jpg" && fileExt != ".jpeg" && fileExt != ".png" && fileExt != ".gif" {
			http.Error(w, "ERROR: Unsupported file type. Only JPEG, PNG, and GIF are allowed.", http.StatusBadRequest)
			return
		}

		// Fotoğrafı uploads dizinine kaydet
		uploadDir := "./uploads"
		if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
			err = os.MkdirAll(uploadDir, os.ModePerm)
			if err != nil {
				http.Error(w, "Unable to create upload directory", http.StatusInternalServerError)
				return
			}
		}

		tempFile, err := os.Create(fmt.Sprintf("%s/%s", uploadDir, handler.Filename))
		if err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}
		defer tempFile.Close()
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Unable to copy file content", http.StatusInternalServerError)
			return
		}
		PhotoPath = fmt.Sprintf("/uploads/%s", handler.Filename)
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

	result, errEx := db.Exec(`INSERT INTO POSTS (UserID, UserName, Title, Content, PhotoPath) VALUES (?, ?, ?, ?, ?)`, userId, userName, title, content, PhotoPath)
	if errEx != nil {
		http.Error(w, "ERROR: Post did not add to the database", http.StatusBadRequest)
		return
	}

	postID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "ERROR: Could not retrieve post ID", http.StatusBadRequest)
		return
	}

	categoryValues := GetCategoryValues(r)
	_, err = db.Exec(`INSERT INTO CATEGORIES (USERID, PostID, GO, HTML, CSS, PHP, PYTHON, C, "CPP", "CSHARP", JS, ASSEMBLY, REACT, FLUTTER, RUST) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		userId, postID, categoryValues["go"], categoryValues["html"], categoryValues["css"], categoryValues["php"],
		categoryValues["python"], categoryValues["c"], categoryValues["cpp"], categoryValues["csharp"],
		categoryValues["js"], categoryValues["assembly"], categoryValues["react"], categoryValues["flutter"], categoryValues["rust"])
	if err != nil {
		http.Error(w, "ERROR: Could not add categories to the database", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Post successfully created")
}

func IsForumPostValid(title, content string) bool {
	return title != "" && content != ""
}

func GetCategoryValues(r *http.Request) map[string]int {
	categories := []string{"go", "html", "css", "php", "python", "c", "cpp", "csharp", "js", "assembly", "react", "flutter", "rust"}
	categoryValues := make(map[string]int)

	for _, category := range categories {
		if r.FormValue(category) == "true" {
			categoryValues[category] = 1
		} else {
			categoryValues[category] = 0
		}
	}
	return categoryValues
}
