package tfbgae

import (
	"appengine"
	"appengine/user"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
	"time"
)

var templates = template.Must(template.ParseFiles("view/posts/list.html",
	"view/posts/addpost.html", "view/header.html", "view/footer.html", "view/posts/show.html", "view/posts/edit.html", "view/dist.html", "view/posts/post.html",
	"view/posts/showadmin.html"))

type Post struct {
	Key         string `datastore:"-"`
	Title       string
	Description string
	Body        string `datastore:",noindex"`
	Date        time.Time
	IsActive    bool
}

type Resp struct {
	NextPage int
	PrevPage int
	Title string
	Posts    []Post
}

const (
	PAGE_SIZE  = 5
	SORT_FIELD = "Date"
)

func List(w http.ResponseWriter, r *http.Request, db *DB) {
	page, _ := strconv.Atoi(mux.Vars(r)["page"])
	if page < 0 {
		page = 0
	}
	offset := page * PAGE_SIZE
	results, keys, err := db.LoadPosts(SORT_FIELD, offset, PAGE_SIZE)
	if err != nil {
		panic(err)
	}
	var posts []Post
	for i := range results {
		r := results[i]
		r.Key = keys[i].Encode()
		posts = append(posts, r)
	}
	nextPage := page + 1
	if len(posts) < PAGE_SIZE {
		nextPage = -1
	}

	err = templates.ExecuteTemplate(w, "list.html", &Resp{Posts: posts,
	 NextPage: nextPage, PrevPage: page - 1, Title:"Latest blog posts"})
	if err != nil {
		panic(err)
	}
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	if isUserAdmin(r) {
		err := templates.ExecuteTemplate(w, "edit.html", nil)
		if err != nil {
			panic(err)
		}
	} else {
		http.Redirect(w, r, "/posts", http.StatusFound)
	}

}

func Show(w http.ResponseWriter, r *http.Request, db *DB) {
	id := mux.Vars(r)["id"]
	post, err := db.LoadPost(id)
	if err != nil {
		panic(err)
	}
	t := "show.html"
	if isUserAdmin(r) {
		t = "showadmin.html"
	}
	err = templates.ExecuteTemplate(w, t, post)
	if err != nil {
		panic(err)
	}
}

func Save(w http.ResponseWriter, r *http.Request, db *DB) {
	if isUserAdmin(r) {
		key := r.FormValue("id")
		title := r.FormValue("title")
		body := r.FormValue("body")
		description := r.FormValue("description")
		p := &Post{Title: title, Body: body, Description: description, Date: time.Now(), Key: key, IsActive: true}
		err := db.SavePost(p)
		if err != nil {
			panic(nil)
		}
	}
	http.Redirect(w, r, "/posts", http.StatusFound)

}

func Edit(w http.ResponseWriter, r *http.Request, db *DB) {
	id := mux.Vars(r)["id"]
	post, err := db.LoadPost(id)
	if err != nil {
		panic(err)
	}
	err = templates.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		panic(err)
	}
}

func Delete(w http.ResponseWriter, r *http.Request, db *DB) {
	if isUserAdmin(r) {
		id := mux.Vars(r)["id"]
		err := db.DeletePost(id)
		if err != nil {
			panic(err)
		}
	}
	http.Redirect(w, r, "/posts", http.StatusFound)
}

func isUserAdmin(r *http.Request) bool {
	c := appengine.NewContext(r)
	u := user.Current(c)
	if u != nil {
		return u.Admin
	}
	return false
}
