/**
 * - Execute 100 transaction in parallel
 */

const { convertPhtToWeiBN, calculateGasCostBN, extractEnvAccountAndPwd, waitFor } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

describe('Workload', () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  let NEW_ACCOUNT_ADDR;
  const NEW_ACCOUNT_PASS = "password";

  it("should create an account for testing purposes, not asserting", async () => {
    NEW_ACCOUNT_ADDR = await web3.eth.personal.newAccount(NEW_ACCOUNT_PASS);
    await web3.eth.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, 1000);
  });

  // // This Test is wasting 0.231 PHT from faucet account per execution
  it("should return 100 tx receipts whose state is 0x1", async () => {
    const weiBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));
    const weiAmountSentBN = convertPhtToWeiBN("0.1");
    const iterations = 100;
    const gasLimit = 21000;
    const sentFundTxReceipt = Array();

    // It runs every txs in parallel
    for ( let i = 0; i < iterations; i++ ) {
      web3.eth.sendTransaction({
        from: ROOT_ACCOUNT,
        to: NEW_ACCOUNT_ADDR,
        value: weiAmountSentBN,
        gas: gasLimit
      }).on('confirmation', function(confirmationNumber, txReceipt) {
        sentFundTxReceipt.push(txReceipt);
      }).on('error', (error) => {
        console.error(error);
        sentFundTxReceipt.push(error);
        assert.equal(true, false, error)
      });
    }

    let maxLoopIt = 10;
    do {
      await waitFor(1);
      --maxLoopIt;
    } while ( sentFundTxReceipt.length < iterations && maxLoopIt > 0 );

    while ( sentFundTxReceipt.length > 0 ) {
      const txReceipt = sentFundTxReceipt.pop();
      assert.equal(txReceipt.status, true, "successful TX status expected");
    }

    const weiActualBalanceBN = web3.utils.toBN(await web3.eth.getBalance(NEW_ACCOUNT_ADDR));
    const weiExpectedBalanceBN = weiAmountSentBN.mul(web3.utils.toBN(iterations));
    assert.equal(weiActualBalanceBN.toString(), weiExpectedBalanceBN.toString(), "incorrect expected balance");

    const weiBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));
    const usedGasCost = await calculateGasCostBN(gasLimit * iterations);
    const weiExpectedBalancePostTxBN = weiBalancePreTxBN.sub(usedGasCost).sub(weiExpectedBalanceBN).toString();
    assert.equal(weiBalancePostTxBN.toString(), weiExpectedBalancePostTxBN.toString(), 'ROOT account balance is incorrect')
  });

  it("should refund received tokens", async () => {
    const gasLimit = 21000;

    // Send back funds to faucet account less gas spent in send tx
    const weiActualBalanceBN = web3.utils.toBN(await web3.eth.getBalance(NEW_ACCOUNT_ADDR));
    const gasCostBN = await calculateGasCostBN(gasLimit);
    const txReceipt = await web3.eth.sendTransaction({
      from: NEW_ACCOUNT_ADDR,
      to: ROOT_ACCOUNT,
      value: weiActualBalanceBN.sub(gasCostBN),
      gas: gasLimit
    });

    assert.equal(txReceipt.status, true, "successful TX status expected");
  })
});
