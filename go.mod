module github.com/lyswifter/beacon-plus

go 1.16

require (
	github.com/drand/drand v1.2.1
	github.com/drand/kyber v1.1.6
	github.com/filecoin-project/go-state-types v0.1.1-0.20210506134452-99b279731c48
	github.com/filecoin-project/lotus v1.6.0
	github.com/gin-gonic/gin v1.7.1
	github.com/go-kit/kit v0.10.0
	github.com/go-playground/validator/v10 v10.5.0 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/mock v1.5.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/ipfs/go-cid v0.0.7
	github.com/ipfs/go-datastore v0.4.5
	github.com/ipfs/go-ds-leveldb v0.4.2
	github.com/ipfs/go-log v1.0.4
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	github.com/libp2p/go-libp2p-pubsub v0.4.2-0.20210212194758-6c1addf493eb
	github.com/mitchellh/go-homedir v1.1.0
	github.com/prometheus/client_golang v1.10.0 // indirect
	github.com/prometheus/common v0.25.0 // indirect
	github.com/raulk/clock v1.1.0
	github.com/syndtr/goleveldb v1.0.0
	github.com/ugorji/go v1.2.5 // indirect
	github.com/whyrusleeping/cbor-gen v0.0.0-20210219115102-f37d292932f2
	go.opencensus.io v0.23.0 // indirect
	go.uber.org/zap v1.16.1-0.20210329175301-c23abee72d19
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/lint v0.0.0-20210508222113-6edffad5e616 // indirect
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420 // indirect
	golang.org/x/sys v0.0.0-20210514084401-e8d321eab015 // indirect
	golang.org/x/tools v0.1.1 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
	google.golang.org/genproto v0.0.0-20210513213006-bf773b8c8384 // indirect
	google.golang.org/grpc v1.38.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/drand/drand => ./extern/drand
