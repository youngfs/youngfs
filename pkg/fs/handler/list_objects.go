package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/youngfs/youngfs/pkg/errors"
	"github.com/youngfs/youngfs/pkg/fs/bucket"
	"github.com/youngfs/youngfs/pkg/fs/fullpath"
	"github.com/youngfs/youngfs/pkg/fs/ui"
	"net/http"
)

func (h *Handler) ListObjectsHandler(c *gin.Context) {
	bkt, fp := bucket.Bucket(c.Param(QueryBucketKey)), fullpath.FullPath(c.Param(QueryPathKey))
	if !bkt.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalBucketName)
		return
	}
	if fp != "/" && fp[len(fp)-1] == '/' {
		fp = fp[:len(fp)-1]
	}
	if !fp.IsLegal() {
		h.errorHandler(c, errors.ErrIllegalObjectName)
		return
	}
	fp = fp.Clean()

	ents, err := h.svr.ListObjects(c, bkt, fp)
	if err != nil {
		h.errorHandler(c, err)
		return
	}

	accepts := c.Request.Header.Values("Accept")
	for _, str := range accepts {
		if str == "application/json" {
			c.JSON(
				http.StatusOK,
				gin.H{
					UUIDKey:   c.Value(UUIDKey),
					"Path":    fp,
					"Bucket":  bkt,
					"Entries": ents,
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
