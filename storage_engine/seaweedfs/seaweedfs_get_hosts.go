package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"icesos/errors"
	"io/ioutil"
	"net/http"
)

type dataNode struct {
	EcShards  uint64 `json:"EcShards"`
	Max       uint64 `json:"Max"`
	PublicUrl string `json:"PublicUrl"`
	Url       string `json:"Url"`
	VolumeIds string `json:"VolumeIds"`
	Volumes   uint64 `json:"Volumes"`
}

type rack struct {
	Id        string     `json:"Id"`
	DataNodes []dataNode `json:"DataNodes"`
}

type dataCenter struct {
	Id    string `json:"Id"`
	Racks []rack `json:"Racks"`
}

type topology struct {
	DataCenters []dataCenter `json:"DataCenters"`
}

type dirStatue struct {
	Topology topology `json:"Topology"`
}

func (se *StorageEngine) GetHosts(ctx context.Context) ([]string, error) {
	resp, err := http.Get("http://" + se.masterServer + "//dir/status")
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrorCodeResponse[errors.ErrSeaweedFSMaster]
	}

	info := &dirStatue{}
	err = jsoniter.Unmarshal(httpBody, info)

	dataCenters := info.Topology.DataCenters
	var racks []rack
	var dataNodes []dataNode
	set := make(map[string]bool)
	ret := make([]string, 0)
	for _, u := range dataCenters {
		if u.Id == "DefaultDataCenter" {
			racks = u.Racks
			break
		}
	}
	for _, u := range racks {
		if u.Id == "DefaultRack" {
			dataNodes = u.DataNodes
			break
		}
	}
	for _, u := range dataNodes {
		if !set[u.Url] {
			ret = append(ret, u.Url)
			set[u.Url] = true
		}
		if !set[u.PublicUrl] {
			ret = append(ret, u.PublicUrl)
			set[u.PublicUrl] = true
		}
	}

	return ret, nil
}
