const Migrations = artifacts.require("./Migrations.sol");
require('dotenv').config({path: `${process.env.PWD}/.env`});

module.exports = (deployer) => {
  const pwd = process.env.PASSPHRASE;
  const from = process.env.ROOT_ACCOUNT;
  web3.personal.unlockAccount(from, pwd, 10000);
  deployer.deploy(Migrations, { from });
};
