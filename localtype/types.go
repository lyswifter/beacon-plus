package localtype

import "github.com/filecoin-project/lotus/chain/types"

type BeaconEntryInfo struct {
	Round uint64
	Entry *types.BeaconEntry
}

type MinerInfo struct {
	Address string `json:"address"`
	Symbol  string `json:"symbol"`
	IP      string `json:"ip"`
}
