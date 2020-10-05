module github.com/celer-network/sgn

go 1.14

require (
	github.com/allegro/bigcache v1.2.1
	github.com/celer-network/goutils v0.1.16
	github.com/cosmos/cosmos-sdk v0.39.1
	github.com/cosmos/go-bip39 v0.0.0-20180819234021-555e2067c45d
	github.com/deckarep/golang-set v1.7.1
	github.com/dvsekhvalnov/jose2go v0.0.0-20200901110807-248326c1351b // indirect
	github.com/ethereum/go-ethereum v1.9.20
	github.com/gammazero/deque v0.0.0-20200721202602-07291166fe33
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/mux v1.7.4
	github.com/iancoleman/strcase v0.1.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/go-amino v0.15.1
	github.com/tendermint/tendermint v0.33.7
	github.com/tendermint/tm-db v0.5.1
)

replace github.com/cosmos/cosmos-sdk => github.com/celer-network/cosmos-sdk v0.39.2-celer.3

replace github.com/dvsekhvalnov/jose2go => github.com/dvsekhvalnov/jose2go v0.0.0-20200901110807-248326c1351b

//replace github.com/cosmos/cosmos-sdk => ../cosmos-sdk
