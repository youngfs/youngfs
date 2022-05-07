package ec_store

import (
	"context"
	"icesos/full_path"
	"icesos/set"
)

type Frag struct {
	full_path.FullPath        // full path
	set.Set                   // set
	Fid                string // fid
	FileSize           uint64 // file size
	OldECid            string // old ecid
}

//Frags must not be []Frag{} (codec will become nill)
type Shard struct {
	Host  string // host
	Frags []Frag // frags
	Md5   []byte // MD5
}

//DataBlocks must not be []Shard{} (codec will become nill)
type Suite struct {
	ECid               string  // erasure code id
	full_path.FullPath         // full path
	set.Set                    // set
	OrigHost           string  // original host
	OrigFid            string  // original fid
	FileSize           uint64  // file size
	BakHost            string  // backup host
	BakFid             string  // backup fid
	Next               string  // next ECid, if it's end,next = ""
	Shards             []Shard // data blocks
}

func (ec *ECStore) getFrags(ctx context.Context, set set.Set, turns int) ([]Frag, error) {
	ecids, err := ec.kvStore.SMembers(ctx, setPlanShardKey(set, turns))
	if err != nil {
		return nil, err
	}

	frags := make([]Frag, len(ecids))
	for i := range frags {
		suite, err := ec.GetSuite(ctx, string(ecids[i]))
		if err != nil {
			return nil, err
		}
		frags[i] = Frag{
			FullPath: suite.FullPath,
			Set:      suite.Set,
			Fid:      suite.OrigFid,
			FileSize: suite.FileSize,
			OldECid:  suite.ECid,
		}
	}

	return frags, nil
}

func (ec *ECStore) InsertSuite(ctx context.Context, suite *Suite) error {
	proto, err := suite.EncodeProto(ctx)
	if err != nil {
		return err
	}

	err = ec.kvStore.KvPut(ctx, ecidKey(suite.ECid), proto)
	if err != nil {
		return err
	}

	return nil
}

func (ec *ECStore) GetSuite(ctx context.Context, ecid string) (*Suite, error) {
	proto, err := ec.kvStore.KvGet(ctx, ecidKey(ecid))
	if err != nil {
		return nil, err
	}

	suite, err := DecodeSuiteProto(ctx, proto)
	if err != nil {
		return nil, err
	}

	return suite, nil
}

func (ec *ECStore) DeleteSuite(ctx context.Context, ecid string) error {
	_, err := ec.kvStore.KvDelete(ctx, ecidKey(ecid))
	if err != nil {
		return err
	}

	return nil
}
