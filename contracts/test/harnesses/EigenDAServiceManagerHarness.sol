// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "../../src/interfaces/IEigenDAServiceManager.sol";
import "../../src/libraries/EigenDAHasher.sol";
import "eigenlayer-middleware/BLSSignatureChecker.sol";
import "eigenlayer-middleware/ServiceManagerBase.sol";

contract EigenDAServiceManagerHarness is IEigenDAServiceManager {
    using EigenDAHasher for BatchHeader;
    using EigenDAHasher for ReducedBatchHeader;

    uint256 public THRESHOLD_DENOMINATOR = 100;
    uint32 public STORE_DURATION_BLOCKS = 2 weeks / 12 seconds;
    uint32 public BLOCK_STALE_MEASURE = 300;

    bytes public quorumAdversaryThresholdPercentages = hex"2121";
    bytes public quorumConfirmationThresholdPercentages = hex"3737";
    bytes public quorumNumbersRequired = hex"0001";

    uint32 public batchId;

    mapping(uint32 => bytes32) public batchIdToBatchMetadataHash;

    function mockConfirmBatch(
        BatchHeader calldata batchHeader,
        bytes32 signatoryRecordHash
    ) external {
        uint32 batchIdMemory = batchId;
        bytes32 batchHeaderHash = batchHeader.hashBatchHeader();
        bytes32 reducedBatchHeaderHash = batchHeader.hashBatchHeaderToReducedBatchHeader();
        batchIdToBatchMetadataHash[batchIdMemory] = EigenDAHasher.hashBatchHashedMetadata(batchHeaderHash, signatoryRecordHash, uint32(block.number));
        batchId = batchIdMemory + 1;
        emit BatchConfirmed(reducedBatchHeaderHash, batchIdMemory);
    }

    function setBatchId(uint32 _batchId) external {
        batchId = _batchId;
    }

    function setBatchIdToBatchMetadataHash(uint32 _batchId, bytes32 _batchMetadataHash) external {
        batchIdToBatchMetadataHash[_batchId] = _batchMetadataHash;
    }

    function setThresholdDenominator(uint256 _thresholdDenominator) external {
        THRESHOLD_DENOMINATOR = _thresholdDenominator;
    }

    function setStoreDurationBlocks(uint32 _storeDurationBlocks) external {
        STORE_DURATION_BLOCKS = _storeDurationBlocks;
    }

    function setBlockStaleMeasure(uint32 _blockStaleMeasure) external {
        BLOCK_STALE_MEASURE = _blockStaleMeasure;
    }

    function setQuorumAdversaryThresholdPercentages(bytes memory _quorumAdversaryThresholdPercentages) external {
        quorumAdversaryThresholdPercentages = _quorumAdversaryThresholdPercentages;
    }

    function setQuorumConfirmationThresholdPercentages(bytes memory _quorumConfirmationThresholdPercentages) external {
        quorumConfirmationThresholdPercentages = _quorumConfirmationThresholdPercentages;
    }

    function setQuorumNumbersRequired(bytes memory _quorumNumbersRequired) external {
        quorumNumbersRequired = _quorumNumbersRequired;
    }

    function confirmBatch(
        BatchHeader calldata batchHeader,
        BLSSignatureChecker.NonSignerStakesAndSignature memory nonSignerStakesAndSignature
    ) external {}

    function setBatchConfirmer(address _batchConfirmer) external {}

    function taskNumber() external view returns (uint32) {}

    function latestServeUntilBlock(uint32 referenceBlockNumber) external view returns (uint32) {}

    function avsDirectory() external view returns (address) {}

    function deregisterOperatorFromAVS(address operator) external {}

    function getOperatorRestakedStrategies(address operator) external view returns (address[] memory) {}

    function getRestakeableStrategies() external view returns (address[] memory) {}

    function payForRange(IPaymentCoordinator.RangePayment[] calldata rangePayments) external {}

    function registerOperatorToAVS(
        address operator,
        ISignatureUtils.SignatureWithSaltAndExpiry memory operatorSignature
    ) external {}

    function updateAVSMetadataURI(string memory _metadataURI) external {}
}