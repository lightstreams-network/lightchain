/**
 * - Create a new wallet, balance is 0
 * - Send funds from one wallet to another
 * - Double spend (one local tendermint node)
 * - Double spend (on two tendermint nodes)
 */

const { isAccountLocked, convertFromWeiBnToPht, convertPhtToWeiBN, fetchTxReceipt, waitFor, extractEnvAccountAndPwd } = require('./utils');

describe('WalletTest', async () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  let NEW_ACCOUNT_ADDR;
  let NEW_ACCOUNT_PASS = "password";

  it("should create a new account with balance 0 and locked", async () => {
    NEW_ACCOUNT_ADDR = await web3.personal.newAccount(NEW_ACCOUNT_PASS);
    const balanceWei = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
    const balance = convertFromWeiBnToPht(balanceWei);
    const isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);

    assert.equal(balance, 0, "New accounts should have balance 0");
    assert.equal(isLocked, true, "New account should be created as locked");
  });

  it("should automatically get locked when session expires", async () => {
    let isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);
    assert.equal(isLocked, true, "Account should be locked");

    const sessionDurationInSec = 1;
    await web3.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, sessionDurationInSec);
    isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);
    assert.equal(isLocked, false, "Account should be unlocked");

    await waitFor(sessionDurationInSec);
    isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);
    assert.equal(isLocked, true, "Account should be locked cause session has expired");
  });

  it('should be able to transfer funds', async function() {
    const weiSentAmountBN = convertPhtToWeiBN(0.1);

    const weiBalancePreTxBN = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
    const txReceiptId = await web3.eth.sendTransaction({
      from: ROOT_ACCOUNT,
      to: NEW_ACCOUNT_ADDR,
      value: weiSentAmountBN
    });

    await fetchTxReceipt(txReceiptId, 15);

    const weiBalancePostTxBN = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
    assert.equal(
        weiBalancePostTxBN.toNumber(),
        weiBalancePreTxBN.toNumber() + weiSentAmountBN.toNumber(),
        "recipient account balance is incorrect"
    );
  });
});
