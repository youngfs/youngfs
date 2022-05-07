package routers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"icesos/command/vars"
	"icesos/log"
	"icesos/util"
	"net/http"
	"strconv"
	"time"
)

func authorizationHeader(user, pw string) string {
	base := user + ":" + pw
	return "Basic " + base64.StdEncoding.EncodeToString(util.StringToBytes(base))
}

func auth(realms ...string) gin.HandlerFunc {
	realm := "Basic realm="
	if len(realms) == 0 {
		realm += "Authorization Required"
	} else {
		for _, u := range realms {
			realm += strconv.Quote(u)
		}
	}

	admin := authorizationHeader("ices", "ices")

	return func(c *gin.Context) {
		if admin != c.Request.Header.Get("Authorization") {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(vars.UserKey, "ices")
	}
}
func Logger(reqs ...string) gin.HandlerFunc {
	req := ""
	for _, u := range reqs {
		req += u
	}

	return func(c *gin.Context) {
		c.Set(vars.UUIDKey, uuid.New().String())
		startTime := time.Now()
		log.Infow("request start",
			vars.UUIDKey, c.Value(vars.UUIDKey),
			requestNameKey, req,
			requestUrlKey, c.Request.URL,
			vars.UserKey, c.Value(vars.UserKey),
			requestStartTimeKey, startTime.Format(timeFormat),
		)
		c.Next()
		endTime := time.Now()
		log.Infow("request finish",
			vars.UUIDKey, c.Value(vars.UUIDKey),
			requestNameKey, req,
			vars.UserKey, c.Value(vars.UserKey),
			requestFinishTimeKey, endTime.Format(timeFormat),
			timeElapsedKey, endTime.Sub(startTime).Seconds(),
			responseHTTPCodeKey, c.Writer.Status(),
			vars.CodeKey, c.Value(vars.CodeKey),
			vars.ErrorKey, c.Value(vars.ErrorKey),
		)
	}
}
