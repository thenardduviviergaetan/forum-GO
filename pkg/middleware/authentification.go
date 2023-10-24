package forum

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	models "forum/pkg/models"
	s "forum/sessions"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Check user credentials
// TODO merge this is some way with the normal Auth.
func AuthGithub(db *sql.DB, w http.ResponseWriter, r *http.Request, user *models.User) error {
	// email, password := r.FormValue("email"), r.FormValue("password")

	var tmpUser models.User

	err := db.QueryRow("SELECT id,username,email, pwd FROM users WHERE email=?", user.Email).Scan(
		&tmpUser.ID,
		&tmpUser.Username,
		&tmpUser.Email,
		&tmpUser.Password)

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return errors.New("Email not found")
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(tmpUser.Password), []byte(user.Password))
	if err != nil {
		return errors.New("Wrong Password")
	}
	user.ID = tmpUser.ID
	s.SetToken(db, w, r, user)
	return nil
}

// Check user credentials
func Auth(db *sql.DB, w http.ResponseWriter, r *http.Request, user *models.User) error {
	email, password := r.FormValue("email"), r.FormValue("password")

	err := db.QueryRow("SELECT id,username,email, pwd FROM users WHERE email=?", email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password)

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		} else {
			return errors.New("Email not found")
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return errors.New("Wrong Password")
	}
	s.SetToken(db, w, r, user)
	return nil
}

// GetGithubToken takes the callback's code sent by github, the ID and Secret from the developper and use this to create a request
// that will be sent back to github to request the clients infos. If everything checks out github returns a body with all the
// information about our user. The function then extracte the access_token and returns it if nothing went wrong.
func GetGithubToken(code, id, secret string) (string, error) {
	reqBodyMap := map[string]string{
		"client_id":     id,
		"client_secret": secret,
		"code":          code,
	}

	requestJSON, _ := json.Marshal(reqBodyMap)

	req, requestErr := http.NewRequest(
		"POST",
		"https://github.com/login/oauth/access_token",
		bytes.NewBuffer(requestJSON),
	)

	if requestErr != nil {
		fmt.Println("[Github Login] -> Error creating new request", requestErr)
		return "", fmt.Errorf("[Github Login] -> Error creating new request\n%v\n", requestErr)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	response, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		fmt.Println("[Github Login] -> Error Sending request to client", respErr)
		return "", fmt.Errorf("[Github Login] -> Error Sending request to client\n%v\n", respErr)
	}

	respBody, errReadBody := io.ReadAll(response.Body)
	if errReadBody != nil {
		fmt.Println("[Github Login] -> Error reading response body", errReadBody)
		return "", fmt.Errorf("[Github Login] -> Error reading response body\n%v\n", respErr)
	}

	var gitResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	errUnmarshal := json.Unmarshal(respBody, &gitResp)
	if errUnmarshal != nil {
		fmt.Println("[Github Login] -> Error trying to Unmarshal response body", errReadBody)
		return "", fmt.Errorf("[Github Login] -> Error trying to Unmarshal response body\n%v\n", respErr)
	}

	return gitResp.AccessToken, nil
}

// GetGithubData takes the user access token and wil create a request with the token embed inside,
// github then returns a json with all the information about the user if no problem is spotted.
func GetGithubData(access_token string) ([]byte, error) {
	request, reqErr := http.NewRequest(
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if reqErr != nil {
		fmt.Println("[Github Data] -> Error trying to create new request", reqErr)
		return nil, fmt.Errorf("[Github Data] -> Error trying to create new request\n%v\n", reqErr)
	}

	request.Header.Set("Authorization", fmt.Sprintf("token %s", access_token))
	response, respErr := http.DefaultClient.Do(request)
	if respErr != nil {
		fmt.Println("[Github Data] -> Error sending request and waiting for response", respErr)
		return nil, fmt.Errorf("[Github Data] -> Error sending request and waiting for response\n%v\n", respErr)
	}

	respBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		fmt.Println("[Github Data] -> Error reading response body", readErr)
		return nil, fmt.Errorf("[Github Data] -> Error reading response body\n%v\n", readErr)
	}

	return respBody, nil
}

func GetGoogleData(access_token, id_token string) ([]byte, error) {
	request, reqErr := http.NewRequest(
		"GET",
		fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", access_token),
		nil,
	)
	if reqErr != nil {
		fmt.Println("[Google Data] -> Error trying to create new request", reqErr)
		return nil, fmt.Errorf("[Google Data] -> Error trying to create new request\n%v\n", reqErr)
	}

	return retrieveData(request)
}

func retrieveData(request *http.Request) ([]byte, error) {
	response, respErr := http.DefaultClient.Do(request)
	if respErr != nil {
		fmt.Println("[Retrieve Data] -> Error sending request and waiting for response", respErr)
		return nil, fmt.Errorf("[Retrieve Data] -> Error sending request and waiting for response\n%v\n", respErr)
	}

	respBody, readErr := io.ReadAll(response.Body)
	if readErr != nil {
		fmt.Println("[Retrieve Data] -> Error reading response body", readErr)
		return nil, fmt.Errorf("[Retrieve Data] -> Error reading response body\n%v\n", readErr)
	}

	return respBody, nil
}
