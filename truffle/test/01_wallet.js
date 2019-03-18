/**
 * - Create a new wallet, balance is 0
 * - Send funds from one wallet to another
 * - Double spend (one local tendermint node)
 * - Double spend (on two tendermint nodes)
 */

const { isAccountLocked, convertPhtToWeiBN, waitFor, extractEnvAccountAndPwd } = require('./utils');

describe('WalletTest', () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  let NEW_ACCOUNT_ADDR;
  let NEW_ACCOUNT_PASS = "password";

  it("should create a new account with balance 0 and locked", async () => {
    NEW_ACCOUNT_ADDR = await web3.eth.personal.newAccount(NEW_ACCOUNT_PASS);
    const balanceWei = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
    const isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);

    assert.equal(balanceWei, "0", "New accounts should have balance 0");
    assert.equal(isLocked, true, "New account should be created as locked");
  });

  it("should automatically get locked when session expires", async () => {
    let isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);
    assert.equal(isLocked, true, "Account should be locked");

    const sessionDurationInSec = 1;
    await web3.eth.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, sessionDurationInSec);
    isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);
    assert.equal(isLocked, false, "Account should be unlocked");

    await waitFor(sessionDurationInSec);
    isLocked = await isAccountLocked(NEW_ACCOUNT_ADDR);
    assert.equal(isLocked, true, "Account should be locked cause session has expired");
  });

  it('should be able to transfer funds', async function() {
    const weiSentAmountBN = convertPhtToWeiBN("0.1");
    const weiBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(NEW_ACCOUNT_ADDR));

    await web3.eth.sendTransaction({
      from: ROOT_ACCOUNT,
      to: NEW_ACCOUNT_ADDR,
      value: weiSentAmountBN
    });

    const weiBalancePostTxBN = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);

    assert.equal(
        weiBalancePostTxBN.toString(),
        weiBalancePreTxBN.add(weiSentAmountBN).toString(),
        "recipient account balance is incorrect"
    );
  });
});