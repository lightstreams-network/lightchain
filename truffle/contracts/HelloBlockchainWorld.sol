pragma solidity ^0.5.0;

contract HelloBlockchainWorld {
    address public owner;
    uint public helloCount;

    event HelloCountIncremented(uint newCount);

    constructor() public {
        owner = msg.sender;
    }

    modifier onlyOwner {
        require(msg.sender == owner, "Only owner can call this function.");
        _;
    }

    function sayHello() public pure returns (string memory message) {
        return "hello";
    }

    function incrementHelloCount() onlyOwner public returns (uint) {
        helloCount = helloCount + 1;

        emit HelloCountIncremented(helloCount);

        return helloCount;
    }
}