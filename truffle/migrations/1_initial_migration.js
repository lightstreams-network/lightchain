const Migrations = artifacts.require("./Migrations.sol");
require('dotenv').config({path: `${process.env.PWD}/.env`});

module.exports = (deployer, network, accounts) => {
  const pwd = process.env.PASSPHRASE;
  web3.personal.unlockAccount(accounts[1], pwd, 10000);
  deployer.deploy(Migrations);
};
