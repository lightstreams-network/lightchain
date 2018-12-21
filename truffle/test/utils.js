module.exports.convertFromWeiBnToEth = function (bn) {
  return Number(web3._extend.utils.fromWei(bn.toNumber(), 'ether'));
};

module.exports.convertEtherToWeiBN = function (ether) {
  const etherInWei = web3._extend.utils.toWei(ether, 'ether');
  return web3._extend.utils.toBigNumber(etherInWei);
};

module.exports.fetchTxReceipt = function(txReceiptId, timeoutInSec) {
    const startTime = new Date();
    const retryInSec = 2;
    
    const waitTime = function(waitInSec) {
        return new Promise((resolve) => {
            setTimeout(resolve, waitInSec * 1000)
        })
    };

    console.log(`Fetching tx ${txReceiptId}. Timeout in ${timeoutInSec}s`);
    
    return new Promise(async (resolve, reject) => {
        while (true) {
            let txReceipt = web3.eth.getTransactionReceipt(txReceiptId);
            if (txReceipt != null && typeof txReceipt !== 'undefined'){
                console.log('Receipt found');
                resolve(txReceipt);
                return;
            }
    
            const now = new Date();
            if (now.getTime() - startTime.getTime() > timeoutInSec * 1000) {
                console.log('Receipt not found');
                reject("Timeout after 5 seconds");
                return;
            }
            
            console.log(`Now: ${now}, StartTime: ${startTime}...retrying in ${retryInSec}`);
            await waitTime(retryInSec)
        }
    });
};
