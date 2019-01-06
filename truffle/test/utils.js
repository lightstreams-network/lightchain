module.exports.convertFromWeiBnToPht = function(bn) {
  return Number(web3._extend.utils.fromWei(bn.toNumber(), 'ether'));
};

module.exports.convertPhtToWeiBN = function(ether) {
  const etherInWei = web3._extend.utils.toWei(ether, 'ether');
  return web3._extend.utils.toBigNumber(etherInWei);
};

module.exports.gasToWei = function(gasAmount) {
  return gasAmount * web3.eth.gasPrice.toNumber();
};

module.exports.fetchTxReceipt = function(txReceiptId, timeoutInSec = 30) {
  const startTime = new Date();
  const retryInSec = 2;

  const waitTime = function(waitInSec) {
    return new Promise((resolve) => {
      setTimeout(resolve, waitInSec * 1000)
    })
  };

  console.log(`Fetching tx ${txReceiptId}`);

  return new Promise(async (resolve, reject) => {
    while ( true ) {
      let txReceipt = web3.eth.getTransactionReceipt(txReceiptId);
      if (txReceipt != null && typeof txReceipt !== 'undefined') {
        console.log('Receipt found');
        resolve(txReceipt);
        return;
      }

      const now = new Date();
      if (now.getTime() - startTime.getTime() > timeoutInSec * 1000) {
        console.log('Receipt not found');
        reject(`Timeout after ${timeoutInSec} seconds`);
        return;
      }

      console.log(`Now: ${now.toISOString()}, StartTime: ${startTime.toISOString()}...retrying in ${retryInSec} seconds`);
      await waitTime(retryInSec)
    }
  });
};

module.exports.waitFor = (waitInSeconds) => {
  return new Promise((resolve) => {
    setTimeout(resolve, waitInSeconds * 1000);
  });
};
