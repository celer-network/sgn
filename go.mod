module github.com/celer-network/sgn

go 1.13

require (
	github.com/allegro/bigcache v1.2.1
	github.com/celer-network/goutils v0.1.13
	github.com/cosmos/cosmos-sdk v0.38.1
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/ethereum/go-ethereum v1.9.6
	github.com/gammazero/deque v0.0.0-20190521012701-46e4ffb7a622
	github.com/golang/protobuf v1.3.4
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/spf13/cobra v0.0.6
	github.com/spf13/viper v1.6.2
	github.com/stretchr/testify v1.5.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.3
	github.com/tendermint/tm-db v0.5.0
)

replace github.com/cosmos/cosmos-sdk => github.com/celer-network/cosmos-sdk v0.38.3-0

// replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk
