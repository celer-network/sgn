package app

import (
	"encoding/json"
	"os"

	"github.com/celer-network/goutils/log"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/monitor"
	"github.com/celer-network/sgn/x/cron"
	"github.com/celer-network/sgn/x/global"
	"github.com/celer-network/sgn/x/gov"
	govclient "github.com/celer-network/sgn/x/gov/client"
	"github.com/celer-network/sgn/x/guard"
	"github.com/celer-network/sgn/x/slash"
	"github.com/celer-network/sgn/x/sync"
	"github.com/celer-network/sgn/x/validator"
	"github.com/cosmos/cosmos-sdk/baseapp"
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
	"github.com/cosmos/cosmos-sdk/x/upgrade"
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

	// DefaultNodeHome sets the folder where the application data and configuration will be stored
	DefaultNodeHome = os.ExpandEnv("$HOME/.sgnd")

	// ModuleBasicManager is in charge of setting up basic module elemnets
	ModuleBasics = module.NewBasicManager(
		genutil.AppModuleBasic{},
		auth.AppModuleBasic{},
		bank.AppModuleBasic{},
		staking.AppModuleBasic{},
		params.AppModuleBasic{},
		supply.AppModuleBasic{},
		upgrade.AppModuleBasic{},

		cron.AppModule{},
		global.AppModule{},
		gov.NewAppModuleBasic(govclient.ParamProposalHandler, govclient.UpgradeProposalHandler),
		slash.AppModule{},
		guard.AppModule{},
		sync.AppModule{},
		validator.AppModuleBasic{},
	)
	// account permissions
	maccPerms = map[string][]string{
		auth.FeeCollectorName:     nil,
		staking.BondedPoolName:    {supply.Burner, supply.Staking},
		staking.NotBondedPoolName: {supply.Burner, supply.Staking},
	}
)

// MakeCodec generates the necessary codecs for Amino
func MakeCodec() *codec.Codec {
	var cdc = codec.New()
	ModuleBasics.RegisterCodec(cdc)
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	return cdc
}

type SgnApp struct {
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
	keyUpgrade   *sdk.KVStoreKey
	keyCron      *sdk.KVStoreKey
	keyGlobal    *sdk.KVStoreKey
	keyGov       *sdk.KVStoreKey
	keySlash     *sdk.KVStoreKey
	keyGuard     *sdk.KVStoreKey
	keySync      *sdk.KVStoreKey
	keyValidator *sdk.KVStoreKey

	// Keepers
	accountKeeper   auth.AccountKeeper
	bankKeeper      bank.Keeper
	stakingKeeper   staking.Keeper
	supplyKeeper    supply.Keeper
	paramsKeeper    params.Keeper
	upgradeKeeper   upgrade.Keeper
	cronKeeper      cron.Keeper
	globalKeeper    global.Keeper
	govKeeper       gov.Keeper
	slashKeeper     slash.Keeper
	guardKeeper     guard.Keeper
	syncKeeper      sync.Keeper
	validatorKeeper validator.Keeper

	// Module Manager
	mm *module.Manager
}

// NewSgnApp is a constructor function for sgnApp
func NewSgnApp(
	logger tlog.Logger,
	db dbm.DB,
	height int64,
	skipUpgradeHeights map[int64]bool,
	baseAppOptions ...func(*bam.BaseApp),
) *SgnApp {
	viper.SetDefault(common.FlagEthPollInterval, 15)
	viper.SetDefault(common.FlagEthBlockDelay, 5)
	viper.SetDefault(common.FlagSgnCheckIntervalGuardQueue, 60)
	viper.SetDefault(common.FlagSgnCheckIntervalSlashQueue, 60)

	err := common.SetupUserPassword()
	if err != nil {
		tmos.Exit(err.Error())
	}

	loglevel := viper.GetString(common.FlagLogLevel)
	log.SetLevelByName(loglevel)
	if loglevel == "trace" {
		baseAppOptions = append(baseAppOptions, baseapp.SetTrace(true))
	}
	if viper.GetBool(common.FlagLogColor) {
		log.EnableColor()
	}

	// First define the top level codec that will be shared by the different modules
	cdc := MakeCodec()

	// BaseApp handles interactions with Tendermint through the ABCI protocol
	bApp := bam.NewBaseApp(appName, logger, db, auth.DefaultTxDecoder(cdc), baseAppOptions...)

	// Here you initialize your application with the store keys it requires
	var app = &SgnApp{
		BaseApp: bApp,
		cdc:     cdc,

		keyMain:      sdk.NewKVStoreKey(bam.MainStoreKey),
		keyAccount:   sdk.NewKVStoreKey(auth.StoreKey),
		keySupply:    sdk.NewKVStoreKey(supply.StoreKey),
		keyStaking:   sdk.NewKVStoreKey(staking.StoreKey),
		tkeyStaking:  sdk.NewTransientStoreKey(staking.TStoreKey),
		keyParams:    sdk.NewKVStoreKey(params.StoreKey),
		tkeyParams:   sdk.NewTransientStoreKey(params.TStoreKey),
		keyUpgrade:   sdk.NewKVStoreKey(upgrade.StoreKey),
		keyCron:      sdk.NewKVStoreKey(cron.StoreKey),
		keyGov:       sdk.NewKVStoreKey(gov.StoreKey),
		keyGlobal:    sdk.NewKVStoreKey(global.StoreKey),
		keySlash:     sdk.NewKVStoreKey(slash.StoreKey),
		keyGuard:     sdk.NewKVStoreKey(guard.StoreKey),
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
	guardSubspace := app.paramsKeeper.Subspace(guard.DefaultParamspace)
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
		app.cdc,
		app.keyGlobal,
		globalSubspace,
	)

	app.validatorKeeper = validator.NewKeeper(
		app.keyValidator,
		app.cdc,
		app.accountKeeper,
		app.bankKeeper,
		app.stakingKeeper,
		validatorSubspace,
	)

	app.upgradeKeeper = upgrade.NewKeeper(skipUpgradeHeights, app.keyUpgrade, app.cdc)

	app.slashKeeper = slash.NewKeeper(
		app.keySlash,
		app.cdc,
		app.globalKeeper,
		app.validatorKeeper,
		slashSubspace,
	)

	app.guardKeeper = guard.NewKeeper(
		app.keyGuard,
		app.cdc,
		app.validatorKeeper,
		guardSubspace,
	)

	app.cronKeeper = cron.NewKeeper(
		app.keyCron,
		app.cdc,
		app.bankKeeper,
		app.validatorKeeper,
	)

	govRouter := gov.NewRouter()
	govRouter.AddRoute(gov.RouterKey, gov.ProposalHandler).
		AddRoute(params.RouterKey, gov.NewParamChangeProposalHandler(app.paramsKeeper)).
		AddRoute(upgrade.RouterKey, gov.NewUpgradeProposalHandler(app.upgradeKeeper))
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
		app.paramsKeeper,
		app.bankKeeper,
		app.globalKeeper,
		app.slashKeeper,
		app.stakingKeeper,
		app.guardKeeper,
		app.validatorKeeper,
	)

	app.mm = module.NewManager(
		genutil.NewAppModule(app.accountKeeper, app.stakingKeeper, app.BaseApp.DeliverTx),
		auth.NewAppModule(app.accountKeeper),
		bank.NewAppModule(app.bankKeeper, app.accountKeeper),
		supply.NewAppModule(app.supplyKeeper, app.accountKeeper),
		staking.NewAppModule(app.stakingKeeper, app.accountKeeper, app.supplyKeeper),
		upgrade.NewAppModule(app.upgradeKeeper),
		cron.NewAppModule(app.cronKeeper),
		global.NewAppModule(app.globalKeeper),
		slash.NewAppModule(app.slashKeeper),
		guard.NewAppModule(app.guardKeeper),
		validator.NewAppModule(app.validatorKeeper),
		gov.NewAppModule(app.govKeeper, app.accountKeeper),
		sync.NewAppModule(app.syncKeeper),
	)

	app.mm.SetOrderBeginBlockers(upgrade.ModuleName, slash.ModuleName)
	app.mm.SetOrderEndBlockers(guard.ModuleName, validator.ModuleName, cron.ModuleName, gov.ModuleName, sync.ModuleName)

	// Sets the order of Genesis - Order matters, genutil is to always come last
	app.mm.SetOrderInitGenesis(
		staking.ModuleName,
		auth.ModuleName,
		bank.ModuleName,
		supply.ModuleName,
		genutil.ModuleName,
		cron.ModuleName,
		global.ModuleName,
		guard.ModuleName,
		slash.ModuleName,
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
		app.keyUpgrade,
		app.keyCron,
		app.keyGlobal,
		app.keySlash,
		app.keyGuard,
		app.keyValidator,
		app.keyGov,
		app.keySync,
	)

	if height == -1 {
		err = app.LoadLatestVersion()
	} else {
		err = app.LoadHeight(height)
	}
	if err != nil {
		tmos.Exit("Failed to load height:" + err.Error())
	}

	enableTxSender := viper.GetBool(common.FlagEthTxEnableTxSender)
	go app.startMonitor(db, enableTxSender)

	return app
}

// GenesisState represents chain state at the start of the chain. Any initial state (account balances) are stored here.
type GenesisState map[string]json.RawMessage

func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}

func (app *SgnApp) InitChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	var genesisState GenesisState

	err := app.cdc.UnmarshalJSON(req.AppStateBytes, &genesisState)
	if err != nil {
		panic(err)
	}

	return app.mm.InitGenesis(ctx, genesisState)
}

func (app *SgnApp) BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock) abci.ResponseBeginBlock {
	return app.mm.BeginBlock(ctx, req)
}

func (app *SgnApp) EndBlocker(ctx sdk.Context, req abci.RequestEndBlock) abci.ResponseEndBlock {
	return app.mm.EndBlock(ctx, req)
}

func (app *SgnApp) LoadLatestVersion() error {
	return app.BaseApp.LoadLatestVersion(app.keyMain)
}

func (app *SgnApp) LoadHeight(height int64) error {
	return app.LoadVersion(height, app.keyMain)
}

//_________________________________________________________

func (app *SgnApp) ExportAppStateAndValidators(forZeroHeight bool, jailWhiteList []string,
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
func (app *SgnApp) ModuleAccountAddrs() map[string]bool {
	modAccAddrs := make(map[string]bool)
	for acc := range maccPerms {
		modAccAddrs[supply.NewModuleAddress(acc).String()] = true
	}

	return modAccAddrs
}

func (app *SgnApp) startMonitor(db dbm.DB, useTxSender bool) {
	var txDB dbm.DB
	if useTxSender {
		txDB = db
	}
	operator, err := monitor.NewOperator(app.cdc, viper.GetString(common.FlagCLIHome), txDB)
	if err != nil {
		tmos.Exit(err.Error())
	}

	err = common.WaitTillHeight(operator.Transactor.CliCtx, 1)
	if err != nil {
		tmos.Exit(err.Error())
	}

	monitor.NewMonitor(operator, db)
}
