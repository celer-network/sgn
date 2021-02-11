package ethtx

import (
	"context"
	"math/big"
	"net/url"
	"strings"
	"time"

	esClient "github.com/celer-network/eth-services/client"
	esLogger "github.com/celer-network/eth-services/logger"
	esStoreTM "github.com/celer-network/eth-services/store/tendermint"
	esTxManager "github.com/celer-network/eth-services/txmanager"
	esTypes "github.com/celer-network/eth-services/types"
	"github.com/celer-network/sgn/common"
	"github.com/celer-network/sgn/mainchain"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	tmDB "github.com/tendermint/tm-db"
	"go.uber.org/zap"
)

const (
	prefixTxManager = "txm"
)

type TxSender struct {
	ethClient  esClient.Client
	from       gethCommon.Address
	txManager  esTxManager.TxManager
	ledgerInfo *ContractInfo
}

func NewTxSender(
	db tmDB.DB,
	keyDir string,
	from gethCommon.Address,
	ledgerAddress gethCommon.Address,
) (*TxSender, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	rpcURL := viper.GetString(common.FlagEthTxRPCURL)
	parsedRPCURL, err := url.Parse(rpcURL)
	if err != nil {
		return nil, err
	}
	secondaryRPCURLStrs := viper.GetStringSlice(common.FlagEthTxSecondaryRPCURLs)
	parsedSecondaryRPCURLs := make([]*url.URL, len(secondaryRPCURLStrs))
	for i, urlStr := range secondaryRPCURLStrs {
		parsed, parseErr := url.Parse(urlStr)
		if parseErr != nil {
			return nil, parseErr
		}
		parsedSecondaryRPCURLs[i] = parsed
	}
	config := &esTypes.Config{
		Logger:           esLogger.NewZapLogger(logger.Sugar()),
		BlockTime:        time.Duration(viper.GetInt64(common.FlagEthTxBlockTime)) * time.Second,
		RPCURL:           parsedRPCURL,
		SecondaryRPCURLs: parsedSecondaryRPCURLs,
		ChainID:          big.NewInt(viper.GetInt64(common.FlagEthChainID)),

		HeadTrackerHistoryDepth:  viper.GetInt64(common.FlagEthTxHeadTrackerHistoryDepth),
		HeadTrackerMaxBufferSize: viper.GetInt(common.FlagEthTxHeadTrackerMaxBufferSize),
		FinalityDepth:            viper.GetInt64(common.FlagEthTxHeadTrackerFinalityDepth),

		DBPollInterval: time.Duration(viper.GetInt64(common.FlagEthTxDBPollInterval)) * time.Second,

		DefaultGasPrice: big.NewInt(viper.GetInt64(common.FlagEthTxDefaultGasPrice)),
		MaxGasPrice:     big.NewInt(viper.GetInt64(common.FlagEthTxMaxGasPrice)),
		GasBumpPercent:  viper.GetUint64(common.FlagEthTxGasBumpPercent),
		GasBumpWei:      big.NewInt(viper.GetInt64(common.FlagEthTxGasBumpWei)),

		GasBumpThreshold: viper.GetInt64(common.FlagEthTxGasBumpThreshold),
		GasBumpTxDepth:   viper.GetInt(common.FlagEthTxGasBumpTxDepth),
	}
	ethClient, err := esClient.NewImpl(config)
	if err != nil {
		return nil, err
	}
	err = ethClient.Dial(context.Background())
	if err != nil {
		return nil, err
	}

	txManagerDB := tmDB.NewPrefixDB(db, []byte(prefixTxManager))
	store := esStoreTM.NewTMStore(txManagerDB)
	keyStore := esClient.NewKeyStore(keyDir, esClient.ScryptParams{N: keystore.StandardScryptN, P: keystore.StandardScryptP})
	err = keyStore.Unlock(viper.GetString(common.FlagEthPassphrase))
	if err != nil {
		return nil, err
	}

	txManager, err := esTxManager.NewTxManager(ethClient, store, keyStore, config)
	if err != nil {
		return nil, err
	}
	err = txManager.RegisterAccount(from)
	if err != nil {
		return nil, err
	}
	ledgerABI, err := abi.JSON(strings.NewReader(mainchain.CelerLedgerABI))
	if err != nil {
		return nil, err
	}
	ledgerInfo := &ContractInfo{
		ABI:     ledgerABI,
		Address: ledgerAddress,
	}
	return &TxSender{
		ethClient:  ethClient,
		from:       from,
		txManager:  txManager,
		ledgerInfo: ledgerInfo,
	}, nil
}

func (s *TxSender) Start() error {
	return s.txManager.Start()
}

func (s *TxSender) sendTx(
	from gethCommon.Address,
	contract *ContractInfo,
	value *big.Int,
	encodedPayload []byte,
	handler esTxManager.JobHandler,
) error {
	to := contract.Address
	gasLimit, err := s.estimateGas(from, &to, value, encodedPayload)
	if err != nil {
		return err
	}
	txID, err := s.txManager.AddTx(from, to, value, encodedPayload, gasLimit)
	if err != nil {
		return err
	}
	jobID, err := s.txManager.AddJob(txID, nil)
	if err != nil {
		return err
	}
	s.txManager.MonitorJob(jobID, handler)
	return nil
}

func (s *TxSender) estimateGas(
	from gethCommon.Address,
	to *gethCommon.Address,
	value *big.Int,
	data []byte,
) (uint64, error) {
	callMsg := ethereum.CallMsg{
		From:     from,
		To:       to,
		Gas:      0,
		GasPrice: big.NewInt(0),
		Value:    value,
		Data:     data,
	}
	gasLimit, err := s.ethClient.EstimateGas(context.Background(), callMsg)
	return gasLimit, errors.Wrap(err, "estimateGas failed")
}
