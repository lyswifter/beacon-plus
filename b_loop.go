package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/lyswifter/beacon-plus/clock"
	ltypes "github.com/lyswifter/beacon-plus/localtype"
	"github.com/lyswifter/beacon-plus/log"
)

func setupBeaconLoop() {
	log.Infof("SetupBeaconLoop")

	for {
		time.Sleep(10 * time.Millisecond)

		fmt.Println()
		log.Infof("loop again")

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
			log.Infof("entry-next round: %d data: %v", entryNext.Round, hex.EncodeToString(entryNext.Data))
			log.Infof("entry-next-next round: %d data: %v", entryNextNext.Round, hex.EncodeToString(entryNextNext.Data))
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
