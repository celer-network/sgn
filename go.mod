module github.com/celer-network/sgn

go 1.15

require (
	github.com/allegro/bigcache v1.2.1
	github.com/celer-network/eth-services v0.0.0-20210210185550-26d7e6f75ead
	github.com/celer-network/goutils v0.1.16
	github.com/celer-network/sgn-contract v0.2.8
	github.com/cosmos/cosmos-sdk v0.39.2
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/deckarep/golang-set v1.7.1
	github.com/ethereum/go-ethereum v1.9.20
	github.com/gammazero/deque v0.0.0-20200721202602-07291166fe33
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.4
	github.com/gorilla/mux v1.7.4
	github.com/iancoleman/strcase v0.1.0
	github.com/levigross/grequests v0.0.0-20190908174114-253788527a1a
	github.com/mitchellh/mapstructure v1.4.0 // indirect
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.9.0 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.9
	github.com/tendermint/tm-db v0.5.2
	go.uber.org/zap v1.16.0
	google.golang.org/grpc v1.31.1 // indirect
)

replace github.com/cosmos/cosmos-sdk => github.com/celer-network/cosmos-sdk v0.39.3-celer.1

// replace github.com/celer-network/eth-services => ../eth-services

//replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk
