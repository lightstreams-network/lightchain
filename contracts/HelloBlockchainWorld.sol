pragma solidity ^0.4.24;

contract HelloBlockchainWorld {
    address public owner;

    constructor() public {
        owner = msg.sender;
    }

    modifier onlyOwner {
        require(msg.sender == owner, "Only owner can call this function.");
        _;
    }

    function sayHello() onlyOwner public view returns (string memory message) {
        return "hello";
    }
}