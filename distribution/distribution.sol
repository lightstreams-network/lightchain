pragma solidity ^0.5.1;

contract Distribution {
    address public owner;
    mapping(address => uint) deposits;

    constructor() public {
        owner = msg.sender;
    }

    modifier onlyOwner {
        require(msg.sender == owner, "Only owner can call this function.");
        _;
    }

    function () external payable {
    }

    function deposit(address _beneficiary) payable public onlyOwner {
        require(_beneficiary != address(0));
        require(deposits[_beneficiary] == 0);
        require(msg.value > 0);

        deposits[_beneficiary] = msg.value;
    }

    function withdraw() payable public {
        require(deposits[msg.sender] > 0);

        uint beneficiaryAmount = deposits[msg.sender];
        deposits[msg.sender] = 0;

        msg.sender.transfer(beneficiaryAmount);
    }

    function changeDepositBeneficiary(address _oldBeneficiary, address _newBeneficiary) public onlyOwner {
        require(_newBeneficiary != address(0));
        require(deposits[_newBeneficiary] == 0);

        uint beneficiaryAmount = deposits[_oldBeneficiary];
        deposits[_oldBeneficiary] = 0;
        deposits[_newBeneficiary] = beneficiaryAmount;
    }
}