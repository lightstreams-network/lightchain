/**
 * - Send funds from one wallet to another cost fixed gas
 * - Part of gas is spent after successful trx
 * - Entire sent gas is wasted after a invalid trx
 * - Validate transaction receipts : Status, From, To, cumulativeGasUsed, Logs (BlockNumber, BlockHash, ...)
 * - Validate Encoded data
 */

const { isAccountLocked, convertPhtToWeiBN, calculateGasCostBN, minimumGasPriceBN, extractEnvAccountAndPwd, toBN } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

describe('TestTransaction', () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  let NEW_ACCOUNT_ADDR;
  const NEW_ACCOUNT_PASS = "password";

  it('should fail transfer on insufficient funds', async function() {
    const expectedWeb3ErrMsg = 'Returned error: insufficient funds for gas * price + value';
    const instance = await HelloBlockchainWorld.deployed();
    NEW_ACCOUNT_ADDR = await web3.eth.personal.newAccount(NEW_ACCOUNT_PASS);
    await web3.eth.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, 10);

    try {
      await web3.eth.sendTransaction({
        from: NEW_ACCOUNT_ADDR,
        to: instance.address,
        value: convertPhtToWeiBN("1")
      });
    } catch (e) {
      assert.equal(e.message, expectedWeb3ErrMsg);
    }
  });

  it('should fail transaction because gas limit is set too low', async function() {
    const expectedWeb3ErrMsg = 'Returned error: intrinsic gas too low';
    const weiBalancePreTx = await web3.eth.getBalance(ROOT_ACCOUNT);
    const weiAmountSentBN = convertPhtToWeiBN("0.1");
    const gasLimit = 20999; // Min required is 21000

    try {
      await web3.eth.sendTransaction({
        from: ROOT_ACCOUNT,
        to: NEW_ACCOUNT_ADDR,
        value: weiAmountSentBN,
        gas: gasLimit
      });
    } catch (e) {
      assert.equal(e.message, expectedWeb3ErrMsg);
    }

    const weiBalancePostTx = await web3.eth.getBalance(ROOT_ACCOUNT);
    assert.equal(weiBalancePostTx, weiBalancePreTx, 'No gas was spent as TX failed');
  });

  // Uncomment once #70
  // it('should fail TX because gas price is set too low', async function() {
  //   const expectedWeb3ErrMsg = 'Timeout after 15 seconds';
  //   const gasLimit = 21000;
  //   const recipient = NEW_ACCOUNT_ADDR;
  //   const weiAmountSentBN = convertPhtToWeiBN(0.1);
  //   const requiredGasPrice = minimumGasPriceBN();
  //   const lowGasPrice = requiredGasPrice.sub(1);
  //   const txGasCost = calculateGasCostBN(gasLimit);
  //   const weiBalancePreTxBN = await web3.eth.getBalance(ROOT_ACCOUNT);
  //   const recipientBalancePreTxBN = await web3.eth.getBalance(recipient);
  //
  //   try {
  //     const txReceiptId = await web3.eth.sendTransaction({
  //       from: ROOT_ACCOUNT,
  //       to: recipient,
  //       value: weiAmountSentBN,
  //       gas: gasLimit,
  //       gasPrice: lowGasPrice,
  //     });
  //
  //     console.log(txReceipt);
  //
  //     assert.equal(undefined, txReceipt, "TX with insufficient gas price should not be added to a block");
  //   } catch (e) {
  //     assert.equal(e.message, expectedWeb3ErrMsg)
  //   }
  //
  //   const weiBalancePostTxBN = await web3.eth.getBalance(ROOT_ACCOUNT);
  //   const recipientBalancePostTxBN = await web3.eth.getBalance(recipient);
  //   assert.equal(recipientBalancePostTxBN.toNumber(), recipientBalancePreTxBN.toNumber(), "recipient should not have received the funds");
  //   assert.equal(weiBalancePostTxBN.toNumber(), weiBalancePreTxBN.toNumber(), 'No gas should be spent on invalid TX')
  // });

  it('should transfer 5 PHT and gas is spent in the transaction', async function() {
    const amountToSendBN = convertPhtToWeiBN("0.1");
    const sender = ROOT_ACCOUNT;
    const recipient = NEW_ACCOUNT_ADDR;
    
    const senderBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(sender));
    const recipientBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(recipient));

    const txReceipt = await web3.eth.sendTransaction({
      from: sender,
      to: recipient,
      value: amountToSendBN
    });

    const expectedGasUsed = 21000;
    const expectedStatus = '0x1';

    assert.equal(txReceipt.gasUsed, expectedGasUsed, 'transfer should consume fixed amount of gas for security purposes');
    assert.equal(txReceipt.status, expectedStatus, 'tx receipt should return a successful status');

    const expectedRecipientBalanceBN = recipientBalancePreTxBN.add(amountToSendBN);
    const recipientBalancePostTxBN = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
    assert.equal(recipientBalancePostTxBN.toString(), expectedRecipientBalanceBN.toString(), "recipient account balance is incorrect");

    const gasUsedCostBN = await calculateGasCostBN(txReceipt.gasUsed);
    const expectedSenderBalanceBN = senderBalancePreTxBN.sub(amountToSendBN).sub(gasUsedCostBN);
    const senderBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(sender));

    assert.equal(senderBalancePostTxBN.toString(), expectedSenderBalanceBN.toString(), "from account balance is incorrect")
  });
});