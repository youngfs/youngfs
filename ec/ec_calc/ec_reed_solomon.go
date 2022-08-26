package ec_calc

import (
	"bytes"
	"context"
	"crypto/md5"
	"github.com/klauspost/reedsolomon"
	"hash"
	"icesfs/command/vars"
	"icesfs/ec/ec_store"
	"icesfs/errors"
	"icesfs/log"
	"icesfs/util"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

func (calc *ECCalc) reedSolomon(ctx context.Context, suite *ec_store.Suite) error {
	fragCnt := 0 // effective frag

	// init len(suite.Shards) == setRules.DataShards + setRules.ParityShards
	for i, shard := range suite.Shards {
		for j, frag := range shard.Frags {
			_, err := calc.storageEngine.GetFidUrl(ctx, frag.Fid)
			if err != nil {
				fragSuite, err := calc.ECStore.GetSuite(ctx, frag.OldECid)
				if err == nil {
					_, err := calc.storageEngine.GetFidUrl(ctx, fragSuite.BakFid)
					if err == nil {
						suite.Shards[i].Frags[j].Fid = fragSuite.BakFid
						_ = calc.storageEngine.AddLink(ctx, fragSuite.BakFid)
						continue
					}
				}
				suite.Shards[i].Frags[j].Fid = ""
				suite.Shards[i].Frags[j].OldECid = ""
				suite.Shards[i].Frags[j].FileSize = 0
				continue
			} else {
				_ = calc.storageEngine.AddLink(ctx, frag.Fid)
				fragCnt++
			}
		}
	}

	size := make([]uint64, len(suite.Shards))
	for i, shard := range suite.Shards {
		for _, frag := range shard.Frags {
			size[i] += frag.FileSize
		}
	}

	mx := uint64(0)
	for _, u := range size {
		mx = util.Max(mx, u)
	}

	md5Hashs := make([]hash.Hash, len(suite.Shards))
	for i := range md5Hashs {
		md5Hashs[i] = md5.New()
	}

	ecReadClosers := make([]io.ReadCloser, suite.DataShards)
	ecReaders := make([]io.Reader, suite.DataShards)
	for i := 0; i < int(suite.DataShards); i++ {
		ecReadClosers[i] = NewECReadCloser(suite.Shards[i].Frags, calc.storageEngine)
		ecReaders[i] = io.TeeReader(ecReadClosers[i], md5Hashs[i])
	}
	defer func() {
		for i := range ecReadClosers {
			_ = ecReadClosers[i].Close()
		}
	}()

	enc, err := reedsolomon.New(int(suite.DataShards), len(suite.Shards)-int(suite.DataShards))
	if err != nil {
		log.Errorw("reed solomon : new reed solomon", vars.ErrorKey, err, "ecid", suite.ECid)
		return errors.GetAPIErr(errors.ErrServer)
	}

	for i := uint64(0); i < mx; i += ecBlockSize {
		length := util.Min(ecBlockSize, mx-i)
		data := make([][]byte, len(suite.Shards))
		for i := range data {
			data[i] = make([]byte, length)
		}
		for i := range ecReaders {
			_, _ = ecReaders[i].Read(data[i])
		}
		err := enc.Encode(data)
		if err != nil {
			return errors.GetAPIErr(errors.ErrServer)
		}

		for i := int(suite.DataShards); i < len(suite.Shards); i++ {
			file := io.TeeReader(bytes.NewReader(data[i]), md5Hashs[i])
			fid, err := calc.storageEngine.PutObject(ctx, length, file, "", true, suite.Shards[i].Host)
			if err != nil {
				log.Errorw("reed solomon : put object", vars.ErrorKey, err, "ecid", suite.ECid)
				return err
			}
			suite.Shards[i].Frags = append(suite.Shards[i].Frags, ec_store.Frag{
				Fid:      fid,
				FileSize: length,
			})
		}
	}

	for i := range md5Hashs {
		suite.Shards[i].Md5 = md5Hashs[i].Sum(nil)
	}

	err = calc.ECStore.InsertSuite(ctx, suite)
	if err != nil {
		log.Errorw("reed solomon : insert suite", vars.ErrorKey, err, "ecid", suite.ECid)
		return err
	}

	err = calc.ECStore.SetECidLink(ctx, suite.ECid, int64(fragCnt))
	if err != nil {
		log.Errorw("reed solomon : set ecid link", vars.ErrorKey, err, "ecid", suite.ECid)
		return errors.GetAPIErr(errors.ErrServer)
	}

	for _, shard := range suite.Shards {
		for _, frag := range shard.Frags {
			if frag.OldECid != "" {
				oldSuite, err := calc.ECStore.GetSuite(ctx, frag.OldECid)
				if err != nil {
					log.Errorw("reed solomon : get suite", vars.ErrorKey, err, "ecid", suite.ECid)
					return err
				}
				if oldSuite.BakFid != "" {
					// delete backup
					err = calc.storageEngine.DeleteObject(ctx, oldSuite.BakFid)
				}
				oldSuite.Next = suite.ECid
				err = calc.ECStore.InsertSuite(ctx, oldSuite)
				if err != nil {
					log.Errorw("reed solomon : insert suite", vars.ErrorKey, err, "ecid", suite.ECid)
					return err
				}
			}
		}
	}

	log.Infow("reed solomon successful", "ecid", suite.ECid)

	return nil
}

func ecFileKey(fid string, cnt int) string {
	return ecFilePrefix + fid + "_" + strconv.Itoa(cnt)
}

func (calc *ECCalc) reedSolomonRecover(ctx context.Context, suite *ec_store.Suite) ([]ec_store.Frag, error) {
	ecReadClosers := make([]io.ReadCloser, len(suite.Shards))
	for i, shard := range suite.Shards {
		ecReadClosers[i] = NewECReadCloser(shard.Frags, calc.storageEngine)
	}
	defer func() {
		for i := range ecReadClosers {
			_ = ecReadClosers[i].Close()
		}
	}()

	size := make([]uint64, len(suite.Shards))
	for i, shard := range suite.Shards {
		for _, frag := range shard.Frags {
			size[i] += frag.FileSize
		}
	}

	mx := uint64(0)
	for _, u := range size {
		mx = util.Max(mx, u)
	}

	corruptionMap := make(map[int]bool)
	corruptionList := make([]int, 0)
	for i, shard := range suite.Shards {
		md5Hash := md5.New()
		reader := io.TeeReader(ecReadClosers[i], md5Hash)
		for i := uint64(0); i < mx; i += ecBlockSize {
			length := util.Min(ecBlockSize, mx-i)
			data := make([]byte, length)
			_, _ = reader.Read(data)
		}

		md5Ret := md5Hash.Sum(nil)
		if !util.BytesIsEqual(md5Ret, shard.Md5) {
			corruptionMap[i] = true
			corruptionList = append(corruptionList, i)
		}
	}

	if len(corruptionList) > len(suite.Shards)-int(suite.DataShards) {
		log.Errorw("reed solomon recover object: corrupt too much", "ecid", suite.ECid)
		return nil, errors.GetAPIErr(errors.ErrRecoverFailed)
	}

	if len(corruptionList) == 0 {
		log.Errorw("reed solomon recover object: request ec shard not corrupt", "ecid", suite.ECid)
		return nil, errors.GetAPIErr(errors.ErrRecoverFailed)
	}

	for i, shard := range suite.Shards {
		ecReadClosers[i] = NewECReadCloser(shard.Frags, calc.storageEngine)
	}

	enc, err := reedsolomon.New(int(suite.DataShards), len(suite.Shards)-int(suite.DataShards))
	if err != nil {
		log.Errorw("reed solomon recover object: init read solomon", vars.ErrorKey, err, "ecid", suite.ECid)
		return nil, errors.GetAPIErr(errors.ErrServer)
	}

	frags := make([][]ec_store.Frag, len(suite.Shards))
	for _, u := range corruptionList {
		frags[u] = make([]ec_store.Frag, 0)
	}

	err = os.MkdirAll(ecFilePrefix, os.ModePerm)
	if err != nil {
		log.Errorw("reed solomon recover object: mkdir ec_data", vars.ErrorKey, err, "ecid", suite.ECid)
		return nil, errors.GetAPIErr(errors.ErrServer)
	}

	filesName := make([]string, 0)
	filesCnt := 0
	defer func() {
		for _, fileName := range filesName {
			_ = os.Remove(fileName)
		}
	}()

	for i := uint64(0); i < mx; i += ecBlockSize {
		length := util.Min(ecBlockSize, mx-i)
		data := make([][]byte, len(suite.Shards))
		for i := range data {
			if corruptionMap[i] {
				data[i] = nil
			} else {
				data[i] = make([]byte, length)
			}
		}
		for i := range ecReadClosers {
			if !corruptionMap[i] {
				_, _ = ecReadClosers[i].Read(data[i])
			}
		}
		err := enc.Reconstruct(data)
		if err != nil {
			log.Errorw("reed solomon recover object: reconstruct", vars.ErrorKey, err, "ecid", suite.ECid)
			return nil, errors.GetAPIErr(errors.ErrServer)
		}

		for _, u := range corruptionList {
			fileName := ecFileKey(suite.ECid, filesCnt)
			filesCnt++
			filesName = append(filesName, fileName)
			err := ioutil.WriteFile(fileName, data[u], os.ModePerm)
			if err != nil {
				log.Errorw("reed solomon recover: write ec temporary data", vars.ErrorKey, err, "ecid", suite.ECid)
				return nil, errors.GetAPIErr(errors.ErrServer)
			}
			frags[u] = append(frags[u], ec_store.Frag{
				Fid: fileName,
			})
		}
	}

	ret := make([]ec_store.Frag, 0)
	for _, u := range corruptionList {
		fileReadCloser := NewFilesReader(frags[u])
		for i, frag := range suite.Shards[u].Frags {
			if frag.Fid != "" {
				// delete original fid
				_ = calc.storageEngine.DeleteObject(ctx, frag.Fid)
				_ = calc.storageEngine.DeleteObject(ctx, frag.Fid)
			}

			fileReadCloser.SetLimit(int(frag.FileSize))
			fid, err := calc.storageEngine.PutObject(ctx, frag.FileSize, fileReadCloser, frag.FullPath.Name(), true)
			if err != nil {
				log.Errorw("reed solomon recover: put recover object", vars.ErrorKey, err, "ecid", suite.ECid)
				return nil, err
			}
			err = calc.storageEngine.AddLink(ctx, fid)
			if err != nil {
				log.Errorw("reed solomon recover: add link object", vars.ErrorKey, err, "ecid", suite.ECid)
				return nil, err
			}
			suite.Shards[u].Frags[i].Fid = fid
			ret = append(ret, ec_store.Frag{
				FullPath: frag.FullPath,
				Set:      frag.Set,
				Fid:      fid,
				FileSize: frag.FileSize,
				OldECid:  frag.OldECid,
			})
		}
		_ = fileReadCloser.Release()
	}

	err = calc.ECStore.InsertSuite(ctx, suite)
	if err != nil {
		log.Errorw("reed solomon recover: insert suite", vars.ErrorKey, err, "ecid", suite.ECid)
		return nil, err
	}

	return ret, nil
}

func (calc *ECCalc) reedSolomonDelete(ctx context.Context, ecid string) error {
	link, err := calc.ECStore.DelECidLink(ctx, ecid)
	if err != nil {
		return err
	}

	if link <= 0 {
		suite, err := calc.ECStore.GetSuite(ctx, ecid)
		if err != nil {
			return err
		}

		err = calc.ECStore.DeleteSuite(ctx, ecid)
		if err != nil {
			return err
		}

		err = calc.ECStore.ClrECidLink(ctx, ecid)
		if err != nil {
			return err
		}

		for _, shard := range suite.Shards {
			for _, frag := range shard.Frags {
				if frag.Fid != "" {
					_ = calc.storageEngine.DeleteObject(ctx, frag.Fid)
					_ = calc.storageEngine.DeleteObject(ctx, frag.Fid)
				}
			}
		}
	}

	return nil
}
