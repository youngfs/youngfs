package seaweedfs

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/youngfs/youngfs/errors"
	"io"
	"net/http"
)

func (se *StorageEngine) GetObject(ctx context.Context, fid string, writer io.Writer) error {
	url, err := se.getFidUrl(fid)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.ErrServer.Wrap("seaweedfs get object: new request")
	}
	req.Header.Add("Accept-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.ErrSeaweedFSVolume.Wrapf("seaweedfs get object: request get error: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			return errors.ErrObjectNotExist
		}
		return errors.ErrSeaweedFSVolume.Wrapf("seaweedfs get object: response http code: %d", resp.StatusCode)
	}

	var reader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
	default:
		reader = resp.Body
	}
	defer func() {
		_ = reader.Close()
	}()

	_, err = io.Copy(writer, reader)
	if err != nil {
		return errors.ErrServer.WrapErr(err)
	}

	return nil
}

func (se *StorageEngine) getFidUrl(fid string) (string, error) {
	volumeId, _, err := se.parseFid(fid)
	if err != nil {
		return "", err
	}

	host, err := se.getVolumeHost(volumeId)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("http://%s/%s", host, fid), nil
}
