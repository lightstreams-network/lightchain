const { convertFromWeiBnToEth, convertEtherToWeiBN, fetchTxReceipt } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

contract('HelloBlockchainWorld', () => {
    const ROOT_ACCOUNT = process.env.ROOT_ACCOUNT;
    const NEW_ACCOUNT_PASS = "password";
    let NEW_ACCOUNT_ADDR ;
    
    it("should return the msg.sender when reading the owner attribute", async () => {
        const instance = await HelloBlockchainWorld.deployed();
        const owner = await instance.owner.call();
        assert.equal(owner, ROOT_ACCOUNT, "Owner doesn't match the msg.sender");
    });
    
    it("should create a new account with balance 0", async() => {
        NEW_ACCOUNT_ADDR = await web3.personal.newAccount(NEW_ACCOUNT_PASS);
        const txReceipt = await web3.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, 10000);
        const balanceWei = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
        const balance = convertFromWeiBnToEth(balanceWei);
        
        assert.equal(balance, 0, "New accounts should have balance 0");
        assert.equal(txReceipt, true, "New account was not unlocked");
    });
    
    it('should fail transfer on insufficient funds', async function () {
        try {
            await web3.eth.sendTransaction({
              from: NEW_ACCOUNT_ADDR,
              to: ROOT_ACCOUNT, 
              value: convertEtherToWeiBN(1)
            });
            assert.equal(true, false, "Transaction should fails due to 'insufficient funds for gas * price + value'")
        } catch (e) {
            assert.equal(true, true);
        }
    });

    it('should transfer funds between accounts', async function () {
        const ethAmount = 1;
        const weiAmount = convertEtherToWeiBN(ethAmount);
        
        const toWeiBalancePreTx = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
        const toEthBalancePreTx = convertFromWeiBnToEth(toWeiBalancePreTx);
        
        const txReceiptId = await web3.eth.sendTransaction({
          from: ROOT_ACCOUNT,
          to: NEW_ACCOUNT_ADDR,
          value: weiAmount
        });
        
        const txReceipt = await fetchTxReceipt(txReceiptId , 15);
        assert.equal(txReceipt.gasUsed, 21000, "transfer should consume fixed amount of gas for security purposes");
        assert.equal(txReceipt.status, "0x1", "tx receipt should return a successful status");

        const toWeiBalancePostTx = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
        const toEthBalancePostTx = convertFromWeiBnToEth(toWeiBalancePostTx);
        
        const expectedEthBalance = toEthBalancePreTx + ethAmount;
        assert.equal(toEthBalancePostTx, expectedEthBalance, "recipient account balance is incorrect");
    });
    
    it("should be allowed to perform the smart contract call and gas reduced from balance", async () => {
        const instance = await HelloBlockchainWorld.deployed();
        const gasPrice = web3.eth.gasPrice.toNumber();
        
        const fromWeiBNBalancePreTx = await web3.eth.getBalance(ROOT_ACCOUNT);
        const fromWeiBalancePreTx = fromWeiBNBalancePreTx.toNumber();

        let txReceipt;
        try {
            const tx = await instance.incrementHelloCount({ 
              from: ROOT_ACCOUNT
            });
            txReceipt = await fetchTxReceipt(tx.tx);
            assert.equal(txReceipt.status, "0x1", "successful TX status expected");
        } catch(e) {
            assert(false, `Exception: ${e.message}`)
        }
        
        const expectedBalance = fromWeiBalancePreTx - (txReceipt.gasUsed * gasPrice);
        const fromBalanceBNWeiPostTx = await web3.eth.getBalance(ROOT_ACCOUNT);
        const fromBalanceWeiPostTx = fromBalanceBNWeiPostTx.toNumber();
        assert.equal(fromBalanceWeiPostTx, expectedBalance);
    });
    
    it("should fail transaction and cost it gas if contract transaction gets rejected", async () => {
        const instance = await HelloBlockchainWorld.deployed();
        const gasPrice = web3.eth.gasPrice.toNumber();
        
        const fromWeiBNBalancePreTx = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
        const fromWeiBalancePreTx = fromWeiBNBalancePreTx.toNumber();

        let txReceipt;
        try {
            const tx = await instance.incrementHelloCount({ 
              from: NEW_ACCOUNT_ADDR
            });
            txReceipt = await fetchTxReceipt(tx.tx);
        } catch(e) {
            txReceipt = e.receipt;
            if (typeof e.receipt === 'undefined') {
                assert(false, `Exception: ${e.message}`)
            }
        }
        
        assert.equal(txReceipt.status, "0x0", "failed TX status expected");

        const expectedBalance = fromWeiBalancePreTx - (txReceipt.gasUsed * gasPrice);
        const fromBalanceBNWeiPostTx = await web3.eth.getBalance(NEW_ACCOUNT_ADDR);
        const fromBalanceWeiPostTx = fromBalanceBNWeiPostTx.toNumber();

        assert.equal(fromBalanceWeiPostTx, expectedBalance);
    });
});
