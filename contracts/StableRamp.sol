// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";
import "@openzeppelin/contracts/security/Pausable.sol";
import "@openzeppelin/contracts/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

// CommitteeArgs holds the current committee state passed into functions.
// Mirrors the Committee type in the Cosmos module.
struct CommitteeArgs {
    address[] members;
    uint256[] powers;
    uint256 nonce;
}

/**
 * @title  StableRamp
 * @notice Bridge contract for the Realio stable-ramp Cosmos module.
 *
 * Deposit  (ETH → Cosmos): user calls sendToRealio → contract locks tokens →
 *   emits SendToRealioEvent with a monotonic eventNonce → relayer submits
 *   MsgDepositClaim on Cosmos using that eventNonce.
 *
 * Withdrawal (Cosmos → ETH): Cosmos module reaches quorum on a withdraw package →
 *   relayer calls submitWithdrawal with committee ECDSA signatures →
 *   contract verifies signatures and releases tokens to the recipient.
 *
 * Committee: mirrors the stable-ramp Committee on Cosmos. Updates require
 *   the current committee to sign off on the new one, same as Peggy valset updates.
 */
contract StableRamp is Initializable, Pausable, ReentrancyGuard, Ownable {
    using SafeERC20 for IERC20;

    // ─────────────────────────────────────────────────────────────────────────
    // State — only append, never reorder (upgradeable-safe)
    // ─────────────────────────────────────────────────────────────────────────

    // Unique identifier for this bridge instance, included in every signed hash
    // to prevent cross-contract replay attacks.
    bytes32 public state_bridgeId;

    // Minimum cumulative power required to approve an action.
    uint256 public state_powerThreshold;

    // Hash of the current committee (members + powers + nonce).
    // Stored instead of the full committee to save gas.
    bytes32 public state_committeeCheckpoint;

    // Nonce of the current committee (monotonically increasing).
    uint256 public state_committeeNonce;

    // Monotonically increasing nonce emitted with every SendToRealioEvent.
    // Used by the Cosmos module as the deposit event_nonce.
    uint256 public state_lastEventNonce;

    // Withdrawal replay protection: Cosmos package nonce → executed.
    // Once set to true, that withdrawal can never be re-executed.
    mapping(uint256 => bool) public state_executedWithdrawals;

    // ─────────────────────────────────────────────────────────────────────────
    // Events
    // ─────────────────────────────────────────────────────────────────────────

    // Emitted when a user deposits tokens to bridge to Cosmos.
    // The relayer uses eventNonce as the MsgDepositClaim.event_nonce.
    event SendToRealioEvent(
        address indexed tokenContract,
        address indexed sender,
        string  cosmosReceiver,
        uint256 amount,
        uint256 eventNonce
    );

    // Emitted when a committee-approved Cosmos withdrawal is executed.
    event WithdrawalExecutedEvent(
        uint256 indexed withdrawalNonce,
        address indexed tokenContract,
        address indexed recipient,
        uint256 amount,
        uint256 eventNonce
    );

    // Emitted when the committee is updated.
    event CommitteeUpdatedEvent(
        uint256 indexed newCommitteeNonce,
        uint256 eventNonce,
        address[] members,
        uint256[] powers
    );

    // ─────────────────────────────────────────────────────────────────────────
    // Initializer
    // ─────────────────────────────────────────────────────────────────────────

    function initialize(
        bytes32  _bridgeId,
        uint256  _powerThreshold,
        address[] calldata _members,
        uint256[] calldata _powers
    ) external initializer {
        require(_members.length == _powers.length, "Malformed committee");
        _requireSufficientPower(_members, _powers, _powerThreshold);

        state_bridgeId          = _bridgeId;
        state_powerThreshold    = _powerThreshold;
        state_committeeCheckpoint = _makeCheckpoint(CommitteeArgs(_members, _powers, 0));
        state_committeeNonce    = 0;
        state_lastEventNonce    = 1;

        emit CommitteeUpdatedEvent(0, state_lastEventNonce, _members, _powers);
    }

    // ─────────────────────────────────────────────────────────────────────────
    // User-facing: deposit ETH → Cosmos
    // ─────────────────────────────────────────────────────────────────────────

    /**
     * @notice Lock ERC20 tokens in this contract and notify the Cosmos module.
     * @param _tokenContract  ERC20 token to deposit.
     * @param _cosmosReceiver Bech32 address on the Cosmos side.
     * @param _amount         Token amount (before fee-on-transfer deduction).
     */
    function sendToRealio(
        address _tokenContract,
        string  calldata _cosmosReceiver,
        uint256 _amount
    ) external nonReentrant whenNotPaused {
        require(_amount > 0,                         "Amount must be positive");
        require(bytes(_cosmosReceiver).length > 0,   "Cosmos receiver required");

        // Use balance delta to support fee-on-transfer tokens.
        uint256 balanceBefore = IERC20(_tokenContract).balanceOf(address(this));
        IERC20(_tokenContract).safeTransferFrom(msg.sender, address(this), _amount);
        uint256 actualAmount = IERC20(_tokenContract).balanceOf(address(this)) - balanceBefore;

        state_lastEventNonce++;
        emit SendToRealioEvent(
            _tokenContract,
            msg.sender,
            _cosmosReceiver,
            actualAmount,
            state_lastEventNonce
        );
    }

    // ─────────────────────────────────────────────────────────────────────────
    // Relayer-facing: execute Cosmos → ETH withdrawal
    // ─────────────────────────────────────────────────────────────────────────

    /**
     * @notice Execute a withdrawal approved by the Cosmos committee.
     *
     * The relayer observes EventTypeWithdrawClaimConfirm on Cosmos, extracts
     * the package nonce, token denom (mapped to tokenContract), recipient,
     * amount, and the committee member signatures collected in TrackedPackage,
     * then calls this function.
     *
     * Committee members sign:
     *   keccak256(bridgeId, "withdrawal", withdrawalNonce, tokenContract,
     *             recipient, amount, timeout)
     *
     * @param _currentCommittee  Must match state_committeeCheckpoint.
     * @param _v / _r / _s       ECDSA signature components per member.
     * @param _withdrawalNonce   Cosmos package nonce (replay guard).
     * @param _tokenContract     ERC20 to release.
     * @param _recipient         Ethereum address to receive the tokens.
     * @param _amount            Amount to release.
     * @param _timeout           Block number deadline; submission fails after this.
     */
    function submitWithdrawal(
        CommitteeArgs calldata _currentCommittee,
        uint8[]  calldata _v,
        bytes32[] calldata _r,
        bytes32[] calldata _s,
        uint256 _withdrawalNonce,
        address _tokenContract,
        address _recipient,
        uint256 _amount,
        uint256 _timeout
    ) external nonReentrant whenNotPaused {
        // ── Checks ────────────────────────────────────────────────────────────

        require(block.number < _timeout,
            "Withdrawal submission has timed out");

        require(!state_executedWithdrawals[_withdrawalNonce],
            "Withdrawal already executed");

        require(
            _currentCommittee.members.length == _currentCommittee.powers.length &&
            _currentCommittee.members.length == _v.length &&
            _currentCommittee.members.length == _r.length &&
            _currentCommittee.members.length == _s.length,
            "Malformed committee or signatures"
        );

        require(
            _makeCheckpoint(_currentCommittee) == state_committeeCheckpoint,
            "Supplied committee does not match stored checkpoint"
        );

        bytes32 withdrawalHash = keccak256(abi.encode(
            state_bridgeId,
            bytes32("withdrawal"),
            _withdrawalNonce,
            _tokenContract,
            _recipient,
            _amount,
            _timeout
        ));

        _checkSignatures(
            _currentCommittee.members,
            _currentCommittee.powers,
            _v, _r, _s,
            withdrawalHash
        );

        // ── Actions ───────────────────────────────────────────────────────────

        state_executedWithdrawals[_withdrawalNonce] = true;
        IERC20(_tokenContract).safeTransfer(_recipient, _amount);

        // ── Logs ──────────────────────────────────────────────────────────────

        state_lastEventNonce++;
        emit WithdrawalExecutedEvent(
            _withdrawalNonce,
            _tokenContract,
            _recipient,
            _amount,
            state_lastEventNonce
        );
    }

    // ─────────────────────────────────────────────────────────────────────────
    // Committee update
    // ─────────────────────────────────────────────────────────────────────────

    /**
     * @notice Replace the committee. The current committee must sign the
     *         checkpoint of the new one, mirroring how Cosmos updates its
     *         Committee via MsgUpdateCommittee + governance.
     */
    function updateCommittee(
        CommitteeArgs calldata _newCommittee,
        CommitteeArgs calldata _currentCommittee,
        uint8[]  calldata _v,
        bytes32[] calldata _r,
        bytes32[] calldata _s
    ) external whenNotPaused {
        require(
            _newCommittee.nonce > _currentCommittee.nonce,
            "New committee nonce must be greater than current"
        );
        require(
            _makeCheckpoint(_currentCommittee) == state_committeeCheckpoint,
            "Current committee mismatch"
        );
        require(
            _currentCommittee.members.length == _v.length &&
            _currentCommittee.members.length == _r.length &&
            _currentCommittee.members.length == _s.length,
            "Malformed signatures"
        );

        _requireSufficientPower(
            _newCommittee.members,
            _newCommittee.powers,
            state_powerThreshold
        );

        bytes32 newCheckpoint = _makeCheckpoint(_newCommittee);
        _checkSignatures(
            _currentCommittee.members,
            _currentCommittee.powers,
            _v, _r, _s,
            newCheckpoint
        );

        state_committeeCheckpoint = newCheckpoint;
        state_committeeNonce      = _newCommittee.nonce;

        state_lastEventNonce++;
        emit CommitteeUpdatedEvent(
            _newCommittee.nonce,
            state_lastEventNonce,
            _newCommittee.members,
            _newCommittee.powers
        );
    }

    // ─────────────────────────────────────────────────────────────────────────
    // Emergency
    // ─────────────────────────────────────────────────────────────────────────

    function pause()   external onlyOwner { _pause(); }
    function unpause() external onlyOwner { _unpause(); }

    // ─────────────────────────────────────────────────────────────────────────
    // Internal helpers
    // ─────────────────────────────────────────────────────────────────────────

    function _makeCheckpoint(
        CommitteeArgs memory _committee
    ) private pure returns (bytes32) {
        return keccak256(abi.encode(
            bytes32("committee"),
            _committee.nonce,
            _committee.members,
            _committee.powers
        ));
    }

    function _verifySig(
        address _signer,
        bytes32 _hash,
        uint8 _v, bytes32 _r, bytes32 _s
    ) private pure returns (bool) {
        bytes32 digest = keccak256(abi.encodePacked(
            "\x19Ethereum Signed Message:\n32",
            _hash
        ));
        return _signer == ecrecover(digest, _v, _r, _s);
    }

    function _checkSignatures(
        address[] memory _members,
        uint256[] memory _powers,
        uint8[]   memory _v,
        bytes32[] memory _r,
        bytes32[] memory _s,
        bytes32 _hash
    ) private view {
        uint256 cumulative;
        for (uint256 i = 0; i < _members.length; i++) {
            if (_v[i] != 0) {
                require(
                    _verifySig(_members[i], _hash, _v[i], _r[i], _s[i]),
                    "Invalid committee member signature"
                );
                cumulative += _powers[i];
                if (cumulative > state_powerThreshold) break;
            }
        }
        require(cumulative > state_powerThreshold, "Insufficient signature power");
    }

    function _requireSufficientPower(
        address[] memory _members,
        uint256[] memory _powers,
        uint256 _threshold
    ) private pure {
        uint256 cumulative;
        for (uint256 i = 0; i < _members.length; i++) {
            cumulative += _powers[i];
            if (cumulative > _threshold) break;
        }
        require(cumulative > _threshold, "Committee has insufficient total power");
    }
}
