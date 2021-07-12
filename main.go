package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-datastore"
	"github.com/lyswifter/beacon-plus/clock"
	ltypes "github.com/lyswifter/beacon-plus/localtype"
	"github.com/lyswifter/beacon-plus/log"
	"github.com/mitchellh/go-homedir"
)

const repoPath = "~/.beaconplus"

func main() {
	log.Infof("This is the beacon plus server")

	repodir, err := homedir.Expand(repoPath)
	if err != nil {
		return
	}

	ldb, err := setupLevelDs(repodir, false)
	if err != nil {
		log.Infof("BeaconDB: err %s", err)
		return
	}
	BeaconDB = ldb
	log.Infof("BeaconDB: %v", BeaconDB)

	go func() {
		ds := BuiltinDrandConfig()
		log.Infof("BuiltinDrandConfig %+v", ds)

		be, err := RandomSchedule(ds)
		if err != nil {
			log.Infof("RandomSchedule %s", err)
			return
		}
		BeaconSche = be

		setupBeaconLoop()
	}()

	setupBeaconServer()
}

func setupBeaconLoop() {
	log.Infof("SetupBeaconLoop")

	for {
		fmt.Println()
		log.Infof("loop again")
		time.Sleep(10 * time.Millisecond)

		pctx := context.TODO()

		nullround := 0
		curTimestamp := clock.Clock.Now().Unix()
		baseEpoch := (curTimestamp - GenesisTimeStamp) / int64(BlockDelaySecs)
		baseTimestamp := GenesisTimeStamp + baseEpoch*int64(BlockDelaySecs)

		// get beacon for the round

		ctx, cancel := context.WithTimeout(pctx, time.Second*3)
		defer cancel()
		nextEpoch := abi.ChainEpoch(baseEpoch) + 1
		entryNext, err := BeaconGetEntry(ctx, nextEpoch)
		if err != nil {
			log.Infof("BeaconGetEntry: %s", err)
			continue
		}

		nextnextEpoch := abi.ChainEpoch(baseEpoch) + 2
		entryNextNext, err := BeaconGetEntry(ctx, nextnextEpoch)
		if err != nil {
			log.Infof("BeaconGetEntry: %s", err)
			continue
		}

		if entryNext != nil && entryNextNext != nil {
			log.Infof("entry-next round: %d data: %v", entryNext.Round, entryNext.Data)
			log.Infof("entry-next-next round: %d data: %v", entryNextNext.Round, entryNextNext.Data)
		}

		entryInfoNext := ltypes.BeaconEntryInfo{
			Round: entryNext.Round,
			Entry: entryNext,
		}

		entryInfoNextNext := ltypes.BeaconEntryInfo{
			Round: entryNextNext.Round,
			Entry: entryNextNext,
		}

		err = saveBeacon(nextEpoch, entryInfoNext)
		if err != nil {
			return
		}

		err = saveBeacon(nextnextEpoch, entryInfoNextNext)
		if err != nil {
			return
		}

		nullround++
		nextRound := time.Unix(int64(baseTimestamp+int64(BlockDelaySecs)*int64(nullround)), 0)

		log.Infof("sleep to next round: %s nullround: %d", nextRound.String(), nullround)

		select {
		case <-clock.Clock.After(clock.Clock.Until(nextRound)):
		case <-pctx.Done():
			return
		}
	}
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

func setupBeaconServer() {
	r := gin.Default()
	r.GET("/public/:epoch", func(c *gin.Context) {
		start := time.Now()
		epoch := c.Param("epoch")
		key := datastore.NewKey(epoch)

		ishas, err := BeaconDB.Has(key)
		if err != nil {
			log.Infof("entrys: has %s", err)
			return
		}

		if !ishas {
			c.JSON(500, gin.H{
				"status":  "Err",
				"epoch":   epoch,
				"message": fmt.Sprintf("get beacon for epoch: %s, is not exist", epoch),
			})
		}

		qt, err := BeaconDB.Get(key)
		if err != nil {
			log.Infof("entrys: list %s", err)
			return
		}

		entrys := ltypes.BeaconEntryInfo{}
		err = json.Unmarshal(qt, &entrys)
		if err != nil {
			return
		}

		c.JSON(200, gin.H{
			"status":  "Ok",
			"epoch":   epoch,
			"message": string(qt),
		})

		log.Infof("Request from client: %s round: %d data: %v took: %s", c.ClientIP(), entrys.Round, entrys.Entry.Data, time.Since(start).String())
	})

	r.Run("0.0.0.0:9090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
