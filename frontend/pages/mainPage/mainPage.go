package mainpage

import (
	"net/http"
	"text/template"

	"forum/backend/auth"
	"forum/backend/database"
	"forum/backend/requests"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
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

	var tmpl *template.Template
	var err error

	authenticated, _, _ := auth.IsAuthenticated(r, db)
	if !authenticated {
		tmpl, err = template.ParseFiles("frontend/pages/mainPage/sessionless/main.html")
		if err != nil {
			http.Error(w, "ERROR: Unable to parse template", http.StatusInternalServerError)
			return
		}
	} else {
		tmpl, err = template.ParseFiles("frontend/pages/mainPage/session/main.html")
		if err != nil {
			http.Error(w, "ERROR: Unable to parse template", http.StatusInternalServerError)
			return
		}
	}

	posts, err := requests.GetDataForServe("http://localhost:8080/api/allposts")
	if err != nil {
		http.Error(w, "Could not fetch post data", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, posts)
	if err != nil {
		http.Error(w, "ERROR: Unable to execute template", http.StatusInternalServerError)
		return
	}
}
