package forum

import (
	"database/sql"
	//"errors"
	"strconv"
	"net/http"
	"log"
	//"time"
	models "forum/pkg/models"
)

func AddCategory(db *sql.DB, r *http.Request) error {

	_, err := db.Exec("INSERT INTO categories(title, description, time) VALUES (?,?,datetime())",
						r.FormValue("catitle"), 
						r.FormValue("catdescription"))
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

func FetchCat(db *sql.DB) []models.Categories {
	rows, err := db.Query("SELECT id, title, description FROM categories")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var categorylst []models.Categories
	for rows.Next() {
		var categories models.Categories
        err = rows.Scan(&categories.ID, &categories.Title, &categories.Description)
        if err != nil {
            log.Fatal(err)
        }
		categorylst = append(categorylst, categories)
	}
	return categorylst
}
