require('dotenv').config({ path: `${process.env.PWD}/.env` });

const Migrations = artifacts.require("./Migrations.sol");

module.exports = (deployer) => {
  console.log("Deploying `Migrations.sol`...");
  deployer.deploy(Migrations).then(() => {
    console.log("Deployment of `Migrations.sol` completed");
  });
};
