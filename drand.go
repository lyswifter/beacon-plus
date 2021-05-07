package main

import (
	"github.com/filecoin-project/go-state-types/abi"
)

const GenesisTimeStamp = 1598306400

var DSchedule = map[abi.ChainEpoch]DrandEnum{
	51000: DrandMainnet,
}

// const GenesisTimeStamp = 1613776200

// var DSchedule = map[abi.ChainEpoch]DrandEnum{
// 	0: DrandMainnet,
// }
