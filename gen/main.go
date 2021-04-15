package main

import (
	"fmt"
	"os"

	"github.com/lyswifter/beacon-plus/localtype"
	gen "github.com/whyrusleeping/cbor-gen"
)

func main() {
	err := gen.WriteMapEncodersToFile("../localtype/cbor_gen.go", "localtype",
		localtype.BeaconEntryInfo{},
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
