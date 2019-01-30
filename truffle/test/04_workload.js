/**
 * - Deploy a new smart contract
 * - Cannot call private method
 * - Payable methods receive funds correctly
 * - Only owner access to protected methods
 * - Test latest protection to popular attacks
 */

const { convertFromWeiBnToPht, convertPhtToWeiBN, fetchTxReceipt, calculateGasCostBN, extractEnvAccountAndPwd } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

contract('Workload', async () => {
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

  it("should return 100 tx receipts whose state is 0x1", async () => {
    const weiAmountSentBN = convertPhtToWeiBN(0.1);
    const iterations = 100;
    const gasLimit = 21000;
    const sentFundTxReceiptPromises = Array();
    const sentFundTxReceiptIds = Array();
    
    // It runs every txs in parallel
    for(let i=0; i < iterations; i++) {
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

    // Send back funds to faucet account less gas spent in txs
    const usedGasCost = calculateGasCostBN(gasLimit * (iterations+1));
    const txReceiptId = await web3.eth.sendTransaction({
      from: NEW_ACCOUNT_ADDR,
      to: ROOT_ACCOUNT,
      value: weiAmountSentBN.mul(iterations).sub(usedGasCost),
      gas: gasLimit
    });
    
    const txReceipt = await fetchTxReceipt(txReceiptId, 15);
    assert.equal(txReceipt.status, "0x1", "successful TX status expected");
    
    console.log(`Gas used in total: ${convertFromWeiBnToPht(usedGasCost)} PHT`);
  });
});
