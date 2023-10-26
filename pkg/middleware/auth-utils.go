package forum

import (
	"fmt"
	"io"
	"net/http"
)

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
