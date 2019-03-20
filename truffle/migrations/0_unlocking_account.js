require('dotenv').config({ path: `${process.env.PWD}/.env` });

const { extractEnvAccountAndPwd } = require('../test/utils');

module.exports = (deployer) => {
  process.env.NETWORK = deployer.network;
  const { from, pwd } = extractEnvAccountAndPwd(deployer.network);

  deployer.then(function() {
    return web3.eth.personal.unlockAccount(from, pwd, 1000)
      .then(console.log('Account unlocked!'))
      .catch((err) => {
        console.log(err);
      });
  });
};
