package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/bucket"
	"youngfs/fs/fullpath"
	"youngfs/fs/server"
	"youngfs/fs/ui"
	"youngfs/vars"
)

func ListObjectsHandler(c *gin.Context) {
	bkt, fp := bucket.Bucket(c.Param("bucket")), fullpath.FullPath(c.Param("path"))
	if !bkt.IsLegal() {
		errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if fp != "/" && fp[len(fp)-1] == '/' {
		fp = fp[:len(fp)-1]
	}
	if !fp.IsLegal() {
		errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	ents, err := server.ListObjects(c, bkt, fp)
	if err != nil {
		errorHandler(c, err)
		return
	}

	accepts := c.Request.Header.Values("Accept")
	for _, str := range accepts {
		if str == "application/json" {
			c.JSON(
				http.StatusOK,
				gin.H{
					vars.UUIDKey: c.Value(vars.UUIDKey),
					"Path":       fp,
					"Bucket":     bkt,
					"Entries":    ents,
				},
			)
			return
		}
	}

	c.HTML(
		http.StatusOK,
		ui.FSName,
		gin.H{
			"FullPath":  string(fp),
			"Bucket":    string(bkt),
			"PathLinks": fp.ToPathLink(),
			"Entries":   ents,
		},
	)
	return
}
