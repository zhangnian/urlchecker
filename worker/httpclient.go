package worker

import (
	"net/http"
)

func HttpGet(url string) (statusCode int, contentLength int64, err error) {
	response, err := http.Get(url)
	if err != nil {
		return
	}
	response.Body.Close()

	statusCode = response.StatusCode
	contentLength = response.ContentLength

	return
}

func HttpPost(url string) {

}
