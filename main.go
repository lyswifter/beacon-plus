package main

import (
	"github.com/lyswifter/beacon-plus/log"
)

const repoPath = "~/.beaconplus"

func main() {
	log.Infof("This is the beacon plus server")

	DataStores()

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
