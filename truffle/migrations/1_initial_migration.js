require('dotenv').config({path: `${process.env.PWD}/.env`});

const { waitFor, extractEnvAccountAndPwd } = require('../test/utils');
const Migrations = artifacts.require("./Migrations.sol");

module.exports = async (deployer) => {
    const { from, pwd } = await extractEnvAccountAndPwd(deployer.network);

    let isUnlock = web3.personal.unlockAccount(from, pwd, 1000);
    if (!isUnlock) {
        console.error(`Account ${from} could not be unlock`);
        process.exit(1);
    }
    console.error(`Account ${from} was unlocked successfully.`);
    console.log("Deploying `Migrations.sol` ...");

    await waitFor(3);
    deployer.deploy(Migrations, { from });
};