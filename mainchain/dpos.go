// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mainchain

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// DPoSABI is the input ABI used to generate the binding from.
const DPoSABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addWhitelisted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextParamProposalId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"removeWhitelisted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"nextSidechainProposalId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_record\",\"type\":\"uint256\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"createParamProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_sidechainAddr\",\"type\":\"address\"}],\"name\":\"isSidechainRegistered\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"checkedValidators\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dposGoLiveTime\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isWhitelisted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isPauser\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"registeredSidechains\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceWhitelistAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"governToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\"},{\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getParamProposalVote\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_record\",\"type\":\"uint256\"}],\"name\":\"getUIntValue\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"UIntStorage\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renouncePauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"miningPool\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addWhitelistAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"paramProposals\",\"outputs\":[{\"name\":\"proposer\",\"type\":\"address\"},{\"name\":\"deposit\",\"type\":\"uint256\"},{\"name\":\"voteDeadline\",\"type\":\"uint256\"},{\"name\":\"record\",\"type\":\"uint256\"},{\"name\":\"newValue\",\"type\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"addPauser\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\"},{\"name\":\"_voter\",\"type\":\"address\"}],\"name\":\"getSidechainProposalVote\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"redeemedMiningReward\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"registerSidechain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"account\",\"type\":\"address\"}],\"name\":\"isWhitelistAdmin\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"COMMISSION_RATE_BASE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"celerToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"enableWhitelist\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceWhitelisted\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sidechainAddr\",\"type\":\"address\"},{\"name\":\"_registered\",\"type\":\"bool\"}],\"name\":\"createSidechainProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"sidechainProposals\",\"outputs\":[{\"name\":\"proposer\",\"type\":\"address\"},{\"name\":\"deposit\",\"type\":\"uint256\"},{\"name\":\"voteDeadline\",\"type\":\"uint256\"},{\"name\":\"sidechainAddr\",\"type\":\"address\"},{\"name\":\"registered\",\"type\":\"bool\"},{\"name\":\"status\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"usedPenaltyNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_celerTokenAddress\",\"type\":\"address\"},{\"name\":\"_governProposalDeposit\",\"type\":\"uint256\"},{\"name\":\"_governVoteTimeout\",\"type\":\"uint256\"},{\"name\":\"_blameTimeout\",\"type\":\"uint256\"},{\"name\":\"_minValidatorNum\",\"type\":\"uint256\"},{\"name\":\"_maxValidatorNum\",\"type\":\"uint256\"},{\"name\":\"_minStakeInPool\",\"type\":\"uint256\"},{\"name\":\"_increaseRateWaitTime\",\"type\":\"uint256\"},{\"name\":\"_dposGoLiveTimeout\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"WhitelistedAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"WhitelistedRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"WhitelistAdminAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"WhitelistAdminRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Paused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"Unpaused\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"voteDeadline\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"record\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"newValue\",\"type\":\"uint256\"}],\"name\":\"CreateParamProposal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"voteType\",\"type\":\"uint8\"}],\"name\":\"VoteParam\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"passed\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"record\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"newValue\",\"type\":\"uint256\"}],\"name\":\"ConfirmParamProposal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"deposit\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"voteDeadline\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sidechainAddr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"registered\",\"type\":\"bool\"}],\"name\":\"CreateSidechainProposal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"voter\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"voteType\",\"type\":\"uint8\"}],\"name\":\"VoteSidechain\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"proposalId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"passed\",\"type\":\"bool\"},{\"indexed\":false,\"name\":\"sidechainAddr\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"registered\",\"type\":\"bool\"}],\"name\":\"ConfirmSidechainProposal\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"account\",\"type\":\"address\"}],\"name\":\"PauserRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"minSelfStake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"rateLockEndTime\",\"type\":\"uint256\"}],\"name\":\"InitializeCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"announcedRate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"announcedLockEndTime\",\"type\":\"uint256\"}],\"name\":\"CommissionRateAnnouncement\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newRate\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"newLockEndTime\",\"type\":\"uint256\"}],\"name\":\"UpdateCommissionRate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newStake\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"stakingPool\",\"type\":\"uint256\"}],\"name\":\"Delegate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"ethAddr\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"changeType\",\"type\":\"uint8\"}],\"name\":\"ValidatorChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"WithdrawFromUnbondedCandidate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"withdrawAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"proposedTime\",\"type\":\"uint256\"}],\"name\":\"IntendWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ConfirmWithdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"delegator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Punish\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"indemnitee\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Indemnify\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"candidate\",\"type\":\"address\"}],\"name\":\"CandidateUnbonded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"receiver\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"reward\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"miningPool\",\"type\":\"uint256\"}],\"name\":\"RedeemMiningReward\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"contributor\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"contribution\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"miningPoolSize\",\"type\":\"uint256\"}],\"name\":\"MiningPoolContribution\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"_enable\",\"type\":\"bool\"}],\"name\":\"updateEnableWhitelist\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"drainToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\"},{\"name\":\"_vote\",\"type\":\"uint8\"}],\"name\":\"voteParam\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\"}],\"name\":\"confirmParamProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\"},{\"name\":\"_vote\",\"type\":\"uint8\"}],\"name\":\"voteSidechain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_proposalId\",\"type\":\"uint256\"}],\"name\":\"confirmSidechainProposal\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"contributeToMiningPool\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_receiver\",\"type\":\"address\"},{\"name\":\"_cumulativeReward\",\"type\":\"uint256\"}],\"name\":\"redeemMiningReward\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_minSelfStake\",\"type\":\"uint256\"},{\"name\":\"_commissionRate\",\"type\":\"uint256\"},{\"name\":\"_rateLockEndTime\",\"type\":\"uint256\"}],\"name\":\"initializeCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newRate\",\"type\":\"uint256\"},{\"name\":\"_newLockEndTime\",\"type\":\"uint256\"}],\"name\":\"nonIncreaseCommissionRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newRate\",\"type\":\"uint256\"},{\"name\":\"_newLockEndTime\",\"type\":\"uint256\"}],\"name\":\"announceIncreaseCommissionRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"confirmIncreaseCommissionRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"delegate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"claimValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"confirmUnbondedCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFromUnbondedCandidate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"intendWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"confirmWithdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_penaltyRequest\",\"type\":\"bytes\"}],\"name\":\"punish\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_request\",\"type\":\"bytes\"}],\"name\":\"validateMultiSigMessage\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isValidDPoS\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidatorNum\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinStakingPool\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"}],\"name\":\"getCandidateInfo\",\"outputs\":[{\"name\":\"initialized\",\"type\":\"bool\"},{\"name\":\"minSelfStake\",\"type\":\"uint256\"},{\"name\":\"stakingPool\",\"type\":\"uint256\"},{\"name\":\"status\",\"type\":\"uint256\"},{\"name\":\"unbondTime\",\"type\":\"uint256\"},{\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"name\":\"rateLockEndTime\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_candidateAddr\",\"type\":\"address\"},{\"name\":\"_delegatorAddr\",\"type\":\"address\"}],\"name\":\"getDelegatorInfo\",\"outputs\":[{\"name\":\"delegatedStake\",\"type\":\"uint256\"},{\"name\":\"undelegatingStake\",\"type\":\"uint256\"},{\"name\":\"intentAmounts\",\"type\":\"uint256[]\"},{\"name\":\"intentProposedTimes\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getMinQuorumStakingPool\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getTotalValidatorStakingPool\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// DPoSBin is the compiled bytecode used for deploying new contracts.
var DPoSBin = "0x60806040523480156200001157600080fd5b50604051620050413803806200504183398181016040526101208110156200003857600080fd5b508051602082015160408301516060840151608085015160a086015160c087015160e088015161010090980151969795969495939492939192909190888888888888888862000090336001600160e01b036200025e16565b600180546001600160a01b0319163317908190556040516001600160a01b0391909116906000907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908290a36001805460ff60a01b19169055620000fd336001600160e01b03620002b016565b600480546001600160a01b03199081166001600160a01b039a8b1617909155600560209081527f05b8ccbb9d4d8fb16ea74ce3c29a41f1b461fbdaff4714a0d9a8eb05499746bc989098557f1471eb6eb2c5e789fc3de43f8ce62938c7d1836ec861730447e2ada8fd81017b969096557f89832631fb3c3307a103ba2c84ab569c64d6182a18893dcd163f0f1c2090733a949094557fa9bc9a3a348c357ba16b37005d7e6b3236198c0e939f4af8c5f19b8deeb8ebc0929092557f3eec716f11ba9e820c81ca75eb978ffb45831ef8b7a53e5e422c26008e1ca6d5557f458b30c2d72bfd2c6317304a4594ecbafe5f729d3111b65fdc3a33bd48e5432d5560066000527f069400f22b28c6c362558d92f66163cec5671cba50b61abd2eecfcd0eaeac5185560108054909116928c16929092179091556200024b904390839062000302811b6200319617901c565b60115550620003ab975050505050505050565b620002798160006200031c60201b620045e11790919060201c565b6040516001600160a01b038216907f6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f890600090a250565b620002cb8160026200031c60201b620045e11790919060201c565b6040516001600160a01b038216907f22380c05984257a1cb900161c713dd71d39e74820f1aea43bd3f1bdd2096129990600090a250565b6000828201838110156200031557600080fd5b9392505050565b6001600160a01b0381166200033057600080fd5b6200034582826001600160e01b036200037516565b156200035057600080fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b60006001600160a01b0382166200038b57600080fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b614c8680620003bb6000396000f3fe608060405234801561001057600080fd5b50600436106103fc5760003560e01c80637362d9c811610215578063be57959d11610125578063e433c1ca116100b8578063eecefef811610087578063eecefef814610c29578063f2fde38b14610cfe578063f64f33f214610d24578063facd743b14610d47578063fb87874914610d6d576103fc565b8063e433c1ca14610b7b578063e64808f314610ba9578063e97b745214610bc6578063eab2ed8c14610c21576103fc565b8063c7ec2f35116100f4578063c7ec2f3514610b1f578063cdfb2b4e14610b45578063d2bfc1c714610b4d578063d6cd947314610b73576103fc565b8063be57959d14610ac3578063bee8380e14610ae6578063c1e1671814610aee578063c6c21e9d14610b17576103fc565b806389ed7939116101a8578063934a18ec11610177578063934a18ec14610a26578063a3e814b914610a43578063aa09fbae14610a4b578063bb5f747b14610a71578063bb9053d014610a97576103fc565b806389ed793914610a065780638da5cb5b14610a0e5780638e9472a314610a165780638f32d59b14610a1e576103fc565b80638456cb59116101e45780638456cb59146109865780638515b0e21461098e57806385bfe017146109ba57806387e53fef146109e0576103fc565b80637362d9c814610898578063785f8ffd146108be5780637e5fb8f3146108ea57806382dc1ec414610960576103fc565b806339c9563e11610310578063581c53c5116102a35780636e7cf85d116102725780636e7cf85d146108705780636e997565146108785780636ef8d66d14610880578063715018a6146108885780637339759714610890576103fc565b8063581c53c5146107de5780635c975abb1461082e57806364c663951461083657806364ed600a14610853576103fc565b806349444b71116102df57806349444b711461076f5780634b7dba6b146107955780634c5a628c146107b257806351abe57b146107ba576103fc565b806339c9563e146107135780633af32abf1461071b5780633f4ba83a1461074157806346fbf68e14610749576103fc565b806325ed6b35116103935780632cb57c48116103625780632cb57c48146106155780633090c0e914610634578063325820b3146106575780633620d1491461067d5780633702db39146106ed576103fc565b806325ed6b351461056157806328bde1e114610587578063291d9549146105e75780632bf0fe591461060d576103fc565b80631c0efd9d116103cf5780631c0efd9d1461048f5780631cfe4f0b146105135780631f7b08861461052d57806322da792714610559576103fc565b8063026e402b1461040157806310154bad1461042f578063145aa116146104555780631a06f73714610472575b600080fd5b61042d6004803603604081101561041757600080fd5b506001600160a01b038135169060200135610d8a565b005b61042d6004803603602081101561044557600080fd5b50356001600160a01b0316610ecf565b61042d6004803603602081101561046b57600080fd5b5035610eed565b61042d6004803603602081101561048857600080fd5b5035610f31565b6104ff600480360360208110156104a557600080fd5b8101906020810181356401000000008111156104c057600080fd5b8201836020820111156104d257600080fd5b803590602001918460018302840111640100000000831117156104f457600080fd5b509092509050610fe6565b604080519115158252519081900360200190f35b61051b6110c1565b60408051918252519081900360200190f35b61042d6004803603604081101561054357600080fd5b506001600160a01b03813516906020013561110f565b61051b611200565b61042d6004803603604081101561057757600080fd5b508035906020013560ff16611206565b6105ad6004803603602081101561059d57600080fd5b50356001600160a01b0316611271565b6040805197151588526020880196909652868601949094526060860192909252608085015260a084015260c0830152519081900360e00190f35b61042d600480360360208110156105fd57600080fd5b50356001600160a01b03166112db565b61051b6112f6565b61042d6004803603602081101561062b57600080fd5b503515156112fc565b61042d6004803603604081101561064a57600080fd5b5080359060200135611320565b6104ff6004803603602081101561066d57600080fd5b50356001600160a01b031661146e565b61042d6004803603602081101561069357600080fd5b8101906020810181356401000000008111156106ae57600080fd5b8201836020820111156106c057600080fd5b803590602001918460018302840111640100000000831117156106e257600080fd5b509092509050611490565b6104ff6004803603602081101561070357600080fd5b50356001600160a01b03166119a2565b61051b6119b7565b6104ff6004803603602081101561073157600080fd5b50356001600160a01b03166119bd565b61042d6119d0565b6104ff6004803603602081101561075f57600080fd5b50356001600160a01b0316611a3a565b6104ff6004803603602081101561078557600080fd5b50356001600160a01b0316611a4c565b61042d600480360360208110156107ab57600080fd5b5035611a61565b61042d611afa565b6107c2611b05565b604080516001600160a01b039092168252519081900360200190f35b61080a600480360360408110156107f457600080fd5b50803590602001356001600160a01b0316611b14565b6040518082600381111561081a57fe5b60ff16815260200191505060405180910390f35b6104ff611b42565b61051b6004803603602081101561084c57600080fd5b5035611b52565b61051b6004803603602081101561086957600080fd5b5035611b64565b61042d611b76565b61042d611edb565b61042d611fad565b61042d611fb6565b61051b612011565b61042d600480360360208110156108ae57600080fd5b50356001600160a01b0316612017565b61042d600480360360408110156108d457600080fd5b506001600160a01b038135169060200135612032565b6109076004803603602081101561090057600080fd5b5035612152565b60405180876001600160a01b03166001600160a01b0316815260200186815260200185815260200184815260200183815260200182600281111561094757fe5b60ff168152602001965050505050505060405180910390f35b61042d6004803603602081101561097657600080fd5b50356001600160a01b0316612194565b61042d6121af565b61080a600480360360408110156109a457600080fd5b50803590602001356001600160a01b0316612220565b61042d600480360360408110156109d057600080fd5b508035906020013560ff1661224f565b61051b600480360360208110156109f657600080fd5b50356001600160a01b03166122b5565b61051b6122c7565b6107c2612325565b61051b612334565b6104ff612415565b61042d60048036036020811015610a3c57600080fd5b5035612426565b61051b6124d4565b61042d60048036036020811015610a6157600080fd5b50356001600160a01b0316612516565b6104ff60048036036020811015610a8757600080fd5b50356001600160a01b031661254b565b61042d60048036036040811015610aad57600080fd5b506001600160a01b03813516906020013561255e565b61042d60048036036040811015610ad957600080fd5b5080359060200135612659565b61051b612703565b61042d60048036036060811015610b0457600080fd5b5080359060208101359060400135612709565b6107c261284d565b61042d60048036036020811015610b3557600080fd5b50356001600160a01b031661285c565b6104ff6128ee565b61042d60048036036020811015610b6357600080fd5b50356001600160a01b03166128f7565b61042d612aff565b61042d60048036036040811015610b9157600080fd5b506001600160a01b0381351690602001351515612b08565b6107c260048036036020811015610bbf57600080fd5b5035612c5b565b610be360048036036020811015610bdc57600080fd5b5035612c76565b604080516001600160a01b038089168252602082018890529181018690529084166060820152821515608082015260a0810182600281111561094757fe5b6104ff612cbe565b610c5760048036036040811015610c3f57600080fd5b506001600160a01b0381358116916020013516612cec565b604051808581526020018481526020018060200180602001838103835285818151815260200191508051906020019060200280838360005b83811015610ca7578181015183820152602001610c8f565b50505050905001838103825284818151815260200191508051906020019060200280838360005b83811015610ce6578181015183820152602001610cce565b50505050905001965050505050505060405180910390f35b61042d60048036036020811015610d1457600080fd5b50356001600160a01b0316612e2b565b61042d60048036036040811015610d3a57600080fd5b5080359060200135612e45565b6104ff60048036036020811015610d5d57600080fd5b50356001600160a01b0316612f36565b6104ff60048036036020811015610d8357600080fd5b5035612f6b565b600154600160a01b900460ff1615610da157600080fd5b816001600160a01b038116610de9576040805162461bcd60e51b815260206004820152600960248201526830206164647265737360b81b604482015290519081900360640190fd5b6001600160a01b0383166000908152600e60205260409020805460ff16610e45576040805162461bcd60e51b815260206004820152601c6024820152600080516020614c32833981519152604482015290519081900360640190fd5b33610e538282866000612f80565b601054610e71906001600160a01b031682308763ffffffff61302916565b846001600160a01b0316816001600160a01b03167f500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b868560020154604051808381526020018281526020019250505060405180910390a35050505050565b610ed83361254b565b610ee157600080fd5b610eea816130b9565b50565b600154600160a01b900460ff16610f0357600080fd5b610f0b612415565b610f1457600080fd5b601054610eea906001600160a01b0316338363ffffffff61310116565b6000610f3e60045b611b52565b90506000805b82811015610fc65760016000828152600b6020526040902054610f719086906001600160a01b0316612220565b6003811115610f7c57fe5b1415610fbe576000818152600b60209081526040808320546001600160a01b03168352600e909152902060020154610fbb90839063ffffffff61319616565b91505b600101610f44565b506000610fd16124d4565b8210159050610fe084826131af565b50505050565b6000610ff13361146e565b610ffa57600080fd5b611002614b3b565b61104184848080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061335992505050565b9050600081600001516040518082805190602001908083835b602083106110795780518252601f19909201916020918201910161105a565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506110b68183602001516134b0565b925050505b92915050565b6000806110ce6004610f39565b90506000805b82811015611108576000818152600b60205260409020546001600160a01b031615611100576001909101905b6001016110d4565b5091505090565b600154600160a01b900460ff161561112657600080fd5b61112f3361146e565b61113857600080fd5b6001600160a01b0382166000908152600f602052604081205461116290839063ffffffff61370e16565b6001600160a01b0384166000908152600f60205260409020839055601254909150611193908263ffffffff61370e16565b6012556010546111b3906001600160a01b0316848363ffffffff61310116565b60125460408051838152602081019290925280516001600160a01b038616927fc243dafa8ee55923dad771198c225cf6dfcdc5e405eda7d4da42b6c6fa018de792908290030190a2505050565b60075481565b3361121081612f36565b611261576040805162461bcd60e51b815260206004820152601d60248201527f6d73672073656e646572206973206e6f7420612076616c696461746f72000000604482015290519081900360640190fd5b61126c838284613723565b505050565b6001600160a01b0381166000908152600e6020526040812080546001820154600280840154600485015460ff9485169693959194849384938493909216908111156112b857fe5b945080600501549350806006015492508060070154915050919395979092949650565b6112e43361254b565b6112ed57600080fd5b610eea816138f2565b600a5481565b611304612415565b61130d57600080fd5b6013805460ff1916911515919091179055565b60075460008181526006602052604090209061134390600163ffffffff61319616565b6007556000808052600560208190527f05b8ccbb9d4d8fb16ea74ce3c29a41f1b461fbdaff4714a0d9a8eb05499746bc5483546001600160a01b0319163390811785556001808601839055909391926113b79290915b8152602001908152602001600020544361319690919063ffffffff16565b600284015560038301859055600480840185905560058401805460ff19166001179055546113f0906001600160a01b0316833084613029565b6007547f40109a070319d6004f4e4b31dba4b605c97bd3474d49865158f55fe093e3b3399061142690600163ffffffff61370e16565b6002850154604080519283526001600160a01b038616602084015282810185905260608301919091526080820188905260a08201879052519081900360c00190a15050505050565b6001600160a01b03811660009081526008602052604090205460ff165b919050565b600154600160a01b900460ff16156114a757600080fd5b6114af612cbe565b6114f4576040805162461bcd60e51b815260206004820152601160248201527011141bd4c81a5cc81b9bdd081d985b1a59607a1b604482015290519081900360640190fd5b6114fc614b3b565b61153b83838080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061393a92505050565b9050611545614b55565b815161155090613a88565b9050600082600001516040518082805190602001908083835b602083106115885780518252601f199092019160209182019101611569565b6001836020036101000a038019825116818451168082178552505050505050905001915050604051809103902090506115c58184602001516134b0565b611616576040805162461bcd60e51b815260206004820152601c60248201527f4661696c20746f20636865636b2076616c696461746f72207369677300000000604482015290519081900360640190fd5b815167ffffffffffffffff166000908152600c602052604090205460ff161561167b576040805162461bcd60e51b8152602060048201526012602482015271557365642070656e616c7479206e6f6e636560701b604482015290519081900360640190fd5b816020015167ffffffffffffffff1643106116cf576040805162461bcd60e51b815260206004820152600f60248201526e14195b985b1d1e48195e1c1a5c9959608a1b604482015290519081900360640190fd5b815167ffffffffffffffff166000908152600c60209081526040808320805460ff19166001179055808501516001600160a01b03168352600e909152812090805b84606001515181101561184d57611725614b82565b8560600151828151811061173557fe5b6020026020010151905061175681602001518461319690919063ffffffff16565b925080600001516001600160a01b031686604001516001600160a01b03167f111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a83602001516040518082815260200191505060405180910390a380516001600160a01b031660009081526003850160209081526040909120908201518154106117f2576117ed85836000015184602001516001612f80565b611843565b8054602083015160009161180c919063ffffffff61370e16565b6001830154909150611824908263ffffffff61370e16565b826001018190555061184186846000015184600001546001612f80565b505b5050600101611710565b5061185b8460400151613d04565b6000805b85608001515181101561194c57611874614b82565b8660800151828151811061188457fe5b602002602001015190506118a581602001518461319690919063ffffffff16565b925080600001516001600160a01b03167f5739ab1e78cc400665fe049852d266e9ac5995451b5e666b25f3c30b19087d8e82602001516040518082815260200191505060405180910390a280516001600160a01b031661191e5760208101516012546119169163ffffffff61319616565b601255611943565b80516020820151601054611943926001600160a01b039091169163ffffffff61310116565b5060010161185f565b50808214611998576040805162461bcd60e51b8152602060048201526014602482015273082dadeeadce840c8decae6dc4ee840dac2e8c6d60631b604482015290519081900360640190fd5b5050505050505050565b600d6020526000908152604090205460ff1681565b60115481565b60006110bb60038363ffffffff613d9416565b6119d933611a3a565b6119e257600080fd5b600154600160a01b900460ff166119f857600080fd5b6001805460ff60a01b191690556040805133815290517f5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa9181900360200190a1565b60006110bb818363ffffffff613d9416565b60086020526000908152604090205460ff1681565b600154600160a01b900460ff1615611a7857600080fd5b6012543390611a8d908363ffffffff61319616565b601255601054611aae906001600160a01b031682308563ffffffff61302916565b60125460408051848152602081019290925280516001600160a01b038416927f97e19c4040b6c46d4275e0c4fea68f8f92c81138372ffdb089932c211938f76592908290030190a25050565b611b0333613dc9565b565b6004546001600160a01b031681565b60008281526006602081815260408084206001600160a01b0386168552909201905290205460ff1692915050565b600154600160a01b900460ff1690565b60009081526005602052604090205490565b60056020526000908152604090205481565b336000818152600e60205260409020805460ff16611bc9576040805162461bcd60e51b815260206004820152601c6024820152600080516020614c32833981519152604482015290519081900360640190fd5b6000600482015460ff166002811115611bde57fe5b1480611bfc57506002600482015460ff166002811115611bfa57fe5b145b611c0557600080fd5b6000611c116005610f39565b90508082600201541015611c6c576040805162461bcd60e51b815260206004820152601960248201527f496e73756666696369656e74207374616b696e6720706f6f6c00000000000000604482015290519081900360640190fd5b60018201546001600160a01b03841660009081526003840160205260409020541015611cd7576040805162461bcd60e51b81526020600482015260156024820152744e6f7420656e6f7567682073656c66207374616b6560581b604482015290519081900360640190fd5b7fdf7de25b7f1fd6d0b5205f0e18f1f35bd7b8d84cce336588d184533ce43a6f76546001600160a01b039081166000818152600e60209081526040822060020154828052600b909152909290919086161415611d75576040805162461bcd60e51b8152602060048201526018602482015277105b1c9958591e481a5b881d985b1a59185d1bdc881cd95d60421b604482015290519081900360640190fd5b6000611d816004610f39565b905060015b81811015611e5e576000818152600b60205260409020546001600160a01b0388811691161415611df8576040805162461bcd60e51b8152602060048201526018602482015277105b1c9958591e481a5b881d985b1a59185d1bdc881cd95d60421b604482015290519081900360640190fd5b6000818152600b60209081526040808320546001600160a01b03168352600e909152902060020154831115611e56576000818152600b60209081526040808320546001600160a01b03168352600e9091529020600201549093509150825b600101611d86565b5081856002015411611ea15760405162461bcd60e51b8152600401808060200182810382526021815260200180614bd76021913960400191505060405180910390fd5b6000838152600b60205260409020546001600160a01b03168015611ec857611ec884613e11565b611ed28785613edf565b50505050505050565b336000908152600e60205260409020805460ff16611f2e576040805162461bcd60e51b815260206004820152601c6024820152600080516020614c32833981519152604482015290519081900360640190fd5b6000611f3a6006610f39565b90508082600a0154014311611f96576040805162461bcd60e51b815260206004820152601c60248201527f6e65772072617465206861736e27742074616b656e2065666665637400000000604482015290519081900360640190fd5b611fa98283600801548460090154613f78565b5050565b611b0333614123565b611fbe612415565b611fc757600080fd5b6001546040516000916001600160a01b0316907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0908390a3600180546001600160a01b0319169055565b60125481565b6120203361254b565b61202957600080fd5b610eea8161416b565b816001600160a01b03811661207a576040805162461bcd60e51b815260206004820152600960248201526830206164647265737360b81b604482015290519081900360640190fd5b6001600160a01b0383166000908152600e60209081526040808320338085526003820190935292209091906120b28284876001612f80565b60018101546120c7908663ffffffff61319616565b60018201556120d586613d04565b6004810180546000908152600283016020908152604091829020888155436001808301918255855401909455925482518981529182015281516001600160a01b03808b1693908816927f7171946bb2a9ef55fcb2eb8cef679db45e2e3a8cef9a44567d34d202b65ff0b1929081900390910190a350505050505050565b6006602052600090815260409020805460018201546002830154600384015460048501546005909501546001600160a01b039094169492939192909160ff1686565b61219d33611a3a565b6121a657600080fd5b610eea816141b3565b6121b833611a3a565b6121c157600080fd5b600154600160a01b900460ff16156121d857600080fd5b6001805460ff60a01b1916600160a01b1790556040805133815290517f62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a2589181900360200190a1565b60008281526009602090815260408083206001600160a01b038516845260040190915290205460ff1692915050565b3361225981612f36565b6122aa576040805162461bcd60e51b815260206004820152601d60248201527f6d73672073656e646572206973206e6f7420612076616c696461746f72000000604482015290519081900360640190fd5b61126c8382846141fb565b600f6020526000908152604090205481565b6000806122d46004610f39565b90506000805b82811015611108576000818152600b60209081526040808320546001600160a01b03168352600e90915290206002015461231b90839063ffffffff61319616565b91506001016122da565b6001546001600160a01b031690565b6000806123416004610f39565b90506000805b828110156123a7576000818152600b60205260409020546001600160a01b03166123705761239f565b6000818152600b60209081526040808320546001600160a01b03168352600e90915290206002015491506123a7565b600101612347565b6001015b82811015611108576000818152600b60209081526040808320546001600160a01b03168352600e90915290206002015482111561240d576000818152600b60209081526040808320546001600160a01b03168352600e90915290206002015491505b6001016123ab565b6001546001600160a01b0316331490565b60006124326004610f39565b90506000805b828110156124ba5760016000828152600b60205260409020546124659086906001600160a01b0316611b14565b600381111561247057fe5b14156124b2576000818152600b60209081526040808320546001600160a01b03168352600e9091529020600201546124af90839063ffffffff61319616565b91505b600101612438565b5060006124c56124d4565b8210159050610fe084826143b6565b6000612511600161250560036124f960026124ed6122c7565b9063ffffffff61452916565b9063ffffffff61455016565b9063ffffffff61319616565b905090565b61251e612415565b61252757600080fd5b6001600160a01b03166000908152600860205260409020805460ff19166001179055565b60006110bb60028363ffffffff613d9416565b816001600160a01b0381166125a6576040805162461bcd60e51b815260206004820152600960248201526830206164647265737360b81b604482015290519081900360640190fd5b6001600160a01b0383166000908152600e6020526040812090600482015460ff1660028111156125d257fe5b146125dc57600080fd5b336125ea8282866001612f80565b601054612607906001600160a01b0316828663ffffffff61310116565b846001600160a01b0316816001600160a01b03167f585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8866040518082815260200191505060405180910390a35050505050565b336000908152600e60205260409020805460ff166126ac576040805162461bcd60e51b815260206004820152601c6024820152600080516020614c32833981519152604482015290519081900360640190fd5b80600601548311156126f8576040805162461bcd60e51b815260206004820152601060248201526f496e76616c6964206e6577207261746560801b604482015290519081900360640190fd5b61126c818484613f78565b61271081565b600154600160a01b900460ff161561272057600080fd5b60135460ff161561276f57612734336119bd565b61276f5760405162461bcd60e51b815260040180806020018281038252603a815260200180614bf8603a913960400191505060405180910390fd5b336000908152600e60205260409020805460ff16156127d5576040805162461bcd60e51b815260206004820152601860248201527f43616e64696461746520697320696e697469616c697a65640000000000000000604482015290519081900360640190fd5b6127108311156127e457600080fd5b805460ff191660019081178255810184905560068101839055600781018290556040805185815260208101859052808201849052905133917f453d56a841836718d9e848e968068cbc2af21ca29d1527fbebd231dc46ceffaa919081900360600190a250505050565b6010546001600160a01b031681565b6001600160a01b0381166000908152600e602052604090206002600482015460ff16600281111561288957fe5b1461289357600080fd5b80600501544310156128a457600080fd5b60048101805460ff191690556000600582018190556040516001600160a01b038416917fbe85a9a7aa606febeaa35606e49cd7324c63cf970f4f5fd0c7e983f42b20b21991a25050565b60135460ff1681565b806001600160a01b03811661293f576040805162461bcd60e51b815260206004820152600960248201526830206164647265737360b81b604482015290519081900360640190fd5b6001600160a01b0382166000818152600e60208181526040808420338086526003820184529185209585529290915260049091015490929143918190819060ff16600281111561298b57fe5b149050836003015491505b8360040154821015612a1657600082815260028086016020526040822091906129be90610f39565b905082806129e15750600182015485906129de908363ffffffff61319616565b11155b15612a045750506000828152600285016020526040812081815560010155612a0b565b5050612a16565b600190910190612996565b6003840182905560005b8460040154831015612a5e57600083815260028601602052604090208054612a4f90839063ffffffff61319616565b6001909401939150612a209050565b6001850154600090821015612aa9576001860154612a82908363ffffffff61370e16565b60018701839055601054909150612aa9906001600160a01b0316888363ffffffff61310116565b886001600160a01b0316876001600160a01b03167f08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220836040518082815260200191505060405180910390a3505050505050505050565b611b03336138f2565b600a54600081815260096020526040902090612b2b90600163ffffffff61319616565b600a556000808052600560208190527f05b8ccbb9d4d8fb16ea74ce3c29a41f1b461fbdaff4714a0d9a8eb05499746bc5483546001600160a01b031916339081178555600180860183905590939192612b85929091611399565b60028401556003830180546001600160a01b0319166001600160a01b038781169190911760ff60a01b1916600160a01b871515021760ff60a81b1916600160a81b17909155600454612bda9116833084613029565b600a547fe6970151d691583ac0aecc2e24c67871318a5c7f7574c6df7929b6dd5d54db6890612c1090600163ffffffff61370e16565b6002850154604080519283526001600160a01b0380871660208501528382018690526060840192909252908816608083015286151560a0830152519081900360c00190a15050505050565b600b602052600090815260409020546001600160a01b031681565b60096020526000908152604090208054600182015460028301546003909301546001600160a01b0392831693919281169060ff600160a01b8204811691600160a81b90041686565b600080612ccb6003610f39565b90506011544310158015612ce6575080612ce36110c1565b10155b91505090565b6001600160a01b038083166000908152600e60209081526040808320938516835260039384019091528120918201546004830154919283926060928392918591612d3c919063ffffffff61370e16565b905080604051908082528060200260200182016040528015612d68578160200160208202803883390190505b50935080604051908082528060200260200182016040528015612d95578160200160208202803883390190505b50925060005b81811015612e13576003830154810160009081526002840160205260409020548551869083908110612dc957fe5b60200260200101818152505082600201600084600301548301815260200190815260200160002060010154848281518110612e0057fe5b6020908102919091010152600101612d9b565b50508054600190910154909790965091945092509050565b612e33612415565b612e3c57600080fd5b610eea81614572565b336000908152600e60205260409020805460ff16612e98576040805162461bcd60e51b815260206004820152601c6024820152600080516020614c32833981519152604482015290519081900360640190fd5b82816006015410612ee3576040805162461bcd60e51b815260206004820152601060248201526f496e76616c6964206e6577207261746560801b604482015290519081900360640190fd5b600881018390556009810182905543600a8201556040805184815260208101849052815133927fd1388fca1fdda1adbe79c9535b48b22e71aa7815469abb61cdbab2a7b4ccd28a928290030190a2505050565b600060016001600160a01b0383166000908152600e602052604090206004015460ff166002811115612f6457fe5b1492915050565b600c6020526000908152604090205460ff1681565b6001600160a01b0383166000908152600385016020526040812090826001811115612fa757fe5b1415612fe0576002850154612fc2908463ffffffff61319616565b60028601558054612fd9908463ffffffff61319616565b8155613022565b6001826001811115612fee57fe5b1415613020576002850154613009908463ffffffff61370e16565b60028601558054612fd9908463ffffffff61370e16565bfe5b5050505050565b604080516323b872dd60e01b81526001600160a01b0385811660048301528481166024830152604482018490529151918616916323b872dd916064808201926020929091908290030181600087803b15801561308457600080fd5b505af1158015613098573d6000803e3d6000fd5b505050506040513d60208110156130ae57600080fd5b5051610fe057600080fd5b6130ca60038263ffffffff6145e116565b6040516001600160a01b038216907fee1504a83b6d4a361f4c1dc78ab59bfa30d6a3b6612c403e86bb01ef2984295f90600090a250565b826001600160a01b031663a9059cbb83836040518363ffffffff1660e01b815260040180836001600160a01b03166001600160a01b0316815260200182815260200192505050602060405180830381600087803b15801561316157600080fd5b505af1158015613175573d6000803e3d6000fd5b505050506040513d602081101561318b57600080fd5b505161126c57600080fd5b6000828201838110156131a857600080fd5b9392505050565b600082815260096020526040902060016003820154600160a81b900460ff1660028111156131d957fe5b14613225576040805162461bcd60e51b8152602060048201526017602482015276496e76616c69642070726f706f73616c2073746174757360481b604482015290519081900360640190fd5b806002015443101561327a576040805162461bcd60e51b8152602060048201526019602482015278159bdd1948191958591b1a5b99481b9bdd081c995858da1959603a1b604482015290519081900360640190fd5b60038101805460ff60a81b1916600160a91b17905581156132f457805460018201546004546132bd926001600160a01b039182169291169063ffffffff61310116565b60038101546001600160a01b0381166000908152600860205260409020805460ff1916600160a01b90920460ff1615159190911790555b60038101546040805185815284151560208201526001600160a01b03831681830152600160a01b90920460ff1615156060830152517f2c26ff0b5547eb09df5dde3569782330829ac9ffa9811847beab5d466066801c916080908290030190a1505050565b613361614b3b565b613369614b99565b6133728361462d565b9050606061338782600263ffffffff61464416565b90508060028151811061339657fe5b60200260200101516040519080825280602002602001820160405280156133d157816020015b60608152602001906001900390816133bc5790505b5083602001819052506000816002815181106133e957fe5b6020026020010181815250506000805b613402846146d4565b156134a757613410846146e0565b9092509050816001141561342e576134278461470d565b85526134a2565b8160021415613492576134408461470d565b85602001518460028151811061345257fe5b60200260200101518151811061346457fe5b60200260200101819052508260028151811061347c57fe5b60209081029190910101805160010190526134a2565b6134a2848263ffffffff61479a16565b6133f9565b50505050919050565b6000806134bb6124d4565b905060006134c8856147f7565b9050606084516040519080825280602002602001820160405280156134f7578160200160208202803883390190505b509050600080805b875181101561369a5761352e88828151811061351757fe5b60200260200101518661484890919063ffffffff16565b84828151811061353a57fe5b60200260200101906001600160a01b031690816001600160a01b031681525050600d600085838151811061356a57fe5b6020908102919091018101516001600160a01b031682528101919091526040016000205460ff161561359f576001915061369a565b6001600e60008684815181106135b157fe5b6020908102919091018101516001600160a01b031682528101919091526040016000206004015460ff1660028111156135e657fe5b146135f057613692565b61363e600e600086848151811061360357fe5b60200260200101516001600160a01b03166001600160a01b03168152602001908152602001600020600201548461319690919063ffffffff16565b92506001600d600086848151811061365257fe5b60200260200101516001600160a01b03166001600160a01b0316815260200190815260200160002060006101000a81548160ff0219169083151502179055505b6001016134ff565b5060005b87518110156136f3576000600d60008684815181106136b957fe5b6020908102919091018101516001600160a01b03168252810191909152604001600020805460ff191691151591909117905560010161369e565b50801580156137025750848210155b98975050505050505050565b60008282111561371d57600080fd5b50900390565b60008381526006602052604090206001600582015460ff16600281111561374657fe5b14613792576040805162461bcd60e51b8152602060048201526017602482015276496e76616c69642070726f706f73616c2073746174757360481b604482015290519081900360640190fd5b806002015443106137e2576040805162461bcd60e51b8152602060048201526015602482015274159bdd1948191958591b1a5b99481c995858da1959605a1b604482015290519081900360640190fd5b6001600160a01b038316600090815260068201602052604081205460ff16600381111561380b57fe5b1461384f576040805162461bcd60e51b815260206004820152600f60248201526e159bdd195c881a185cc81d9bdd1959608a1b604482015290519081900360640190fd5b6001600160a01b03831660009081526006820160205260409020805483919060ff1916600183600381111561388057fe5b02179055507f06c7ef6e19454637e93ee60cc680c61fb2ebabb57e58cf36d94141a5036b3d6584848460405180848152602001836001600160a01b03166001600160a01b031681526020018260038111156138d757fe5b60ff168152602001935050505060405180910390a150505050565b61390360038263ffffffff61491916565b6040516001600160a01b038216907f270d9b30cf5b0793bbfd54c9d5b94aeb49462b8148399000265144a8722da6b690600090a250565b613942614b3b565b61394a614b99565b6139538361462d565b9050606061396882600263ffffffff61464416565b90508060028151811061397757fe5b60200260200101516040519080825280602002602001820160405280156139b257816020015b606081526020019060019003908161399d5790505b5083602001819052506000816002815181106139ca57fe5b6020026020010181815250506000805b6139e3846146d4565b156134a7576139f1846146e0565b90925090508160011415613a0f57613a088461470d565b8552613a83565b8160021415613a7357613a218461470d565b856020015184600281518110613a3357fe5b602002602001015181518110613a4557fe5b602002602001018190525082600281518110613a5d57fe5b6020908102919091010180516001019052613a83565b613a83848263ffffffff61479a16565b6139da565b613a90614b55565b613a98614b99565b613aa18361462d565b90506060613ab682600563ffffffff61464416565b905080600481518110613ac557fe5b6020026020010151604051908082528060200260200182016040528015613b0657816020015b613af3614b82565b815260200190600190039081613aeb5790505b508360600181905250600081600481518110613b1e57fe5b60200260200101818152505080600581518110613b3757fe5b6020026020010151604051908082528060200260200182016040528015613b7857816020015b613b65614b82565b815260200190600190039081613b5d5790505b508360800181905250600081600581518110613b9057fe5b6020026020010181815250506000805b613ba9846146d4565b156134a757613bb7846146e0565b90925090508160011415613bdf57613bce84614961565b67ffffffffffffffff168552613cff565b8160021415613c0557613bf184614961565b67ffffffffffffffff166020860152613cff565b8160031415613c3257613c1f613c1a8561470d565b6149bc565b6001600160a01b03166040860152613cff565b8160041415613c9e57613c4c613c478561470d565b6149c7565b856060015184600481518110613c5e57fe5b602002602001015181518110613c7057fe5b602002602001018190525082600481518110613c8857fe5b6020908102919091010180516001019052613cff565b8160051415613cef57613cb3613c478561470d565b856080015184600581518110613cc557fe5b602002602001015181518110613cd757fe5b602002602001018190525082600581518110613c8857fe5b613cff848263ffffffff61479a16565b613ba0565b6001600160a01b0381166000908152600e602052604090206001600482015460ff166002811115613d3157fe5b14613d3c5750610eea565b60018101546001600160a01b03831660009081526003830160205260408120549190911090613d6b6005610f39565b600284015490915081118280613d7e5750805b1561302257613022613d8f86614a68565b613e11565b60006001600160a01b038216613da957600080fd5b506001600160a01b03166000908152602091909152604090205460ff1690565b613dda60028263ffffffff61491916565b6040516001600160a01b038216907f0a8eb35e5ca14b3d6f28e4abf2f128dbab231a58b56e89beb5d636115001e16590600090a250565b6000818152600b60205260409020546001600160a01b031680613e345750610eea565b6000828152600b6020908152604080832080546001600160a01b03191690556001600160a01b0384168352600e9091528120600401805460ff19166002908117909155613e8090610f39565b9050613e92438263ffffffff61319616565b6001600160a01b0383166000818152600e60205260408082206005019390935591516001927f63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c91a3505050565b6000818152600b60205260409020546001600160a01b031615613f0157600080fd5b6000818152600b6020908152604080832080546001600160a01b0319166001600160a01b038716908117909155808452600e90925280832060048101805460ff19166001179055600501839055517f63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c908390a35050565b612710821115613fc2576040805162461bcd60e51b815260206004820152601060248201526f496e76616c6964206e6577207261746560801b604482015290519081900360640190fd5b438111614016576040805162461bcd60e51b815260206004820152601a60248201527f4f75746461746564206e6577206c6f636b20656e642074696d65000000000000604482015290519081900360640190fd5b826006015482141561406957826007015481116140645760405162461bcd60e51b8152600401808060200182810382526023815260200180614bb46023913960400191505060405180910390fd5b6140c1565b826007015443116140c1576040805162461bcd60e51b815260206004820152601960248201527f436f6d6d697373696f6e2072617465206973206c6f636b656400000000000000604482015290519081900360640190fd5b600683018290556007830181905560006008840181905560098401819055600a8401556040805183815260208101839052815133927f37954fc2aa8b4424ad16c75da2ea4d51ba08ef9e07907e37ccae54a0b4ce1e9e928290030190a2505050565b61413460008263ffffffff61491916565b6040516001600160a01b038216907fcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e90600090a250565b61417c60028263ffffffff6145e116565b6040516001600160a01b038216907f22380c05984257a1cb900161c713dd71d39e74820f1aea43bd3f1bdd2096129990600090a250565b6141c460008263ffffffff6145e116565b6040516001600160a01b038216907f6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f890600090a250565b600083815260096020526040902060016003820154600160a81b900460ff16600281111561422557fe5b14614271576040805162461bcd60e51b8152602060048201526017602482015276496e76616c69642070726f706f73616c2073746174757360481b604482015290519081900360640190fd5b806002015443106142c1576040805162461bcd60e51b8152602060048201526015602482015274159bdd1948191958591b1a5b99481c995858da1959605a1b604482015290519081900360640190fd5b6001600160a01b038316600090815260048201602052604081205460ff1660038111156142ea57fe5b1461432e576040805162461bcd60e51b815260206004820152600f60248201526e159bdd195c881a185cc81d9bdd1959608a1b604482015290519081900360640190fd5b6001600160a01b03831660009081526004820160205260409020805483919060ff1916600183600381111561435f57fe5b02179055507f7686976924e1fdb79b36f7445ada20b6e9d3377d85b34d5162116e675c39d34c84848460405180848152602001836001600160a01b03166001600160a01b031681526020018260038111156138d757fe5b60008281526006602052604090206001600582015460ff1660028111156143d957fe5b14614425576040805162461bcd60e51b8152602060048201526017602482015276496e76616c69642070726f706f73616c2073746174757360481b604482015290519081900360640190fd5b806002015443101561447a576040805162461bcd60e51b8152602060048201526019602482015278159bdd1948191958591b1a5b99481b9bdd081c995858da1959603a1b604482015290519081900360640190fd5b60058101805460ff1916600217905581156144d157805460018201546004546144b7926001600160a01b039182169291169063ffffffff61310116565b600481015460038201546000908152600560205260409020555b60038101546004820154604080518681528515156020820152808201939093526060830191909152517f106f43a560e53395081c0423504b476d1a2cfed9d56ff972bf77ae43ff7d4ba49181900360800190a1505050565b600082614538575060006110bb565b8282028284828161454557fe5b04146131a857600080fd5b600080821161455e57600080fd5b600082848161456957fe5b04949350505050565b6001600160a01b03811661458557600080fd5b6001546040516001600160a01b038084169216907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a3600180546001600160a01b0319166001600160a01b0392909216919091179055565b6001600160a01b0381166145f457600080fd5b6145fe8282613d94565b1561460857600080fd5b6001600160a01b0316600090815260209190915260409020805460ff19166001179055565b614635614b99565b60208101919091526000815290565b815160408051600184018082526020808202830101909252606092918015614676578160200160208202803883390190505b5091506000805b614686866146d4565b156146cb57614694866146e0565b809250819350505060018483815181106146aa57fe5b6020026020010181815101915081815250506146c6868261479a565b61467d565b50509092525090565b60208101515190511090565b60008060006146ee84614961565b905060088104925080600716600581111561470557fe5b915050915091565b6060600061471a83614961565b835160208501515191925082019081111561473457600080fd5b816040519080825280601f01601f19166020018201604052801561475f576020820181803883390190505b50602080860151865192955091818601919083010160005b8581101561478f578181015183820152602001614777565b505050935250919050565b60008160058111156147a857fe5b14156147bd576147b782614961565b50611fa9565b60028160058111156147cb57fe5b14156103fc5760006147dc83614961565b8351810180855260208501515191925011156147b757600080fd5b604080517f19457468657265756d205369676e6564204d6573736167653a0a333200000000602080830191909152603c8083019490945282518083039094018452605c909101909152815191012090565b600080600080845160411461486357600093505050506110bb565b50505060208201516040830151606084015160001a601b81101561488557601b015b8060ff16601b1415801561489d57508060ff16601c14155b156148ae57600093505050506110bb565b6040805160008152602080820180845289905260ff8416828401526060820186905260808201859052915160019260a0808401939192601f1981019281900390910190855afa158015614905573d6000803e3d6000fd5b5050506020604051035193505050506110bb565b6001600160a01b03811661492c57600080fd5b6149368282613d94565b61493f57600080fd5b6001600160a01b0316600090815260209190915260409020805460ff19169055565b602080820151825181019091015160009182805b600a8110156149b65783811a91508060070282607f16901b8517945081608016600014156149ae5785510160010185525061148b915050565b600101614975565b50600080fd5b60006110bb82614af9565b6149cf614b82565b6149d7614b99565b6149e08361462d565b90506000805b6149ef836146d4565b15614a60576149fd836146e0565b90925090508160011415614a2757614a17613c1a8461470d565b6001600160a01b03168452614a5b565b8160021415614a4b57614a41614a3c8461470d565b614b18565b6020850152614a5b565b614a5b838263ffffffff61479a16565b6149e6565b505050919050565b600080614a756004610f39565b905060005b81811015614ab5576000818152600b60205260409020546001600160a01b0385811691161415614aad57915061148b9050565b600101614a7a565b506040805162461bcd60e51b815260206004820152601360248201527227379039bab1b41030903b30b634b230ba37b960691b604482015290519081900360640190fd5b60008151601414614b0957600080fd5b5060200151600160601b900490565b6000602082511115614b2957600080fd5b50602081810151915160089103021c90565b604051806040016040528060608152602001606081525090565b6040805160a081018252600080825260208201819052918101919091526060808201819052608082015290565b604080518082019091526000808252602082015290565b60405180604001604052806000815260200160608152509056fe4e6577206c6f636b20656e642074696d65206973206e6f7420696e6372656173696e675374616b65206973206c657373207468616e20616c6c2076616c696461746f727357686974656c6973746564526f6c653a2063616c6c657220646f6573206e6f742068617665207468652057686974656c697374656420726f6c6543616e646964617465206973206e6f7420696e697469616c697a656400000000a265627a7a723058205df5c2259a9e3568ec5ca45be692bd6e6852a250f5c8fc6cc648a0aa9758b14464736f6c634300050a0032"

// DeployDPoS deploys a new Ethereum contract, binding an instance of DPoS to it.
func DeployDPoS(auth *bind.TransactOpts, backend bind.ContractBackend, _celerTokenAddress common.Address, _governProposalDeposit *big.Int, _governVoteTimeout *big.Int, _blameTimeout *big.Int, _minValidatorNum *big.Int, _maxValidatorNum *big.Int, _minStakeInPool *big.Int, _increaseRateWaitTime *big.Int, _dposGoLiveTimeout *big.Int) (common.Address, *types.Transaction, *DPoS, error) {
	parsed, err := abi.JSON(strings.NewReader(DPoSABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(DPoSBin), backend, _celerTokenAddress, _governProposalDeposit, _governVoteTimeout, _blameTimeout, _minValidatorNum, _maxValidatorNum, _minStakeInPool, _increaseRateWaitTime, _dposGoLiveTimeout)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DPoS{DPoSCaller: DPoSCaller{contract: contract}, DPoSTransactor: DPoSTransactor{contract: contract}, DPoSFilterer: DPoSFilterer{contract: contract}}, nil
}

// DPoS is an auto generated Go binding around an Ethereum contract.
type DPoS struct {
	DPoSCaller     // Read-only binding to the contract
	DPoSTransactor // Write-only binding to the contract
	DPoSFilterer   // Log filterer for contract events
}

// DPoSCaller is an auto generated read-only Go binding around an Ethereum contract.
type DPoSCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DPoSTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DPoSTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DPoSFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DPoSFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DPoSSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DPoSSession struct {
	Contract     *DPoS             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DPoSCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DPoSCallerSession struct {
	Contract *DPoSCaller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// DPoSTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DPoSTransactorSession struct {
	Contract     *DPoSTransactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DPoSRaw is an auto generated low-level Go binding around an Ethereum contract.
type DPoSRaw struct {
	Contract *DPoS // Generic contract binding to access the raw methods on
}

// DPoSCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DPoSCallerRaw struct {
	Contract *DPoSCaller // Generic read-only contract binding to access the raw methods on
}

// DPoSTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DPoSTransactorRaw struct {
	Contract *DPoSTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDPoS creates a new instance of DPoS, bound to a specific deployed contract.
func NewDPoS(address common.Address, backend bind.ContractBackend) (*DPoS, error) {
	contract, err := bindDPoS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DPoS{DPoSCaller: DPoSCaller{contract: contract}, DPoSTransactor: DPoSTransactor{contract: contract}, DPoSFilterer: DPoSFilterer{contract: contract}}, nil
}

// NewDPoSCaller creates a new read-only instance of DPoS, bound to a specific deployed contract.
func NewDPoSCaller(address common.Address, caller bind.ContractCaller) (*DPoSCaller, error) {
	contract, err := bindDPoS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DPoSCaller{contract: contract}, nil
}

// NewDPoSTransactor creates a new write-only instance of DPoS, bound to a specific deployed contract.
func NewDPoSTransactor(address common.Address, transactor bind.ContractTransactor) (*DPoSTransactor, error) {
	contract, err := bindDPoS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DPoSTransactor{contract: contract}, nil
}

// NewDPoSFilterer creates a new log filterer instance of DPoS, bound to a specific deployed contract.
func NewDPoSFilterer(address common.Address, filterer bind.ContractFilterer) (*DPoSFilterer, error) {
	contract, err := bindDPoS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DPoSFilterer{contract: contract}, nil
}

// bindDPoS binds a generic wrapper to an already deployed contract.
func bindDPoS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DPoSABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DPoS *DPoSRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DPoS.Contract.DPoSCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DPoS *DPoSRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.Contract.DPoSTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DPoS *DPoSRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DPoS.Contract.DPoSTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DPoS *DPoSCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _DPoS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DPoS *DPoSTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DPoS *DPoSTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DPoS.Contract.contract.Transact(opts, method, params...)
}

// COMMISSIONRATEBASE is a free data retrieval call binding the contract method 0xbee8380e.
//
// Solidity: function COMMISSION_RATE_BASE() constant returns(uint256)
func (_DPoS *DPoSCaller) COMMISSIONRATEBASE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "COMMISSION_RATE_BASE")
	return *ret0, err
}

// COMMISSIONRATEBASE is a free data retrieval call binding the contract method 0xbee8380e.
//
// Solidity: function COMMISSION_RATE_BASE() constant returns(uint256)
func (_DPoS *DPoSSession) COMMISSIONRATEBASE() (*big.Int, error) {
	return _DPoS.Contract.COMMISSIONRATEBASE(&_DPoS.CallOpts)
}

// COMMISSIONRATEBASE is a free data retrieval call binding the contract method 0xbee8380e.
//
// Solidity: function COMMISSION_RATE_BASE() constant returns(uint256)
func (_DPoS *DPoSCallerSession) COMMISSIONRATEBASE() (*big.Int, error) {
	return _DPoS.Contract.COMMISSIONRATEBASE(&_DPoS.CallOpts)
}

// UIntStorage is a free data retrieval call binding the contract method 0x64ed600a.
//
// Solidity: function UIntStorage(uint256 ) constant returns(uint256)
func (_DPoS *DPoSCaller) UIntStorage(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "UIntStorage", arg0)
	return *ret0, err
}

// UIntStorage is a free data retrieval call binding the contract method 0x64ed600a.
//
// Solidity: function UIntStorage(uint256 ) constant returns(uint256)
func (_DPoS *DPoSSession) UIntStorage(arg0 *big.Int) (*big.Int, error) {
	return _DPoS.Contract.UIntStorage(&_DPoS.CallOpts, arg0)
}

// UIntStorage is a free data retrieval call binding the contract method 0x64ed600a.
//
// Solidity: function UIntStorage(uint256 ) constant returns(uint256)
func (_DPoS *DPoSCallerSession) UIntStorage(arg0 *big.Int) (*big.Int, error) {
	return _DPoS.Contract.UIntStorage(&_DPoS.CallOpts, arg0)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_DPoS *DPoSCaller) CelerToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "celerToken")
	return *ret0, err
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_DPoS *DPoSSession) CelerToken() (common.Address, error) {
	return _DPoS.Contract.CelerToken(&_DPoS.CallOpts)
}

// CelerToken is a free data retrieval call binding the contract method 0xc6c21e9d.
//
// Solidity: function celerToken() constant returns(address)
func (_DPoS *DPoSCallerSession) CelerToken() (common.Address, error) {
	return _DPoS.Contract.CelerToken(&_DPoS.CallOpts)
}

// CheckedValidators is a free data retrieval call binding the contract method 0x3702db39.
//
// Solidity: function checkedValidators(address ) constant returns(bool)
func (_DPoS *DPoSCaller) CheckedValidators(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "checkedValidators", arg0)
	return *ret0, err
}

// CheckedValidators is a free data retrieval call binding the contract method 0x3702db39.
//
// Solidity: function checkedValidators(address ) constant returns(bool)
func (_DPoS *DPoSSession) CheckedValidators(arg0 common.Address) (bool, error) {
	return _DPoS.Contract.CheckedValidators(&_DPoS.CallOpts, arg0)
}

// CheckedValidators is a free data retrieval call binding the contract method 0x3702db39.
//
// Solidity: function checkedValidators(address ) constant returns(bool)
func (_DPoS *DPoSCallerSession) CheckedValidators(arg0 common.Address) (bool, error) {
	return _DPoS.Contract.CheckedValidators(&_DPoS.CallOpts, arg0)
}

// DposGoLiveTime is a free data retrieval call binding the contract method 0x39c9563e.
//
// Solidity: function dposGoLiveTime() constant returns(uint256)
func (_DPoS *DPoSCaller) DposGoLiveTime(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "dposGoLiveTime")
	return *ret0, err
}

// DposGoLiveTime is a free data retrieval call binding the contract method 0x39c9563e.
//
// Solidity: function dposGoLiveTime() constant returns(uint256)
func (_DPoS *DPoSSession) DposGoLiveTime() (*big.Int, error) {
	return _DPoS.Contract.DposGoLiveTime(&_DPoS.CallOpts)
}

// DposGoLiveTime is a free data retrieval call binding the contract method 0x39c9563e.
//
// Solidity: function dposGoLiveTime() constant returns(uint256)
func (_DPoS *DPoSCallerSession) DposGoLiveTime() (*big.Int, error) {
	return _DPoS.Contract.DposGoLiveTime(&_DPoS.CallOpts)
}

// EnableWhitelist is a free data retrieval call binding the contract method 0xcdfb2b4e.
//
// Solidity: function enableWhitelist() constant returns(bool)
func (_DPoS *DPoSCaller) EnableWhitelist(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "enableWhitelist")
	return *ret0, err
}

// EnableWhitelist is a free data retrieval call binding the contract method 0xcdfb2b4e.
//
// Solidity: function enableWhitelist() constant returns(bool)
func (_DPoS *DPoSSession) EnableWhitelist() (bool, error) {
	return _DPoS.Contract.EnableWhitelist(&_DPoS.CallOpts)
}

// EnableWhitelist is a free data retrieval call binding the contract method 0xcdfb2b4e.
//
// Solidity: function enableWhitelist() constant returns(bool)
func (_DPoS *DPoSCallerSession) EnableWhitelist() (bool, error) {
	return _DPoS.Contract.EnableWhitelist(&_DPoS.CallOpts)
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, uint256 stakingPool, uint256 status, uint256 unbondTime, uint256 commissionRate, uint256 rateLockEndTime)
func (_DPoS *DPoSCaller) GetCandidateInfo(opts *bind.CallOpts, _candidateAddr common.Address) (struct {
	Initialized     bool
	MinSelfStake    *big.Int
	StakingPool     *big.Int
	Status          *big.Int
	UnbondTime      *big.Int
	CommissionRate  *big.Int
	RateLockEndTime *big.Int
}, error) {
	ret := new(struct {
		Initialized     bool
		MinSelfStake    *big.Int
		StakingPool     *big.Int
		Status          *big.Int
		UnbondTime      *big.Int
		CommissionRate  *big.Int
		RateLockEndTime *big.Int
	})
	out := ret
	err := _DPoS.contract.Call(opts, out, "getCandidateInfo", _candidateAddr)
	return *ret, err
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, uint256 stakingPool, uint256 status, uint256 unbondTime, uint256 commissionRate, uint256 rateLockEndTime)
func (_DPoS *DPoSSession) GetCandidateInfo(_candidateAddr common.Address) (struct {
	Initialized     bool
	MinSelfStake    *big.Int
	StakingPool     *big.Int
	Status          *big.Int
	UnbondTime      *big.Int
	CommissionRate  *big.Int
	RateLockEndTime *big.Int
}, error) {
	return _DPoS.Contract.GetCandidateInfo(&_DPoS.CallOpts, _candidateAddr)
}

// GetCandidateInfo is a free data retrieval call binding the contract method 0x28bde1e1.
//
// Solidity: function getCandidateInfo(address _candidateAddr) constant returns(bool initialized, uint256 minSelfStake, uint256 stakingPool, uint256 status, uint256 unbondTime, uint256 commissionRate, uint256 rateLockEndTime)
func (_DPoS *DPoSCallerSession) GetCandidateInfo(_candidateAddr common.Address) (struct {
	Initialized     bool
	MinSelfStake    *big.Int
	StakingPool     *big.Int
	Status          *big.Int
	UnbondTime      *big.Int
	CommissionRate  *big.Int
	RateLockEndTime *big.Int
}, error) {
	return _DPoS.Contract.GetCandidateInfo(&_DPoS.CallOpts, _candidateAddr)
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 delegatedStake, uint256 undelegatingStake, uint256[] intentAmounts, uint256[] intentProposedTimes)
func (_DPoS *DPoSCaller) GetDelegatorInfo(opts *bind.CallOpts, _candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	DelegatedStake      *big.Int
	UndelegatingStake   *big.Int
	IntentAmounts       []*big.Int
	IntentProposedTimes []*big.Int
}, error) {
	ret := new(struct {
		DelegatedStake      *big.Int
		UndelegatingStake   *big.Int
		IntentAmounts       []*big.Int
		IntentProposedTimes []*big.Int
	})
	out := ret
	err := _DPoS.contract.Call(opts, out, "getDelegatorInfo", _candidateAddr, _delegatorAddr)
	return *ret, err
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 delegatedStake, uint256 undelegatingStake, uint256[] intentAmounts, uint256[] intentProposedTimes)
func (_DPoS *DPoSSession) GetDelegatorInfo(_candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	DelegatedStake      *big.Int
	UndelegatingStake   *big.Int
	IntentAmounts       []*big.Int
	IntentProposedTimes []*big.Int
}, error) {
	return _DPoS.Contract.GetDelegatorInfo(&_DPoS.CallOpts, _candidateAddr, _delegatorAddr)
}

// GetDelegatorInfo is a free data retrieval call binding the contract method 0xeecefef8.
//
// Solidity: function getDelegatorInfo(address _candidateAddr, address _delegatorAddr) constant returns(uint256 delegatedStake, uint256 undelegatingStake, uint256[] intentAmounts, uint256[] intentProposedTimes)
func (_DPoS *DPoSCallerSession) GetDelegatorInfo(_candidateAddr common.Address, _delegatorAddr common.Address) (struct {
	DelegatedStake      *big.Int
	UndelegatingStake   *big.Int
	IntentAmounts       []*big.Int
	IntentProposedTimes []*big.Int
}, error) {
	return _DPoS.Contract.GetDelegatorInfo(&_DPoS.CallOpts, _candidateAddr, _delegatorAddr)
}

// GetMinQuorumStakingPool is a free data retrieval call binding the contract method 0xa3e814b9.
//
// Solidity: function getMinQuorumStakingPool() constant returns(uint256)
func (_DPoS *DPoSCaller) GetMinQuorumStakingPool(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getMinQuorumStakingPool")
	return *ret0, err
}

// GetMinQuorumStakingPool is a free data retrieval call binding the contract method 0xa3e814b9.
//
// Solidity: function getMinQuorumStakingPool() constant returns(uint256)
func (_DPoS *DPoSSession) GetMinQuorumStakingPool() (*big.Int, error) {
	return _DPoS.Contract.GetMinQuorumStakingPool(&_DPoS.CallOpts)
}

// GetMinQuorumStakingPool is a free data retrieval call binding the contract method 0xa3e814b9.
//
// Solidity: function getMinQuorumStakingPool() constant returns(uint256)
func (_DPoS *DPoSCallerSession) GetMinQuorumStakingPool() (*big.Int, error) {
	return _DPoS.Contract.GetMinQuorumStakingPool(&_DPoS.CallOpts)
}

// GetMinStakingPool is a free data retrieval call binding the contract method 0x8e9472a3.
//
// Solidity: function getMinStakingPool() constant returns(uint256)
func (_DPoS *DPoSCaller) GetMinStakingPool(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getMinStakingPool")
	return *ret0, err
}

// GetMinStakingPool is a free data retrieval call binding the contract method 0x8e9472a3.
//
// Solidity: function getMinStakingPool() constant returns(uint256)
func (_DPoS *DPoSSession) GetMinStakingPool() (*big.Int, error) {
	return _DPoS.Contract.GetMinStakingPool(&_DPoS.CallOpts)
}

// GetMinStakingPool is a free data retrieval call binding the contract method 0x8e9472a3.
//
// Solidity: function getMinStakingPool() constant returns(uint256)
func (_DPoS *DPoSCallerSession) GetMinStakingPool() (*big.Int, error) {
	return _DPoS.Contract.GetMinStakingPool(&_DPoS.CallOpts)
}

// GetParamProposalVote is a free data retrieval call binding the contract method 0x581c53c5.
//
// Solidity: function getParamProposalVote(uint256 _proposalId, address _voter) constant returns(uint8)
func (_DPoS *DPoSCaller) GetParamProposalVote(opts *bind.CallOpts, _proposalId *big.Int, _voter common.Address) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getParamProposalVote", _proposalId, _voter)
	return *ret0, err
}

// GetParamProposalVote is a free data retrieval call binding the contract method 0x581c53c5.
//
// Solidity: function getParamProposalVote(uint256 _proposalId, address _voter) constant returns(uint8)
func (_DPoS *DPoSSession) GetParamProposalVote(_proposalId *big.Int, _voter common.Address) (uint8, error) {
	return _DPoS.Contract.GetParamProposalVote(&_DPoS.CallOpts, _proposalId, _voter)
}

// GetParamProposalVote is a free data retrieval call binding the contract method 0x581c53c5.
//
// Solidity: function getParamProposalVote(uint256 _proposalId, address _voter) constant returns(uint8)
func (_DPoS *DPoSCallerSession) GetParamProposalVote(_proposalId *big.Int, _voter common.Address) (uint8, error) {
	return _DPoS.Contract.GetParamProposalVote(&_DPoS.CallOpts, _proposalId, _voter)
}

// GetSidechainProposalVote is a free data retrieval call binding the contract method 0x8515b0e2.
//
// Solidity: function getSidechainProposalVote(uint256 _proposalId, address _voter) constant returns(uint8)
func (_DPoS *DPoSCaller) GetSidechainProposalVote(opts *bind.CallOpts, _proposalId *big.Int, _voter common.Address) (uint8, error) {
	var (
		ret0 = new(uint8)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getSidechainProposalVote", _proposalId, _voter)
	return *ret0, err
}

// GetSidechainProposalVote is a free data retrieval call binding the contract method 0x8515b0e2.
//
// Solidity: function getSidechainProposalVote(uint256 _proposalId, address _voter) constant returns(uint8)
func (_DPoS *DPoSSession) GetSidechainProposalVote(_proposalId *big.Int, _voter common.Address) (uint8, error) {
	return _DPoS.Contract.GetSidechainProposalVote(&_DPoS.CallOpts, _proposalId, _voter)
}

// GetSidechainProposalVote is a free data retrieval call binding the contract method 0x8515b0e2.
//
// Solidity: function getSidechainProposalVote(uint256 _proposalId, address _voter) constant returns(uint8)
func (_DPoS *DPoSCallerSession) GetSidechainProposalVote(_proposalId *big.Int, _voter common.Address) (uint8, error) {
	return _DPoS.Contract.GetSidechainProposalVote(&_DPoS.CallOpts, _proposalId, _voter)
}

// GetTotalValidatorStakingPool is a free data retrieval call binding the contract method 0x89ed7939.
//
// Solidity: function getTotalValidatorStakingPool() constant returns(uint256)
func (_DPoS *DPoSCaller) GetTotalValidatorStakingPool(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getTotalValidatorStakingPool")
	return *ret0, err
}

// GetTotalValidatorStakingPool is a free data retrieval call binding the contract method 0x89ed7939.
//
// Solidity: function getTotalValidatorStakingPool() constant returns(uint256)
func (_DPoS *DPoSSession) GetTotalValidatorStakingPool() (*big.Int, error) {
	return _DPoS.Contract.GetTotalValidatorStakingPool(&_DPoS.CallOpts)
}

// GetTotalValidatorStakingPool is a free data retrieval call binding the contract method 0x89ed7939.
//
// Solidity: function getTotalValidatorStakingPool() constant returns(uint256)
func (_DPoS *DPoSCallerSession) GetTotalValidatorStakingPool() (*big.Int, error) {
	return _DPoS.Contract.GetTotalValidatorStakingPool(&_DPoS.CallOpts)
}

// GetUIntValue is a free data retrieval call binding the contract method 0x64c66395.
//
// Solidity: function getUIntValue(uint256 _record) constant returns(uint256)
func (_DPoS *DPoSCaller) GetUIntValue(opts *bind.CallOpts, _record *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getUIntValue", _record)
	return *ret0, err
}

// GetUIntValue is a free data retrieval call binding the contract method 0x64c66395.
//
// Solidity: function getUIntValue(uint256 _record) constant returns(uint256)
func (_DPoS *DPoSSession) GetUIntValue(_record *big.Int) (*big.Int, error) {
	return _DPoS.Contract.GetUIntValue(&_DPoS.CallOpts, _record)
}

// GetUIntValue is a free data retrieval call binding the contract method 0x64c66395.
//
// Solidity: function getUIntValue(uint256 _record) constant returns(uint256)
func (_DPoS *DPoSCallerSession) GetUIntValue(_record *big.Int) (*big.Int, error) {
	return _DPoS.Contract.GetUIntValue(&_DPoS.CallOpts, _record)
}

// GetValidatorNum is a free data retrieval call binding the contract method 0x1cfe4f0b.
//
// Solidity: function getValidatorNum() constant returns(uint256)
func (_DPoS *DPoSCaller) GetValidatorNum(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "getValidatorNum")
	return *ret0, err
}

// GetValidatorNum is a free data retrieval call binding the contract method 0x1cfe4f0b.
//
// Solidity: function getValidatorNum() constant returns(uint256)
func (_DPoS *DPoSSession) GetValidatorNum() (*big.Int, error) {
	return _DPoS.Contract.GetValidatorNum(&_DPoS.CallOpts)
}

// GetValidatorNum is a free data retrieval call binding the contract method 0x1cfe4f0b.
//
// Solidity: function getValidatorNum() constant returns(uint256)
func (_DPoS *DPoSCallerSession) GetValidatorNum() (*big.Int, error) {
	return _DPoS.Contract.GetValidatorNum(&_DPoS.CallOpts)
}

// GovernToken is a free data retrieval call binding the contract method 0x51abe57b.
//
// Solidity: function governToken() constant returns(address)
func (_DPoS *DPoSCaller) GovernToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "governToken")
	return *ret0, err
}

// GovernToken is a free data retrieval call binding the contract method 0x51abe57b.
//
// Solidity: function governToken() constant returns(address)
func (_DPoS *DPoSSession) GovernToken() (common.Address, error) {
	return _DPoS.Contract.GovernToken(&_DPoS.CallOpts)
}

// GovernToken is a free data retrieval call binding the contract method 0x51abe57b.
//
// Solidity: function governToken() constant returns(address)
func (_DPoS *DPoSCallerSession) GovernToken() (common.Address, error) {
	return _DPoS.Contract.GovernToken(&_DPoS.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DPoS *DPoSCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DPoS *DPoSSession) IsOwner() (bool, error) {
	return _DPoS.Contract.IsOwner(&_DPoS.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_DPoS *DPoSCallerSession) IsOwner() (bool, error) {
	return _DPoS.Contract.IsOwner(&_DPoS.CallOpts)
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(address account) constant returns(bool)
func (_DPoS *DPoSCaller) IsPauser(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isPauser", account)
	return *ret0, err
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(address account) constant returns(bool)
func (_DPoS *DPoSSession) IsPauser(account common.Address) (bool, error) {
	return _DPoS.Contract.IsPauser(&_DPoS.CallOpts, account)
}

// IsPauser is a free data retrieval call binding the contract method 0x46fbf68e.
//
// Solidity: function isPauser(address account) constant returns(bool)
func (_DPoS *DPoSCallerSession) IsPauser(account common.Address) (bool, error) {
	return _DPoS.Contract.IsPauser(&_DPoS.CallOpts, account)
}

// IsSidechainRegistered is a free data retrieval call binding the contract method 0x325820b3.
//
// Solidity: function isSidechainRegistered(address _sidechainAddr) constant returns(bool)
func (_DPoS *DPoSCaller) IsSidechainRegistered(opts *bind.CallOpts, _sidechainAddr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isSidechainRegistered", _sidechainAddr)
	return *ret0, err
}

// IsSidechainRegistered is a free data retrieval call binding the contract method 0x325820b3.
//
// Solidity: function isSidechainRegistered(address _sidechainAddr) constant returns(bool)
func (_DPoS *DPoSSession) IsSidechainRegistered(_sidechainAddr common.Address) (bool, error) {
	return _DPoS.Contract.IsSidechainRegistered(&_DPoS.CallOpts, _sidechainAddr)
}

// IsSidechainRegistered is a free data retrieval call binding the contract method 0x325820b3.
//
// Solidity: function isSidechainRegistered(address _sidechainAddr) constant returns(bool)
func (_DPoS *DPoSCallerSession) IsSidechainRegistered(_sidechainAddr common.Address) (bool, error) {
	return _DPoS.Contract.IsSidechainRegistered(&_DPoS.CallOpts, _sidechainAddr)
}

// IsValidDPoS is a free data retrieval call binding the contract method 0xeab2ed8c.
//
// Solidity: function isValidDPoS() constant returns(bool)
func (_DPoS *DPoSCaller) IsValidDPoS(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isValidDPoS")
	return *ret0, err
}

// IsValidDPoS is a free data retrieval call binding the contract method 0xeab2ed8c.
//
// Solidity: function isValidDPoS() constant returns(bool)
func (_DPoS *DPoSSession) IsValidDPoS() (bool, error) {
	return _DPoS.Contract.IsValidDPoS(&_DPoS.CallOpts)
}

// IsValidDPoS is a free data retrieval call binding the contract method 0xeab2ed8c.
//
// Solidity: function isValidDPoS() constant returns(bool)
func (_DPoS *DPoSCallerSession) IsValidDPoS() (bool, error) {
	return _DPoS.Contract.IsValidDPoS(&_DPoS.CallOpts)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address _addr) constant returns(bool)
func (_DPoS *DPoSCaller) IsValidator(opts *bind.CallOpts, _addr common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isValidator", _addr)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address _addr) constant returns(bool)
func (_DPoS *DPoSSession) IsValidator(_addr common.Address) (bool, error) {
	return _DPoS.Contract.IsValidator(&_DPoS.CallOpts, _addr)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address _addr) constant returns(bool)
func (_DPoS *DPoSCallerSession) IsValidator(_addr common.Address) (bool, error) {
	return _DPoS.Contract.IsValidator(&_DPoS.CallOpts, _addr)
}

// IsWhitelistAdmin is a free data retrieval call binding the contract method 0xbb5f747b.
//
// Solidity: function isWhitelistAdmin(address account) constant returns(bool)
func (_DPoS *DPoSCaller) IsWhitelistAdmin(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isWhitelistAdmin", account)
	return *ret0, err
}

// IsWhitelistAdmin is a free data retrieval call binding the contract method 0xbb5f747b.
//
// Solidity: function isWhitelistAdmin(address account) constant returns(bool)
func (_DPoS *DPoSSession) IsWhitelistAdmin(account common.Address) (bool, error) {
	return _DPoS.Contract.IsWhitelistAdmin(&_DPoS.CallOpts, account)
}

// IsWhitelistAdmin is a free data retrieval call binding the contract method 0xbb5f747b.
//
// Solidity: function isWhitelistAdmin(address account) constant returns(bool)
func (_DPoS *DPoSCallerSession) IsWhitelistAdmin(account common.Address) (bool, error) {
	return _DPoS.Contract.IsWhitelistAdmin(&_DPoS.CallOpts, account)
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x3af32abf.
//
// Solidity: function isWhitelisted(address account) constant returns(bool)
func (_DPoS *DPoSCaller) IsWhitelisted(opts *bind.CallOpts, account common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "isWhitelisted", account)
	return *ret0, err
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x3af32abf.
//
// Solidity: function isWhitelisted(address account) constant returns(bool)
func (_DPoS *DPoSSession) IsWhitelisted(account common.Address) (bool, error) {
	return _DPoS.Contract.IsWhitelisted(&_DPoS.CallOpts, account)
}

// IsWhitelisted is a free data retrieval call binding the contract method 0x3af32abf.
//
// Solidity: function isWhitelisted(address account) constant returns(bool)
func (_DPoS *DPoSCallerSession) IsWhitelisted(account common.Address) (bool, error) {
	return _DPoS.Contract.IsWhitelisted(&_DPoS.CallOpts, account)
}

// MiningPool is a free data retrieval call binding the contract method 0x73397597.
//
// Solidity: function miningPool() constant returns(uint256)
func (_DPoS *DPoSCaller) MiningPool(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "miningPool")
	return *ret0, err
}

// MiningPool is a free data retrieval call binding the contract method 0x73397597.
//
// Solidity: function miningPool() constant returns(uint256)
func (_DPoS *DPoSSession) MiningPool() (*big.Int, error) {
	return _DPoS.Contract.MiningPool(&_DPoS.CallOpts)
}

// MiningPool is a free data retrieval call binding the contract method 0x73397597.
//
// Solidity: function miningPool() constant returns(uint256)
func (_DPoS *DPoSCallerSession) MiningPool() (*big.Int, error) {
	return _DPoS.Contract.MiningPool(&_DPoS.CallOpts)
}

// NextParamProposalId is a free data retrieval call binding the contract method 0x22da7927.
//
// Solidity: function nextParamProposalId() constant returns(uint256)
func (_DPoS *DPoSCaller) NextParamProposalId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "nextParamProposalId")
	return *ret0, err
}

// NextParamProposalId is a free data retrieval call binding the contract method 0x22da7927.
//
// Solidity: function nextParamProposalId() constant returns(uint256)
func (_DPoS *DPoSSession) NextParamProposalId() (*big.Int, error) {
	return _DPoS.Contract.NextParamProposalId(&_DPoS.CallOpts)
}

// NextParamProposalId is a free data retrieval call binding the contract method 0x22da7927.
//
// Solidity: function nextParamProposalId() constant returns(uint256)
func (_DPoS *DPoSCallerSession) NextParamProposalId() (*big.Int, error) {
	return _DPoS.Contract.NextParamProposalId(&_DPoS.CallOpts)
}

// NextSidechainProposalId is a free data retrieval call binding the contract method 0x2bf0fe59.
//
// Solidity: function nextSidechainProposalId() constant returns(uint256)
func (_DPoS *DPoSCaller) NextSidechainProposalId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "nextSidechainProposalId")
	return *ret0, err
}

// NextSidechainProposalId is a free data retrieval call binding the contract method 0x2bf0fe59.
//
// Solidity: function nextSidechainProposalId() constant returns(uint256)
func (_DPoS *DPoSSession) NextSidechainProposalId() (*big.Int, error) {
	return _DPoS.Contract.NextSidechainProposalId(&_DPoS.CallOpts)
}

// NextSidechainProposalId is a free data retrieval call binding the contract method 0x2bf0fe59.
//
// Solidity: function nextSidechainProposalId() constant returns(uint256)
func (_DPoS *DPoSCallerSession) NextSidechainProposalId() (*big.Int, error) {
	return _DPoS.Contract.NextSidechainProposalId(&_DPoS.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DPoS *DPoSCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DPoS *DPoSSession) Owner() (common.Address, error) {
	return _DPoS.Contract.Owner(&_DPoS.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_DPoS *DPoSCallerSession) Owner() (common.Address, error) {
	return _DPoS.Contract.Owner(&_DPoS.CallOpts)
}

// ParamProposals is a free data retrieval call binding the contract method 0x7e5fb8f3.
//
// Solidity: function paramProposals(uint256 ) constant returns(address proposer, uint256 deposit, uint256 voteDeadline, uint256 record, uint256 newValue, uint8 status)
func (_DPoS *DPoSCaller) ParamProposals(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Proposer     common.Address
	Deposit      *big.Int
	VoteDeadline *big.Int
	Record       *big.Int
	NewValue     *big.Int
	Status       uint8
}, error) {
	ret := new(struct {
		Proposer     common.Address
		Deposit      *big.Int
		VoteDeadline *big.Int
		Record       *big.Int
		NewValue     *big.Int
		Status       uint8
	})
	out := ret
	err := _DPoS.contract.Call(opts, out, "paramProposals", arg0)
	return *ret, err
}

// ParamProposals is a free data retrieval call binding the contract method 0x7e5fb8f3.
//
// Solidity: function paramProposals(uint256 ) constant returns(address proposer, uint256 deposit, uint256 voteDeadline, uint256 record, uint256 newValue, uint8 status)
func (_DPoS *DPoSSession) ParamProposals(arg0 *big.Int) (struct {
	Proposer     common.Address
	Deposit      *big.Int
	VoteDeadline *big.Int
	Record       *big.Int
	NewValue     *big.Int
	Status       uint8
}, error) {
	return _DPoS.Contract.ParamProposals(&_DPoS.CallOpts, arg0)
}

// ParamProposals is a free data retrieval call binding the contract method 0x7e5fb8f3.
//
// Solidity: function paramProposals(uint256 ) constant returns(address proposer, uint256 deposit, uint256 voteDeadline, uint256 record, uint256 newValue, uint8 status)
func (_DPoS *DPoSCallerSession) ParamProposals(arg0 *big.Int) (struct {
	Proposer     common.Address
	Deposit      *big.Int
	VoteDeadline *big.Int
	Record       *big.Int
	NewValue     *big.Int
	Status       uint8
}, error) {
	return _DPoS.Contract.ParamProposals(&_DPoS.CallOpts, arg0)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_DPoS *DPoSCaller) Paused(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "paused")
	return *ret0, err
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_DPoS *DPoSSession) Paused() (bool, error) {
	return _DPoS.Contract.Paused(&_DPoS.CallOpts)
}

// Paused is a free data retrieval call binding the contract method 0x5c975abb.
//
// Solidity: function paused() constant returns(bool)
func (_DPoS *DPoSCallerSession) Paused() (bool, error) {
	return _DPoS.Contract.Paused(&_DPoS.CallOpts)
}

// RedeemedMiningReward is a free data retrieval call binding the contract method 0x87e53fef.
//
// Solidity: function redeemedMiningReward(address ) constant returns(uint256)
func (_DPoS *DPoSCaller) RedeemedMiningReward(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "redeemedMiningReward", arg0)
	return *ret0, err
}

// RedeemedMiningReward is a free data retrieval call binding the contract method 0x87e53fef.
//
// Solidity: function redeemedMiningReward(address ) constant returns(uint256)
func (_DPoS *DPoSSession) RedeemedMiningReward(arg0 common.Address) (*big.Int, error) {
	return _DPoS.Contract.RedeemedMiningReward(&_DPoS.CallOpts, arg0)
}

// RedeemedMiningReward is a free data retrieval call binding the contract method 0x87e53fef.
//
// Solidity: function redeemedMiningReward(address ) constant returns(uint256)
func (_DPoS *DPoSCallerSession) RedeemedMiningReward(arg0 common.Address) (*big.Int, error) {
	return _DPoS.Contract.RedeemedMiningReward(&_DPoS.CallOpts, arg0)
}

// RegisteredSidechains is a free data retrieval call binding the contract method 0x49444b71.
//
// Solidity: function registeredSidechains(address ) constant returns(bool)
func (_DPoS *DPoSCaller) RegisteredSidechains(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "registeredSidechains", arg0)
	return *ret0, err
}

// RegisteredSidechains is a free data retrieval call binding the contract method 0x49444b71.
//
// Solidity: function registeredSidechains(address ) constant returns(bool)
func (_DPoS *DPoSSession) RegisteredSidechains(arg0 common.Address) (bool, error) {
	return _DPoS.Contract.RegisteredSidechains(&_DPoS.CallOpts, arg0)
}

// RegisteredSidechains is a free data retrieval call binding the contract method 0x49444b71.
//
// Solidity: function registeredSidechains(address ) constant returns(bool)
func (_DPoS *DPoSCallerSession) RegisteredSidechains(arg0 common.Address) (bool, error) {
	return _DPoS.Contract.RegisteredSidechains(&_DPoS.CallOpts, arg0)
}

// SidechainProposals is a free data retrieval call binding the contract method 0xe97b7452.
//
// Solidity: function sidechainProposals(uint256 ) constant returns(address proposer, uint256 deposit, uint256 voteDeadline, address sidechainAddr, bool registered, uint8 status)
func (_DPoS *DPoSCaller) SidechainProposals(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Proposer      common.Address
	Deposit       *big.Int
	VoteDeadline  *big.Int
	SidechainAddr common.Address
	Registered    bool
	Status        uint8
}, error) {
	ret := new(struct {
		Proposer      common.Address
		Deposit       *big.Int
		VoteDeadline  *big.Int
		SidechainAddr common.Address
		Registered    bool
		Status        uint8
	})
	out := ret
	err := _DPoS.contract.Call(opts, out, "sidechainProposals", arg0)
	return *ret, err
}

// SidechainProposals is a free data retrieval call binding the contract method 0xe97b7452.
//
// Solidity: function sidechainProposals(uint256 ) constant returns(address proposer, uint256 deposit, uint256 voteDeadline, address sidechainAddr, bool registered, uint8 status)
func (_DPoS *DPoSSession) SidechainProposals(arg0 *big.Int) (struct {
	Proposer      common.Address
	Deposit       *big.Int
	VoteDeadline  *big.Int
	SidechainAddr common.Address
	Registered    bool
	Status        uint8
}, error) {
	return _DPoS.Contract.SidechainProposals(&_DPoS.CallOpts, arg0)
}

// SidechainProposals is a free data retrieval call binding the contract method 0xe97b7452.
//
// Solidity: function sidechainProposals(uint256 ) constant returns(address proposer, uint256 deposit, uint256 voteDeadline, address sidechainAddr, bool registered, uint8 status)
func (_DPoS *DPoSCallerSession) SidechainProposals(arg0 *big.Int) (struct {
	Proposer      common.Address
	Deposit       *big.Int
	VoteDeadline  *big.Int
	SidechainAddr common.Address
	Registered    bool
	Status        uint8
}, error) {
	return _DPoS.Contract.SidechainProposals(&_DPoS.CallOpts, arg0)
}

// UsedPenaltyNonce is a free data retrieval call binding the contract method 0xfb878749.
//
// Solidity: function usedPenaltyNonce(uint256 ) constant returns(bool)
func (_DPoS *DPoSCaller) UsedPenaltyNonce(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "usedPenaltyNonce", arg0)
	return *ret0, err
}

// UsedPenaltyNonce is a free data retrieval call binding the contract method 0xfb878749.
//
// Solidity: function usedPenaltyNonce(uint256 ) constant returns(bool)
func (_DPoS *DPoSSession) UsedPenaltyNonce(arg0 *big.Int) (bool, error) {
	return _DPoS.Contract.UsedPenaltyNonce(&_DPoS.CallOpts, arg0)
}

// UsedPenaltyNonce is a free data retrieval call binding the contract method 0xfb878749.
//
// Solidity: function usedPenaltyNonce(uint256 ) constant returns(bool)
func (_DPoS *DPoSCallerSession) UsedPenaltyNonce(arg0 *big.Int) (bool, error) {
	return _DPoS.Contract.UsedPenaltyNonce(&_DPoS.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) constant returns(address)
func (_DPoS *DPoSCaller) ValidatorSet(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _DPoS.contract.Call(opts, out, "validatorSet", arg0)
	return *ret0, err
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) constant returns(address)
func (_DPoS *DPoSSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _DPoS.Contract.ValidatorSet(&_DPoS.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0xe64808f3.
//
// Solidity: function validatorSet(uint256 ) constant returns(address)
func (_DPoS *DPoSCallerSession) ValidatorSet(arg0 *big.Int) (common.Address, error) {
	return _DPoS.Contract.ValidatorSet(&_DPoS.CallOpts, arg0)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(address account) returns()
func (_DPoS *DPoSTransactor) AddPauser(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "addPauser", account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(address account) returns()
func (_DPoS *DPoSSession) AddPauser(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.AddPauser(&_DPoS.TransactOpts, account)
}

// AddPauser is a paid mutator transaction binding the contract method 0x82dc1ec4.
//
// Solidity: function addPauser(address account) returns()
func (_DPoS *DPoSTransactorSession) AddPauser(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.AddPauser(&_DPoS.TransactOpts, account)
}

// AddWhitelistAdmin is a paid mutator transaction binding the contract method 0x7362d9c8.
//
// Solidity: function addWhitelistAdmin(address account) returns()
func (_DPoS *DPoSTransactor) AddWhitelistAdmin(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "addWhitelistAdmin", account)
}

// AddWhitelistAdmin is a paid mutator transaction binding the contract method 0x7362d9c8.
//
// Solidity: function addWhitelistAdmin(address account) returns()
func (_DPoS *DPoSSession) AddWhitelistAdmin(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.AddWhitelistAdmin(&_DPoS.TransactOpts, account)
}

// AddWhitelistAdmin is a paid mutator transaction binding the contract method 0x7362d9c8.
//
// Solidity: function addWhitelistAdmin(address account) returns()
func (_DPoS *DPoSTransactorSession) AddWhitelistAdmin(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.AddWhitelistAdmin(&_DPoS.TransactOpts, account)
}

// AddWhitelisted is a paid mutator transaction binding the contract method 0x10154bad.
//
// Solidity: function addWhitelisted(address account) returns()
func (_DPoS *DPoSTransactor) AddWhitelisted(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "addWhitelisted", account)
}

// AddWhitelisted is a paid mutator transaction binding the contract method 0x10154bad.
//
// Solidity: function addWhitelisted(address account) returns()
func (_DPoS *DPoSSession) AddWhitelisted(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.AddWhitelisted(&_DPoS.TransactOpts, account)
}

// AddWhitelisted is a paid mutator transaction binding the contract method 0x10154bad.
//
// Solidity: function addWhitelisted(address account) returns()
func (_DPoS *DPoSTransactorSession) AddWhitelisted(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.AddWhitelisted(&_DPoS.TransactOpts, account)
}

// AnnounceIncreaseCommissionRate is a paid mutator transaction binding the contract method 0xf64f33f2.
//
// Solidity: function announceIncreaseCommissionRate(uint256 _newRate, uint256 _newLockEndTime) returns()
func (_DPoS *DPoSTransactor) AnnounceIncreaseCommissionRate(opts *bind.TransactOpts, _newRate *big.Int, _newLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "announceIncreaseCommissionRate", _newRate, _newLockEndTime)
}

// AnnounceIncreaseCommissionRate is a paid mutator transaction binding the contract method 0xf64f33f2.
//
// Solidity: function announceIncreaseCommissionRate(uint256 _newRate, uint256 _newLockEndTime) returns()
func (_DPoS *DPoSSession) AnnounceIncreaseCommissionRate(_newRate *big.Int, _newLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.AnnounceIncreaseCommissionRate(&_DPoS.TransactOpts, _newRate, _newLockEndTime)
}

// AnnounceIncreaseCommissionRate is a paid mutator transaction binding the contract method 0xf64f33f2.
//
// Solidity: function announceIncreaseCommissionRate(uint256 _newRate, uint256 _newLockEndTime) returns()
func (_DPoS *DPoSTransactorSession) AnnounceIncreaseCommissionRate(_newRate *big.Int, _newLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.AnnounceIncreaseCommissionRate(&_DPoS.TransactOpts, _newRate, _newLockEndTime)
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x6e7cf85d.
//
// Solidity: function claimValidator() returns()
func (_DPoS *DPoSTransactor) ClaimValidator(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "claimValidator")
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x6e7cf85d.
//
// Solidity: function claimValidator() returns()
func (_DPoS *DPoSSession) ClaimValidator() (*types.Transaction, error) {
	return _DPoS.Contract.ClaimValidator(&_DPoS.TransactOpts)
}

// ClaimValidator is a paid mutator transaction binding the contract method 0x6e7cf85d.
//
// Solidity: function claimValidator() returns()
func (_DPoS *DPoSTransactorSession) ClaimValidator() (*types.Transaction, error) {
	return _DPoS.Contract.ClaimValidator(&_DPoS.TransactOpts)
}

// ConfirmIncreaseCommissionRate is a paid mutator transaction binding the contract method 0x6e997565.
//
// Solidity: function confirmIncreaseCommissionRate() returns()
func (_DPoS *DPoSTransactor) ConfirmIncreaseCommissionRate(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "confirmIncreaseCommissionRate")
}

// ConfirmIncreaseCommissionRate is a paid mutator transaction binding the contract method 0x6e997565.
//
// Solidity: function confirmIncreaseCommissionRate() returns()
func (_DPoS *DPoSSession) ConfirmIncreaseCommissionRate() (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmIncreaseCommissionRate(&_DPoS.TransactOpts)
}

// ConfirmIncreaseCommissionRate is a paid mutator transaction binding the contract method 0x6e997565.
//
// Solidity: function confirmIncreaseCommissionRate() returns()
func (_DPoS *DPoSTransactorSession) ConfirmIncreaseCommissionRate() (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmIncreaseCommissionRate(&_DPoS.TransactOpts)
}

// ConfirmParamProposal is a paid mutator transaction binding the contract method 0x934a18ec.
//
// Solidity: function confirmParamProposal(uint256 _proposalId) returns()
func (_DPoS *DPoSTransactor) ConfirmParamProposal(opts *bind.TransactOpts, _proposalId *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "confirmParamProposal", _proposalId)
}

// ConfirmParamProposal is a paid mutator transaction binding the contract method 0x934a18ec.
//
// Solidity: function confirmParamProposal(uint256 _proposalId) returns()
func (_DPoS *DPoSSession) ConfirmParamProposal(_proposalId *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmParamProposal(&_DPoS.TransactOpts, _proposalId)
}

// ConfirmParamProposal is a paid mutator transaction binding the contract method 0x934a18ec.
//
// Solidity: function confirmParamProposal(uint256 _proposalId) returns()
func (_DPoS *DPoSTransactorSession) ConfirmParamProposal(_proposalId *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmParamProposal(&_DPoS.TransactOpts, _proposalId)
}

// ConfirmSidechainProposal is a paid mutator transaction binding the contract method 0x1a06f737.
//
// Solidity: function confirmSidechainProposal(uint256 _proposalId) returns()
func (_DPoS *DPoSTransactor) ConfirmSidechainProposal(opts *bind.TransactOpts, _proposalId *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "confirmSidechainProposal", _proposalId)
}

// ConfirmSidechainProposal is a paid mutator transaction binding the contract method 0x1a06f737.
//
// Solidity: function confirmSidechainProposal(uint256 _proposalId) returns()
func (_DPoS *DPoSSession) ConfirmSidechainProposal(_proposalId *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmSidechainProposal(&_DPoS.TransactOpts, _proposalId)
}

// ConfirmSidechainProposal is a paid mutator transaction binding the contract method 0x1a06f737.
//
// Solidity: function confirmSidechainProposal(uint256 _proposalId) returns()
func (_DPoS *DPoSTransactorSession) ConfirmSidechainProposal(_proposalId *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmSidechainProposal(&_DPoS.TransactOpts, _proposalId)
}

// ConfirmUnbondedCandidate is a paid mutator transaction binding the contract method 0xc7ec2f35.
//
// Solidity: function confirmUnbondedCandidate(address _candidateAddr) returns()
func (_DPoS *DPoSTransactor) ConfirmUnbondedCandidate(opts *bind.TransactOpts, _candidateAddr common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "confirmUnbondedCandidate", _candidateAddr)
}

// ConfirmUnbondedCandidate is a paid mutator transaction binding the contract method 0xc7ec2f35.
//
// Solidity: function confirmUnbondedCandidate(address _candidateAddr) returns()
func (_DPoS *DPoSSession) ConfirmUnbondedCandidate(_candidateAddr common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmUnbondedCandidate(&_DPoS.TransactOpts, _candidateAddr)
}

// ConfirmUnbondedCandidate is a paid mutator transaction binding the contract method 0xc7ec2f35.
//
// Solidity: function confirmUnbondedCandidate(address _candidateAddr) returns()
func (_DPoS *DPoSTransactorSession) ConfirmUnbondedCandidate(_candidateAddr common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmUnbondedCandidate(&_DPoS.TransactOpts, _candidateAddr)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0xd2bfc1c7.
//
// Solidity: function confirmWithdraw(address _candidateAddr) returns()
func (_DPoS *DPoSTransactor) ConfirmWithdraw(opts *bind.TransactOpts, _candidateAddr common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "confirmWithdraw", _candidateAddr)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0xd2bfc1c7.
//
// Solidity: function confirmWithdraw(address _candidateAddr) returns()
func (_DPoS *DPoSSession) ConfirmWithdraw(_candidateAddr common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmWithdraw(&_DPoS.TransactOpts, _candidateAddr)
}

// ConfirmWithdraw is a paid mutator transaction binding the contract method 0xd2bfc1c7.
//
// Solidity: function confirmWithdraw(address _candidateAddr) returns()
func (_DPoS *DPoSTransactorSession) ConfirmWithdraw(_candidateAddr common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.ConfirmWithdraw(&_DPoS.TransactOpts, _candidateAddr)
}

// ContributeToMiningPool is a paid mutator transaction binding the contract method 0x4b7dba6b.
//
// Solidity: function contributeToMiningPool(uint256 _amount) returns()
func (_DPoS *DPoSTransactor) ContributeToMiningPool(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "contributeToMiningPool", _amount)
}

// ContributeToMiningPool is a paid mutator transaction binding the contract method 0x4b7dba6b.
//
// Solidity: function contributeToMiningPool(uint256 _amount) returns()
func (_DPoS *DPoSSession) ContributeToMiningPool(_amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.ContributeToMiningPool(&_DPoS.TransactOpts, _amount)
}

// ContributeToMiningPool is a paid mutator transaction binding the contract method 0x4b7dba6b.
//
// Solidity: function contributeToMiningPool(uint256 _amount) returns()
func (_DPoS *DPoSTransactorSession) ContributeToMiningPool(_amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.ContributeToMiningPool(&_DPoS.TransactOpts, _amount)
}

// CreateParamProposal is a paid mutator transaction binding the contract method 0x3090c0e9.
//
// Solidity: function createParamProposal(uint256 _record, uint256 _value) returns()
func (_DPoS *DPoSTransactor) CreateParamProposal(opts *bind.TransactOpts, _record *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "createParamProposal", _record, _value)
}

// CreateParamProposal is a paid mutator transaction binding the contract method 0x3090c0e9.
//
// Solidity: function createParamProposal(uint256 _record, uint256 _value) returns()
func (_DPoS *DPoSSession) CreateParamProposal(_record *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.CreateParamProposal(&_DPoS.TransactOpts, _record, _value)
}

// CreateParamProposal is a paid mutator transaction binding the contract method 0x3090c0e9.
//
// Solidity: function createParamProposal(uint256 _record, uint256 _value) returns()
func (_DPoS *DPoSTransactorSession) CreateParamProposal(_record *big.Int, _value *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.CreateParamProposal(&_DPoS.TransactOpts, _record, _value)
}

// CreateSidechainProposal is a paid mutator transaction binding the contract method 0xe433c1ca.
//
// Solidity: function createSidechainProposal(address _sidechainAddr, bool _registered) returns()
func (_DPoS *DPoSTransactor) CreateSidechainProposal(opts *bind.TransactOpts, _sidechainAddr common.Address, _registered bool) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "createSidechainProposal", _sidechainAddr, _registered)
}

// CreateSidechainProposal is a paid mutator transaction binding the contract method 0xe433c1ca.
//
// Solidity: function createSidechainProposal(address _sidechainAddr, bool _registered) returns()
func (_DPoS *DPoSSession) CreateSidechainProposal(_sidechainAddr common.Address, _registered bool) (*types.Transaction, error) {
	return _DPoS.Contract.CreateSidechainProposal(&_DPoS.TransactOpts, _sidechainAddr, _registered)
}

// CreateSidechainProposal is a paid mutator transaction binding the contract method 0xe433c1ca.
//
// Solidity: function createSidechainProposal(address _sidechainAddr, bool _registered) returns()
func (_DPoS *DPoSTransactorSession) CreateSidechainProposal(_sidechainAddr common.Address, _registered bool) (*types.Transaction, error) {
	return _DPoS.Contract.CreateSidechainProposal(&_DPoS.TransactOpts, _sidechainAddr, _registered)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSTransactor) Delegate(opts *bind.TransactOpts, _candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "delegate", _candidateAddr, _amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSSession) Delegate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.Delegate(&_DPoS.TransactOpts, _candidateAddr, _amount)
}

// Delegate is a paid mutator transaction binding the contract method 0x026e402b.
//
// Solidity: function delegate(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSTransactorSession) Delegate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.Delegate(&_DPoS.TransactOpts, _candidateAddr, _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0x145aa116.
//
// Solidity: function drainToken(uint256 _amount) returns()
func (_DPoS *DPoSTransactor) DrainToken(opts *bind.TransactOpts, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "drainToken", _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0x145aa116.
//
// Solidity: function drainToken(uint256 _amount) returns()
func (_DPoS *DPoSSession) DrainToken(_amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.DrainToken(&_DPoS.TransactOpts, _amount)
}

// DrainToken is a paid mutator transaction binding the contract method 0x145aa116.
//
// Solidity: function drainToken(uint256 _amount) returns()
func (_DPoS *DPoSTransactorSession) DrainToken(_amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.DrainToken(&_DPoS.TransactOpts, _amount)
}

// InitializeCandidate is a paid mutator transaction binding the contract method 0xc1e16718.
//
// Solidity: function initializeCandidate(uint256 _minSelfStake, uint256 _commissionRate, uint256 _rateLockEndTime) returns()
func (_DPoS *DPoSTransactor) InitializeCandidate(opts *bind.TransactOpts, _minSelfStake *big.Int, _commissionRate *big.Int, _rateLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "initializeCandidate", _minSelfStake, _commissionRate, _rateLockEndTime)
}

// InitializeCandidate is a paid mutator transaction binding the contract method 0xc1e16718.
//
// Solidity: function initializeCandidate(uint256 _minSelfStake, uint256 _commissionRate, uint256 _rateLockEndTime) returns()
func (_DPoS *DPoSSession) InitializeCandidate(_minSelfStake *big.Int, _commissionRate *big.Int, _rateLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.InitializeCandidate(&_DPoS.TransactOpts, _minSelfStake, _commissionRate, _rateLockEndTime)
}

// InitializeCandidate is a paid mutator transaction binding the contract method 0xc1e16718.
//
// Solidity: function initializeCandidate(uint256 _minSelfStake, uint256 _commissionRate, uint256 _rateLockEndTime) returns()
func (_DPoS *DPoSTransactorSession) InitializeCandidate(_minSelfStake *big.Int, _commissionRate *big.Int, _rateLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.InitializeCandidate(&_DPoS.TransactOpts, _minSelfStake, _commissionRate, _rateLockEndTime)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x785f8ffd.
//
// Solidity: function intendWithdraw(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSTransactor) IntendWithdraw(opts *bind.TransactOpts, _candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "intendWithdraw", _candidateAddr, _amount)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x785f8ffd.
//
// Solidity: function intendWithdraw(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSSession) IntendWithdraw(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.IntendWithdraw(&_DPoS.TransactOpts, _candidateAddr, _amount)
}

// IntendWithdraw is a paid mutator transaction binding the contract method 0x785f8ffd.
//
// Solidity: function intendWithdraw(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSTransactorSession) IntendWithdraw(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.IntendWithdraw(&_DPoS.TransactOpts, _candidateAddr, _amount)
}

// NonIncreaseCommissionRate is a paid mutator transaction binding the contract method 0xbe57959d.
//
// Solidity: function nonIncreaseCommissionRate(uint256 _newRate, uint256 _newLockEndTime) returns()
func (_DPoS *DPoSTransactor) NonIncreaseCommissionRate(opts *bind.TransactOpts, _newRate *big.Int, _newLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "nonIncreaseCommissionRate", _newRate, _newLockEndTime)
}

// NonIncreaseCommissionRate is a paid mutator transaction binding the contract method 0xbe57959d.
//
// Solidity: function nonIncreaseCommissionRate(uint256 _newRate, uint256 _newLockEndTime) returns()
func (_DPoS *DPoSSession) NonIncreaseCommissionRate(_newRate *big.Int, _newLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.NonIncreaseCommissionRate(&_DPoS.TransactOpts, _newRate, _newLockEndTime)
}

// NonIncreaseCommissionRate is a paid mutator transaction binding the contract method 0xbe57959d.
//
// Solidity: function nonIncreaseCommissionRate(uint256 _newRate, uint256 _newLockEndTime) returns()
func (_DPoS *DPoSTransactorSession) NonIncreaseCommissionRate(_newRate *big.Int, _newLockEndTime *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.NonIncreaseCommissionRate(&_DPoS.TransactOpts, _newRate, _newLockEndTime)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_DPoS *DPoSTransactor) Pause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "pause")
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_DPoS *DPoSSession) Pause() (*types.Transaction, error) {
	return _DPoS.Contract.Pause(&_DPoS.TransactOpts)
}

// Pause is a paid mutator transaction binding the contract method 0x8456cb59.
//
// Solidity: function pause() returns()
func (_DPoS *DPoSTransactorSession) Pause() (*types.Transaction, error) {
	return _DPoS.Contract.Pause(&_DPoS.TransactOpts)
}

// Punish is a paid mutator transaction binding the contract method 0x3620d149.
//
// Solidity: function punish(bytes _penaltyRequest) returns()
func (_DPoS *DPoSTransactor) Punish(opts *bind.TransactOpts, _penaltyRequest []byte) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "punish", _penaltyRequest)
}

// Punish is a paid mutator transaction binding the contract method 0x3620d149.
//
// Solidity: function punish(bytes _penaltyRequest) returns()
func (_DPoS *DPoSSession) Punish(_penaltyRequest []byte) (*types.Transaction, error) {
	return _DPoS.Contract.Punish(&_DPoS.TransactOpts, _penaltyRequest)
}

// Punish is a paid mutator transaction binding the contract method 0x3620d149.
//
// Solidity: function punish(bytes _penaltyRequest) returns()
func (_DPoS *DPoSTransactorSession) Punish(_penaltyRequest []byte) (*types.Transaction, error) {
	return _DPoS.Contract.Punish(&_DPoS.TransactOpts, _penaltyRequest)
}

// RedeemMiningReward is a paid mutator transaction binding the contract method 0x1f7b0886.
//
// Solidity: function redeemMiningReward(address _receiver, uint256 _cumulativeReward) returns()
func (_DPoS *DPoSTransactor) RedeemMiningReward(opts *bind.TransactOpts, _receiver common.Address, _cumulativeReward *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "redeemMiningReward", _receiver, _cumulativeReward)
}

// RedeemMiningReward is a paid mutator transaction binding the contract method 0x1f7b0886.
//
// Solidity: function redeemMiningReward(address _receiver, uint256 _cumulativeReward) returns()
func (_DPoS *DPoSSession) RedeemMiningReward(_receiver common.Address, _cumulativeReward *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.RedeemMiningReward(&_DPoS.TransactOpts, _receiver, _cumulativeReward)
}

// RedeemMiningReward is a paid mutator transaction binding the contract method 0x1f7b0886.
//
// Solidity: function redeemMiningReward(address _receiver, uint256 _cumulativeReward) returns()
func (_DPoS *DPoSTransactorSession) RedeemMiningReward(_receiver common.Address, _cumulativeReward *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.RedeemMiningReward(&_DPoS.TransactOpts, _receiver, _cumulativeReward)
}

// RegisterSidechain is a paid mutator transaction binding the contract method 0xaa09fbae.
//
// Solidity: function registerSidechain(address _addr) returns()
func (_DPoS *DPoSTransactor) RegisterSidechain(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "registerSidechain", _addr)
}

// RegisterSidechain is a paid mutator transaction binding the contract method 0xaa09fbae.
//
// Solidity: function registerSidechain(address _addr) returns()
func (_DPoS *DPoSSession) RegisterSidechain(_addr common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.RegisterSidechain(&_DPoS.TransactOpts, _addr)
}

// RegisterSidechain is a paid mutator transaction binding the contract method 0xaa09fbae.
//
// Solidity: function registerSidechain(address _addr) returns()
func (_DPoS *DPoSTransactorSession) RegisterSidechain(_addr common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.RegisterSidechain(&_DPoS.TransactOpts, _addr)
}

// RemoveWhitelisted is a paid mutator transaction binding the contract method 0x291d9549.
//
// Solidity: function removeWhitelisted(address account) returns()
func (_DPoS *DPoSTransactor) RemoveWhitelisted(opts *bind.TransactOpts, account common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "removeWhitelisted", account)
}

// RemoveWhitelisted is a paid mutator transaction binding the contract method 0x291d9549.
//
// Solidity: function removeWhitelisted(address account) returns()
func (_DPoS *DPoSSession) RemoveWhitelisted(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.RemoveWhitelisted(&_DPoS.TransactOpts, account)
}

// RemoveWhitelisted is a paid mutator transaction binding the contract method 0x291d9549.
//
// Solidity: function removeWhitelisted(address account) returns()
func (_DPoS *DPoSTransactorSession) RemoveWhitelisted(account common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.RemoveWhitelisted(&_DPoS.TransactOpts, account)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DPoS *DPoSTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DPoS *DPoSSession) RenounceOwnership() (*types.Transaction, error) {
	return _DPoS.Contract.RenounceOwnership(&_DPoS.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_DPoS *DPoSTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _DPoS.Contract.RenounceOwnership(&_DPoS.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_DPoS *DPoSTransactor) RenouncePauser(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "renouncePauser")
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_DPoS *DPoSSession) RenouncePauser() (*types.Transaction, error) {
	return _DPoS.Contract.RenouncePauser(&_DPoS.TransactOpts)
}

// RenouncePauser is a paid mutator transaction binding the contract method 0x6ef8d66d.
//
// Solidity: function renouncePauser() returns()
func (_DPoS *DPoSTransactorSession) RenouncePauser() (*types.Transaction, error) {
	return _DPoS.Contract.RenouncePauser(&_DPoS.TransactOpts)
}

// RenounceWhitelistAdmin is a paid mutator transaction binding the contract method 0x4c5a628c.
//
// Solidity: function renounceWhitelistAdmin() returns()
func (_DPoS *DPoSTransactor) RenounceWhitelistAdmin(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "renounceWhitelistAdmin")
}

// RenounceWhitelistAdmin is a paid mutator transaction binding the contract method 0x4c5a628c.
//
// Solidity: function renounceWhitelistAdmin() returns()
func (_DPoS *DPoSSession) RenounceWhitelistAdmin() (*types.Transaction, error) {
	return _DPoS.Contract.RenounceWhitelistAdmin(&_DPoS.TransactOpts)
}

// RenounceWhitelistAdmin is a paid mutator transaction binding the contract method 0x4c5a628c.
//
// Solidity: function renounceWhitelistAdmin() returns()
func (_DPoS *DPoSTransactorSession) RenounceWhitelistAdmin() (*types.Transaction, error) {
	return _DPoS.Contract.RenounceWhitelistAdmin(&_DPoS.TransactOpts)
}

// RenounceWhitelisted is a paid mutator transaction binding the contract method 0xd6cd9473.
//
// Solidity: function renounceWhitelisted() returns()
func (_DPoS *DPoSTransactor) RenounceWhitelisted(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "renounceWhitelisted")
}

// RenounceWhitelisted is a paid mutator transaction binding the contract method 0xd6cd9473.
//
// Solidity: function renounceWhitelisted() returns()
func (_DPoS *DPoSSession) RenounceWhitelisted() (*types.Transaction, error) {
	return _DPoS.Contract.RenounceWhitelisted(&_DPoS.TransactOpts)
}

// RenounceWhitelisted is a paid mutator transaction binding the contract method 0xd6cd9473.
//
// Solidity: function renounceWhitelisted() returns()
func (_DPoS *DPoSTransactorSession) RenounceWhitelisted() (*types.Transaction, error) {
	return _DPoS.Contract.RenounceWhitelisted(&_DPoS.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DPoS *DPoSTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DPoS *DPoSSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.TransferOwnership(&_DPoS.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DPoS *DPoSTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DPoS.Contract.TransferOwnership(&_DPoS.TransactOpts, newOwner)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_DPoS *DPoSTransactor) Unpause(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "unpause")
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_DPoS *DPoSSession) Unpause() (*types.Transaction, error) {
	return _DPoS.Contract.Unpause(&_DPoS.TransactOpts)
}

// Unpause is a paid mutator transaction binding the contract method 0x3f4ba83a.
//
// Solidity: function unpause() returns()
func (_DPoS *DPoSTransactorSession) Unpause() (*types.Transaction, error) {
	return _DPoS.Contract.Unpause(&_DPoS.TransactOpts)
}

// UpdateEnableWhitelist is a paid mutator transaction binding the contract method 0x2cb57c48.
//
// Solidity: function updateEnableWhitelist(bool _enable) returns()
func (_DPoS *DPoSTransactor) UpdateEnableWhitelist(opts *bind.TransactOpts, _enable bool) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "updateEnableWhitelist", _enable)
}

// UpdateEnableWhitelist is a paid mutator transaction binding the contract method 0x2cb57c48.
//
// Solidity: function updateEnableWhitelist(bool _enable) returns()
func (_DPoS *DPoSSession) UpdateEnableWhitelist(_enable bool) (*types.Transaction, error) {
	return _DPoS.Contract.UpdateEnableWhitelist(&_DPoS.TransactOpts, _enable)
}

// UpdateEnableWhitelist is a paid mutator transaction binding the contract method 0x2cb57c48.
//
// Solidity: function updateEnableWhitelist(bool _enable) returns()
func (_DPoS *DPoSTransactorSession) UpdateEnableWhitelist(_enable bool) (*types.Transaction, error) {
	return _DPoS.Contract.UpdateEnableWhitelist(&_DPoS.TransactOpts, _enable)
}

// ValidateMultiSigMessage is a paid mutator transaction binding the contract method 0x1c0efd9d.
//
// Solidity: function validateMultiSigMessage(bytes _request) returns(bool)
func (_DPoS *DPoSTransactor) ValidateMultiSigMessage(opts *bind.TransactOpts, _request []byte) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "validateMultiSigMessage", _request)
}

// ValidateMultiSigMessage is a paid mutator transaction binding the contract method 0x1c0efd9d.
//
// Solidity: function validateMultiSigMessage(bytes _request) returns(bool)
func (_DPoS *DPoSSession) ValidateMultiSigMessage(_request []byte) (*types.Transaction, error) {
	return _DPoS.Contract.ValidateMultiSigMessage(&_DPoS.TransactOpts, _request)
}

// ValidateMultiSigMessage is a paid mutator transaction binding the contract method 0x1c0efd9d.
//
// Solidity: function validateMultiSigMessage(bytes _request) returns(bool)
func (_DPoS *DPoSTransactorSession) ValidateMultiSigMessage(_request []byte) (*types.Transaction, error) {
	return _DPoS.Contract.ValidateMultiSigMessage(&_DPoS.TransactOpts, _request)
}

// VoteParam is a paid mutator transaction binding the contract method 0x25ed6b35.
//
// Solidity: function voteParam(uint256 _proposalId, uint8 _vote) returns()
func (_DPoS *DPoSTransactor) VoteParam(opts *bind.TransactOpts, _proposalId *big.Int, _vote uint8) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "voteParam", _proposalId, _vote)
}

// VoteParam is a paid mutator transaction binding the contract method 0x25ed6b35.
//
// Solidity: function voteParam(uint256 _proposalId, uint8 _vote) returns()
func (_DPoS *DPoSSession) VoteParam(_proposalId *big.Int, _vote uint8) (*types.Transaction, error) {
	return _DPoS.Contract.VoteParam(&_DPoS.TransactOpts, _proposalId, _vote)
}

// VoteParam is a paid mutator transaction binding the contract method 0x25ed6b35.
//
// Solidity: function voteParam(uint256 _proposalId, uint8 _vote) returns()
func (_DPoS *DPoSTransactorSession) VoteParam(_proposalId *big.Int, _vote uint8) (*types.Transaction, error) {
	return _DPoS.Contract.VoteParam(&_DPoS.TransactOpts, _proposalId, _vote)
}

// VoteSidechain is a paid mutator transaction binding the contract method 0x85bfe017.
//
// Solidity: function voteSidechain(uint256 _proposalId, uint8 _vote) returns()
func (_DPoS *DPoSTransactor) VoteSidechain(opts *bind.TransactOpts, _proposalId *big.Int, _vote uint8) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "voteSidechain", _proposalId, _vote)
}

// VoteSidechain is a paid mutator transaction binding the contract method 0x85bfe017.
//
// Solidity: function voteSidechain(uint256 _proposalId, uint8 _vote) returns()
func (_DPoS *DPoSSession) VoteSidechain(_proposalId *big.Int, _vote uint8) (*types.Transaction, error) {
	return _DPoS.Contract.VoteSidechain(&_DPoS.TransactOpts, _proposalId, _vote)
}

// VoteSidechain is a paid mutator transaction binding the contract method 0x85bfe017.
//
// Solidity: function voteSidechain(uint256 _proposalId, uint8 _vote) returns()
func (_DPoS *DPoSTransactorSession) VoteSidechain(_proposalId *big.Int, _vote uint8) (*types.Transaction, error) {
	return _DPoS.Contract.VoteSidechain(&_DPoS.TransactOpts, _proposalId, _vote)
}

// WithdrawFromUnbondedCandidate is a paid mutator transaction binding the contract method 0xbb9053d0.
//
// Solidity: function withdrawFromUnbondedCandidate(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSTransactor) WithdrawFromUnbondedCandidate(opts *bind.TransactOpts, _candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.contract.Transact(opts, "withdrawFromUnbondedCandidate", _candidateAddr, _amount)
}

// WithdrawFromUnbondedCandidate is a paid mutator transaction binding the contract method 0xbb9053d0.
//
// Solidity: function withdrawFromUnbondedCandidate(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSSession) WithdrawFromUnbondedCandidate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.WithdrawFromUnbondedCandidate(&_DPoS.TransactOpts, _candidateAddr, _amount)
}

// WithdrawFromUnbondedCandidate is a paid mutator transaction binding the contract method 0xbb9053d0.
//
// Solidity: function withdrawFromUnbondedCandidate(address _candidateAddr, uint256 _amount) returns()
func (_DPoS *DPoSTransactorSession) WithdrawFromUnbondedCandidate(_candidateAddr common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _DPoS.Contract.WithdrawFromUnbondedCandidate(&_DPoS.TransactOpts, _candidateAddr, _amount)
}

// DPoSCandidateUnbondedIterator is returned from FilterCandidateUnbonded and is used to iterate over the raw logs and unpacked data for CandidateUnbonded events raised by the DPoS contract.
type DPoSCandidateUnbondedIterator struct {
	Event *DPoSCandidateUnbonded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSCandidateUnbondedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSCandidateUnbonded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSCandidateUnbonded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSCandidateUnbondedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSCandidateUnbondedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSCandidateUnbonded represents a CandidateUnbonded event raised by the DPoS contract.
type DPoSCandidateUnbonded struct {
	Candidate common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterCandidateUnbonded is a free log retrieval operation binding the contract event 0xbe85a9a7aa606febeaa35606e49cd7324c63cf970f4f5fd0c7e983f42b20b219.
//
// Solidity: event CandidateUnbonded(address indexed candidate)
func (_DPoS *DPoSFilterer) FilterCandidateUnbonded(opts *bind.FilterOpts, candidate []common.Address) (*DPoSCandidateUnbondedIterator, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "CandidateUnbonded", candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSCandidateUnbondedIterator{contract: _DPoS.contract, event: "CandidateUnbonded", logs: logs, sub: sub}, nil
}

// WatchCandidateUnbonded is a free log subscription operation binding the contract event 0xbe85a9a7aa606febeaa35606e49cd7324c63cf970f4f5fd0c7e983f42b20b219.
//
// Solidity: event CandidateUnbonded(address indexed candidate)
func (_DPoS *DPoSFilterer) WatchCandidateUnbonded(opts *bind.WatchOpts, sink chan<- *DPoSCandidateUnbonded, candidate []common.Address) (event.Subscription, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "CandidateUnbonded", candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSCandidateUnbonded)
				if err := _DPoS.contract.UnpackLog(event, "CandidateUnbonded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCandidateUnbonded is a log parse operation binding the contract event 0xbe85a9a7aa606febeaa35606e49cd7324c63cf970f4f5fd0c7e983f42b20b219.
//
// Solidity: event CandidateUnbonded(address indexed candidate)
func (_DPoS *DPoSFilterer) ParseCandidateUnbonded(log types.Log) (*DPoSCandidateUnbonded, error) {
	event := new(DPoSCandidateUnbonded)
	if err := _DPoS.contract.UnpackLog(event, "CandidateUnbonded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSCommissionRateAnnouncementIterator is returned from FilterCommissionRateAnnouncement and is used to iterate over the raw logs and unpacked data for CommissionRateAnnouncement events raised by the DPoS contract.
type DPoSCommissionRateAnnouncementIterator struct {
	Event *DPoSCommissionRateAnnouncement // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSCommissionRateAnnouncementIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSCommissionRateAnnouncement)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSCommissionRateAnnouncement)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSCommissionRateAnnouncementIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSCommissionRateAnnouncementIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSCommissionRateAnnouncement represents a CommissionRateAnnouncement event raised by the DPoS contract.
type DPoSCommissionRateAnnouncement struct {
	Candidate            common.Address
	AnnouncedRate        *big.Int
	AnnouncedLockEndTime *big.Int
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterCommissionRateAnnouncement is a free log retrieval operation binding the contract event 0xd1388fca1fdda1adbe79c9535b48b22e71aa7815469abb61cdbab2a7b4ccd28a.
//
// Solidity: event CommissionRateAnnouncement(address indexed candidate, uint256 announcedRate, uint256 announcedLockEndTime)
func (_DPoS *DPoSFilterer) FilterCommissionRateAnnouncement(opts *bind.FilterOpts, candidate []common.Address) (*DPoSCommissionRateAnnouncementIterator, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "CommissionRateAnnouncement", candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSCommissionRateAnnouncementIterator{contract: _DPoS.contract, event: "CommissionRateAnnouncement", logs: logs, sub: sub}, nil
}

// WatchCommissionRateAnnouncement is a free log subscription operation binding the contract event 0xd1388fca1fdda1adbe79c9535b48b22e71aa7815469abb61cdbab2a7b4ccd28a.
//
// Solidity: event CommissionRateAnnouncement(address indexed candidate, uint256 announcedRate, uint256 announcedLockEndTime)
func (_DPoS *DPoSFilterer) WatchCommissionRateAnnouncement(opts *bind.WatchOpts, sink chan<- *DPoSCommissionRateAnnouncement, candidate []common.Address) (event.Subscription, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "CommissionRateAnnouncement", candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSCommissionRateAnnouncement)
				if err := _DPoS.contract.UnpackLog(event, "CommissionRateAnnouncement", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCommissionRateAnnouncement is a log parse operation binding the contract event 0xd1388fca1fdda1adbe79c9535b48b22e71aa7815469abb61cdbab2a7b4ccd28a.
//
// Solidity: event CommissionRateAnnouncement(address indexed candidate, uint256 announcedRate, uint256 announcedLockEndTime)
func (_DPoS *DPoSFilterer) ParseCommissionRateAnnouncement(log types.Log) (*DPoSCommissionRateAnnouncement, error) {
	event := new(DPoSCommissionRateAnnouncement)
	if err := _DPoS.contract.UnpackLog(event, "CommissionRateAnnouncement", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSConfirmParamProposalIterator is returned from FilterConfirmParamProposal and is used to iterate over the raw logs and unpacked data for ConfirmParamProposal events raised by the DPoS contract.
type DPoSConfirmParamProposalIterator struct {
	Event *DPoSConfirmParamProposal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSConfirmParamProposalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSConfirmParamProposal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSConfirmParamProposal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSConfirmParamProposalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSConfirmParamProposalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSConfirmParamProposal represents a ConfirmParamProposal event raised by the DPoS contract.
type DPoSConfirmParamProposal struct {
	ProposalId *big.Int
	Passed     bool
	Record     *big.Int
	NewValue   *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterConfirmParamProposal is a free log retrieval operation binding the contract event 0x106f43a560e53395081c0423504b476d1a2cfed9d56ff972bf77ae43ff7d4ba4.
//
// Solidity: event ConfirmParamProposal(uint256 proposalId, bool passed, uint256 record, uint256 newValue)
func (_DPoS *DPoSFilterer) FilterConfirmParamProposal(opts *bind.FilterOpts) (*DPoSConfirmParamProposalIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "ConfirmParamProposal")
	if err != nil {
		return nil, err
	}
	return &DPoSConfirmParamProposalIterator{contract: _DPoS.contract, event: "ConfirmParamProposal", logs: logs, sub: sub}, nil
}

// WatchConfirmParamProposal is a free log subscription operation binding the contract event 0x106f43a560e53395081c0423504b476d1a2cfed9d56ff972bf77ae43ff7d4ba4.
//
// Solidity: event ConfirmParamProposal(uint256 proposalId, bool passed, uint256 record, uint256 newValue)
func (_DPoS *DPoSFilterer) WatchConfirmParamProposal(opts *bind.WatchOpts, sink chan<- *DPoSConfirmParamProposal) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "ConfirmParamProposal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSConfirmParamProposal)
				if err := _DPoS.contract.UnpackLog(event, "ConfirmParamProposal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfirmParamProposal is a log parse operation binding the contract event 0x106f43a560e53395081c0423504b476d1a2cfed9d56ff972bf77ae43ff7d4ba4.
//
// Solidity: event ConfirmParamProposal(uint256 proposalId, bool passed, uint256 record, uint256 newValue)
func (_DPoS *DPoSFilterer) ParseConfirmParamProposal(log types.Log) (*DPoSConfirmParamProposal, error) {
	event := new(DPoSConfirmParamProposal)
	if err := _DPoS.contract.UnpackLog(event, "ConfirmParamProposal", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSConfirmSidechainProposalIterator is returned from FilterConfirmSidechainProposal and is used to iterate over the raw logs and unpacked data for ConfirmSidechainProposal events raised by the DPoS contract.
type DPoSConfirmSidechainProposalIterator struct {
	Event *DPoSConfirmSidechainProposal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSConfirmSidechainProposalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSConfirmSidechainProposal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSConfirmSidechainProposal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSConfirmSidechainProposalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSConfirmSidechainProposalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSConfirmSidechainProposal represents a ConfirmSidechainProposal event raised by the DPoS contract.
type DPoSConfirmSidechainProposal struct {
	ProposalId    *big.Int
	Passed        bool
	SidechainAddr common.Address
	Registered    bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterConfirmSidechainProposal is a free log retrieval operation binding the contract event 0x2c26ff0b5547eb09df5dde3569782330829ac9ffa9811847beab5d466066801c.
//
// Solidity: event ConfirmSidechainProposal(uint256 proposalId, bool passed, address sidechainAddr, bool registered)
func (_DPoS *DPoSFilterer) FilterConfirmSidechainProposal(opts *bind.FilterOpts) (*DPoSConfirmSidechainProposalIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "ConfirmSidechainProposal")
	if err != nil {
		return nil, err
	}
	return &DPoSConfirmSidechainProposalIterator{contract: _DPoS.contract, event: "ConfirmSidechainProposal", logs: logs, sub: sub}, nil
}

// WatchConfirmSidechainProposal is a free log subscription operation binding the contract event 0x2c26ff0b5547eb09df5dde3569782330829ac9ffa9811847beab5d466066801c.
//
// Solidity: event ConfirmSidechainProposal(uint256 proposalId, bool passed, address sidechainAddr, bool registered)
func (_DPoS *DPoSFilterer) WatchConfirmSidechainProposal(opts *bind.WatchOpts, sink chan<- *DPoSConfirmSidechainProposal) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "ConfirmSidechainProposal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSConfirmSidechainProposal)
				if err := _DPoS.contract.UnpackLog(event, "ConfirmSidechainProposal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfirmSidechainProposal is a log parse operation binding the contract event 0x2c26ff0b5547eb09df5dde3569782330829ac9ffa9811847beab5d466066801c.
//
// Solidity: event ConfirmSidechainProposal(uint256 proposalId, bool passed, address sidechainAddr, bool registered)
func (_DPoS *DPoSFilterer) ParseConfirmSidechainProposal(log types.Log) (*DPoSConfirmSidechainProposal, error) {
	event := new(DPoSConfirmSidechainProposal)
	if err := _DPoS.contract.UnpackLog(event, "ConfirmSidechainProposal", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSConfirmWithdrawIterator is returned from FilterConfirmWithdraw and is used to iterate over the raw logs and unpacked data for ConfirmWithdraw events raised by the DPoS contract.
type DPoSConfirmWithdrawIterator struct {
	Event *DPoSConfirmWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSConfirmWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSConfirmWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSConfirmWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSConfirmWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSConfirmWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSConfirmWithdraw represents a ConfirmWithdraw event raised by the DPoS contract.
type DPoSConfirmWithdraw struct {
	Delegator common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterConfirmWithdraw is a free log retrieval operation binding the contract event 0x08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220.
//
// Solidity: event ConfirmWithdraw(address indexed delegator, address indexed candidate, uint256 amount)
func (_DPoS *DPoSFilterer) FilterConfirmWithdraw(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*DPoSConfirmWithdrawIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "ConfirmWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSConfirmWithdrawIterator{contract: _DPoS.contract, event: "ConfirmWithdraw", logs: logs, sub: sub}, nil
}

// WatchConfirmWithdraw is a free log subscription operation binding the contract event 0x08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220.
//
// Solidity: event ConfirmWithdraw(address indexed delegator, address indexed candidate, uint256 amount)
func (_DPoS *DPoSFilterer) WatchConfirmWithdraw(opts *bind.WatchOpts, sink chan<- *DPoSConfirmWithdraw, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "ConfirmWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSConfirmWithdraw)
				if err := _DPoS.contract.UnpackLog(event, "ConfirmWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConfirmWithdraw is a log parse operation binding the contract event 0x08d0283ea9a2e520a2f09611cf37ca6eb70f62e9a807e53756047dd2dc027220.
//
// Solidity: event ConfirmWithdraw(address indexed delegator, address indexed candidate, uint256 amount)
func (_DPoS *DPoSFilterer) ParseConfirmWithdraw(log types.Log) (*DPoSConfirmWithdraw, error) {
	event := new(DPoSConfirmWithdraw)
	if err := _DPoS.contract.UnpackLog(event, "ConfirmWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSCreateParamProposalIterator is returned from FilterCreateParamProposal and is used to iterate over the raw logs and unpacked data for CreateParamProposal events raised by the DPoS contract.
type DPoSCreateParamProposalIterator struct {
	Event *DPoSCreateParamProposal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSCreateParamProposalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSCreateParamProposal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSCreateParamProposal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSCreateParamProposalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSCreateParamProposalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSCreateParamProposal represents a CreateParamProposal event raised by the DPoS contract.
type DPoSCreateParamProposal struct {
	ProposalId   *big.Int
	Proposer     common.Address
	Deposit      *big.Int
	VoteDeadline *big.Int
	Record       *big.Int
	NewValue     *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterCreateParamProposal is a free log retrieval operation binding the contract event 0x40109a070319d6004f4e4b31dba4b605c97bd3474d49865158f55fe093e3b339.
//
// Solidity: event CreateParamProposal(uint256 proposalId, address proposer, uint256 deposit, uint256 voteDeadline, uint256 record, uint256 newValue)
func (_DPoS *DPoSFilterer) FilterCreateParamProposal(opts *bind.FilterOpts) (*DPoSCreateParamProposalIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "CreateParamProposal")
	if err != nil {
		return nil, err
	}
	return &DPoSCreateParamProposalIterator{contract: _DPoS.contract, event: "CreateParamProposal", logs: logs, sub: sub}, nil
}

// WatchCreateParamProposal is a free log subscription operation binding the contract event 0x40109a070319d6004f4e4b31dba4b605c97bd3474d49865158f55fe093e3b339.
//
// Solidity: event CreateParamProposal(uint256 proposalId, address proposer, uint256 deposit, uint256 voteDeadline, uint256 record, uint256 newValue)
func (_DPoS *DPoSFilterer) WatchCreateParamProposal(opts *bind.WatchOpts, sink chan<- *DPoSCreateParamProposal) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "CreateParamProposal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSCreateParamProposal)
				if err := _DPoS.contract.UnpackLog(event, "CreateParamProposal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreateParamProposal is a log parse operation binding the contract event 0x40109a070319d6004f4e4b31dba4b605c97bd3474d49865158f55fe093e3b339.
//
// Solidity: event CreateParamProposal(uint256 proposalId, address proposer, uint256 deposit, uint256 voteDeadline, uint256 record, uint256 newValue)
func (_DPoS *DPoSFilterer) ParseCreateParamProposal(log types.Log) (*DPoSCreateParamProposal, error) {
	event := new(DPoSCreateParamProposal)
	if err := _DPoS.contract.UnpackLog(event, "CreateParamProposal", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSCreateSidechainProposalIterator is returned from FilterCreateSidechainProposal and is used to iterate over the raw logs and unpacked data for CreateSidechainProposal events raised by the DPoS contract.
type DPoSCreateSidechainProposalIterator struct {
	Event *DPoSCreateSidechainProposal // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSCreateSidechainProposalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSCreateSidechainProposal)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSCreateSidechainProposal)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSCreateSidechainProposalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSCreateSidechainProposalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSCreateSidechainProposal represents a CreateSidechainProposal event raised by the DPoS contract.
type DPoSCreateSidechainProposal struct {
	ProposalId    *big.Int
	Proposer      common.Address
	Deposit       *big.Int
	VoteDeadline  *big.Int
	SidechainAddr common.Address
	Registered    bool
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCreateSidechainProposal is a free log retrieval operation binding the contract event 0xe6970151d691583ac0aecc2e24c67871318a5c7f7574c6df7929b6dd5d54db68.
//
// Solidity: event CreateSidechainProposal(uint256 proposalId, address proposer, uint256 deposit, uint256 voteDeadline, address sidechainAddr, bool registered)
func (_DPoS *DPoSFilterer) FilterCreateSidechainProposal(opts *bind.FilterOpts) (*DPoSCreateSidechainProposalIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "CreateSidechainProposal")
	if err != nil {
		return nil, err
	}
	return &DPoSCreateSidechainProposalIterator{contract: _DPoS.contract, event: "CreateSidechainProposal", logs: logs, sub: sub}, nil
}

// WatchCreateSidechainProposal is a free log subscription operation binding the contract event 0xe6970151d691583ac0aecc2e24c67871318a5c7f7574c6df7929b6dd5d54db68.
//
// Solidity: event CreateSidechainProposal(uint256 proposalId, address proposer, uint256 deposit, uint256 voteDeadline, address sidechainAddr, bool registered)
func (_DPoS *DPoSFilterer) WatchCreateSidechainProposal(opts *bind.WatchOpts, sink chan<- *DPoSCreateSidechainProposal) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "CreateSidechainProposal")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSCreateSidechainProposal)
				if err := _DPoS.contract.UnpackLog(event, "CreateSidechainProposal", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCreateSidechainProposal is a log parse operation binding the contract event 0xe6970151d691583ac0aecc2e24c67871318a5c7f7574c6df7929b6dd5d54db68.
//
// Solidity: event CreateSidechainProposal(uint256 proposalId, address proposer, uint256 deposit, uint256 voteDeadline, address sidechainAddr, bool registered)
func (_DPoS *DPoSFilterer) ParseCreateSidechainProposal(log types.Log) (*DPoSCreateSidechainProposal, error) {
	event := new(DPoSCreateSidechainProposal)
	if err := _DPoS.contract.UnpackLog(event, "CreateSidechainProposal", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSDelegateIterator is returned from FilterDelegate and is used to iterate over the raw logs and unpacked data for Delegate events raised by the DPoS contract.
type DPoSDelegateIterator struct {
	Event *DPoSDelegate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSDelegateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSDelegate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSDelegate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSDelegateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSDelegateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSDelegate represents a Delegate event raised by the DPoS contract.
type DPoSDelegate struct {
	Delegator   common.Address
	Candidate   common.Address
	NewStake    *big.Int
	StakingPool *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDelegate is a free log retrieval operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 stakingPool)
func (_DPoS *DPoSFilterer) FilterDelegate(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*DPoSDelegateIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "Delegate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSDelegateIterator{contract: _DPoS.contract, event: "Delegate", logs: logs, sub: sub}, nil
}

// WatchDelegate is a free log subscription operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 stakingPool)
func (_DPoS *DPoSFilterer) WatchDelegate(opts *bind.WatchOpts, sink chan<- *DPoSDelegate, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "Delegate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSDelegate)
				if err := _DPoS.contract.UnpackLog(event, "Delegate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDelegate is a log parse operation binding the contract event 0x500599802164a08023e87ffc3eed0ba3ae60697b3083ba81d046683679d81c6b.
//
// Solidity: event Delegate(address indexed delegator, address indexed candidate, uint256 newStake, uint256 stakingPool)
func (_DPoS *DPoSFilterer) ParseDelegate(log types.Log) (*DPoSDelegate, error) {
	event := new(DPoSDelegate)
	if err := _DPoS.contract.UnpackLog(event, "Delegate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSIndemnifyIterator is returned from FilterIndemnify and is used to iterate over the raw logs and unpacked data for Indemnify events raised by the DPoS contract.
type DPoSIndemnifyIterator struct {
	Event *DPoSIndemnify // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSIndemnifyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSIndemnify)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSIndemnify)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSIndemnifyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSIndemnifyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSIndemnify represents a Indemnify event raised by the DPoS contract.
type DPoSIndemnify struct {
	Indemnitee common.Address
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterIndemnify is a free log retrieval operation binding the contract event 0x5739ab1e78cc400665fe049852d266e9ac5995451b5e666b25f3c30b19087d8e.
//
// Solidity: event Indemnify(address indexed indemnitee, uint256 amount)
func (_DPoS *DPoSFilterer) FilterIndemnify(opts *bind.FilterOpts, indemnitee []common.Address) (*DPoSIndemnifyIterator, error) {

	var indemniteeRule []interface{}
	for _, indemniteeItem := range indemnitee {
		indemniteeRule = append(indemniteeRule, indemniteeItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "Indemnify", indemniteeRule)
	if err != nil {
		return nil, err
	}
	return &DPoSIndemnifyIterator{contract: _DPoS.contract, event: "Indemnify", logs: logs, sub: sub}, nil
}

// WatchIndemnify is a free log subscription operation binding the contract event 0x5739ab1e78cc400665fe049852d266e9ac5995451b5e666b25f3c30b19087d8e.
//
// Solidity: event Indemnify(address indexed indemnitee, uint256 amount)
func (_DPoS *DPoSFilterer) WatchIndemnify(opts *bind.WatchOpts, sink chan<- *DPoSIndemnify, indemnitee []common.Address) (event.Subscription, error) {

	var indemniteeRule []interface{}
	for _, indemniteeItem := range indemnitee {
		indemniteeRule = append(indemniteeRule, indemniteeItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "Indemnify", indemniteeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSIndemnify)
				if err := _DPoS.contract.UnpackLog(event, "Indemnify", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIndemnify is a log parse operation binding the contract event 0x5739ab1e78cc400665fe049852d266e9ac5995451b5e666b25f3c30b19087d8e.
//
// Solidity: event Indemnify(address indexed indemnitee, uint256 amount)
func (_DPoS *DPoSFilterer) ParseIndemnify(log types.Log) (*DPoSIndemnify, error) {
	event := new(DPoSIndemnify)
	if err := _DPoS.contract.UnpackLog(event, "Indemnify", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSInitializeCandidateIterator is returned from FilterInitializeCandidate and is used to iterate over the raw logs and unpacked data for InitializeCandidate events raised by the DPoS contract.
type DPoSInitializeCandidateIterator struct {
	Event *DPoSInitializeCandidate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSInitializeCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSInitializeCandidate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSInitializeCandidate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSInitializeCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSInitializeCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSInitializeCandidate represents a InitializeCandidate event raised by the DPoS contract.
type DPoSInitializeCandidate struct {
	Candidate       common.Address
	MinSelfStake    *big.Int
	CommissionRate  *big.Int
	RateLockEndTime *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterInitializeCandidate is a free log retrieval operation binding the contract event 0x453d56a841836718d9e848e968068cbc2af21ca29d1527fbebd231dc46ceffaa.
//
// Solidity: event InitializeCandidate(address indexed candidate, uint256 minSelfStake, uint256 commissionRate, uint256 rateLockEndTime)
func (_DPoS *DPoSFilterer) FilterInitializeCandidate(opts *bind.FilterOpts, candidate []common.Address) (*DPoSInitializeCandidateIterator, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "InitializeCandidate", candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSInitializeCandidateIterator{contract: _DPoS.contract, event: "InitializeCandidate", logs: logs, sub: sub}, nil
}

// WatchInitializeCandidate is a free log subscription operation binding the contract event 0x453d56a841836718d9e848e968068cbc2af21ca29d1527fbebd231dc46ceffaa.
//
// Solidity: event InitializeCandidate(address indexed candidate, uint256 minSelfStake, uint256 commissionRate, uint256 rateLockEndTime)
func (_DPoS *DPoSFilterer) WatchInitializeCandidate(opts *bind.WatchOpts, sink chan<- *DPoSInitializeCandidate, candidate []common.Address) (event.Subscription, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "InitializeCandidate", candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSInitializeCandidate)
				if err := _DPoS.contract.UnpackLog(event, "InitializeCandidate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitializeCandidate is a log parse operation binding the contract event 0x453d56a841836718d9e848e968068cbc2af21ca29d1527fbebd231dc46ceffaa.
//
// Solidity: event InitializeCandidate(address indexed candidate, uint256 minSelfStake, uint256 commissionRate, uint256 rateLockEndTime)
func (_DPoS *DPoSFilterer) ParseInitializeCandidate(log types.Log) (*DPoSInitializeCandidate, error) {
	event := new(DPoSInitializeCandidate)
	if err := _DPoS.contract.UnpackLog(event, "InitializeCandidate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSIntendWithdrawIterator is returned from FilterIntendWithdraw and is used to iterate over the raw logs and unpacked data for IntendWithdraw events raised by the DPoS contract.
type DPoSIntendWithdrawIterator struct {
	Event *DPoSIntendWithdraw // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSIntendWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSIntendWithdraw)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSIntendWithdraw)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSIntendWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSIntendWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSIntendWithdraw represents a IntendWithdraw event raised by the DPoS contract.
type DPoSIntendWithdraw struct {
	Delegator      common.Address
	Candidate      common.Address
	WithdrawAmount *big.Int
	ProposedTime   *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterIntendWithdraw is a free log retrieval operation binding the contract event 0x7171946bb2a9ef55fcb2eb8cef679db45e2e3a8cef9a44567d34d202b65ff0b1.
//
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 proposedTime)
func (_DPoS *DPoSFilterer) FilterIntendWithdraw(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*DPoSIntendWithdrawIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "IntendWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSIntendWithdrawIterator{contract: _DPoS.contract, event: "IntendWithdraw", logs: logs, sub: sub}, nil
}

// WatchIntendWithdraw is a free log subscription operation binding the contract event 0x7171946bb2a9ef55fcb2eb8cef679db45e2e3a8cef9a44567d34d202b65ff0b1.
//
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 proposedTime)
func (_DPoS *DPoSFilterer) WatchIntendWithdraw(opts *bind.WatchOpts, sink chan<- *DPoSIntendWithdraw, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "IntendWithdraw", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSIntendWithdraw)
				if err := _DPoS.contract.UnpackLog(event, "IntendWithdraw", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseIntendWithdraw is a log parse operation binding the contract event 0x7171946bb2a9ef55fcb2eb8cef679db45e2e3a8cef9a44567d34d202b65ff0b1.
//
// Solidity: event IntendWithdraw(address indexed delegator, address indexed candidate, uint256 withdrawAmount, uint256 proposedTime)
func (_DPoS *DPoSFilterer) ParseIntendWithdraw(log types.Log) (*DPoSIntendWithdraw, error) {
	event := new(DPoSIntendWithdraw)
	if err := _DPoS.contract.UnpackLog(event, "IntendWithdraw", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSMiningPoolContributionIterator is returned from FilterMiningPoolContribution and is used to iterate over the raw logs and unpacked data for MiningPoolContribution events raised by the DPoS contract.
type DPoSMiningPoolContributionIterator struct {
	Event *DPoSMiningPoolContribution // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSMiningPoolContributionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSMiningPoolContribution)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSMiningPoolContribution)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSMiningPoolContributionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSMiningPoolContributionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSMiningPoolContribution represents a MiningPoolContribution event raised by the DPoS contract.
type DPoSMiningPoolContribution struct {
	Contributor    common.Address
	Contribution   *big.Int
	MiningPoolSize *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterMiningPoolContribution is a free log retrieval operation binding the contract event 0x97e19c4040b6c46d4275e0c4fea68f8f92c81138372ffdb089932c211938f765.
//
// Solidity: event MiningPoolContribution(address indexed contributor, uint256 contribution, uint256 miningPoolSize)
func (_DPoS *DPoSFilterer) FilterMiningPoolContribution(opts *bind.FilterOpts, contributor []common.Address) (*DPoSMiningPoolContributionIterator, error) {

	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "MiningPoolContribution", contributorRule)
	if err != nil {
		return nil, err
	}
	return &DPoSMiningPoolContributionIterator{contract: _DPoS.contract, event: "MiningPoolContribution", logs: logs, sub: sub}, nil
}

// WatchMiningPoolContribution is a free log subscription operation binding the contract event 0x97e19c4040b6c46d4275e0c4fea68f8f92c81138372ffdb089932c211938f765.
//
// Solidity: event MiningPoolContribution(address indexed contributor, uint256 contribution, uint256 miningPoolSize)
func (_DPoS *DPoSFilterer) WatchMiningPoolContribution(opts *bind.WatchOpts, sink chan<- *DPoSMiningPoolContribution, contributor []common.Address) (event.Subscription, error) {

	var contributorRule []interface{}
	for _, contributorItem := range contributor {
		contributorRule = append(contributorRule, contributorItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "MiningPoolContribution", contributorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSMiningPoolContribution)
				if err := _DPoS.contract.UnpackLog(event, "MiningPoolContribution", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMiningPoolContribution is a log parse operation binding the contract event 0x97e19c4040b6c46d4275e0c4fea68f8f92c81138372ffdb089932c211938f765.
//
// Solidity: event MiningPoolContribution(address indexed contributor, uint256 contribution, uint256 miningPoolSize)
func (_DPoS *DPoSFilterer) ParseMiningPoolContribution(log types.Log) (*DPoSMiningPoolContribution, error) {
	event := new(DPoSMiningPoolContribution)
	if err := _DPoS.contract.UnpackLog(event, "MiningPoolContribution", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the DPoS contract.
type DPoSOwnershipTransferredIterator struct {
	Event *DPoSOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSOwnershipTransferred represents a OwnershipTransferred event raised by the DPoS contract.
type DPoSOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DPoS *DPoSFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DPoSOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DPoSOwnershipTransferredIterator{contract: _DPoS.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DPoS *DPoSFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DPoSOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSOwnershipTransferred)
				if err := _DPoS.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_DPoS *DPoSFilterer) ParseOwnershipTransferred(log types.Log) (*DPoSOwnershipTransferred, error) {
	event := new(DPoSOwnershipTransferred)
	if err := _DPoS.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSPausedIterator is returned from FilterPaused and is used to iterate over the raw logs and unpacked data for Paused events raised by the DPoS contract.
type DPoSPausedIterator struct {
	Event *DPoSPaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSPausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSPaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSPaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSPausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSPausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSPaused represents a Paused event raised by the DPoS contract.
type DPoSPaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPaused is a free log retrieval operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DPoS *DPoSFilterer) FilterPaused(opts *bind.FilterOpts) (*DPoSPausedIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return &DPoSPausedIterator{contract: _DPoS.contract, event: "Paused", logs: logs, sub: sub}, nil
}

// WatchPaused is a free log subscription operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DPoS *DPoSFilterer) WatchPaused(opts *bind.WatchOpts, sink chan<- *DPoSPaused) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "Paused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSPaused)
				if err := _DPoS.contract.UnpackLog(event, "Paused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePaused is a log parse operation binding the contract event 0x62e78cea01bee320cd4e420270b5ea74000d11b0c9f74754ebdbfc544b05a258.
//
// Solidity: event Paused(address account)
func (_DPoS *DPoSFilterer) ParsePaused(log types.Log) (*DPoSPaused, error) {
	event := new(DPoSPaused)
	if err := _DPoS.contract.UnpackLog(event, "Paused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSPauserAddedIterator is returned from FilterPauserAdded and is used to iterate over the raw logs and unpacked data for PauserAdded events raised by the DPoS contract.
type DPoSPauserAddedIterator struct {
	Event *DPoSPauserAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSPauserAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSPauserAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSPauserAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSPauserAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSPauserAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSPauserAdded represents a PauserAdded event raised by the DPoS contract.
type DPoSPauserAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPauserAdded is a free log retrieval operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: event PauserAdded(address indexed account)
func (_DPoS *DPoSFilterer) FilterPauserAdded(opts *bind.FilterOpts, account []common.Address) (*DPoSPauserAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "PauserAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &DPoSPauserAddedIterator{contract: _DPoS.contract, event: "PauserAdded", logs: logs, sub: sub}, nil
}

// WatchPauserAdded is a free log subscription operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: event PauserAdded(address indexed account)
func (_DPoS *DPoSFilterer) WatchPauserAdded(opts *bind.WatchOpts, sink chan<- *DPoSPauserAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "PauserAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSPauserAdded)
				if err := _DPoS.contract.UnpackLog(event, "PauserAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePauserAdded is a log parse operation binding the contract event 0x6719d08c1888103bea251a4ed56406bd0c3e69723c8a1686e017e7bbe159b6f8.
//
// Solidity: event PauserAdded(address indexed account)
func (_DPoS *DPoSFilterer) ParsePauserAdded(log types.Log) (*DPoSPauserAdded, error) {
	event := new(DPoSPauserAdded)
	if err := _DPoS.contract.UnpackLog(event, "PauserAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSPauserRemovedIterator is returned from FilterPauserRemoved and is used to iterate over the raw logs and unpacked data for PauserRemoved events raised by the DPoS contract.
type DPoSPauserRemovedIterator struct {
	Event *DPoSPauserRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSPauserRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSPauserRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSPauserRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSPauserRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSPauserRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSPauserRemoved represents a PauserRemoved event raised by the DPoS contract.
type DPoSPauserRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterPauserRemoved is a free log retrieval operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: event PauserRemoved(address indexed account)
func (_DPoS *DPoSFilterer) FilterPauserRemoved(opts *bind.FilterOpts, account []common.Address) (*DPoSPauserRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "PauserRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &DPoSPauserRemovedIterator{contract: _DPoS.contract, event: "PauserRemoved", logs: logs, sub: sub}, nil
}

// WatchPauserRemoved is a free log subscription operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: event PauserRemoved(address indexed account)
func (_DPoS *DPoSFilterer) WatchPauserRemoved(opts *bind.WatchOpts, sink chan<- *DPoSPauserRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "PauserRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSPauserRemoved)
				if err := _DPoS.contract.UnpackLog(event, "PauserRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePauserRemoved is a log parse operation binding the contract event 0xcd265ebaf09df2871cc7bd4133404a235ba12eff2041bb89d9c714a2621c7c7e.
//
// Solidity: event PauserRemoved(address indexed account)
func (_DPoS *DPoSFilterer) ParsePauserRemoved(log types.Log) (*DPoSPauserRemoved, error) {
	event := new(DPoSPauserRemoved)
	if err := _DPoS.contract.UnpackLog(event, "PauserRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSPunishIterator is returned from FilterPunish and is used to iterate over the raw logs and unpacked data for Punish events raised by the DPoS contract.
type DPoSPunishIterator struct {
	Event *DPoSPunish // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSPunishIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSPunish)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSPunish)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSPunishIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSPunishIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSPunish represents a Punish event raised by the DPoS contract.
type DPoSPunish struct {
	Validator common.Address
	Delegator common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterPunish is a free log retrieval operation binding the contract event 0x111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a.
//
// Solidity: event Punish(address indexed validator, address indexed delegator, uint256 amount)
func (_DPoS *DPoSFilterer) FilterPunish(opts *bind.FilterOpts, validator []common.Address, delegator []common.Address) (*DPoSPunishIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "Punish", validatorRule, delegatorRule)
	if err != nil {
		return nil, err
	}
	return &DPoSPunishIterator{contract: _DPoS.contract, event: "Punish", logs: logs, sub: sub}, nil
}

// WatchPunish is a free log subscription operation binding the contract event 0x111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a.
//
// Solidity: event Punish(address indexed validator, address indexed delegator, uint256 amount)
func (_DPoS *DPoSFilterer) WatchPunish(opts *bind.WatchOpts, sink chan<- *DPoSPunish, validator []common.Address, delegator []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "Punish", validatorRule, delegatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSPunish)
				if err := _DPoS.contract.UnpackLog(event, "Punish", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePunish is a log parse operation binding the contract event 0x111897aba775ed4cb659e35805c453dcd8f0024cc414f560f9677bdfae30952a.
//
// Solidity: event Punish(address indexed validator, address indexed delegator, uint256 amount)
func (_DPoS *DPoSFilterer) ParsePunish(log types.Log) (*DPoSPunish, error) {
	event := new(DPoSPunish)
	if err := _DPoS.contract.UnpackLog(event, "Punish", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSRedeemMiningRewardIterator is returned from FilterRedeemMiningReward and is used to iterate over the raw logs and unpacked data for RedeemMiningReward events raised by the DPoS contract.
type DPoSRedeemMiningRewardIterator struct {
	Event *DPoSRedeemMiningReward // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSRedeemMiningRewardIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSRedeemMiningReward)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSRedeemMiningReward)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSRedeemMiningRewardIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSRedeemMiningRewardIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSRedeemMiningReward represents a RedeemMiningReward event raised by the DPoS contract.
type DPoSRedeemMiningReward struct {
	Receiver   common.Address
	Reward     *big.Int
	MiningPool *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterRedeemMiningReward is a free log retrieval operation binding the contract event 0xc243dafa8ee55923dad771198c225cf6dfcdc5e405eda7d4da42b6c6fa018de7.
//
// Solidity: event RedeemMiningReward(address indexed receiver, uint256 reward, uint256 miningPool)
func (_DPoS *DPoSFilterer) FilterRedeemMiningReward(opts *bind.FilterOpts, receiver []common.Address) (*DPoSRedeemMiningRewardIterator, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "RedeemMiningReward", receiverRule)
	if err != nil {
		return nil, err
	}
	return &DPoSRedeemMiningRewardIterator{contract: _DPoS.contract, event: "RedeemMiningReward", logs: logs, sub: sub}, nil
}

// WatchRedeemMiningReward is a free log subscription operation binding the contract event 0xc243dafa8ee55923dad771198c225cf6dfcdc5e405eda7d4da42b6c6fa018de7.
//
// Solidity: event RedeemMiningReward(address indexed receiver, uint256 reward, uint256 miningPool)
func (_DPoS *DPoSFilterer) WatchRedeemMiningReward(opts *bind.WatchOpts, sink chan<- *DPoSRedeemMiningReward, receiver []common.Address) (event.Subscription, error) {

	var receiverRule []interface{}
	for _, receiverItem := range receiver {
		receiverRule = append(receiverRule, receiverItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "RedeemMiningReward", receiverRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSRedeemMiningReward)
				if err := _DPoS.contract.UnpackLog(event, "RedeemMiningReward", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRedeemMiningReward is a log parse operation binding the contract event 0xc243dafa8ee55923dad771198c225cf6dfcdc5e405eda7d4da42b6c6fa018de7.
//
// Solidity: event RedeemMiningReward(address indexed receiver, uint256 reward, uint256 miningPool)
func (_DPoS *DPoSFilterer) ParseRedeemMiningReward(log types.Log) (*DPoSRedeemMiningReward, error) {
	event := new(DPoSRedeemMiningReward)
	if err := _DPoS.contract.UnpackLog(event, "RedeemMiningReward", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSUnpausedIterator is returned from FilterUnpaused and is used to iterate over the raw logs and unpacked data for Unpaused events raised by the DPoS contract.
type DPoSUnpausedIterator struct {
	Event *DPoSUnpaused // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSUnpausedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSUnpaused)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSUnpaused)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSUnpausedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSUnpausedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSUnpaused represents a Unpaused event raised by the DPoS contract.
type DPoSUnpaused struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterUnpaused is a free log retrieval operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DPoS *DPoSFilterer) FilterUnpaused(opts *bind.FilterOpts) (*DPoSUnpausedIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return &DPoSUnpausedIterator{contract: _DPoS.contract, event: "Unpaused", logs: logs, sub: sub}, nil
}

// WatchUnpaused is a free log subscription operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DPoS *DPoSFilterer) WatchUnpaused(opts *bind.WatchOpts, sink chan<- *DPoSUnpaused) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "Unpaused")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSUnpaused)
				if err := _DPoS.contract.UnpackLog(event, "Unpaused", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUnpaused is a log parse operation binding the contract event 0x5db9ee0a495bf2e6ff9c91a7834c1ba4fdd244a5e8aa4e537bd38aeae4b073aa.
//
// Solidity: event Unpaused(address account)
func (_DPoS *DPoSFilterer) ParseUnpaused(log types.Log) (*DPoSUnpaused, error) {
	event := new(DPoSUnpaused)
	if err := _DPoS.contract.UnpackLog(event, "Unpaused", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSUpdateCommissionRateIterator is returned from FilterUpdateCommissionRate and is used to iterate over the raw logs and unpacked data for UpdateCommissionRate events raised by the DPoS contract.
type DPoSUpdateCommissionRateIterator struct {
	Event *DPoSUpdateCommissionRate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSUpdateCommissionRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSUpdateCommissionRate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSUpdateCommissionRate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSUpdateCommissionRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSUpdateCommissionRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSUpdateCommissionRate represents a UpdateCommissionRate event raised by the DPoS contract.
type DPoSUpdateCommissionRate struct {
	Candidate      common.Address
	NewRate        *big.Int
	NewLockEndTime *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommissionRate is a free log retrieval operation binding the contract event 0x37954fc2aa8b4424ad16c75da2ea4d51ba08ef9e07907e37ccae54a0b4ce1e9e.
//
// Solidity: event UpdateCommissionRate(address indexed candidate, uint256 newRate, uint256 newLockEndTime)
func (_DPoS *DPoSFilterer) FilterUpdateCommissionRate(opts *bind.FilterOpts, candidate []common.Address) (*DPoSUpdateCommissionRateIterator, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "UpdateCommissionRate", candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSUpdateCommissionRateIterator{contract: _DPoS.contract, event: "UpdateCommissionRate", logs: logs, sub: sub}, nil
}

// WatchUpdateCommissionRate is a free log subscription operation binding the contract event 0x37954fc2aa8b4424ad16c75da2ea4d51ba08ef9e07907e37ccae54a0b4ce1e9e.
//
// Solidity: event UpdateCommissionRate(address indexed candidate, uint256 newRate, uint256 newLockEndTime)
func (_DPoS *DPoSFilterer) WatchUpdateCommissionRate(opts *bind.WatchOpts, sink chan<- *DPoSUpdateCommissionRate, candidate []common.Address) (event.Subscription, error) {

	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "UpdateCommissionRate", candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSUpdateCommissionRate)
				if err := _DPoS.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpdateCommissionRate is a log parse operation binding the contract event 0x37954fc2aa8b4424ad16c75da2ea4d51ba08ef9e07907e37ccae54a0b4ce1e9e.
//
// Solidity: event UpdateCommissionRate(address indexed candidate, uint256 newRate, uint256 newLockEndTime)
func (_DPoS *DPoSFilterer) ParseUpdateCommissionRate(log types.Log) (*DPoSUpdateCommissionRate, error) {
	event := new(DPoSUpdateCommissionRate)
	if err := _DPoS.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSValidatorChangeIterator is returned from FilterValidatorChange and is used to iterate over the raw logs and unpacked data for ValidatorChange events raised by the DPoS contract.
type DPoSValidatorChangeIterator struct {
	Event *DPoSValidatorChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSValidatorChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSValidatorChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSValidatorChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSValidatorChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSValidatorChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSValidatorChange represents a ValidatorChange event raised by the DPoS contract.
type DPoSValidatorChange struct {
	EthAddr    common.Address
	ChangeType uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterValidatorChange is a free log retrieval operation binding the contract event 0x63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c.
//
// Solidity: event ValidatorChange(address indexed ethAddr, uint8 indexed changeType)
func (_DPoS *DPoSFilterer) FilterValidatorChange(opts *bind.FilterOpts, ethAddr []common.Address, changeType []uint8) (*DPoSValidatorChangeIterator, error) {

	var ethAddrRule []interface{}
	for _, ethAddrItem := range ethAddr {
		ethAddrRule = append(ethAddrRule, ethAddrItem)
	}
	var changeTypeRule []interface{}
	for _, changeTypeItem := range changeType {
		changeTypeRule = append(changeTypeRule, changeTypeItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "ValidatorChange", ethAddrRule, changeTypeRule)
	if err != nil {
		return nil, err
	}
	return &DPoSValidatorChangeIterator{contract: _DPoS.contract, event: "ValidatorChange", logs: logs, sub: sub}, nil
}

// WatchValidatorChange is a free log subscription operation binding the contract event 0x63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c.
//
// Solidity: event ValidatorChange(address indexed ethAddr, uint8 indexed changeType)
func (_DPoS *DPoSFilterer) WatchValidatorChange(opts *bind.WatchOpts, sink chan<- *DPoSValidatorChange, ethAddr []common.Address, changeType []uint8) (event.Subscription, error) {

	var ethAddrRule []interface{}
	for _, ethAddrItem := range ethAddr {
		ethAddrRule = append(ethAddrRule, ethAddrItem)
	}
	var changeTypeRule []interface{}
	for _, changeTypeItem := range changeType {
		changeTypeRule = append(changeTypeRule, changeTypeItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "ValidatorChange", ethAddrRule, changeTypeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSValidatorChange)
				if err := _DPoS.contract.UnpackLog(event, "ValidatorChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseValidatorChange is a log parse operation binding the contract event 0x63f783ba869265648de5e70add96be9f4914e3bde064fdc19fd7e6a8ebf2f46c.
//
// Solidity: event ValidatorChange(address indexed ethAddr, uint8 indexed changeType)
func (_DPoS *DPoSFilterer) ParseValidatorChange(log types.Log) (*DPoSValidatorChange, error) {
	event := new(DPoSValidatorChange)
	if err := _DPoS.contract.UnpackLog(event, "ValidatorChange", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSVoteParamIterator is returned from FilterVoteParam and is used to iterate over the raw logs and unpacked data for VoteParam events raised by the DPoS contract.
type DPoSVoteParamIterator struct {
	Event *DPoSVoteParam // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSVoteParamIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSVoteParam)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSVoteParam)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSVoteParamIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSVoteParamIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSVoteParam represents a VoteParam event raised by the DPoS contract.
type DPoSVoteParam struct {
	ProposalId *big.Int
	Voter      common.Address
	VoteType   uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteParam is a free log retrieval operation binding the contract event 0x06c7ef6e19454637e93ee60cc680c61fb2ebabb57e58cf36d94141a5036b3d65.
//
// Solidity: event VoteParam(uint256 proposalId, address voter, uint8 voteType)
func (_DPoS *DPoSFilterer) FilterVoteParam(opts *bind.FilterOpts) (*DPoSVoteParamIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "VoteParam")
	if err != nil {
		return nil, err
	}
	return &DPoSVoteParamIterator{contract: _DPoS.contract, event: "VoteParam", logs: logs, sub: sub}, nil
}

// WatchVoteParam is a free log subscription operation binding the contract event 0x06c7ef6e19454637e93ee60cc680c61fb2ebabb57e58cf36d94141a5036b3d65.
//
// Solidity: event VoteParam(uint256 proposalId, address voter, uint8 voteType)
func (_DPoS *DPoSFilterer) WatchVoteParam(opts *bind.WatchOpts, sink chan<- *DPoSVoteParam) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "VoteParam")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSVoteParam)
				if err := _DPoS.contract.UnpackLog(event, "VoteParam", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoteParam is a log parse operation binding the contract event 0x06c7ef6e19454637e93ee60cc680c61fb2ebabb57e58cf36d94141a5036b3d65.
//
// Solidity: event VoteParam(uint256 proposalId, address voter, uint8 voteType)
func (_DPoS *DPoSFilterer) ParseVoteParam(log types.Log) (*DPoSVoteParam, error) {
	event := new(DPoSVoteParam)
	if err := _DPoS.contract.UnpackLog(event, "VoteParam", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSVoteSidechainIterator is returned from FilterVoteSidechain and is used to iterate over the raw logs and unpacked data for VoteSidechain events raised by the DPoS contract.
type DPoSVoteSidechainIterator struct {
	Event *DPoSVoteSidechain // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSVoteSidechainIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSVoteSidechain)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSVoteSidechain)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSVoteSidechainIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSVoteSidechainIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSVoteSidechain represents a VoteSidechain event raised by the DPoS contract.
type DPoSVoteSidechain struct {
	ProposalId *big.Int
	Voter      common.Address
	VoteType   uint8
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVoteSidechain is a free log retrieval operation binding the contract event 0x7686976924e1fdb79b36f7445ada20b6e9d3377d85b34d5162116e675c39d34c.
//
// Solidity: event VoteSidechain(uint256 proposalId, address voter, uint8 voteType)
func (_DPoS *DPoSFilterer) FilterVoteSidechain(opts *bind.FilterOpts) (*DPoSVoteSidechainIterator, error) {

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "VoteSidechain")
	if err != nil {
		return nil, err
	}
	return &DPoSVoteSidechainIterator{contract: _DPoS.contract, event: "VoteSidechain", logs: logs, sub: sub}, nil
}

// WatchVoteSidechain is a free log subscription operation binding the contract event 0x7686976924e1fdb79b36f7445ada20b6e9d3377d85b34d5162116e675c39d34c.
//
// Solidity: event VoteSidechain(uint256 proposalId, address voter, uint8 voteType)
func (_DPoS *DPoSFilterer) WatchVoteSidechain(opts *bind.WatchOpts, sink chan<- *DPoSVoteSidechain) (event.Subscription, error) {

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "VoteSidechain")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSVoteSidechain)
				if err := _DPoS.contract.UnpackLog(event, "VoteSidechain", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseVoteSidechain is a log parse operation binding the contract event 0x7686976924e1fdb79b36f7445ada20b6e9d3377d85b34d5162116e675c39d34c.
//
// Solidity: event VoteSidechain(uint256 proposalId, address voter, uint8 voteType)
func (_DPoS *DPoSFilterer) ParseVoteSidechain(log types.Log) (*DPoSVoteSidechain, error) {
	event := new(DPoSVoteSidechain)
	if err := _DPoS.contract.UnpackLog(event, "VoteSidechain", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSWhitelistAdminAddedIterator is returned from FilterWhitelistAdminAdded and is used to iterate over the raw logs and unpacked data for WhitelistAdminAdded events raised by the DPoS contract.
type DPoSWhitelistAdminAddedIterator struct {
	Event *DPoSWhitelistAdminAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSWhitelistAdminAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSWhitelistAdminAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSWhitelistAdminAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSWhitelistAdminAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSWhitelistAdminAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSWhitelistAdminAdded represents a WhitelistAdminAdded event raised by the DPoS contract.
type DPoSWhitelistAdminAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWhitelistAdminAdded is a free log retrieval operation binding the contract event 0x22380c05984257a1cb900161c713dd71d39e74820f1aea43bd3f1bdd20961299.
//
// Solidity: event WhitelistAdminAdded(address indexed account)
func (_DPoS *DPoSFilterer) FilterWhitelistAdminAdded(opts *bind.FilterOpts, account []common.Address) (*DPoSWhitelistAdminAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "WhitelistAdminAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &DPoSWhitelistAdminAddedIterator{contract: _DPoS.contract, event: "WhitelistAdminAdded", logs: logs, sub: sub}, nil
}

// WatchWhitelistAdminAdded is a free log subscription operation binding the contract event 0x22380c05984257a1cb900161c713dd71d39e74820f1aea43bd3f1bdd20961299.
//
// Solidity: event WhitelistAdminAdded(address indexed account)
func (_DPoS *DPoSFilterer) WatchWhitelistAdminAdded(opts *bind.WatchOpts, sink chan<- *DPoSWhitelistAdminAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "WhitelistAdminAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSWhitelistAdminAdded)
				if err := _DPoS.contract.UnpackLog(event, "WhitelistAdminAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWhitelistAdminAdded is a log parse operation binding the contract event 0x22380c05984257a1cb900161c713dd71d39e74820f1aea43bd3f1bdd20961299.
//
// Solidity: event WhitelistAdminAdded(address indexed account)
func (_DPoS *DPoSFilterer) ParseWhitelistAdminAdded(log types.Log) (*DPoSWhitelistAdminAdded, error) {
	event := new(DPoSWhitelistAdminAdded)
	if err := _DPoS.contract.UnpackLog(event, "WhitelistAdminAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSWhitelistAdminRemovedIterator is returned from FilterWhitelistAdminRemoved and is used to iterate over the raw logs and unpacked data for WhitelistAdminRemoved events raised by the DPoS contract.
type DPoSWhitelistAdminRemovedIterator struct {
	Event *DPoSWhitelistAdminRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSWhitelistAdminRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSWhitelistAdminRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSWhitelistAdminRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSWhitelistAdminRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSWhitelistAdminRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSWhitelistAdminRemoved represents a WhitelistAdminRemoved event raised by the DPoS contract.
type DPoSWhitelistAdminRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWhitelistAdminRemoved is a free log retrieval operation binding the contract event 0x0a8eb35e5ca14b3d6f28e4abf2f128dbab231a58b56e89beb5d636115001e165.
//
// Solidity: event WhitelistAdminRemoved(address indexed account)
func (_DPoS *DPoSFilterer) FilterWhitelistAdminRemoved(opts *bind.FilterOpts, account []common.Address) (*DPoSWhitelistAdminRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "WhitelistAdminRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &DPoSWhitelistAdminRemovedIterator{contract: _DPoS.contract, event: "WhitelistAdminRemoved", logs: logs, sub: sub}, nil
}

// WatchWhitelistAdminRemoved is a free log subscription operation binding the contract event 0x0a8eb35e5ca14b3d6f28e4abf2f128dbab231a58b56e89beb5d636115001e165.
//
// Solidity: event WhitelistAdminRemoved(address indexed account)
func (_DPoS *DPoSFilterer) WatchWhitelistAdminRemoved(opts *bind.WatchOpts, sink chan<- *DPoSWhitelistAdminRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "WhitelistAdminRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSWhitelistAdminRemoved)
				if err := _DPoS.contract.UnpackLog(event, "WhitelistAdminRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWhitelistAdminRemoved is a log parse operation binding the contract event 0x0a8eb35e5ca14b3d6f28e4abf2f128dbab231a58b56e89beb5d636115001e165.
//
// Solidity: event WhitelistAdminRemoved(address indexed account)
func (_DPoS *DPoSFilterer) ParseWhitelistAdminRemoved(log types.Log) (*DPoSWhitelistAdminRemoved, error) {
	event := new(DPoSWhitelistAdminRemoved)
	if err := _DPoS.contract.UnpackLog(event, "WhitelistAdminRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSWhitelistedAddedIterator is returned from FilterWhitelistedAdded and is used to iterate over the raw logs and unpacked data for WhitelistedAdded events raised by the DPoS contract.
type DPoSWhitelistedAddedIterator struct {
	Event *DPoSWhitelistedAdded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSWhitelistedAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSWhitelistedAdded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSWhitelistedAdded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSWhitelistedAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSWhitelistedAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSWhitelistedAdded represents a WhitelistedAdded event raised by the DPoS contract.
type DPoSWhitelistedAdded struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWhitelistedAdded is a free log retrieval operation binding the contract event 0xee1504a83b6d4a361f4c1dc78ab59bfa30d6a3b6612c403e86bb01ef2984295f.
//
// Solidity: event WhitelistedAdded(address indexed account)
func (_DPoS *DPoSFilterer) FilterWhitelistedAdded(opts *bind.FilterOpts, account []common.Address) (*DPoSWhitelistedAddedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "WhitelistedAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return &DPoSWhitelistedAddedIterator{contract: _DPoS.contract, event: "WhitelistedAdded", logs: logs, sub: sub}, nil
}

// WatchWhitelistedAdded is a free log subscription operation binding the contract event 0xee1504a83b6d4a361f4c1dc78ab59bfa30d6a3b6612c403e86bb01ef2984295f.
//
// Solidity: event WhitelistedAdded(address indexed account)
func (_DPoS *DPoSFilterer) WatchWhitelistedAdded(opts *bind.WatchOpts, sink chan<- *DPoSWhitelistedAdded, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "WhitelistedAdded", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSWhitelistedAdded)
				if err := _DPoS.contract.UnpackLog(event, "WhitelistedAdded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWhitelistedAdded is a log parse operation binding the contract event 0xee1504a83b6d4a361f4c1dc78ab59bfa30d6a3b6612c403e86bb01ef2984295f.
//
// Solidity: event WhitelistedAdded(address indexed account)
func (_DPoS *DPoSFilterer) ParseWhitelistedAdded(log types.Log) (*DPoSWhitelistedAdded, error) {
	event := new(DPoSWhitelistedAdded)
	if err := _DPoS.contract.UnpackLog(event, "WhitelistedAdded", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSWhitelistedRemovedIterator is returned from FilterWhitelistedRemoved and is used to iterate over the raw logs and unpacked data for WhitelistedRemoved events raised by the DPoS contract.
type DPoSWhitelistedRemovedIterator struct {
	Event *DPoSWhitelistedRemoved // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSWhitelistedRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSWhitelistedRemoved)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSWhitelistedRemoved)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSWhitelistedRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSWhitelistedRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSWhitelistedRemoved represents a WhitelistedRemoved event raised by the DPoS contract.
type DPoSWhitelistedRemoved struct {
	Account common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterWhitelistedRemoved is a free log retrieval operation binding the contract event 0x270d9b30cf5b0793bbfd54c9d5b94aeb49462b8148399000265144a8722da6b6.
//
// Solidity: event WhitelistedRemoved(address indexed account)
func (_DPoS *DPoSFilterer) FilterWhitelistedRemoved(opts *bind.FilterOpts, account []common.Address) (*DPoSWhitelistedRemovedIterator, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "WhitelistedRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return &DPoSWhitelistedRemovedIterator{contract: _DPoS.contract, event: "WhitelistedRemoved", logs: logs, sub: sub}, nil
}

// WatchWhitelistedRemoved is a free log subscription operation binding the contract event 0x270d9b30cf5b0793bbfd54c9d5b94aeb49462b8148399000265144a8722da6b6.
//
// Solidity: event WhitelistedRemoved(address indexed account)
func (_DPoS *DPoSFilterer) WatchWhitelistedRemoved(opts *bind.WatchOpts, sink chan<- *DPoSWhitelistedRemoved, account []common.Address) (event.Subscription, error) {

	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "WhitelistedRemoved", accountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSWhitelistedRemoved)
				if err := _DPoS.contract.UnpackLog(event, "WhitelistedRemoved", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWhitelistedRemoved is a log parse operation binding the contract event 0x270d9b30cf5b0793bbfd54c9d5b94aeb49462b8148399000265144a8722da6b6.
//
// Solidity: event WhitelistedRemoved(address indexed account)
func (_DPoS *DPoSFilterer) ParseWhitelistedRemoved(log types.Log) (*DPoSWhitelistedRemoved, error) {
	event := new(DPoSWhitelistedRemoved)
	if err := _DPoS.contract.UnpackLog(event, "WhitelistedRemoved", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DPoSWithdrawFromUnbondedCandidateIterator is returned from FilterWithdrawFromUnbondedCandidate and is used to iterate over the raw logs and unpacked data for WithdrawFromUnbondedCandidate events raised by the DPoS contract.
type DPoSWithdrawFromUnbondedCandidateIterator struct {
	Event *DPoSWithdrawFromUnbondedCandidate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *DPoSWithdrawFromUnbondedCandidateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DPoSWithdrawFromUnbondedCandidate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(DPoSWithdrawFromUnbondedCandidate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *DPoSWithdrawFromUnbondedCandidateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DPoSWithdrawFromUnbondedCandidateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DPoSWithdrawFromUnbondedCandidate represents a WithdrawFromUnbondedCandidate event raised by the DPoS contract.
type DPoSWithdrawFromUnbondedCandidate struct {
	Delegator common.Address
	Candidate common.Address
	Amount    *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterWithdrawFromUnbondedCandidate is a free log retrieval operation binding the contract event 0x585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8.
//
// Solidity: event WithdrawFromUnbondedCandidate(address indexed delegator, address indexed candidate, uint256 amount)
func (_DPoS *DPoSFilterer) FilterWithdrawFromUnbondedCandidate(opts *bind.FilterOpts, delegator []common.Address, candidate []common.Address) (*DPoSWithdrawFromUnbondedCandidateIterator, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.FilterLogs(opts, "WithdrawFromUnbondedCandidate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return &DPoSWithdrawFromUnbondedCandidateIterator{contract: _DPoS.contract, event: "WithdrawFromUnbondedCandidate", logs: logs, sub: sub}, nil
}

// WatchWithdrawFromUnbondedCandidate is a free log subscription operation binding the contract event 0x585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8.
//
// Solidity: event WithdrawFromUnbondedCandidate(address indexed delegator, address indexed candidate, uint256 amount)
func (_DPoS *DPoSFilterer) WatchWithdrawFromUnbondedCandidate(opts *bind.WatchOpts, sink chan<- *DPoSWithdrawFromUnbondedCandidate, delegator []common.Address, candidate []common.Address) (event.Subscription, error) {

	var delegatorRule []interface{}
	for _, delegatorItem := range delegator {
		delegatorRule = append(delegatorRule, delegatorItem)
	}
	var candidateRule []interface{}
	for _, candidateItem := range candidate {
		candidateRule = append(candidateRule, candidateItem)
	}

	logs, sub, err := _DPoS.contract.WatchLogs(opts, "WithdrawFromUnbondedCandidate", delegatorRule, candidateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DPoSWithdrawFromUnbondedCandidate)
				if err := _DPoS.contract.UnpackLog(event, "WithdrawFromUnbondedCandidate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseWithdrawFromUnbondedCandidate is a log parse operation binding the contract event 0x585e40624b400c05be4193af453d2fd2e69facd17163bda6afd44546f3dbbaa8.
//
// Solidity: event WithdrawFromUnbondedCandidate(address indexed delegator, address indexed candidate, uint256 amount)
func (_DPoS *DPoSFilterer) ParseWithdrawFromUnbondedCandidate(log types.Log) (*DPoSWithdrawFromUnbondedCandidate, error) {
	event := new(DPoSWithdrawFromUnbondedCandidate)
	if err := _DPoS.contract.UnpackLog(event, "WithdrawFromUnbondedCandidate", log); err != nil {
		return nil, err
	}
	return event, nil
}
