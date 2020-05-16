package lib

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	timeToken = "timeToken"
)

// Options stores configurations for a CSRF middleware.
type TimeTokenOptions struct {
	IgnoreMethods []string
	ErrorFunc     gin.HandlerFunc
}

func getTime(c *gin.Context) (*time.Time, bool) {
	session := sessions.Default(c)

	flashes := session.Flashes(timeToken)
	session.Save()
	if len(flashes) == 0 {
		return nil, false
	}

	timeStr := flashes[0].(string)
	epoch, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return nil, false
	}

	t := time.Unix(0, epoch)
	return &t, true
}

// Middleware validates One time token.
func TimeTokenMiddleware(options TimeTokenOptions) gin.HandlerFunc {
	ignoreMethods := options.IgnoreMethods
	errorFunc := options.ErrorFunc

	if ignoreMethods == nil {
		ignoreMethods = []string{"GET", "HEAD", "OPTIONS"}
	}

	if errorFunc == nil {
		errorFunc = func(c *gin.Context) {
			panic(errors.New("request Timeout"))
		}
	}

	return func(c *gin.Context) {
		if inArray(ignoreMethods, c.Request.Method) {
			c.Next()
			return
		}

		t, ok := getTime(c)

		if !ok || !time.Now().Before(*t) {
			errorFunc(c)
			return
		}

		session := sessions.Default(c)
		session.AddFlash(getUnixNanoAsStr(*t), timeToken)
		session.Save()

		c.Next()
	}
}

func BeginOtt(c *gin.Context, validSec time.Duration) {
	session := sessions.Default(c)

	session.Flashes(timeToken)
	t := time.Now().Add(time.Second * validSec)
	session.AddFlash(getUnixNanoAsStr(t), timeToken)

	session.Save()
}

func getUnixNanoAsStr(t time.Time) string {
	epoch := t.UnixNano()
	return strconv.FormatInt(epoch, 10)
}
