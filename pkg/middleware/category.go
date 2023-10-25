package forum

import (
	"database/sql"
	//"errors"
	"log"
	"net/http"
	"strconv"
	//"fmt"
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
	rows, err := db.Query("SELECT id FROM post")
	if err != nil {
		return err
	}
	for rows.Next() {
		var temp int
		err = rows.Scan(&temp)
		if err != nil {
			return err
		}
		var exist bool
		err := db.QueryRow("SELECT EXISTS( SELECT * FROM linkcatpost WHERE postid = ?) AS exist", temp).Scan(&exist)
		if err != nil {
			return err
		}
		if !exist {
			_, err = db.Exec("DELETE FROM post WHERE id=?", temp)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func FetchCat(db *sql.DB, current []int) []models.Categories {
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
		for _, v := range current {
			if v == int(categories.ID) {
				categories.Ifcurtentcat = true
			}
		}
		categorylst = append(categorylst, categories)
		//categories.Ifcurtentcat = current == categories.ID
	}
	return categorylst
}
