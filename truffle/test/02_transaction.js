/**
 * - Send funds from one wallet to another cost fixed gas
 * - Part of gas is spent after successful trx
 * - Entire sent gas is wasted after a invalid trx
 * - Validate transaction receipts : Status, From, To, cumulativeGasUsed, Logs (BlockNumber, BlockHash, ...)
 * - Validate Encoded data
 */

const { convertFromWeiBnToPht, convertPhtToWeiBN, fetchTxReceipt, calculateGasCostBN, extractEnvAccountAndPwd } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

contract('TestTransaction', async () => {
  let ROOT_ACCOUNT;
  let NEW_ACCOUNT_ADDR;
  const NEW_ACCOUNT_PASS = "password";

  it("setup", async () => {
      const account = await extractEnvAccountAndPwd(process.env.NETWORK);
      ROOT_ACCOUNT = account.from;
  });

  it('should fail transfer on insufficient funds', async function() {
    const instance = await HelloBlockchainWorld.deployed();
    NEW_ACCOUNT_ADDR = await web3.personal.newAccount(NEW_ACCOUNT_PASS);
    await web3.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, 10);

    const errMsg = 'insufficient funds for gas * price + value';
    try {
      await web3.eth.sendTransaction({
        from: NEW_ACCOUNT_ADDR,
        to: instance.address,
        value: convertPhtToWeiBN(1)
      });
      assert.equal(true, false, `Transaction should fails due to "${errMsg}"`)
    } catch ( e ) {
      assert.equal(e.message, errMsg);
    }
  });

  it('should fail transaction because gas limit is set too low', async function() {
    const weiBalancePreTxBN = await web3.eth.getBalance(ROOT_ACCOUNT);
    const weiAmountSentBN = convertPhtToWeiBN(0.1);
    const errMsg = 'intrinsic gas too low';
    const gasLimit = 20999; // Min required is 21000
    try {
      const txReceiptId = await web3.eth.sendTransaction({
        from: ROOT_ACCOUNT,
        to: NEW_ACCOUNT_ADDR,
        value: weiAmountSentBN,
        gas: gasLimit
      });
      await fetchTxReceipt(txReceiptId, 15);
      assert(false, `Transaction should fails due to "${errMsg}"`)
    } catch ( e ) {
      assert.equal(e.message, errMsg);
    }

    const weiBalancePostTxBN = await web3.eth.getBalance(ROOT_ACCOUNT);
    assert.equal(weiBalancePostTxBN.toNumber(), weiBalancePreTxBN.toNumber(), 'No gas is wasted as trx failed')
  });

  it('should transfer 0.1 PHT and gas is spent in the transaction', async function() {
    const amountToSend = convertPhtToWeiBN(0.1);
    const sender = ROOT_ACCOUNT;
    const senderBalancePreTxBN = await web3.eth.getBalance(sender);

    const txReceiptId = await web3.eth.sendTransaction({
      from: sender,
      to: NEW_ACCOUNT_ADDR,
      value: amountToSend
    });

    const txReceipt = await fetchTxReceipt(txReceiptId, 15);
    const expectedGasUsed = 21000;
    const expectedStatus = '0x1';

    assert.equal(txReceipt.gasUsed, expectedGasUsed, 'transfer should consume fixed amount of gas for security purposes');
    assert.equal(txReceipt.status, expectedStatus, 'tx receipt should return a successful status');

    const recipientBalancePostTxBN = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
    assert.equal(recipientBalancePostTxBN.toNumber(), amountToSend.toNumber(), "recipient account balance is incorrect");

    const expectedSenderBalance = senderBalancePreTxBN.sub(amountToSend.add(calculateGasCostBN(txReceipt.gasUsed)));
    const senderBalancePostTxBN = await web3.eth.getBalance(sender);

    assert.equal(senderBalancePostTxBN.toNumber(), expectedSenderBalance.toNumber(), "from account balance is incorrect")
  });
});
