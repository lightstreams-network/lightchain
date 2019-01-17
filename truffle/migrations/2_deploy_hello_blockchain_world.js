require('dotenv').config({path: `${process.env.PWD}/.env`});

const HelloBlockchainWorld = artifacts.require("./HelloBlockchainWorld.sol");

module.exports = (deployer) => {
  const from = process.env.ROOT_ACCOUNT;
  console.log("Deploying `HelloBlockchainWorld.sol` ...");
  deployer.deploy(HelloBlockchainWorld, { from });
};
