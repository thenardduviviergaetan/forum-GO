package forum

import (
	"fmt"
	"io"
	"net/http"
)

// performRequest receive a request as argument and sends it then proceed to retrieve the response
// and then extracts the body to send it back as a return. If anything went wrong, return a nil array
// and an error.
func performRequest(request *http.Request) ([]byte, error) {
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
