package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginPageHandler)
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/ui/static/"))))

	fmt.Println("Server running on port: 8080")
	http.ListenAndServe(":8080", mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title              string
		TotalStudents      int
		TotalCourses       int
		TotalBills         int
		AverageScansPerDay int
	}{
		Title:              "Student Management System",
		TotalStudents:      1220,
		TotalCourses:       45,
		AverageScansPerDay: 169,
	}
	tmpl := template.Must(template.ParseFiles("ui/html/templates/index.html"))
	if err := tmpl.Execute(w, data); err != nil {
		fmt.Println(err)
	}
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("ui/html/templates/login.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println(err)
	}
}
