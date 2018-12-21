var HelloBlockchainWorld = artifacts.require("./HelloBlockchainWorld.sol");

module.exports = function(deployer) {
  deployer.deploy(HelloBlockchainWorld);
};
