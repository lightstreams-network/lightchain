pragma solidity ^0.5.0;

import "./Ownable.sol";

contract ValidatorSet is Ownable {

  mapping(bytes32 => address) internal _vAccAddress;
  string[] public _vPubKeyAddresses;

  bool public _freeze = false;
  address public _nextVersion;

  event ValidatorAdded(bytes32 key, address vAddress, string vPubKeyAddress);
  event ValidatorRemoved(bytes32 key, address vAddress, string vPubKeyAddress);
  event SetNextVersion(address vAddress);
  event Freeze();

  function addValidator(string memory vPubKeyAddress, address vAddress) onlyOwner public {
    require(_freeze == false);
    require(vAddress != address(0x0));
    require(bytes(vPubKeyAddress).length == 40);

    bytes32 _vItemKey = calculateValidatorItemKey(vPubKeyAddress);
    require(_vAccAddress[_vItemKey] == address(0x0));
    
    _vAccAddress[_vItemKey] = vAddress;
    _vPubKeyAddresses.push(vPubKeyAddress);

    emit ValidatorAdded(_vItemKey, vAddress, vPubKeyAddress);
  }
  
  function removeValidator(string memory vPubKeyAddress, address vAddress) onlyOwner public {
    require(_freeze == false);
    require(vAddress != address(0x0));
    require(bytes(vPubKeyAddress).length == 40);

    bytes32 vItemKey = calculateValidatorItemKey(vPubKeyAddress);
    require(_vAccAddress[vItemKey] == vAddress);

    _vAccAddress[vItemKey] = address(0x0);
    bool _foundDeletedItem = false;
    for (uint i = 0; i < _vPubKeyAddresses.length-1; i++) {
      if (calculateValidatorItemKey(_vPubKeyAddresses[i]) == vItemKey) {
        _foundDeletedItem = true;
      }
      if (_foundDeletedItem) {
        _vPubKeyAddresses[i] = _vPubKeyAddresses[i+1];
      }
    }

    delete _vPubKeyAddresses[uint(_vPubKeyAddresses.length-1)];
    _vPubKeyAddresses.length--;
    emit ValidatorRemoved(vItemKey, vAddress, vPubKeyAddress);
  }
  
  function validatorAddress(string memory vPubKeyAddress) public view returns (address) {
    bytes32 vItemKey = calculateValidatorItemKey(vPubKeyAddress);
    return _vAccAddress[vItemKey];
  }
  
  function validatorPubKey(uint index) public view returns (string memory) {
    return _vPubKeyAddresses[index];
  }

  function setFreezeStatus(bool value) onlyOwner public {
    _freeze = value;
    emit Freeze();
  }
  
  
  function calculateValidatorItemKey(string memory vPubKeyAddress) pure internal returns (bytes32) {
    return sha256(bytes(vPubKeyAddress));
  }

  function validatorSetSize() public view returns (uint) {
    return uint(_vPubKeyAddresses.length);
  }

  function _setNextVersionAddress(address value) onlyOwner public {
    _nextVersion = value;
    emit SetNextVersion(value);
  }
}
