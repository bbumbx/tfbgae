package tfbgae

import(
	"net/http"
)

func InjectDb(fn func (http.ResponseWriter, *http.Request, *DB)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db := NewDB(r)
        fn(w, r, db)
	}
}

func Render() {
	
}