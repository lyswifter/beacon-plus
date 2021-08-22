package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	levelds "github.com/ipfs/go-ds-leveldb"
	ltypes "github.com/lyswifter/beacon-plus/localtype"
	"github.com/lyswifter/beacon-plus/log"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
	"go.uber.org/multierr"
	"golang.org/x/xerrors"
)

var BeaconDB datastore.Batching
var MinerInfoDB datastore.Batching

func setupLevelDs(path string, readonly bool) (datastore.Batching, error) {
	if _, err := os.ReadDir(path); err != nil {
		if os.IsNotExist(err) {
			//mkdir
			err = os.Mkdir(path, 0777)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := levelds.NewDatastore(path, &levelds.Options{
		Compression: ldbopts.NoCompression,
		NoSync:      false,
		Strict:      ldbopts.StrictAll,
		ReadOnly:    readonly,
	})
	if err != nil {
		fmt.Printf("NewDatastore: %s\n", err)
		return nil, err
	}

	return db, err
}

func DataStores() {
	repodir, err := homedir.Expand(repoPath)
	if err != nil {
		return
	}

	ldb, err := setupLevelDs(repodir, false)
	if err != nil {
		log.Infof("setup beacondb: err %s", err)
		return
	}
	BeaconDB = ldb
	log.Infof("BeaconDB: %+v", BeaconDB)

	idb, err := setupLevelDs(path.Join(repodir, "minfo"), false)
	if err != nil {
		log.Infof("setup infodb: err %s", err)
		return
	}
	MinerInfoDB = idb
	log.Infof("MinerInfoDB: %+v", MinerInfoDB)
}

func saveBeacon(epoch abi.ChainEpoch, info ltypes.BeaconEntryInfo) error {
	key := datastore.NewKey(fmt.Sprintf("%d", epoch))
	ishas, err := BeaconDB.Has(key)
	if err != nil {
		log.Infof("entrys: has %s", err)
		return err
	}

	if !ishas {
		in, err := json.Marshal(info)
		if err != nil {
			return err
		}

		err = BeaconDB.Put(key, in)
		if err != nil {
			log.Infof("entrys: begin %s", err)
			return err
		}

		log.Infof("write beacon for epoch: %s", key.String())
	}

	return nil
}

func saveMInfo(addr string, info ltypes.MinerInfo) error {
	key := datastore.NewKey(addr)

	isHas, err := MinerInfoDB.Has(key)
	if err != nil {
		log.Infof("minfo: has %s", err)
		return err
	}

	if !isHas {
		in, err := json.Marshal(info)
		if err != nil {
			return err
		}

		err = MinerInfoDB.Put(key, in)
		if err != nil {
			log.Infof("minfo: begin %s", err)
			return err
		}

		log.Infof("write minfo for addr: %s val %v", key.String(), info)
	}

	return nil
}

func readMInfos() ([]ltypes.MinerInfo, error) {
	res, err := MinerInfoDB.Query(query.Query{})
	if err != nil {
		return nil, err
	}

	defer res.Close()

	minfos := []ltypes.MinerInfo{}

	var errs error

	for {
		res, ok := res.NextSync()
		if !ok {
			break
		}

		if res.Error != nil {
			return nil, res.Error
		}

		minfo := &ltypes.MinerInfo{}
		err := json.Unmarshal(res.Value, minfo)
		if err != nil {
			errs = multierr.Append(errs, xerrors.Errorf("decoding state for key '%s': %w", res.Key, err))
			continue
		}

		minfos = append(minfos, *minfo)
	}

	log.Infof("read minfos ok, len %d", len(minfos))

	return minfos, nil
}

func readmInfo(maddr string) (*ltypes.MinerInfo, error) {
	key := datastore.NewKey(maddr)
	isHas, err := MinerInfoDB.Has(key)
	if err != nil {
		log.Infof("minfo: has %s", err)
		return nil, xerrors.Errorf("has err for %s err %s", key.String(), err.Error())
	}

	if !isHas {
		return nil, xerrors.Errorf("minfo not exist: %s", key.String())
	}

	res, err := MinerInfoDB.Get(key)
	if err != nil {
		return nil, err
	}

	minfo := &ltypes.MinerInfo{}
	err = json.Unmarshal(res, minfo)
	if err != nil {
		return nil, xerrors.Errorf("unmarsal err %s", err.Error())
	}

	log.Infof("read minfo(%s) ok", key.String())

	return minfo, nil
}
