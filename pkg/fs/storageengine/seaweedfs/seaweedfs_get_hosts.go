package seaweedfs

import (
	"context"
	jsoniter "github.com/json-iterator/go"
	"github.com/youngfs/youngfs/pkg/log"
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

func (se *StorageEngine) updateHosts() {
	resp, err := http.Get("http://" + se.masterServer + "/dir/status")
	if err != nil {
		log.Errorw("seaweedfs get hosts : http get error")
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorw("seaweedfs get hosts: get http body error")
		return
	}

	info := &dirStatue{}
	err = jsoniter.Unmarshal(httpBody, info)
	if err != nil {
		log.Errorw("seaweedfs get hosts: http body unmarshal error")
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

	se.hostsMutex.Lock()
	defer se.hostsMutex.Unlock()
	se.hosts = ret
}

func (se *StorageEngine) scheduledUpdateHosts() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		se.updateHosts()
	}
}

func (se *StorageEngine) GetHosts(ctx context.Context) ([]string, error) {
	se.hostsMutex.RLock()
	defer se.hostsMutex.RUnlock()
	ret := make([]string, len(se.hosts))
	copy(ret, se.hosts)
	return ret, nil
}
