require('dotenv').config({ path: `${process.env.PWD}/.env` });

const { extractEnvAccountAndPwd } = require('../test/utils');
const Migrations = artifacts.require("./Migrations.sol");

module.exports = (deployer) => {
  process.env.NETWORK = deployer.network;
  const { from, pwd } = extractEnvAccountAndPwd(deployer.network);

  console.log(`Unlocking  ${from} account...`);
  web3.eth.personal.unlockAccount(from, pwd, 20)
      .then(() => {
          console.log("Deploying `Migrations.sol`...");
          deployer.deploy(Migrations, {overwrite: false})
              .then(() => {
                  console.log("Deployment of `Migrations.sol` completed");
              });
      })
      .catch((err) => {
          console.log(err);
          process.exit(1);
      });
};