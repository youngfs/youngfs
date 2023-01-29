package ec_calc

import (
	"context"
	"crypto/md5"
	"io"
	"net/http"
	"youngfs/errors"
	"youngfs/fs/ec/ec_store"
	"youngfs/fs/entry"
	"youngfs/log"
	"youngfs/util"
)

func (calc *ECCalc) backup(ctx context.Context, suite *ec_store.Suite) error {
	// todo: storage engine need a full interface
	url, err := calc.storageEngine.GetFidUrl(ctx, suite.OrigFid)
	if err != nil {
		return errors.Wrap(err, "backup failed, ecid: "+suite.ECid)
	}

	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "backup failed, ecid: "+suite.ECid)
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("backup failed, ecid: " + suite.ECid)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	bakFid, err := calc.storageEngine.PutObject(ctx, suite.FileSize, resp.Body, suite.FullPath.Name(), true, suite.BakHost)
	if err != nil {
		return errors.WithMessage(err, "backup failed, ecid: "+suite.ECid)
	}

	suite.BakFid = bakFid
	err = calc.ECStore.InsertSuite(ctx, suite)
	if err != nil {
		return err
	}

	log.Infow("back up successful", "ecid", suite.ECid)

	return nil
}

// return len([]ec_store.Frag) == 1
func (calc *ECCalc) backupRecover(ctx context.Context, suite *ec_store.Suite, ent *entry.Entry) ([]ec_store.Frag, error) {
	// todo: storage engine need a full interface
	url, err := calc.storageEngine.GetFidUrl(ctx, suite.BakFid)
	if err != nil {
		return nil, errors.Wrap(err, "backup recover failed, ecid: "+suite.ECid)
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "backup recover failed, ecid: "+suite.ECid)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("backup recover failed, ecid: " + suite.ECid)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	md5Hash := md5.New()
	file := io.TeeReader(resp.Body, md5Hash)

	fid, err := calc.storageEngine.PutObject(ctx, suite.FileSize, file, suite.FullPath.Name(), true)
	if err != nil {
		return nil, errors.WithMessage(err, "backup recover failed: put object, ecid: "+suite.ECid)
	}

	md5Ret := md5Hash.Sum(nil)
	if len(ent.Md5) != 0 && !util.BytesIsEqual(md5Ret, ent.Md5) {
		suite.BakFid = ""
		err := calc.storageEngine.DeleteObject(ctx, fid)
		return nil, errors.ErrRecoverFailed.WrapErr(err)
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
		OldECid:  suite.ECid,
	}

	return frag, nil
}
