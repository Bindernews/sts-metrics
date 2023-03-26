package web

import (
	"errors"
	"net/url"
	"os"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4"
)

func AbortMsg(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, c.Error(err).JSON())
}

// Returns the value of the session key, or an empty string if
// the value doesn't exist, or is not a string.
func SessGetString(s sessions.Session, key string) string {
	val := s.Get(key)
	if val == nil {
		return ""
	} else if val2, ok := val.(string); ok {
		return val2
	} else {
		return ""
	}
}

// Scan database rows into an array of T, returning the first
// error that occurs.
func ScanRows[T any](rows pgx.Rows, dest *[]T) error {
	for rows.Next() {
		o := new(T)
		if err := rows.Scan(o); err != nil {
			return err
		}
		*dest = append(*dest, *o)
	}
	return nil
}

func TryEach[T any](list []T, fn func(T) error) error {
	for _, v := range list {
		if err := fn(v); err != nil {
			return err
		}
	}
	return nil
}

// Takes a set of handler functions, some of which may be nil,
// and returns a slice of the non-nil ones in order.
func HandlerChain(chain ...gin.HandlerFunc) []gin.HandlerFunc {
	out := make([]gin.HandlerFunc, 0, len(chain))
	for _, f := range chain {
		if f != nil {
			out = append(out, f)
		}
	}
	return out
}

// Open a file for writing, failing if the file already existed.
func CreateNewFile(name string) (*os.File, error) {
	fd, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}
	if info, err := fd.Stat(); err != nil {
		fd.Close()
		return nil, err
	} else if info.Size() > 0 {
		fd.Close()
		return nil, errors.New("duplicate play id")
	}
	return fd, nil
}

// Strips suffix from the request URI before passing on to the next handler.
func StripRequestPrefix(suffix string) gin.HandlerFunc {
	// If suffix is empty, return a nop
	if suffix == "" {
		return func(_ *gin.Context) {}
	}
	return func(c *gin.Context) {
		oldurl := *c.Request.URL // make a copy
		oldurl.Path = strings.TrimPrefix(oldurl.Path, suffix)
		oldurl.RawPath = ""
		if newurl, err := url.Parse(oldurl.String()); err != nil {
			c.AbortWithError(404, err)
		} else {
			c.Request.URL = newurl
		}
	}
}

// Cast an any-slice to the destination type
func CastSlice[D, S any](src []S) []D {
	dst := make([]D, len(src))
	for i, v := range src {
		dst[i] = any(v).(D)
	}
	return dst
}

// Set all values from 'values' to true in 'st'.
func SetAdd[T comparable](st map[T]bool, values ...T) {
	for _, v := range values {
		st[v] = true
	}
}
