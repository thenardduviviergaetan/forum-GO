package forum

// // FIXME
// func CheckSessionToken(db *sql.DB, w http.ResponseWriter, r *http.Request) {
// 	c, err := middle.GetCookie(w, r)
// 	if err != nil {
// 		return
// 	}
// 	sessionToken := c.Value

// 	// userSession, exists := models.Sessions[sessionToken]
// 	// if !exists {
// 	// 	w.WriteHeader(http.StatusUnauthorized)
// 	// 	return
// 	// }
// 	// userSession := db.QueryRow("SELECT session_token, expires_at FROM users WHERE session_token = ?", sessionToken)

// 	var userSession string
// 	err = db.QueryRow("SELECT session_token FROM users WHERE session_token=?", sessionToken).Scan(&userSession)
// 	if err != nil {
// 		if err != sql.ErrNoRows {
// 			http.Redirect(w, r, "/logout", http.StatusUnauthorized)
// 		}
// 	}
// if userSession.IsExpired() {
// 	// delete(models.Sessions, sessionToken)
// 	w.WriteHeader(http.StatusUnauthorized)
// 	return
// }

// fmt.Println(userSession.Username)
// }

// func refreshToken(w http.ResponseWriter, r *http.Request) {
// 	c, err := middle.GetCookie(w, r)
// 	if err != nil {
// 		return
// 	}
// 	sessionToken := c.Value

// 	userSession, exists := models.Sessions[sessionToken]
// 	if !exists {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}
// 	if userSession.IsExpired() {
// 		delete(models.Sessions, sessionToken)
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	middle.SetToken(w, r, userSession.Username)

// }
