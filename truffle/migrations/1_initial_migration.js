require('dotenv').config({ path: `${process.env.PWD}/.env` });

const { waitForAccountToUnlock, extractEnvAccountAndPwd } = require('../test/utils');
const Migrations = artifacts.require("./Migrations.sol");

module.exports = (deployer) => {
  process.env.NETWORK = deployer.network;
  const { from, pwd } = extractEnvAccountAndPwd(deployer.network);

  let isUnlock = web3.personal.unlockAccount(from, pwd, 1000);
  if (!isUnlock) {
    console.error(`Account ${from} could not be unlock`);
    process.exit(1);
  }

  // waitForAccountToUnlock(from).then(() => {
  console.log(`Account ${from} was unlocked successfully.`);
  console.log("Deploying `Migrations.sol` ...");
  deployer.deploy(Migrations, { overwrite: false }).then(() => {
    console.log("Deployment of `Migrations.sol` completed");
  });
  // }).catch((err) => {
  //   console.error(err)
  // });
};
