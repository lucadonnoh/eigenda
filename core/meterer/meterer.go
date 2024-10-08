package meterer

import (
	"context"
	"fmt"
	"math/big"
	"slices"
	"time"

	"github.com/Layr-Labs/eigenda/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"github.com/ethereum/go-ethereum/common"
)

type TimeoutConfig struct {
	ChainReadTimeout    time.Duration
	ChainWriteTimeout   time.Duration
	ChainStateTimeout   time.Duration
	TxnBroadcastTimeout time.Duration
}

// network parameters (this should be published on-chain and read through contracts)
type Config struct {
	GlobalBytesPerSecond uint64 // 2^64 bytes ~= 18 exabytes per second; if we use uint32, that's ~4GB/s
	PricePerChargeable   uint32 // 2^64 gwei ~= 18M Eth; uint32 => ~4ETH
	MinChargeableSize    uint32
	ReservationWindow    uint32
}

// disperser API server will receive requests from clients. these requests will be with a blobHeader with payments information (CumulativePayments, BinIndex, and Signature)
// Disperser will pass the blob header to the meterer, which will check if the payments information is valid. if it is, it will be added to the meterer's state.
// To check if the payment is valid, the meterer will:
//  1. check if the signature is valid
//     (against the CumulativePayments and BinIndex fields ;
//     maybe need something else to secure against using this appraoch for reservations when rev request comes in same bin interval; say that nonce is signed over as well)
//  2. For reservations, check offchain bin state as demonstrated in pseudocode, also check onchain state before rejecting (since onchain data is pulled)
//  3. For on-demand, check against payments and the global rates, similar to the reservation case
//
// If the payment is valid, the meterer will add the blob header to its state and return a success response to the disperser API server.
// if any of the checks fail, the meterer will return a failure response to the disperser API server.
var OnDemandQuorumNumbers = []uint8{0, 1}

type Meterer struct {
	Config
	TimeoutConfig

	ChainState    *OnchainPaymentState
	OffchainStore *OffchainStore

	logger logging.Logger
}

func NewMeterer(
	config Config,
	timeoutConfig TimeoutConfig,
	paymentChainState *OnchainPaymentState,
	offchainStore *OffchainStore,
	logger logging.Logger,
) (*Meterer, error) {
	// TODO: create a separate thread to pull from the chain and update chain state
	return &Meterer{
		Config:        config,
		TimeoutConfig: timeoutConfig,

		ChainState:    paymentChainState,
		OffchainStore: offchainStore,

		logger: logger.With("component", "Meterer"),
	}, nil
}

// MeterRequest validates a blob header and adds it to the meterer's state
// TODO: return error if there's a rejection (with reasoning) or internal error (should be very rare)
func (m *Meterer) MeterRequest(ctx context.Context, header BlobHeader) error {
	// TODO: validate signing
	// if err := m.ValidateSignature(ctx, header); err != nil {
	// 	return fmt.Errorf("invalid signature: %w", err)
	// }

	blockNumber, err := m.ChainState.GetCurrentBlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current block number: %w", err)
	}

	// Validate against the payment method
	if header.CumulativePayment == 0 {
		reservation, err := m.ChainState.GetActiveReservationByAccount(ctx, blockNumber, header.AccountID)
		if err != nil {
			return fmt.Errorf("failed to get active reservation by account: %w", err)
		}
		if err := m.ServeReservationRequest(ctx, header, reservation); err != nil {
			return fmt.Errorf("invalid reservation: %w", err)
		}
	} else {
		onDemandPayment, err := m.ChainState.GetOnDemandPaymentByAccount(ctx, blockNumber, header.AccountID)
		if err != nil {
			return fmt.Errorf("failed to get on-demand payment by account: %w", err)
		}
		if err := m.ServeOnDemandRequest(ctx, header, onDemandPayment); err != nil {
			return fmt.Errorf("invalid on-demand request: %w", err)
		}
	}

	return nil
}

// TODO: mocked EIP712 domain, change to the real thing when available
// ValidateSignature checks if the signature is valid against all other fields in the header
// Assuming the signature is an eip712 signature
func (m *Meterer) ValidateSignature(ctx context.Context, header BlobHeader) error {
	// Create the EIP712Signer
	//TODO: update the chainID and verifyingContract
	signer := NewEIP712Signer(big.NewInt(17000), common.HexToAddress("0x1234000000000000000000000000000000000000"))

	recoveredAddress, err := signer.RecoverSender(&header)
	if err != nil {
		return fmt.Errorf("failed to recover sender: %w", err)
	}

	accountAddress := common.HexToAddress(header.AccountID)

	if recoveredAddress != accountAddress {
		return fmt.Errorf("invalid signature: recovered address %s does not match account ID %s", recoveredAddress.Hex(), accountAddress.Hex())
	}

	return nil
}

// ServeReservationRequest handles the rate limiting logic for incoming requests
func (m *Meterer) ServeReservationRequest(ctx context.Context, blobHeader BlobHeader, reservation *ActiveReservation) error {
	if err := m.ValidateQuorum(blobHeader, reservation.QuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for reservation: %w", err)
	}
	if !m.ValidateBinIndex(blobHeader, reservation) {
		return fmt.Errorf("invalid bin index for reservation")
	}

	// Update bin usage atomically and check against reservation's data rate as the bin limit
	if err := m.IncrementBinUsage(ctx, blobHeader, reservation); err != nil {
		return fmt.Errorf("bin overflows: %w", err)
	}

	return nil
}

func (m *Meterer) ValidateQuorum(blobHeader BlobHeader, allowedQuoroms []uint8) error {
	if len(blobHeader.QuorumNumbers) == 0 {
		return fmt.Errorf("no quorum params in blob header")
	}

	// check that all the quorum ids are in ActiveReservation's
	for _, q := range blobHeader.QuorumNumbers {
		if !slices.Contains(allowedQuoroms, q) {
			// fail the entire request if there's a quorum number mismatch
			return fmt.Errorf("quorum number mismatch: %d", q)
		}
	}
	return nil
}

// ValidateBinIndex checks if the provided bin index is valid
func (m *Meterer) ValidateBinIndex(blobHeader BlobHeader, reservation *ActiveReservation) bool {
	now := uint64(time.Now().Unix())
	currentBinIndex := GetBinIndex(now, m.ReservationWindow)
	// Valid bin indexes are either the current bin or the previous bin
	if (blobHeader.BinIndex != currentBinIndex && blobHeader.BinIndex != (currentBinIndex-1)) || (GetBinIndex(reservation.StartTimestamp, m.ReservationWindow) > blobHeader.BinIndex || blobHeader.BinIndex > GetBinIndex(reservation.EndTimestamp, m.ReservationWindow)) {
		return false
	}
	return true
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
// TODO: Bin limit should be direct write to the Store
func (m *Meterer) IncrementBinUsage(ctx context.Context, blobHeader BlobHeader, reservation *ActiveReservation) error {
	recordedSize := uint64(max(blobHeader.DataLength, uint32(m.MinChargeableSize)))
	newUsage, err := m.OffchainStore.UpdateReservationBin(ctx, blobHeader.AccountID, uint64(blobHeader.BinIndex), recordedSize)
	if err != nil {
		return fmt.Errorf("failed to increment bin usage: %w", err)
	}

	// metered usage stays within the bin limit
	if newUsage <= reservation.DataRate {
		return nil
	} else if newUsage-recordedSize >= reservation.DataRate {
		// metered usage before updating the size already exceeded the limit
		return fmt.Errorf("bin has already been filled")
	}
	if newUsage <= 2*reservation.DataRate && blobHeader.BinIndex+2 <= GetBinIndex(reservation.EndTimestamp, m.ReservationWindow) {
		m.OffchainStore.UpdateReservationBin(ctx, blobHeader.AccountID, uint64(blobHeader.BinIndex+2), newUsage-reservation.DataRate)
		return nil
	}
	return fmt.Errorf("overflow usage exceeds bin limit")
}

// GetBinIndex returns the current bin index by chunking time by the bin interval;
// bin interval used by the disperser should be public information
func GetBinIndex(timestamp uint64, binInterval uint32) uint32 {
	return uint32(timestamp) / binInterval
}

// ServeOnDemandRequest handles the rate limiting logic for incoming requests
func (m *Meterer) ServeOnDemandRequest(ctx context.Context, blobHeader BlobHeader, onDemandPayment *OnDemandPayment) error {
	if err := m.ValidateQuorum(blobHeader, OnDemandQuorumNumbers); err != nil {
		return fmt.Errorf("invalid quorum for On-Demand Request: %w", err)
	}
	// update blob header to use the miniumum chargeable size
	blobSizeCharged := m.BlobSizeCharged(blobHeader.DataLength)
	err := m.OffchainStore.AddOnDemandPayment(ctx, blobHeader, blobSizeCharged)
	if err != nil {
		return fmt.Errorf("failed to update cumulative payment: %w", err)
	}
	// Validate payments attached
	err = m.ValidatePayment(ctx, blobHeader, onDemandPayment)
	if err != nil {
		// No tolerance for incorrect payment amounts; no rollbacks
		return fmt.Errorf("invalid on-demand payment: %w", err)
	}

	// Update bin usage atomically and check against bin capacity
	if err := m.IncrementGlobalBinUsage(ctx, blobSizeCharged); err != nil {
		//TODO: conditionally remove the payment based on the error type (maybe if the error is store-op related)
		m.OffchainStore.RemoveOnDemandPayment(ctx, blobHeader.AccountID, blobHeader.CumulativePayment)
		return fmt.Errorf("failed global rate limiting")
	}

	return nil
}

// ValidatePayment checks if the provided payment header is valid against the local accounting
// prevPmt is the largest  cumulative payment strictly less    than blobHeader.cumulativePayment if exists
// nextPmt is the smallest cumulative payment strictly greater than blobHeader.cumulativePayment if exists
// nextPmtDataLength is the dataLength of corresponding to nextPmt if exists
// prevPmt + blobHeader.DataLength * m.FixedFeePerByte
// <= blobHeader.CumulativePayment
// <= nextPmt - nextPmtDataLength * m.FixedFeePerByte > nextPmt
func (m *Meterer) ValidatePayment(ctx context.Context, blobHeader BlobHeader, onDemandPayment *OnDemandPayment) error {
	if blobHeader.CumulativePayment > uint64(onDemandPayment.CumulativePayment) {
		return fmt.Errorf("request claims a cumulative payment greater than the on-chain deposit")
	}

	prevPmt, nextPmt, nextPmtDataLength, err := m.OffchainStore.GetRelevantOnDemandRecords(ctx, blobHeader.AccountID, blobHeader.CumulativePayment) // zero if DNE
	if err != nil {
		return fmt.Errorf("failed to get relevant on-demand records: %w", err)
	}
	// the current request must increment cumulative payment by a magnitude sufficient to cover the blob size
	if prevPmt+m.PaymentCharged(blobHeader.DataLength) > blobHeader.CumulativePayment {
		return fmt.Errorf("insufficient cumulative payment increment")
	}
	// the current request must not break the payment magnitude for the next payment if the two requests were delivered out-of-order
	if nextPmt != 0 && blobHeader.CumulativePayment+m.PaymentCharged(nextPmtDataLength) > nextPmt {
		return fmt.Errorf("breaking cumulative payment invariants")
	}
	// check passed: blob can be safely inserted into the set of payments
	return nil
}

// PaymentCharged returns the chargeable price for a given data length
func (m *Meterer) PaymentCharged(dataLength uint32) uint64 {
	return uint64(core.RoundUpDivide(uint(m.BlobSizeCharged(dataLength)*m.PricePerChargeable), uint(m.MinChargeableSize)))
}

// BlobSizeCharged returns the chargeable data length for a given data length
func (m *Meterer) BlobSizeCharged(dataLength uint32) uint32 {
	return max(dataLength, uint32(m.MinChargeableSize))
}

// ValidateBinIndex checks if the provided bin index is valid
func (m *Meterer) ValidateGlobalBinIndex(blobHeader BlobHeader) (uint32, error) {
	// Deterministic function: local clock -> index (1second intervals)
	currentBinIndex := uint32(time.Now().Unix())

	// Valid bin indexes are either the current bin or the previous bin (allow this second or prev sec)
	if blobHeader.BinIndex != currentBinIndex && blobHeader.BinIndex != (currentBinIndex-1) {
		return 0, fmt.Errorf("invalid bin index for on-demand request")
	}
	return currentBinIndex, nil
}

// IncrementBinUsage increments the bin usage atomically and checks for overflow
func (m *Meterer) IncrementGlobalBinUsage(ctx context.Context, blobSizeCharged uint32) error {
	globalIndex := uint64(time.Now().Unix())
	newUsage, err := m.OffchainStore.UpdateGlobalBin(ctx, globalIndex, blobSizeCharged)
	if err != nil {
		return fmt.Errorf("failed to increment global bin usage: %w", err)
	}
	if newUsage > m.GlobalBytesPerSecond {
		return fmt.Errorf("global bin usage overflows")
	}
	return nil
}
