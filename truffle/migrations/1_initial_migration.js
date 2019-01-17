const { waitFor } = require('../test/utils');
require('dotenv').config({path: `${process.env.PWD}/.env`});

const Migrations = artifacts.require("./Migrations.sol");

module.exports = async (deployer) => {
  const pwd = process.env.PASSPHRASE;
  const from = process.env.ROOT_ACCOUNT;

  let isUnlock = web3.personal.unlockAccount(from, pwd, 1000);
  if (!isUnlock) {
    console.error(`Account ${from} could not be unlock`);
    process.exit(1);
  }
  console.error(`Account ${from} was unlock successfully`);
  console.log("Deploying `Migrations.sol` ...");
  await waitFor(1);
  deployer.deploy(Migrations, { from });
};
