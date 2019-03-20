const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld.sol");

module.exports = (deployer) => {
  console.log("Deploying `HelloBlockchainWorld.sol`...");
  deployer.deploy(HelloBlockchainWorld).then(() => {
    console.log("HelloBlockchainWorld.sol successfully deployed!");
  });
};
