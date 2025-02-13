// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts (last updated v4.6.0) (token/ERC20/IERC20.sol)

pragma solidity ^0.8.0;

/**
 * @dev Interface of the ERC20 standard as defined in the EIP.
 */
interface IAccessControl {

    /**
     * @dev Set role for an address
     */
   function grantRole(bytes32 role, address addr) external returns (bool);

   /**
     * @dev Unset role for an address
     */
   function revokeRole(bytes32 role, address addr) external returns (bool);
}
