package routers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/youngfs/youngfs/log"
	"github.com/youngfs/youngfs/vars"
	"net/http"
	"strconv"
)

func authorizationHeader(user, pw string) string {
	base := user + ":" + pw
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
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

	admin := authorizationHeader("young", "young")

	return func(c *gin.Context) {
		if admin != c.Request.Header.Get("Authorization") {
			// Credentials doesn't match, we return 401 and abort handlers chain.
			c.Header("WWW-Authenticate", realm)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// The user credentials was found, set user's id to key AuthUserKey in this context, the user's id can be read later using
		// c.MustGet(gin.AuthUserKey).
		c.Set(vars.UserKey, "young")
	}
}
func Logger(reqs ...string) gin.HandlerFunc {
	req := ""
	for _, u := range reqs {
		req += u
	}

	return func(c *gin.Context) {
		c.Set(vars.UUIDKey, uuid.New().String())
		log.Infow("request start",
			vars.UUIDKey, c.Value(vars.UUIDKey),
			requestNameKey, req,
			requestUrlKey, c.Request.URL,
			vars.UserKey, c.Value(vars.UserKey),
		)
		c.Next()
		log.Infow("request finish",
			vars.UUIDKey, c.Value(vars.UUIDKey),
			requestNameKey, req,
			vars.UserKey, c.Value(vars.UserKey),
			responseHTTPCodeKey, c.Writer.Status(),
			vars.CodeKey, c.Value(vars.CodeKey),
			vars.ErrorKey, c.Value(vars.ErrorKey),
		)
	}
}
