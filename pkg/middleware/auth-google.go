package forum

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GetGoogleToken takes the callback's code sent by Google, the ID and Secret from the developper and use this to create a request
// that will be sent back to Google to request the clients infos. If everything checks out github returns a body with all the
// information about our user. The function then extracte the access_token and returns it if nothing went wrong.
func GetGoogleToken(code, id, secret, origin string) (string, string, error) {
	reqURL := fmt.Sprintf("%s?%s&%s&%s&%s&%s",
		"https://www.googleapis.com/oauth2/v4/token",
		"client_id="+id,
		"client_secret="+secret,
		"grant_type=authorization_code",
		"code="+code,
		"redirect_uri=https://localhost:8080/google/callback/"+origin)

	request, errReq := http.NewRequest(
		"POST",
		reqURL,
		nil,
	)

	if errReq != nil {
		fmt.Println("[Google Login][Token] -> Error creating new request", errReq)
		return "", "", fmt.Errorf("[Google Login][Token] -> Error creating new request\n%v\n", errReq)
	}

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")

	respBody, errPerfReq := performRequest(request)
	if errPerfReq != nil {
		fmt.Println("[Google Login] -> Error Sending request to client", errPerfReq)
		return "", "", fmt.Errorf("[Github Login] -> Error Sending request to client\n%v\n", errPerfReq)
	}

	var googleTokenResp map[string]interface{}

	errUnmarshal := json.Unmarshal(respBody, &googleTokenResp)
	if errUnmarshal != nil || googleTokenResp == nil || googleTokenResp["error"] != nil {
		fmt.Println("[Google login][Token] -> error trying to unmarshal response body", errUnmarshal)
		return "", "", fmt.Errorf("[Google login][Token] -> error trying to unmarshal response body\n%v\n", errUnmarshal)
	}

	return googleTokenResp["access_token"].(string), googleTokenResp["id_token"].(string), nil
}

// GetGoogleData takes the user access token and wil create a request with the token embed inside,
// Google then returns a json with all the information about the user if no problem is spotted.
func GetGoogleData(access_token, id_token string) ([]byte, error) {
	request, reqErr := http.NewRequest(
		"GET",
		fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", access_token),
		nil,
	)
	if reqErr != nil {
		fmt.Println("[Google Data][Data] -> Error trying to create new request", reqErr)
		return nil, fmt.Errorf("[Google Data][Data] -> Error trying to create new request\n%v\n", reqErr)
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", id_token))

	return performRequest(request)
}
