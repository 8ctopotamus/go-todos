package main

import (

	// "strings"
	"fmt"
	"html/template"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model
	ID    uint
	Title string
	Done  bool
}

type TodoPageData struct {
	PageTitle string
	Todos     []Todo
}

const PORT int = 8080

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Todo{})

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			if err := r.ParseForm(); err != nil {
				fmt.Fprintf(w, "ParseForm() err: %v", err)
				return
			}
			todo := r.FormValue("todo")
			db.Create(&Todo{Title: todo, Done: false})
		}

		var todos []Todo
		db.Find(&todos)
		data := TodoPageData{
			PageTitle: "My Go TODOs list",
			Todos:     todos,
		}

		tmpl.Execute(w, data)
	})

	fmt.Printf("Listening on http://localhost:%v\n", PORT)
	http.ListenAndServe(fmt.Sprintf(":%v", PORT), nil)
}
