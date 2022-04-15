package routers

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"icesos/util"
	"net/http"
	"strconv"
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
		c.Set(UserKey, "ices")
	}
}
