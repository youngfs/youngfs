package ec_server

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/go-playground/assert/v2"
	"io"
	"math/rand"
	"net/http"
	"os"
	"testing"
	"time"
	"youngfs/errors"
	"youngfs/fs/ec/ec_calc"
	"youngfs/fs/ec/ec_store"
	"youngfs/fs/entry"
	"youngfs/fs/full_path"
	fs_set "youngfs/fs/set"
	"youngfs/fs/storage_engine/seaweedfs"
	"youngfs/kv/redis"
	"youngfs/log"
	"youngfs/util"
	"youngfs/vars"
)

func TestECServer_Backup(t *testing.T) {
	vars.Debug = true
	log.InitLogger()
	defer func() {
		_ = os.Remove(vars.ServerName + ".log")
		_ = os.RemoveAll(ec_calc.EcFilePrefix)
	}()
	defer log.Sync()

	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewEC(kvStore, se)
	ecCalc := ec_calc.NewECCalc(ecStore, se)
	client := NewECServer(ecStore, ecCalc)

	ctx := context.Background()

	hosts, err := se.GetHosts(ctx)
	assert.Equal(t, err, nil)
	set := fs_set.Set("ec_test")
	ec_calc.EcBlockSize = 16 * 1024 * 1024
	size := uint64(40 * 1024 * 1024)

	if len(hosts) < 2 {
		fmt.Printf("Can't do backup unit test")
		return
	}

	setRules := &fs_set.SetRules{
		Set:             set,
		Hosts:           hosts,
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    0,
		ECMode:          false,
		ReplicationMode: true,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	for i := 0; i < 4; i++ {
		ent := &entry.Entry{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set,
			FileSize: size,
		}

		b := util.RandByte(size)
		file := io.Reader(bytes.NewReader(b))
		md5Hash := md5.New()
		file = io.TeeReader(file, md5Hash)

		host, ecid, err := client.InsertObject(ctx, ent)
		assert.Equal(t, err, nil)

		fid, err := se.PutObject(ctx, size, file, "", true, host)
		assert.Equal(t, err, nil)

		ent.Fid = fid
		ent.ECid = ecid
		ent.Md5 = md5Hash.Sum(nil)

		err = client.ConfirmEC(ctx, ent)
		assert.Equal(t, err, nil)

		err = client.ExecEC(ctx, ent.ECid)
		assert.Equal(t, err, nil)

		time.Sleep(3 * time.Second)

		url, err := se.GetFidUrl(ctx, fid)
		assert.Equal(t, err, nil)

		err = se.DeleteObject(ctx, fid)
		assert.Equal(t, err, nil)

		time.Sleep(3 * time.Second)

		resp1, err := http.Get(url)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp1.StatusCode, http.StatusNotFound)
		defer func() {
			_ = resp1.Body.Close()
		}()

		frag, err := client.RecoverObject(ctx, ent)
		assert.Equal(t, err, nil)

		url, err = se.GetFidUrl(ctx, frag[0].Fid)
		assert.Equal(t, err, nil)

		resp2, err := http.Get(url)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp2.StatusCode, http.StatusOK)
		defer func() {
			_ = resp2.Body.Close()
		}()

		httpBody, err := io.ReadAll(resp2.Body)
		assert.Equal(t, err, nil)
		assert.Equal(t, httpBody, b)

		err = se.DeleteObject(ctx, frag[0].Fid)
		assert.Equal(t, err, nil)

		err = client.DeleteObject(ctx, ent.ECid)
		assert.Equal(t, err, nil)
	}

	err = client.DeleteSetRules(ctx, set)
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)
}

func TestECServer_NoEC(t *testing.T) {
	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewEC(kvStore, se)
	ecCalc := ec_calc.NewECCalc(ecStore, se)
	client := NewECServer(ecStore, ecCalc)

	ctx := context.Background()

	hosts, err := se.GetHosts(ctx)
	assert.Equal(t, err, nil)
	set := fs_set.Set("ec_test")
	ec_calc.EcBlockSize = 16 * 1024 * 1024
	size := uint64(40 * 1024 * 1024)

	setRules := &fs_set.SetRules{
		Set:             set,
		Hosts:           hosts[:1],
		DataShards:      0,
		ParityShards:    0,
		MAXShardSize:    0,
		ECMode:          false,
		ReplicationMode: false,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	for i := 0; i < 16; i++ {
		ent := &entry.Entry{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set,
			FileSize: size,
		}

		b := util.RandByte(size)
		file := io.Reader(bytes.NewReader(b))
		md5Hash := md5.New()
		file = io.TeeReader(file, md5Hash)

		host, ecid, err := client.InsertObject(ctx, ent)
		assert.Equal(t, err, nil)
		assert.Equal(t, ecid, "")
		assert.Equal(t, host, hosts[0])

		fid, err := se.PutObject(ctx, size, file, "", true, host)
		assert.Equal(t, err, nil)

		ent.Fid = fid
		ent.ECid = ecid
		ent.Md5 = md5Hash.Sum(nil)

		err = client.ConfirmEC(ctx, ent)
		assert.Equal(t, err, nil)

		err = client.ExecEC(ctx, ent.ECid)
		assert.Equal(t, err, nil)

		time.Sleep(1 * time.Second)

		url, err := se.GetFidUrl(ctx, fid)
		assert.Equal(t, err, nil)

		err = se.DeleteObject(ctx, fid)
		assert.Equal(t, err, nil)

		time.Sleep(3 * time.Second)

		resp1, err := http.Get(url)
		assert.Equal(t, err, nil)
		assert.Equal(t, resp1.StatusCode, http.StatusNotFound)
		defer func() {
			_ = resp1.Body.Close()
		}()

		frag, err := client.RecoverObject(ctx, ent)
		assert.Equal(t, errors.Is(err, errors.ErrRecoverFailed), true)
		assert.Equal(t, frag, nil)
	}

	err = client.DeleteSetRules(ctx, set)
	assert.Equal(t, err, nil)
}

func TestECServer_ReedSolomon(t *testing.T) {
	vars.Debug = true
	log.InitLogger()
	defer func() {
		_ = os.Remove(vars.ServerName + ".log")
		_ = os.RemoveAll(ec_calc.EcFilePrefix)
	}()
	defer log.Sync()

	kvStore := redis.NewKvStore(vars.RedisSocket, vars.RedisPassword, vars.RedisDatabase)
	se := seaweedfs.NewStorageEngine(vars.SeaweedFSMaster, kvStore)
	ecStore := ec_store.NewEC(kvStore, se)
	ecCalc := ec_calc.NewECCalc(ecStore, se)
	client := NewECServer(ecStore, ecCalc)

	ctx := context.Background()

	hosts, err := se.GetHosts(ctx)
	assert.Equal(t, err, nil)
	set := fs_set.Set("ec_test")
	ec_calc.EcBlockSize = 16 * 1024 * 1024
	size := uint64(40 * 1024 * 1024)

	if len(hosts) < 3 {
		fmt.Printf("Can't do reed solomon unit test")
		return
	}

	setRules := &fs_set.SetRules{
		Set:             set,
		Hosts:           hosts[:3],
		DataShards:      2,
		ParityShards:    1,
		MAXShardSize:    96 * 1024 * 1024,
		ECMode:          true,
		ReplicationMode: true,
	}

	err = client.InsertSetRules(ctx, setRules)
	assert.Equal(t, err, nil)

	ents := make([]*entry.Entry, 9)
	for i := 0; i < 5; i++ {
		sz := size - 128 + uint64(rand.Intn(256))
		ent := &entry.Entry{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set,
			FileSize: sz,
		}

		b := util.RandByte(sz)
		file := io.Reader(bytes.NewReader(b))
		md5Hash := md5.New()
		file = io.TeeReader(file, md5Hash)

		host, ecid, err := client.InsertObject(ctx, ent)
		assert.Equal(t, err, nil)

		fid, err := se.PutObject(ctx, size, file, "", true, host)
		assert.Equal(t, err, nil)

		ent.Fid = fid
		ent.ECid = ecid
		ent.Md5 = md5Hash.Sum(nil)

		err = client.ConfirmEC(ctx, ent)
		assert.Equal(t, err, nil)

		err = client.ExecEC(ctx, ent.ECid)
		assert.Equal(t, err, nil)

		ents[i] = ent
	}

	time.Sleep(60 * time.Second)

	url, err := se.GetFidUrl(ctx, ents[0].Fid)
	assert.Equal(t, err, nil)

	err = se.DeleteObject(ctx, ents[0].Fid)
	assert.Equal(t, err, nil)

	err = se.DeleteObject(ctx, ents[0].Fid) // because do reed solomon add link, so delete twice
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp1, err := http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp1.StatusCode, http.StatusNotFound)
	defer func() {
		_ = resp1.Body.Close()
	}()

	frags, err := client.RecoverObject(ctx, ents[0])
	assert.Equal(t, err, nil)

	for _, frag := range frags {
		if frag.OldECid == ents[0].ECid {
			url, err := se.GetFidUrl(ctx, frag.Fid)
			assert.Equal(t, err, nil)

			resp, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.StatusCode, http.StatusOK)
			defer func() {
				_ = resp.Body.Close()
			}()

			httpBody, err := io.ReadAll(resp.Body)
			assert.Equal(t, err, nil)

			md5Hash := md5.New()
			md5Hash.Write(httpBody)
			md5Ret := md5Hash.Sum(nil)
			assert.Equal(t, md5Ret, ents[0].Md5)

			ents[0].Fid = frag.Fid
		}
	}

	url, err = se.GetFidUrl(ctx, ents[0].Fid)
	assert.Equal(t, err, nil)

	err = se.DeleteObject(ctx, ents[0].Fid)
	assert.Equal(t, err, nil)

	err = se.DeleteObject(ctx, ents[0].Fid) // because do reed solomon add link, so delete twice
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp2, err := http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp2.StatusCode, http.StatusNotFound)
	defer func() {
		_ = resp2.Body.Close()
	}()

	frags, err = client.RecoverObject(ctx, ents[0])
	assert.Equal(t, err, nil)
	assert.Equal(t, len(frags), 2)

	for _, frag := range frags {
		if frag.OldECid == ents[0].ECid {
			url, err := se.GetFidUrl(ctx, frag.Fid)
			assert.Equal(t, err, nil)

			resp, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.StatusCode, http.StatusOK)
			defer func() {
				_ = resp.Body.Close()
			}()

			httpBody, err := io.ReadAll(resp.Body)
			assert.Equal(t, err, nil)

			md5Hash := md5.New()
			md5Hash.Write(httpBody)
			md5Ret := md5Hash.Sum(nil)
			assert.Equal(t, md5Ret, ents[0].Md5)

			ents[0].Fid = frag.Fid
		}
	}

	for _, frag := range frags {
		if frag.OldECid == ents[1].ECid {
			url, err := se.GetFidUrl(ctx, frag.Fid)
			assert.Equal(t, err, nil)

			resp, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.StatusCode, http.StatusOK)
			defer func() {
				_ = resp.Body.Close()
			}()

			httpBody, err := io.ReadAll(resp.Body)
			assert.Equal(t, err, nil)

			md5Hash := md5.New()
			md5Hash.Write(httpBody)
			md5Ret := md5Hash.Sum(nil)
			assert.Equal(t, md5Ret, ents[1].Md5)

			ents[1].Fid = frag.Fid
		}
	}

	for i := 5; i < 9; i++ {
		sz := size - 128 + uint64(rand.Intn(256))
		ent := &entry.Entry{
			FullPath: full_path.FullPath(util.RandString(16)),
			Set:      set,
			FileSize: sz,
		}

		b := util.RandByte(sz)
		file := io.Reader(bytes.NewReader(b))
		md5Hash := md5.New()
		file = io.TeeReader(file, md5Hash)

		host, ecid, err := client.InsertObject(ctx, ent)
		assert.Equal(t, err, nil)

		fid, err := se.PutObject(ctx, size, file, "", true, host)
		assert.Equal(t, err, nil)

		ent.Fid = fid
		ent.ECid = ecid
		ent.Md5 = md5Hash.Sum(nil)

		err = client.ConfirmEC(ctx, ent)
		assert.Equal(t, err, nil)

		ents[i] = ent

		err = client.ExecEC(ctx, ent.ECid)
		assert.Equal(t, err, nil)

		if i == 5 {
			time.Sleep(30 * time.Second)

			url, err = se.GetFidUrl(ctx, ents[i].Fid)
			assert.Equal(t, err, nil)

			err = se.DeleteObject(ctx, ents[i].Fid)
			assert.Equal(t, err, nil)

			time.Sleep(3 * time.Second)

			url2, err := se.GetFidUrl(ctx, ents[i].Fid)
			assert.Equal(t, errors.Is(err, errors.ErrObjectNotExist), true)
			assert.Equal(t, url2, "")

			resp1, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp1.StatusCode, http.StatusNotFound)
			defer func() {
				_ = resp1.Body.Close()
			}()

			frags, err = client.RecoverObject(ctx, ents[i])
			assert.Equal(t, err, nil)

			url, err = se.GetFidUrl(ctx, frags[0].Fid)
			assert.Equal(t, err, nil)

			resp2, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp2.StatusCode, http.StatusOK)
			defer func() {
				_ = resp2.Body.Close()
			}()

			b2, err := io.ReadAll(resp2.Body)
			assert.Equal(t, err, nil)

			md5Hash2 := md5.New()
			md5Hash2.Write(b2)
			assert.Equal(t, md5Hash2.Sum(nil), md5Hash.Sum(nil))

			ents[i].Fid = frags[0].Fid
		}
	}

	time.Sleep(60 * time.Second)

	url, err = se.GetFidUrl(ctx, ents[6].Fid)
	assert.Equal(t, err, nil)

	err = se.DeleteObject(ctx, ents[6].Fid)
	assert.Equal(t, err, nil)

	err = se.DeleteObject(ctx, ents[6].Fid) // because do reed solomon add link, so delete twice
	assert.Equal(t, err, nil)

	time.Sleep(3 * time.Second)

	resp3, err := http.Get(url)
	assert.Equal(t, err, nil)
	assert.Equal(t, resp3.StatusCode, http.StatusNotFound)
	defer func() {
		_ = resp3.Body.Close()
	}()

	frags, err = client.RecoverObject(ctx, ents[6])
	assert.Equal(t, err, nil)
	assert.Equal(t, len(frags), 2)

	for _, frag := range frags {
		if frag.OldECid == ents[6].ECid {
			url, err := se.GetFidUrl(ctx, frag.Fid)
			assert.Equal(t, err, nil)

			resp, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.StatusCode, http.StatusOK)
			defer func() {
				_ = resp.Body.Close()
			}()

			httpBody, err := io.ReadAll(resp.Body)
			assert.Equal(t, err, nil)

			md5Hash := md5.New()
			md5Hash.Write(httpBody)
			md5Ret := md5Hash.Sum(nil)
			assert.Equal(t, md5Ret, ents[6].Md5)

			ents[6].Fid = frag.Fid
		}
	}

	for _, frag := range frags {
		if frag.OldECid == ents[7].ECid {
			url, err := se.GetFidUrl(ctx, frag.Fid)
			assert.Equal(t, err, nil)

			resp, err := http.Get(url)
			assert.Equal(t, err, nil)
			assert.Equal(t, resp.StatusCode, http.StatusOK)
			defer func() {
				_ = resp.Body.Close()
			}()

			httpBody, err := io.ReadAll(resp.Body)
			assert.Equal(t, err, nil)

			md5Hash := md5.New()
			md5Hash.Write(httpBody)
			md5Ret := md5Hash.Sum(nil)
			assert.Equal(t, md5Ret, ents[7].Md5)

			ents[7].Fid = frag.Fid
		}
	}

	err = client.DeleteSetRules(ctx, set)
	assert.Equal(t, err, nil)

	for i := 0; i < 8; i++ {
		err := se.DeleteObject(ctx, ents[i].Fid)
		assert.Equal(t, err, nil)

		err = client.DeleteObject(ctx, ents[i].ECid)
		assert.Equal(t, err, nil)
	}

	err = se.DeleteObject(ctx, ents[8].Fid)
	assert.Equal(t, err, nil)

	err = client.DeleteObject(ctx, ents[8].ECid)
	assert.Equal(t, err, nil)

	time.Sleep(10 * time.Second)
}
