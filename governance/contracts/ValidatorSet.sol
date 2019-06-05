pragma solidity ^0.5.0;

import "./Ownable.sol";

contract ValidatorSet is Ownable {

  mapping(bytes32 => address) public _validators;
  
  event ValidatorAdded(address _address, string _key);
  event ValidatorRemoved(address _address, string _key);

  
  function addValidator(string memory _key, address _address) onlyOwner public {
    require(_address != address(0x0));
    bytes32 _bkey = _stringToBytes32(_key);
    require(_validators[_bkey] == address(0x0));
    _validators[_bkey] = _address;
    
    emit ValidatorAdded(_address, _key);
  }
  
  function removeValidator(string memory _key, address _address) onlyOwner public {
    bytes32 _bkey = _stringToBytes32(_key);
    require(_validators[_bkey] == _address);
    _validators[_bkey] = address(0x0);
    
    emit ValidatorRemoved(_address, _key);
  }
  
  function validatorAddress(string memory _key) public view returns (address) {
      return _validators[_stringToBytes32(_key)];
  }
  
  function _stringToBytes32(string memory source) pure internal returns (bytes32 result) {
        bytes memory tempEmptyStringTest = bytes(source);
        if (tempEmptyStringTest.length == 0) {
            return 0x0;
        }

        assembly {
            result := mload(add(source, 32))
        }
    }
}
