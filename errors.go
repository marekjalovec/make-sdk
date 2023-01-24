package makesdk

import (
	"fmt"
	"net/http"
)

type HttpError struct {
	StatusCode int
	Err        error
}

func (e HttpError) Error() string {
	return fmt.Sprintf("HTTP status code %d, err: %v", e.StatusCode, e.Err)
}

func NewHttpError(url string, resp *http.Response) error {
	return &HttpError{
		StatusCode: resp.StatusCode,
		Err:        fmt.Errorf("HTTP request failure [%s] on %s", resp.Status, url),
	}
}

func getHttpError(err error) *HttpError {
	httpErr, ok := err.(*HttpError)
	if ok {
		return httpErr
	}

	return nil
}
