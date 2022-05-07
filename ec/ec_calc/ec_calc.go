package ec_calc

import (
	"context"
	"crypto/md5"
	"icesos/ec/ec_store"
	"icesos/entry"
	"icesos/errors"
	"icesos/log"
	"icesos/storage_engine"
	"icesos/util"
	"io"
	"net/http"
)

type ECCalc struct {
	ECStore       *ec_store.ECStore
	storageEngine storage_engine.StorageEngine
}

func NewECCalc(ecStore *ec_store.ECStore, storageEngine storage_engine.StorageEngine) *ECCalc {
	return &ECCalc{
		ECStore:       ecStore,
		storageEngine: storageEngine,
	}
}

func (calc *ECCalc) ExecEC(ctx context.Context, ecid string) error {
	suite, err := calc.ECStore.GetSuite(ctx, ecid)
	if err != nil {
		return err
	}

	for suite.Next != "" {
		suite, err = calc.ECStore.GetSuite(ctx, ecid)
		if err != nil {
			return err
		}
	}

	if suite.BakHost != "" && suite.OrigFid != "" {
		err := calc.backup(ctx, suite)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}

func (calc *ECCalc) backup(ctx context.Context, suite *ec_store.Suite) error {
	url, err := calc.storageEngine.GetFidUrl(ctx, suite.OrigFid)
	if err != nil {
		log.Errorw("backup err: get fid url", "ecid", suite.ECid, "suite", suite)
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Errorw("backup err: http get", "ecid", suite.ECid, "suite", suite, "http code", resp.StatusCode)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorw("backup err: http get", "ecid", suite.ECid, "suite", suite, "http code", resp.StatusCode)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bakFid, err := calc.storageEngine.PutObject(ctx, util.GetContentLength(resp.Header), resp.Body, suite.BakHost)
	if err != nil {
		log.Errorw("backup err: put backup", "ecid", suite.ECid, "suite", suite)
		return err
	}

	suite.BakFid = bakFid
	err = calc.ECStore.InsertSuite(ctx, suite)
	if err != nil {
		log.Errorw("backup err: insert suite", "ecid", suite.ECid, "suite", suite)
		return err
	}

	log.Infow("back up successful", "suite", suite)

	return nil
}

func (calc *ECCalc) RecoverObject(ctx context.Context, ent *entry.Entry) ([]ec_store.Frag, error) {
	if ent.ECid == "" {
		return nil, errors.GetAPIErr(errors.ErrRecoverFailed)
	}

	suite, err := calc.ECStore.GetSuite(ctx, ent.ECid)
	if err != nil {
		return nil, err
	}

	if suite.BakFid != "" {
		url, err := calc.storageEngine.GetFidUrl(ctx, suite.BakFid)
		if err != nil {
			log.Errorw("resume object: get fid url", "ecid", suite.ECid, "suite", suite)
			return nil, err
		}

		resp, err := http.Get(url)
		if err != nil {
			log.Errorw("resume object: http get", "ecid", suite.ECid, "suite", suite, "http code", resp.StatusCode)
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			log.Errorw("resume object: http get", "ecid", suite.ECid, "suite", suite, "http code", resp.StatusCode)
			return nil, err
		}
		defer func() {
			_ = resp.Body.Close()
		}()

		md5Hash := md5.New()
		file := io.TeeReader(resp.Body, md5Hash)

		fid, err := calc.storageEngine.PutObject(ctx, util.GetContentLength(resp.Header), file)
		if err != nil {
			log.Errorw("backup err: put backup", "ecid", suite.ECid, "suite", suite)
			return nil, err
		}

		md5Ret := md5Hash.Sum(nil)
		if !util.BytesIsEqual(md5Ret, ent.Md5) {
			suite.BakFid = ""
			err := calc.storageEngine.DeleteObject(ctx, fid)
			if err != nil {
				log.Errorw("backup err: delete ")
			}
			return nil, errors.GetAPIErr(errors.ErrRecoverFailed)
		}

		suite.OrigFid = fid
		err = calc.ECStore.InsertSuite(ctx, suite)
		if err != nil {
			return nil, err
		}

		frag := make([]ec_store.Frag, 1)
		frag[0] = ec_store.Frag{
			FullPath: ent.FullPath,
			Set:      ent.Set,
			Fid:      fid,
		}

		return frag, nil
	}

	return nil, errors.GetAPIErr(errors.ErrRecoverFailed)
}

func (calc *ECCalc) DeleteObject(ctx context.Context, ecid string) error {
	if ecid == "" {
		return nil
	}

	suite, err := calc.ECStore.GetSuite(ctx, ecid)
	if err != nil {
		return nil
	}

	if suite.BakFid != "" {
		err := calc.storageEngine.DeleteObject(ctx, suite.BakFid)
		if err != nil {
			return err
		}
	}

	err = calc.ECStore.DeleteSuite(ctx, ecid)
	if err != nil {
		return err
	}

	return nil
}
