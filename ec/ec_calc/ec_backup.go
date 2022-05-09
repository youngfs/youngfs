package ec_calc

import (
	"context"
	"crypto/md5"
	"icesos/ec/ec_store"
	"icesos/entry"
	"icesos/errors"
	"icesos/log"
	"icesos/util"
	"io"
	"net/http"
)

func (calc *ECCalc) backup(ctx context.Context, suite *ec_store.Suite) error {
	url, err := calc.storageEngine.GetFidUrl(ctx, suite.OrigFid)
	if err != nil {
		log.Errorw("backup err: get fid url", "ecid", suite.ECid)
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Errorw("backup err: http get", "ecid", suite.ECid, "http code", resp.StatusCode)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorw("backup err: http get", "ecid", suite.ECid, "http code", resp.StatusCode)
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bakFid, err := calc.storageEngine.PutObject(ctx, util.GetContentLength(resp.Header), resp.Body, suite.BakHost)
	if err != nil {
		log.Errorw("backup err: put backup", "ecid", suite.ECid)
		return err
	}

	suite.BakFid = bakFid
	err = calc.ECStore.InsertSuite(ctx, suite)
	if err != nil {
		log.Errorw("backup err: insert suite", "ecid", suite.ECid)
		return err
	}

	log.Infow("back up successful", "ecid", suite.ECid)

	return nil
}

// return len([]ec_store.Frag) == 1
func (calc *ECCalc) backupRecover(ctx context.Context, suite *ec_store.Suite, ent *entry.Entry) ([]ec_store.Frag, error) {
	url, err := calc.storageEngine.GetFidUrl(ctx, suite.BakFid)
	if err != nil {
		log.Errorw("backup recover object: get fid url", "ecid", suite.ECid)
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Errorw("backup recover object: http get", "ecid", suite.ECid, "http code", resp.StatusCode)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorw("backup recover object: http get", "ecid", suite.ECid, "http code", resp.StatusCode)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	md5Hash := md5.New()
	file := io.TeeReader(resp.Body, md5Hash)

	fid, err := calc.storageEngine.PutObject(ctx, util.GetContentLength(resp.Header), file)
	if err != nil {
		log.Errorw("backup recover object: put backup", "ecid", suite.ECid)
		return nil, err
	}

	md5Ret := md5Hash.Sum(nil)
	if len(ent.Md5) != 0 && !util.BytesIsEqual(md5Ret, ent.Md5) {
		suite.BakFid = ""
		err := calc.storageEngine.DeleteObject(ctx, fid)
		if err != nil {
			log.Errorw("backup recover object: delete ")
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
