package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"time"
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

func (e *Engine) updateEndpoints() {
	resp, err := http.Get("http://" + e.masterEndpoint + "/dir/status")
	if err != nil {
		e.logger.Errorf("seaweedfs update endpoints: %s", err.Error())
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		e.logger.Errorf("seaweedfs update endpoints: %s", err.Error())
		return
	}

	info := &dirStatue{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		e.logger.Errorf("seaweedfs update endpoints: %s", err.Error())
		return
	}

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

	e.endpointsMutex.Lock()
	defer e.endpointsMutex.Unlock()
	e.endpoints = ret
}

func (e *Engine) scheduledUpdateEndpoints() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		e.updateEndpoints()
	}
}

func (e *Engine) GetEndpoints(ctx context.Context) ([]string, error) {
	e.endpointsMutex.RLock()
	defer e.endpointsMutex.RUnlock()
	ret := make([]string, len(e.endpoints))
	copy(ret, e.endpoints)
	return ret, nil
}
