require('dotenv').config({ path: `${process.env.PWD}/.env` });

const { waitForAccountToUnlock, extractEnvAccountAndPwd } = require('../test/utils');
const Migrations = artifacts.require("./Migrations.sol");

module.exports = async (deployer) => {
  process.env.NETWORK = deployer.network;
  const { from, pwd } = extractEnvAccountAndPwd(deployer.network);

  let isUnlock = await web3.eth.personal.unlockAccount(from, pwd, 1000);
  if (!isUnlock) {
    console.error(`Account ${from} could not be unlock`);
    process.exit(1);
  }

  console.log(`Account ${from} was unlocked successfully.`);
  console.log("Deploying `Migrations.sol` ...");

  deployer.deploy(Migrations);
};