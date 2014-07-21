package tfbgae

import(
	"appengine"
	"appengine/datastore"
	"net/http"
)

const(
	postKind = "Post"
)

type DB struct {
	context appengine.Context
}

func NewDB(r *http.Request) *DB{
	return &DB{context: appengine.NewContext(r)}	
}

func (db *DB) LoadPosts(sortField string, offset int,  limit int) ([]Post,[] *datastore.Key, error) {
	q := datastore.NewQuery(postKind).
		Filter("IsActive=", true).
		Project("Title", "Date", "Description").
		Offset(offset).
		Limit(limit).
        Order("-" + sortField)

    var result []Post            
    keys, err := q.GetAll(db.context, &result)
	if err != nil{
		return nil, nil, err
	}
	return result, keys, nil		
}

func (db *DB) LoadPost(id string) (*Post, error) {
	key, _ := datastore.DecodeKey(id)
	var result Post
	err := datastore.Get(db.context, key, &result)
	result.Key = key.Encode()
	if err != nil{
		return nil, err
	}
	return &result, nil		
}

 func (db *DB) SavePost(p *Post) error {
 	if p.Key != "" {
 		key, _ := datastore.DecodeKey(p.Key)
 		_, err := datastore.Put(db.context, key, p)
 		return err
 	} else {
 		key := datastore.NewIncompleteKey(db.context, postKind, nil)	
 		_, err := datastore.Put(db.context, key, p)
 		return err
 	} 				
 }

 func (db *DB) DeletePost(id string) error {
	p, err := db.LoadPost(id)
	if err != nil{
		return err
	}
	p.IsActive = false
	err = db.SavePost(p)
	return err	
}



