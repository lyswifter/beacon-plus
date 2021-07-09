package main

import (
	"context"
	"fmt"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/chain/beacon"
	"github.com/filecoin-project/lotus/chain/beacon/drand"
	"github.com/filecoin-project/lotus/chain/types"
	"github.com/filecoin-project/lotus/node/modules/dtypes"
	"golang.org/x/xerrors"
)

var BeaconSche beacon.Schedule

func BuiltinDrandConfig() dtypes.DrandSchedule {
	return DrandConfigSchedule()
}

func RandomSchedule(DrandConfig dtypes.DrandSchedule) (beacon.Schedule, error) {
	shd := beacon.Schedule{}
	for _, dc := range DrandConfig {
		bc, err := drand.NewDrandBeacon(GenesisTimeStamp, BlockDelaySecs, nil, dtypes.DrandConfig{
			Servers:       dc.Config.Servers,
			Relays:        dc.Config.Relays,
			ChainInfoJSON: dc.Config.ChainInfoJSON,
		})
		if err != nil {
			return nil, xerrors.Errorf("creating drand beacon: %w", err)
		}
		shd = append(shd, beacon.BeaconPoint{Start: dc.Start, Beacon: bc})
	}

	return shd, nil
}

func BeaconGetEntry(ctx context.Context, epoch abi.ChainEpoch) (*types.BeaconEntry, error) {
	b := BeaconSche.BeaconForEpoch(epoch)
	rr := b.MaxBeaconRoundForEpoch(epoch)
	e := b.Entry(ctx, rr)

	go b.Watch(ctx)

	select {
	case be, ok := <-e:
		if !ok {
			return nil, fmt.Errorf("beacon get returned no value")
		}
		if be.Err != nil {
			return nil, be.Err
		}
		return &be.Entry, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
