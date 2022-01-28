// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.7.0 <0.9.0;
// simple contract for depositing erc-20 tokens with cosmos address, and emitting event

interface IERC20 {

    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function allowance(address owner, address spender) external view returns (uint256);

    function transfer(address recipient, uint256 amount) external returns (bool);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address sender, address recipient, uint256 amount) external returns (bool);

    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}

library SafeMath {
    function sub(uint256 a, uint256 b) internal pure returns (uint256) {
      assert(b <= a);
      return a - b;
    }

    function add(uint256 a, uint256 b) internal pure returns (uint256) {
      uint256 c = a + b;
      assert(c >= a);
      return c;
    }
}

contract BaseledgerTest {
  event SendToCosmosEvent(
		address indexed _tokenContract,
		address indexed _sender,
		string _destination,
		uint256 _amount,
		uint256 _eventNonce
	);

  event ValidatorPowerChangeEvent(
		address indexed _tokenContract,
		address indexed _sender,
		string _destination,
		uint256 _amount,
		uint256 _eventNonce
	);

  IERC20 private token;
  address private tokenAddress;
  
  constructor(address erc20Address) {
    tokenAddress = erc20Address;
    token = IERC20(erc20Address);
  }

  // event nonce zero is reserved by the Cosmos module as a special
	// value indicating that no events have yet been submitted
	uint256 public state_lastEventNonce = 0;

  function deposit(uint256 amount, string calldata destination) public {
    require(amount > 0, "Deposit should be greater than zero.");
    uint256 allowance = token.allowance(msg.sender, address(this));
    require(allowance >= amount, "Check the token allowance");
    token.transferFrom(msg.sender, address(this), amount);
    state_lastEventNonce = state_lastEventNonce + 1;

		emit SendToCosmosEvent(
			tokenAddress,
			msg.sender,
			destination,
			amount,
			state_lastEventNonce
		);
  }

  // dummy method, implementation same as above, just to test emitting and catching power change event
  function powerChange(uint256 amount, string calldata destination) public {
    require(amount > 0, "Deposit should be greater than zero.");
    state_lastEventNonce = state_lastEventNonce + 1;

		emit ValidatorPowerChangeEvent(
			tokenAddress,
			msg.sender,
			destination,
			amount,
			state_lastEventNonce
		);
  }
}