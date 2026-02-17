package pattern

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPWikiLoader struct {
	client  *http.Client
	maxSize int64
}

func NewHTTPWikiLoader(timeout time.Duration, maxSize int64) HTTPWikiLoader {
	return HTTPWikiLoader{
		client:  &http.Client{Timeout: timeout},
		maxSize: maxSize,
	}
}

func (l HTTPWikiLoader) Load(url string) (string, error) {
	resp, err := l.client.Get(url)
	if err != nil {
		return "", RecoverableError{Message: err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", RecoverableError{Message: fmt.Sprintf("unexpected http status: %d", resp.StatusCode)}
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/") {
		return "", RecoverableError{Message: "unsupported content-type"}
	}

	limited := io.LimitReader(resp.Body, l.maxSize+1)
	body, err := io.ReadAll(limited)
	if err != nil {
		return "", RecoverableError{Message: err.Error()}
	}
	if int64(len(body)) > l.maxSize {
		return "", RecoverableError{Message: "response size exceeds limit"}
	}

	return string(body), nil
}
