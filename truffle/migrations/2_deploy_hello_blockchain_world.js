const HelloBlockchainWorld = artifacts.require("./HelloBlockchainWorld.sol");

module.exports = (deployer) => {
  console.log("Deploying `HelloBlockchainWorld.sol` ...");
  deployer.deploy(HelloBlockchainWorld, { overwrite: false }).then(() => {
    console.log("Deployment of `HelloBlockchainWorld.sol` completed");
  });
};
