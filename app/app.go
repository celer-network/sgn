package app

import (
	"encoding/json"
	"os"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	"github.com/celer-network/sgn/monitor"
	"github.com/celer-network/sgn/transactor"
	"github.com/celer-network/sgn/x/cron"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/gov"
	govclient "github.com/celer-network/sgn/x/gov/client"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/subscribe"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/genutil"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	tlog "github.com/tendermint/tendermint/libs/log"
	tmos "github.com/tendermint/tendermint/libs/os"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

const appName = "sgn"

var (
	// default home directories for the application CLI
	DefaultCLIHome = os.ExpandEnv("$HOME/.sgncli")

	// DefaultNodeHome sets the folder where the applcation data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.sgn")

	// ModuleBasicManager is in charge of setting up basic module elemnets
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},

		cron.AppModule{},
		global.AppModule{},
		gov.NewAppModuleBasic(govclient.ParamProposalHandler),
		slash.AppModule{},
		subscribe.AppModule{},
		sync.AppModule{},
		validator.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}

	monitored = false
	ethClient *mainchain.EthClient
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type sgnApp struct {
	*bam.BaseApp
	cdc *codec.Codec

	// Keys to access the substores
	tkeyStaking  *sdk.TransientStoreKey
	tkeyParams   *sdk.TransientStoreKey
	keyMain      *sdk.KVStoreKey
	keyAccount   *sdk.KVStoreKey
	keySupply    *sdk.KVStoreKey
	keyStaking   *sdk.KVStoreKey
	keyParams    *sdk.KVStoreKey
	keyCron      *sdk.KVStoreKey
	keyGlobal    *sdk.KVStoreKey
	keyGov       *sdk.KVStoreKey
	keySlash     *sdk.KVStoreKey
	keySubscribe *sdk.KVStoreKey
	keySync      *sdk.KVStoreKey
	keyValidator *sdk.KVStoreKey

	// Keepers
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	stakingKeeper   staking.Keeper
	supplyKeeper    supply.Keeper
	paramsKeeper    params.Keeper
	cronKeeper      cron.Keeper
	globalKeeper    global.Keeper
	govKeeper       gov.Keeper
	slashKeeper     slash.Keeper
	subscribeKeeper subscribe.Keeper
	syncKeeper      sync.Keeper
	validatorKeeper validator.Keeper

	// Module Manager
	mm *module.Manager
}

// NewSgnApp is a constructor function for sgnApp
func NewSgnApp(logger tlog.Logger, db dbm.DB, baseAppOptions ...func(*bam.BaseApp)) *sgnApp {
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		tmos.Exit(err.Error())
	}
	viper.SetDefault(common.FlagStartMonitor, true)
	viper.SetDefault(common.FlagEthPollInterval, 5)
	viper.SetDefault(common.FlagEthConfirmCount, 7)

	ethClient, err = mainchain.NewEthClient(
		viper.GetString(common.FlagEthInstance),
		viper.GetString(common.FlagEthDPoSAddress),
		viper.GetString(common.FlagEthSGNAddress),
		viper.GetString(common.FlagEthLedgerAddress),
		viper.GetString(common.FlagEthKeystore),
		viper.GetString(common.FlagEthPassphrase),
	)
	if err != nil {
		tmos.Exit(err.Error())
	}

	log.SetLevelByName(viper.GetString(common.FlagLogLevel))
	if viper.GetBool(common.FlagLogColor) {
		log.EnableColor()
	}
	if viper.GetBool(common.FlagLogLongFile) {
		common.EnableLogLongFile()
	}

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	// Here you initialize your application with the store keys it requires
	var app = &sgnApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:      sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:   sdk.NewKVStoreKey(auth.StoreKey),
		keySupply:    sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:   sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:  sdk.NewTransientStoreKey(staking.TStoreKey),
		keyParams:    sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:   sdk.NewTransientStoreKey(params.TStoreKey),
		keyCron:      sdk.NewKVStoreKey(cron.StoreKey),
		keyGlobal:    sdk.NewKVStoreKey(global.StoreKey),
		keyGov:       sdk.NewKVStoreKey(gov.StoreKey),
		keySlash:     sdk.NewKVStoreKey(slash.StoreKey),
		keySubscribe: sdk.NewKVStoreKey(subscribe.StoreKey),
		keySync:      sdk.NewKVStoreKey(sync.StoreKey),
		keyValidator: sdk.NewKVStoreKey(validator.StoreKey),
	}

	// The ParamsKeeper handles parameter storage for the application
	app.paramsKeeper = params.NewKeeper(app.cdc, app.keyParams, app.tkeyParams)
	// Set specific subspaces
	authSubspace := app.paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSupspace := app.paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := app.paramsKeeper.Subspace(staking.DefaultParamspace)
	globalSubspace := app.paramsKeeper.Subspace(global.DefaultParamspace)
	govSubspace := app.paramsKeeper.Subspace(gov.DefaultParamspace).WithKeyTable(gov.ParamKeyTable())
	validatorSubspace := app.paramsKeeper.Subspace(validator.DefaultParamspace)
	slashSubspace := app.paramsKeeper.Subspace(slash.DefaultParamspace)
	subscribeSubspace := app.paramsKeeper.Subspace(subscribe.DefaultParamspace)
	syncSubspace := app.paramsKeeper.Subspace(sync.DefaultParamspace).WithKeyTable(sync.ParamKeyTable())

	// The AccountKeeper handles address -> account lookups
	app.accountKeeper = auth.NewAccountKeeper(
		app.cdc,
		app.keyAccount,
		authSubspace,
		auth.ProtoBaseAccount,
	)

	// The BankKeeper allows you perform sdk.Coins interactions
	app.bankKeeper = bank.NewBaseKeeper(
		app.accountKeeper,
		bankSupspace,
		app.ModuleAccountAddrs(),
	)

	// The SupplyKeeper collects transaction fees and renders them to the fee distribution module
	app.supplyKeeper = supply.NewKeeper(
		app.cdc,
		app.keySupply,
		app.accountKeeper,
		app.bankKeeper,
		maccPerms)

	// The staking keeper
	stakingKeeper := staking.NewKeeper(
		app.cdc,
		app.keyStaking,
		app.supplyKeeper,
		stakingSubspace,
	)

	// register the staking hooks
	// NOTE: stakingKeeper above is passed by reference, so that it will contain these hooks
	app.stakingKeeper = *stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(),
	)

	app.globalKeeper = global.NewKeeper(
		app.keyGlobal,
		app.cdc,
		globalSubspace,
	)

	app.validatorKeeper = validator.NewKeeper(
		app.keyValidator,
		app.cdc,
		app.globalKeeper,
		app.accountKeeper,
		app.stakingKeeper,
		validatorSubspace,
	)

	app.slashKeeper = slash.NewKeeper(
		app.keySlash,
		app.cdc,
		app.validatorKeeper,
		slashSubspace,
	)

	app.subscribeKeeper = subscribe.NewKeeper(
		app.keySubscribe,
		app.cdc,
		app.globalKeeper,
		app.validatorKeeper,
		subscribeSubspace,
	)

	app.cronKeeper = cron.NewKeeper(
		app.keyCron,
		app.cdc,
		app.bankKeeper,
		app.validatorKeeper,
	)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, gov.NewParamChangeProposalHandler(app.paramsKeeper))
	app.govKeeper = gov.NewKeeper(
		app.cdc,
		app.keyGov,
		govSubspace,
		app.validatorKeeper,
		app.slashKeeper,
		govRouter,
	)

	app.syncKeeper = sync.NewKeeper(
		app.cdc,
		app.keySync,
		syncSubspace,
		app.globalKeeper,
		app.paramsKeeper,
		app.slashKeeper,
		app.stakingKeeper,
		app.subscribeKeeper,
		app.validatorKeeper,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		cron.NewAppModule(app.cronKeeper),
		global.NewAppModule(app.globalKeeper),
		slash.NewAppModule(app.slashKeeper),
		subscribe.NewAppModule(app.subscribeKeeper),
		validator.NewAppModule(app.validatorKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper),
		sync.NewAppModule(app.syncKeeper),
	)

	app.mm.SetOrderBeginBlockers(slash.ModuleName)
	app.mm.SetOrderEndBlockers(subscribe.ModuleName, validator.ModuleName, cron.ModuleName, gov.ModuleName, sync.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		genutil.ModuleName,
		cron.ModuleName,
		global.ModuleName,
		slash.ModuleName,
		subscribe.ModuleName,
		validator.ModuleName,
		gov.ModuleName,
		sync.ModuleName,
	)

	// register all module routes and module queriers
	app.mm.RegisterRoutes(app.Router(), app.QueryRouter())

	// The initChainer handles translating the genesis.json file into initial state for the network
	app.SetInitChainer(app.InitChainer)
	app.SetBeginBlocker(app.BeginBlocker)
	app.SetEndBlocker(app.EndBlocker)

	// The AnteHandler handles signature verification and transaction pre-processing
	app.SetAnteHandler(
		auth.NewAnteHandler(
			app.accountKeeper,
			app.supplyKeeper,
			auth.DefaultSigVerificationGasConsumer,
		),
	)

	app.MountStores(
		app.tkeyParams,
		app.tkeyStaking,
		app.keyMain,
		app.keyAccount,
		app.keySupply,
		app.keyStaking,
		app.keyParams,
		app.keyCron,
		app.keyGlobal,
		app.keySlash,
		app.keySubscribe,
		app.keyValidator,
		app.keyGov,
		app.keySync,
	)

	err = app.LoadLatestVersion(app.keyMain)
	if err != nil {
		tmos.Exit(err.Error())
	}

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *sgnApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *sgnApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	if !monitored {
		monitored = true
		go app.startMonitor(ctx)
	}
	return app.mm.BeginBlock(ctx, req)
}

func (app *sgnApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *sgnApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

//_________________________________________________________

func (app *sgnApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
) (appState json.RawMessage, validators []tmtypes.GenesisValidator, err error) {

	// as if they could withdraw from the start of the next block
	ctx := app.NewContext(true, abci.Header{Height: app.LastBlockHeight()})

	genState := app.mm.ExportGenesis(ctx)
	appState, err = codec.MarshalJSONIndent(app.cdc, genState)
	if err != nil {
		return nil, nil, err
	}

	validators = staking.WriteValidators(ctx, app.stakingKeeper)

	return appState, validators, nil
}

// ModuleAccountAddrs returns all the app's module account addresses.
func (app *sgnApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *sgnApp) startMonitor(ctx sdk.Context) {
	if !viper.GetBool(common.FlagStartMonitor) {
		return
	}

	operator, err := transactor.NewTransactor(
		viper.GetString(common.FlagCLIHome),
		ctx.ChainID(),
		viper.GetString(common.FlagSgnNodeURI),
		viper.GetString(common.FlagSgnOperator),
		viper.GetString(common.FlagSgnPassphrase),
		app.cdc,
		transactor.NewGasPriceEstimator(viper.GetString(common.FlagSgnNodeURI)),
	)
	if err != nil {
		tmos.Exit(err.Error())
	}

	monitor.NewEthMonitor(ethClient, operator)
}
