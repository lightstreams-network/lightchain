/**
 * - Execute 100 transaction in parallel
 */

const { convertFromWeiBnToPht, convertPhtToWeiBN, fetchTxReceipt, calculateGasCostBN, extractEnvAccountAndPwd } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

describe('Workload', async () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  let NEW_ACCOUNT_ADDR;
  const NEW_ACCOUNT_PASS = "password";
  
  it("should create an account for testing purposes, not asserting", async () => {
    NEW_ACCOUNT_ADDR = await web3.personal.newAccount(NEW_ACCOUNT_PASS);
    await web3.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, 1000);

    const txReceiptId = await web3.eth.sendTransaction({
      from: ROOT_ACCOUNT,
      to: NEW_ACCOUNT_ADDR,
      value: convertPhtToWeiBN(1)
    });

    await fetchTxReceipt(txReceiptId, 15);
  });

  // This Test is wasting 0.231 PHT from faucet account per execution
  it("should return 100 tx receipts whose state is 0x1", async () => {
    const weiBalancePreTxBN = await web3.eth.getBalance(ROOT_ACCOUNT);
    const weiAmountSentBN = convertPhtToWeiBN(0.1);
    const iterations = 100;
    const gasLimit = 21000;
    const sentFundTxReceiptPromises = Array();
    const sentFundTxReceiptIds = Array();
    
    // It runs every txs in parallel
    for (let i=0; i < iterations; i++) {
      const txReceiptPromise = web3.eth.sendTransaction({
        from: ROOT_ACCOUNT,
        to: NEW_ACCOUNT_ADDR,
        value: weiAmountSentBN,
        gas: gasLimit
      });
      sentFundTxReceiptPromises.push(txReceiptPromise);
    }
    
    while(sentFundTxReceiptPromises.length) {
      const txReceiptId = await sentFundTxReceiptPromises.pop();
      sentFundTxReceiptIds.push(txReceiptId)
    }
    
    while(sentFundTxReceiptIds.length) {
      const txReceiptId = sentFundTxReceiptIds.pop();
      const txReceipt = await fetchTxReceipt(txReceiptId, 15);
      assert.equal(txReceipt.status, "0x1", "successful TX status expected");
    }

    // Send back funds to faucet account less gas spent in send tx
    const txReceiptId = await web3.eth.sendTransaction({
      from: NEW_ACCOUNT_ADDR,
      to: ROOT_ACCOUNT,
      value: weiAmountSentBN.mul(iterations).sub(calculateGasCostBN(gasLimit)),
      gas: gasLimit
    });
    
    const txReceipt = await fetchTxReceipt(txReceiptId, 15);
    assert.equal(txReceipt.status, "0x1", "successful TX status expected");
    
    const usedGasCost = calculateGasCostBN(gasLimit * (iterations + 1));
    const weiBalancePostTxBN = await web3.eth.getBalance(ROOT_ACCOUNT);
    assert.equal(weiBalancePostTxBN.toNumber(), weiBalancePreTxBN.sub(usedGasCost).toNumber(), 'from account balance is incorrect')
  });
});
