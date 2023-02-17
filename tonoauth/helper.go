package tonoauth

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func abortErr(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, c.Error(err).JSON())
}

// Returns the value of the session key, or an empty string if
// the value doesn't exist, or is not a string.
func sessGetString(s sessions.Session, key string) string {
	val := s.Get(key)
	if val == nil {
		return ""
	} else if val2, ok := val.(string); ok {
		return val2
	} else {
		return ""
	}
}

// Helper struct to build an HTTP request, perform it, and parse the output as JSON
type EzHttpRequest struct {
	Method  string
	Url     string
	Body    []byte
	Headers map[string]any
}

func (r *EzHttpRequest) Do(client *http.Client, jsonOut any) error {
	var body io.Reader = nil
	if r.Body != nil {
		body = bytes.NewReader(r.Body)
	}
	req, err := http.NewRequest(r.Method, r.Url, body)
	if err != nil {
		return err
	}
	for k, v := range r.Headers {
		req.Header.Set(k, v.(string))
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(jsonOut); err != nil {
		return err
	}
	return nil
}