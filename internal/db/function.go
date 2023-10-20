package forum

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	middle "forum/pkg/middleware"
	models "forum/pkg/models"
)

func linkpost(app *App_db, postid int64) (tablike map[int64]bool, tabdislike map[int64]bool) {
	tablike, tabdislike = make(map[int64]bool), make(map[int64]bool)
	rows, err := app.DB.Query("SELECT userid,likes FROM linkpost WHERE postid = ?", postid)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	var userid int64
	var like bool
	for rows.Next() {
		rows.Scan(&userid, &like)
		if like {
			tablike[userid] = true
		} else {
			tabdislike[userid] = true
		}
	}
	return
}

func mkdirPostAsset(app *App_db, idpost int64, post *models.Post, r *http.Request) {
	dirfile := "web/static/upload/img/post" + strconv.Itoa(int(idpost)) + "/"
	if _, err := os.Stat(dirfile); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dirfile+"comment", os.ModePerm)
		if err != nil {
			fmt.Println("mkfirPostAsset : ", err)
			return
		}
	}
	file, fileheader, err := r.FormFile("myFile")
	if err != nil {
		if err.Error() == "http: no such file" {
			return
		}
		fmt.Println("formfile :", err.Error())
		return
	}
	defer file.Close()
	//.svg,.gif,.pnj,.jpeg
	tab := strings.Split(fileheader.Filename, ".")
	// fileheader.Size in octer
	if len(tab) != 2 || fileheader.Size > (20*1000000) {
		// file.Close()
		fmt.Println("badformat or >20Mo")
		return
	}
	key := tab[1]
	switch key {
	case "svg", "gif", "pnj", "jpeg", "jpg":
	default:
		// file.Close()
		fmt.Println("badformat")
		return
	}
	f, err := os.OpenFile(dirfile+fileheader.Filename, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	io.Copy(f, file)
	post.Img = fileheader.Filename
	middle.UpdateImgPoste(app.DB, idpost, fileheader.Filename)
}

func mkdirCommentAsset(app *App_db, idpost int64, comment *models.Comment, r *http.Request) {
	dirfile := "web/static/upload/img/post" + strconv.Itoa(int(idpost)) + "/comment"
	if _, err := os.Stat(dirfile); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(dirfile+"comment", os.ModePerm)
		if err != nil {
			fmt.Println("mkfirPostAsset : ", err)
			return
		}
	}
	dirfile += "/comment_" + strconv.Itoa(int(comment.ID))
	file, fileheader, err := r.FormFile("myFile")
	if err != nil {
		if err.Error() == "http: no such file" {
			return
		}
		fmt.Println("formfile :", err.Error())
		return
	}
	defer file.Close()
	//.svg,.gif,.pnj,.jpeg
	tab := strings.Split(fileheader.Filename, ".")
	if len(tab) != 2 || fileheader.Size > 20*1000000 {
		// file.Close()
		fmt.Println("badformat or >20Mo")
		return
	}
	key := tab[1]
	switch key {
	case "svg", "gif", "pnj", "jpeg", "jpg":
	default:
		// file.Close()
		fmt.Println("badformat")
		return
	}
	f, err := os.OpenFile(dirfile+"_"+fileheader.Filename, os.O_WRONLY|os.O_CREATE, 0o666)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	io.Copy(f, file)
	comment.Img = "comment_" + strconv.Itoa(int(comment.ID)) + "_" + fileheader.Filename
	middle.UpdateImgComment(app.DB, comment.Postid, comment.ID, comment.Img)
}
