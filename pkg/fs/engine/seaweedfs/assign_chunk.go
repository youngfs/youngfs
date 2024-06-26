package seaweedfs

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/youngfs/youngfs/pkg/errors"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type assignChunkInfo struct {
	Fid       string `json:"fid"` //Fid = VolumeId,Fid
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

func (e *Engine) assignChunk(ctx context.Context, endpoints ...string) (*assignChunkInfo, error) {
	endpointReq, endpoint := "", ""
	if len(endpoints) > 0 {
		endpoint = endpoints[rand.Intn(len(endpoints))]
		endpointReq = "dataCenter=DefaultDataCenter&rack=DefaultRack&dataNode=" + endpoint
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("http://%s/dir/assign?%s", e.masterEndpoint, endpointReq), nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.ErrEngineMaster.WarpErr(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrEngineMaster.WarpErr(err)
	}

	assignFileInfo := &assignChunkInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		return nil, errors.ErrEngineMaster.WarpErr(err)
	}

	if endpoint != "" {
		if endpoint != assignFileInfo.Url && endpoint != assignFileInfo.PublicUrl {
			return nil, errors.ErrEngineMaster.WithMessagef("seaweedfs assign object: endpoint not match, want %s, got %s", endpoint, assignFileInfo.Url)
		}
	}

	return assignFileInfo, nil
}

func (e *Engine) parseFid(fullFid string) (uint64, string, error) {
	ret := strings.Split(fullFid, ",")
	if len(ret) != 2 {
		return 0, "", errors.ErrEngineMaster.WithMessagef("seaweedfs parse fid error: %s", fullFid)
	}
	volumeId, err := strconv.ParseUint(ret[0], 10, 64)
	if err != nil {
		return 0, "", errors.ErrEngineMaster.WithMessagef("seaweedfs parse fid error: %s", fullFid)
	}
	return volumeId, ret[1], nil
}
