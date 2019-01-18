require('dotenv').config({path: `${process.env.PWD}/.env`});

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

module.exports.fetchTxReceipt = function(txReceiptId, timeoutInSec = 30) {
  const startTime = new Date();
  const retryInSec = 2;

  const waitTime = function(waitInSec) {
    return new Promise((resolve) => {
      setTimeout(resolve, waitInSec * 1000)
    })
  };

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

      await waitTime(retryInSec)
    }
  });
};

module.exports.waitFor = (waitInSeconds) => {
  return new Promise((resolve) => {
    setTimeout(resolve, waitInSeconds * 1000);
  });
};

module.exports.extractEnvAccountAndPwd = (network) => {
  return new Promise((resolve, reject) => {
      const account = {
        from: "",
        pwd: "",
      };

      if (network === "sirius") {
          account.from = process.env.SIRIUS_ACCOUNT;
          account.pwd = process.env.SIRIUS_PASSPHRASE;

          resolve(account);
      } else if (network === "standalone") {
          account.from = process.env.STANDALONE_ACCOUNT;
          account.pwd = process.env.STANDALONE_PASSPHRASE;

          resolve(account);
      } else {
          console.log("unknown network " + network);
          reject("undefined network to deploy to");
      }
  });
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
