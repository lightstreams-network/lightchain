const { waitFor, extractEnvAccountAndPwd } = require('../test/utils');
const HelloBlockchainWorld = artifacts.require("./HelloBlockchainWorld.sol");

module.exports = async (deployer) => {
  console.log("Deploying `HelloBlockchainWorld.sol` ...");

  deployer.deploy(HelloBlockchainWorld);
};