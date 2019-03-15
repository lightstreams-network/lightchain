require('dotenv').config({path: `${process.env.PWD}/.env`});

module.exports.isAccountLocked = async (address) => {
    return new Promise((resolve) => {
        web3.eth.sendTransaction({
            from: address,
            to: address,
            value: 0
        }, function(err, txHash) {
            resolve(err !== undefined && err.message === "Returned error: authentication needed: password or unlock")
        });
    });
};

const waitFor = (waitInSeconds) => {
  return new Promise((resolve) => {
    setTimeout(resolve, waitInSeconds * 1000);
  });
};
module.exports.waitFor = waitFor;

module.exports.convertPhtToWeiBN = function(pht) {
  return web3.utils.toBN(web3.utils.toWei(pht, 'ether'));
};

module.exports.calculateGasCostBN = async function(gasAmount) {
  return new Promise(resolve => {
      web3.eth.getGasPrice().then(function (gasPrice) {
          resolve(web3.utils.toBN(gasPrice).mul(web3.utils.toBN(gasAmount)));
      });
  })
};

module.exports.minimumGasPriceBN = function() {
  return web3.utils.toBN("500000000000");
};

module.exports.toBN = function(wei) {
  return web3.utils.toBN(wei);
};

module.exports.extractEnvAccountAndPwd = (network) => {
  if (network === "sirius") {
    return {
      from: process.env.SIRIUS_ACCOUNT,
      pwd: process.env.SIRIUS_PASSPHRASE
    }
  }
  if (network === "standalone") {
    return {
      from: process.env.STANDALONE_ACCOUNT,
      pwd: process.env.STANDALONE_PASSPHRASE
    }
  }
  
  if (network === "mainnet") {
    return {
      from: process.env.STANDALONE_ACCOUNT,
      pwd: process.env.STANDALONE_PASSPHRASE
    }
  }

  console.error("unknown network " + network);
  throw Error("undefined network to deploy to");
};

module.exports.timeTravel = (time) => {
  return new Promise((resolve, reject) => {
    web3.currentProvider.sendAsync({
      jsonrpc: "2.0",
      method: "evm_increaseTime",
      params: [time], // 86400 is num seconds in day
      id: new Date().getTime()
    }, (err, result) => {
      if (err) {
        return reject(err)
      }
      return resolve(result)
    });
  })
};

module.exports.waitForAccountToUnlock = function(address, timeoutInSec = 10) {
  const startTime = new Date();
  const retryInSec = 2;

  return new Promise(async (resolve, reject) => {
    while ( true ) {
      let isUnlock = isAccountLocked(address);
      if (isUnlock) {
        resolve(txReceipt);
        return;
      }

      const now = new Date();
      if (now.getTime() - startTime.getTime() > timeoutInSec * 1000) {
        reject(`Timeout after ${timeoutInSec} seconds`);
        return;
      }

      await waitFor(retryInSec)
    }
  });
};
