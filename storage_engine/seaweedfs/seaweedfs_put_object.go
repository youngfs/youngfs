package seaweedfs

import (
	"bytes"
	"context"
	"fmt"
	"icesos/command/vars"
	"icesos/errors"
	"icesos/log"
	"icesos/util"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
)

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (se *StorageEngine) PutObject(ctx context.Context, size uint64, file io.Reader, fileName string, compress bool, hosts ...string) (string, error) {

	info, err := se.assignObject(ctx, size, hosts...)
	if err != nil {
		log.Errorw("seaweedfs put object: assign object error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "size", size, "hosts", hosts)
		return "", err
	}

	var multiReader io.Reader

	buf := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, util.EscapeQuotes("file"), util.EscapeQuotes(fileName)))
	h.Set("Content-Type", "application/octet-stream")
	if compress {
		h.Set("Content-Encoding", "gzip")
	}

	formFile, err := multiWriter.CreatePart(h)
	if err != nil {
		log.Errorw("seaweedfs put object: get form file", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+info.Url+"/"+info.Fid)
		return "", errors.GetAPIErr(errors.ErrServer)
	}
	if compress {
		_, err := se.gzipWriterPool.GzipStream(formFile, file)
		if err != nil {
			log.Errorw("seaweedfs put object: gzip copy", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "size", size, "hosts", hosts)
			return "", errors.GetAPIErr(errors.ErrServer)
		}

		// add multiWriter.Boundary()
		err = multiWriter.Close()
		if err != nil {
			log.Errorw("seaweedfs put object: multi writer close", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "size", size, "hosts", hosts)
			return "", err
		}

		multiReader = buf
	} else {
		// add multiWriter.Boundary()
		end := strings.NewReader("\r\n--" + multiWriter.Boundary() + "--\r\n")
		multiReader = io.MultiReader(buf, file, end)
	}

	req, err := http.NewRequest("POST", "http://"+info.Url+"/"+info.Fid, multiReader)
	if err != nil {
		log.Errorw("seaweedfs put object: new request put error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+info.Url+"/"+info.Fid, "request", req)
		return "", errors.GetAPIErr(errors.ErrServer)
	}
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorw("seaweedfs put object: do request put error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), vars.ErrorKey, err.Error(), "request url", "http://"+info.Url+"/"+info.Fid, "request", req, "response", resp)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		log.Errorw("seaweedfs put object: request error", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "request url", "http://"+info.Url+"/"+info.Fid, "http code", resp.StatusCode, "request", req, "response", resp)
		return "", errors.GetAPIErr(errors.ErrSeaweedFSVolume)
	}

	err = se.AddLink(ctx, info.Fid)
	if err != nil {
		log.Errorw("seaweedfs put object: add link", vars.UUIDKey, ctx.Value(vars.UUIDKey), vars.UserKey, ctx.Value(vars.UserKey), "request url", "http://"+info.Url+"/"+info.Fid, "http code", resp.StatusCode, "request", req, "response", resp)
		return "", err
	}

	return info.Fid, nil
}
