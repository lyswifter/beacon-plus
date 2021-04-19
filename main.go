package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/gin-gonic/gin"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"
	"github.com/lyswifter/beacon-plus/clock"
	ltypes "github.com/lyswifter/beacon-plus/localtype"
	"github.com/mitchellh/go-homedir"
)

const repoPath = "~/.beaconplus"

var log = logging.Logger("main")

func main() {
	fmt.Print("This is the beacon plus server\n")

	repodir, err := homedir.Expand(repoPath)
	if err != nil {
		return
	}

	ldb, err := setupLevelDs(repodir, false)
	if err != nil {
		fmt.Printf("BeaconDB: err %s\n", err)
		return
	}
	BeaconDB = ldb
	fmt.Printf("BeaconDB: %v\n", BeaconDB)

	go func() {
		ds := BuiltinDrandConfig()
		fmt.Printf("BuiltinDrandConfig %+v\n", ds)

		be, err := RandomSchedule(ds)
		if err != nil {
			fmt.Printf("RandomSchedule %s\n", err)
			return
		}
		BeaconSche = be

		setupBeaconLoop()
	}()

	setupBeaconServer()
}

func setupBeaconLoop() {
	fmt.Print("setupBeaconLoop\n")

	for {
		fmt.Println()
		fmt.Print("loop again\n")

		pctx := context.TODO()

		nullround := 0
		curTimestamp := clock.Clock.Now().Unix()
		baseEpoch := (curTimestamp - GenesisTimeStamp) / int64(BlockDelaySecs)
		baseTimestamp := GenesisTimeStamp + baseEpoch*int64(BlockDelaySecs)
		fmt.Printf("baseEpoch: %d baseTimestamp: %d\n", baseEpoch, baseTimestamp)

		//get beacon for the round

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
			fmt.Printf("entry-next round: %d data: %v\n", entryNext.Round, entryNext.Data)
			fmt.Printf("entry-next-next round: %d data: %v\n", entryNextNext.Round, entryNextNext.Data)
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
	fmt.Printf("write key: %+v\n", key)
	ishas, err := BeaconDB.Has(key)
	if err != nil {
		fmt.Printf("entrys: has %s\n", err)
		return err
	}

	if !ishas {
		in, err := json.Marshal(info)
		if err != nil {
			return err
		}

		err = BeaconDB.Put(key, in)
		if err != nil {
			fmt.Printf("entrys: begin %s\n", err)
			return err
		}
	}

	return nil
}

func setupBeaconServer() {
	r := gin.Default()
	r.GET("/public/:epoch", func(c *gin.Context) {
		epoch := c.Param("epoch")

		key := datastore.NewKey(epoch)
		fmt.Printf("read key: %+v\n", key)
		ishas, err := BeaconDB.Has(key)
		if err != nil {
			fmt.Printf("entrys: has %s\n", err)
			return
		}

		if !ishas {
			// c.String(http.StatusInternalServerError, "Err get beacon for epoch: %s, is not exist", epoch)
			c.JSON(500, gin.H{
				"status":  "Err",
				"epoch":   epoch,
				"message": fmt.Sprintf("get beacon for epoch: %s, is not exist", epoch),
			})
		}

		qt, err := BeaconDB.Get(key)
		if err != nil {
			fmt.Printf("entrys: list %s\n", err)
			return
		}

		entrys := ltypes.BeaconEntryInfo{}
		err = json.Unmarshal(qt, &entrys)
		if err != nil {
			return
		}

		fmt.Printf("entry read round: %d data: %+v\n", entrys.Round, entrys)
		c.JSON(200, gin.H{
			"status":  "Ok",
			"epoch":   epoch,
			"message": string(qt),
		})

		// c.String(http.StatusOK, "Finished get beacon for epoch: %s ret: %v", epoch, string(qt))
	})

	r.Run("0.0.0.0:9090") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func OkFunc(c *gin.Context) {

}

func ErrFunc(c *gin.Context) {

}
