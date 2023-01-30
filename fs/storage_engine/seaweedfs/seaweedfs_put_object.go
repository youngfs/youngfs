package seaweedfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strings"
	"youngfs/errors"
)

type PutObjectInfo struct {
	Size uint64 `json:"size"`
	ETag string `json:"eTag"`
}

func (se *StorageEngine) PutObject(ctx context.Context, size uint64, file io.Reader, fileName string, compress bool, hosts ...string) (string, error) {
	info, err := se.assignObject(ctx, size, hosts...)
	if err != nil {
		return "", err
	}

	var multiReader io.Reader

	buf := &bytes.Buffer{}
	multiWriter := multipart.NewWriter(buf)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, escapeQuotes("file"), escapeQuotes(fileName)))
	h.Set("Content-Type", "application/octet-stream")
	if compress {
		h.Set("Content-Encoding", "gzip")
	}

	formFile, err := multiWriter.CreatePart(h)
	if err != nil {
		return "", errors.ErrServer.Wrap("seaweedfs put object: get form file")
	}
	if compress {
		_, err := se.gzipWriterPool.GzipStream(formFile, file)
		if err != nil {
			return "", errors.ErrServer.Wrap("seaweedfs put object: gzip copy")
		}

		// add multiWriter.Boundary()
		err = multiWriter.Close()
		if err != nil {
			return "", errors.ErrServer.Wrap("seaweedfs put object: multi writer close")
		}

		multiReader = buf
	} else {
		// add multiWriter.Boundary()
		end := strings.NewReader("\r\n--" + multiWriter.Boundary() + "--\r\n")
		multiReader = io.MultiReader(buf, file, end)
	}

	req, err := http.NewRequest("POST", "http://"+info.Url+"/"+info.Fid, multiReader)
	if err != nil {
		return "", errors.ErrServer.Wrap("seaweedfs put object: new request put error")
	}
	req.Header.Set("Content-Type", multiWriter.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.ErrServer.Wrap("seaweedfs put object: do request put error")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusCreated {
		return "", errors.ErrSeaweedFSVolume.Wrap("seaweedfs put object: request error")
	}

	err = se.AddLink(ctx, info.Fid)
	if err != nil {
		return "", err
	}

	return info.Fid, nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
