package forum

// func GetCurrentSession(w http.ResponseWriter, r *http.Request) {
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

// 	fmt.Println(userSession.Username)
// }
