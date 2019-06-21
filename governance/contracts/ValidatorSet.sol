pragma solidity ^0.5.0;

import "./Ownable.sol";

contract ValidatorSet is Ownable {

  mapping(bytes20 => address) internal _vAccountAddresses;
  bytes20[] public _vPubKeyAddresses;

  bool public _freeze = false;
  address public _nextVersion;

  event ValidatorAdded(bytes20 vPubKeyAddress, address vAddress);
  event ValidatorRemoved(bytes20 vPubKeyAddress, address vAddress);
  event SetNextVersion(address nAddress);
  event Freeze();

  function addValidator(bytes20 vPubKeyAddress, address vAddress) onlyOwner public {
    require(_freeze == false);
    require(vAddress != address(0x0));
    require(_vAccountAddresses[vPubKeyAddress] == address(0x0));

    _vAccountAddresses[vPubKeyAddress] = vAddress;
    _vPubKeyAddresses.push(vPubKeyAddress);

    emit ValidatorAdded(vPubKeyAddress, vAddress);
  }
  
  function removeValidator(bytes20 vPubKeyAddress, address vAddress) onlyOwner public {
    require(_freeze == false);
    require(vAddress != address(0x0));
    require(_vAccountAddresses[vPubKeyAddress] == vAddress);

    _vAccountAddresses[vPubKeyAddress] = address(0x0);

    bool _foundDeletedItem = false;
    for (uint i = 0; i < _vPubKeyAddresses.length-1; i++) {
      if (_vPubKeyAddresses[i] == vPubKeyAddress) {
        _foundDeletedItem = true;
      }
      if (_foundDeletedItem) {
        _vPubKeyAddresses[i] = _vPubKeyAddresses[i+1];
      }
    }

    delete _vPubKeyAddresses[uint(_vPubKeyAddresses.length-1)];
    _vPubKeyAddresses.length--;

    emit ValidatorRemoved(vPubKeyAddress, vAddress);
  }
  
  function validatorAddress(bytes20 vPubKeyAddress) public view returns (address) {
    return _vAccountAddresses[vPubKeyAddress];
  }
  
  function validatorPubKey(uint index) public view returns (bytes20) {
    return _vPubKeyAddresses[index];
  }

  function setFreezeStatus(bool value) onlyOwner public {
    _freeze = value;
    emit Freeze();
  }

  function validatorSetSize() public view returns (uint) {
    return uint(_vPubKeyAddresses.length);
  }

  function _setNextVersionAddress(address value) onlyOwner public {
    _nextVersion = value;
    emit SetNextVersion(value);
  }
}
