package storageEngine

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"object-storage-server/command"
)

type AssignFileInfo struct {
	Fid       string `json:"fid"`
	Url       string `json:"url"`
	PublicUrl string `json:"publicUrl"`
	Count     int64  `json:"count"`
}

func AssignFileHandler() (assignFileInfo *AssignFileInfo, err error) {
	resp, err := http.Get("http://" + command.MasterServer + "/dir/assign")
	if err != nil {
		return
	}

	httpBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	assignFileInfo = &AssignFileInfo{}
	err = jsoniter.Unmarshal(httpBody, assignFileInfo)
	if err != nil {
		return
	}

	return
}
