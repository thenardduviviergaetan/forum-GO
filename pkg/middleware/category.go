package forum

import (
	"database/sql"
	//"errors"
	"log"
	"net/http"
	"strconv"

	//"time"
	models "forum/pkg/models"
)

func AddCategory(db *sql.DB, r *http.Request) error {
	_, err := db.Exec("INSERT INTO categories(title, descriptions, creation) VALUES (?,?,datetime())",
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
	_, err = db.Exec("UPDATE categories SET title=?, descriptions=? WHERE id=?", r.FormValue("catitle"), r.FormValue("catdescription"), id)
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

func FetchCat(db *sql.DB, current int64) []models.Categories {
	rows, err := db.Query("SELECT id, title, descriptions FROM categories")
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
		categories.Ifcurtentcat = current == categories.ID
		if categories.ID == 0 {
			continue
		}
		categorylst = append(categorylst, categories)
	}
	return categorylst
}
