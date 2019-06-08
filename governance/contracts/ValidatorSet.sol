pragma solidity ^0.5.0;

import "./Ownable.sol";

contract ValidatorSet is Ownable {

  mapping(bytes32 => address) internal _validatorsAddress;
  string[] public _validatorPubKeys;

  bool public _freeze = false;
  address public _nextVersion;

  event ValidatorAdded(bytes32 _key, address _address, string _pubKey);
  event ValidatorRemoved(bytes32 _key, address _address, string _pubKey);
  event Freeze();
  event SetNextVersion(address _address);

  function addValidator(string memory _pubKey, address _address) onlyOwner public {
    require(_freeze == false);
    require(_address != address(0x0));
    require(bytes(_pubKey).length == 40);

    bytes32 _vItemKey = calculateValidatorItemKey(_pubKey);
    require(_validatorsAddress[_vItemKey] == address(0x0));
    
    _validatorsAddress[_vItemKey] = _address;
    _validatorPubKeys.push(_pubKey);

    emit ValidatorAdded(_vItemKey, _address, _pubKey);
  }
  
  function removeValidator(string memory _pubKey, address _address) onlyOwner public {
    require(_freeze == false);
    require(_address != address(0x0));
    require(bytes(_pubKey).length == 40);

    bytes32 _vItemKey = calculateValidatorItemKey(_pubKey);
    require(_validatorsAddress[_vItemKey] == _address);

    _validatorsAddress[_vItemKey] = address(0x0);
    bool _foundDeletedItem = false;
    for (uint i = 0; i < _validatorPubKeys.length-1; i++) {
      if (calculateValidatorItemKey(_validatorPubKeys[i]) == _vItemKey) {
        _foundDeletedItem = true;
      }
      if (_foundDeletedItem) {
        _validatorPubKeys[i] = _validatorPubKeys[i+1];
      }
    }

    delete _validatorPubKeys[uint(_validatorPubKeys.length-1)];
    _validatorPubKeys.length--;
    emit ValidatorRemoved(_vItemKey, _address, _pubKey);
  }
  
  function validatorAddress(string memory _pubKey) public view returns (address) {
    bytes32 _vItemKey = calculateValidatorItemKey(_pubKey);
    return _validatorsAddress[_vItemKey];
  }
  
  function validatorPubKey(uint index) public view returns (string memory) {
    return _validatorPubKeys[index];
  }

  function setFreezeStatus(bool _value) onlyOwner public {
    _freeze = _value;
    emit Freeze();
  }
  
  
  function calculateValidatorItemKey(string memory _pubKey) pure internal returns (bytes32) {
    return sha256(bytes(_pubKey));
  }

  function validatorSetSize() public view returns (uint) {
    return uint(_validatorPubKeys.length);
  }

  function _setNextVersionAddress(address _value) onlyOwner public {
    _nextVersion = _value;
    emit SetNextVersion(_value);
  }
}
