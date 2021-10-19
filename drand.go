package main

import "github.com/filecoin-project/go-state-types/abi"

// main network
// const GenesisTimeStamp = 1598306400

// var DSchedule = map[abi.ChainEpoch]DrandEnum{
// 	51000: DrandMainnet,
// }

// calibration network
const GenesisTimeStamp = 1624060800

var DSchedule = map[abi.ChainEpoch]DrandEnum{
	0: DrandMainnet,
}
