require('dotenv').config({path: `${process.env.PWD}/.env`});

const isAccountLocked = async (address) => {
  try {
    web3.eth.sendTransaction({
      from: address,
      to: address,
      value: 0
    });
    return false;
  } catch ( err ) {
    return (err.message === "authentication needed: password or unlock");
  }
};
module.exports.isAccountLocked = isAccountLocked;

const waitFor = (waitInSeconds) => {
  return new Promise((resolve) => {
    setTimeout(resolve, waitInSeconds * 1000);
  });
};
module.exports.waitFor = waitFor;

module.exports.convertFromWeiBnToPht = function(bn) {
  return Number(web3._extend.utils.fromWei(bn.toNumber(), 'ether'));
};

module.exports.convertPhtToWeiBN = function(ether) {
  const etherInWei = web3._extend.utils.toWei(ether, 'ether');
  return web3._extend.utils.toBigNumber(etherInWei);
};

module.exports.calculateGasCostBN = function(gasAmount) {
  return web3.eth.gasPrice.mul(gasAmount);
};

module.exports.minimumGasPriceBN = function() {
  return web3._extend.utils.toBigNumber("500000000000");
};

module.exports.toBN = function(wei) {
  return web3._extend.utils.toBigNumber(wei);
};

module.exports.fetchTxReceipt = function(txReceiptId, timeoutInSec = 30) {
  const startTime = new Date();
  const retryInSec = 2;

  return new Promise(async (resolve, reject) => {
    while ( true ) {
      let txReceipt = web3.eth.getTransactionReceipt(txReceiptId);
      if (txReceipt != null && typeof txReceipt !== 'undefined') {
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
