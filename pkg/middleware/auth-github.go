package forum

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	models "forum/pkg/models"
)

// GetGithubToken takes the callback's code sent by Github, the ID and Secret from the developper and use this to create a request
// that will be sent back to Github to request the clients infos. If everything checks out github returns a body with all the
// information about our user. The function then extracte the access_token and returns it if nothing went wrong.
func GetGithubToken(code, id, secret string) (string, error) {
	requestJSON, _ := json.Marshal(map[string]string{
		"client_id":     id,
		"client_secret": secret,
		"code":          code,
	})

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

	respBody, errPerfReq := performRequest(req)
	if errPerfReq != nil {
		fmt.Println("[Github Login] -> Error Sending request to client", errPerfReq)
		return "", fmt.Errorf("[Github Login] -> Error Sending request to client\n%v\n", errPerfReq)
	}

	var gitResp models.GitResp

	errUnmarshal := json.Unmarshal(respBody, &gitResp)
	if errUnmarshal != nil {
		fmt.Println("[Github login] -> error trying to unmarshal response body", errUnmarshal)
		return "", fmt.Errorf("[Github login] -> error trying to unmarshal response body\n%v\n", errUnmarshal)
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

	return performRequest(request)
}
