package storage_engine

import (
	jsoniter "github.com/json-iterator/go"
	"icesos/command"
	"io/ioutil"
	"net/http"
)

type AssignFileInfo struct {
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

func AssignFileHandler() (*AssignFileInfo, error) {
	resp, err := http.Get("http://" + command.MasterServer + "/dir/assign")
	if err != nil {
		return nil, err
	}

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	assignFileInfo := &AssignFileInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		return nil, err
	}

	return assignFileInfo, err
}
