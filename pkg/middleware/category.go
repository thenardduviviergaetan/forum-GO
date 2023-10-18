package forum

import (
	"database/sql"
	//"errors"
	"strconv"
	"net/http"
	//"log"
	"time"
	//models "forum/pkg/models"
)

func AddCategory(db *sql.DB, r *http.Request) error {

	_, err := db.Exec("INSERT INTO categories(title, description, time) VALUES (?,?,?)",
						r.FormValue("catitle"), 
						r.FormValue("catdescription"), 
						time.Now())
	if err != nil {
        return err
    }
	return nil
}

func ModCategory(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("creatcat"))
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE categories SET title=?, description=? WHERE id=?", r.FormValue("catitle"), r.FormValue("catdescription"), id)
	if err != nil {
		return err
	}
	return nil
}

func DelCategory(db *sql.DB, r *http.Request) error {

	id, err := strconv.Atoi(r.FormValue("delcat"))
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM categories WHERE id=?", id)
	if err != nil {
        return err
    }
	return nil
}