package searchedpostspage

import (
	"html/template"
	"net/http"

	"forum/backend/requests"
)

func SearchedPostsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "ERROR: Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	filter := r.FormValue("filter")
	category := r.FormValue("category")
	search := r.FormValue("search")

	data, err := requests.GetSearchedDataForServeWithReq("http://localhost:8080/api/searchedposts", filter, category, search)
	if err != nil {
		http.Error(w, "ERROR: Cannot get post", http.StatusBadRequest)
		return
	}

	tmpl, err := template.ParseFiles("frontend/pages/searchedPostsPage/searchedPostsPage.html")
	if err != nil {
		http.Error(w, "ERROR: Unable to parse template", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "ERROR: Unable to execute template", http.StatusInternalServerError)
		return
	}
}
