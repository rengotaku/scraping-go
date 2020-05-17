package lib

import (
	"crypto/sha256"
	"encoding/hex"
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

func getTime(c *gin.Context) (*time.Time, string, bool) {
	session := sessions.Default(c)

	flashes := session.Flashes(timeToken)
	session.Save()
	if len(flashes) == 0 {
		return nil, "", false
	}

	timeStr := flashes[0].(string)
	epoch, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		return nil, "", false
	}

	t := time.Unix(0, epoch)
	return &t, flashes[1].(string), true
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

		t, id, ok := getTime(c)

		if !ok || !time.Now().Before(*t) || id != getIdentifyHash(c) {
			errorFunc(c)
			return
		}

		session := sessions.Default(c)
		session.AddFlash(getUnixNanoAsStr(*t), timeToken)
		session.AddFlash(id, timeToken)
		session.Save()

		c.Next()
	}
}

func BeginOtt(c *gin.Context, validSec time.Duration) bool {
	session := sessions.Default(c)

	session.Flashes(timeToken)
	t := time.Now().Add(time.Second * validSec)
	id := getIdentifyHash(c)
	if id == "" {
		return false
	}

	session.AddFlash(getUnixNanoAsStr(t), timeToken)
	session.AddFlash(id, timeToken)
	session.Save()

	return true
}

func getIdentifyHash(c *gin.Context) string {
	if c.GetHeader("User-Agent") == "" || c.ClientIP() == "" {
		return ""
	}

	// Mozilla/5.0... + ::1
	seed := []byte(c.GetHeader("User-Agent") + c.ClientIP())
	sha := sha256.Sum256(seed)
	return hex.EncodeToString(sha[:])
}

func getUnixNanoAsStr(t time.Time) string {
	epoch := t.UnixNano()
	return strconv.FormatInt(epoch, 10)
}

func EndOtt(c *gin.Context) string {
	session := sessions.Default(c)
	flashes := session.Flashes(timeToken)
	session.Save()

	// dosen't match
	sha := sha256.Sum256([]byte(flashes[0].(string) + flashes[1].(string)))
	return hex.EncodeToString(sha[:])
}
