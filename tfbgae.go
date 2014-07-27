package tfbgae

import (
    "net/http"
    "html/template"
    // "appengine/user"
    "github.com/gorilla/mux"
)

func init() {
    r := mux.NewRouter()
    r.StrictSlash(true)
    r.HandleFunc("/", rootHandler)
    r.HandleFunc("/posts", InjectDb(List))
    r.HandleFunc("/posts/page/{page}", InjectDb(List))
    r.HandleFunc("/posts/add", AddPost)
    r.HandleFunc("/posts/save", InjectDb(Save))
    r.HandleFunc("/posts/edit/save", InjectDb(Save))
    r.HandleFunc("/posts/{id}", InjectDb(Show))
    r.HandleFunc("/posts/edit/{id}", InjectDb(Edit))
    r.HandleFunc("/posts/delete/{id}", InjectDb(Delete))
    
    http.Handle("/", r)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    t, _ := template.ParseFiles("view/index.html")
    t.Execute(w, nil)
}

