var HelloBlockchainWorld = artifacts.require("./HelloBlockchainWorld.sol");

module.exports = function(deployer) {
  const from = process.env.ROOT_ACCOUNT;
  deployer.deploy(HelloBlockchainWorld, { from });
};
