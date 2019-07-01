/**
 * - Execute 100 transaction in parallel
 */

const { convertPhtToWeiBN, calculateGasCostBN, extractEnvAccountAndPwd, waitFor } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

describe('Workload', () => {
    let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
    let NEW_ACCOUNT_ADDR;
    const NEW_ACCOUNT_PASS = "password";

    it("should create an account for testing purposes", async () => {
        NEW_ACCOUNT_ADDR = await web3.eth.personal.newAccount(NEW_ACCOUNT_PASS);
        await web3.eth.personal.unlockAccount(NEW_ACCOUNT_ADDR, NEW_ACCOUNT_PASS, 1000);

        console.log("Account to receive all workload TXs: ", NEW_ACCOUNT_ADDR);

        const txReceipt = await web3.eth.sendTransaction({
            from: ROOT_ACCOUNT,
            to: NEW_ACCOUNT_ADDR,
            value: convertPhtToWeiBN("0.2")
        });

        assert.equal(txReceipt.status, "0x1", 'tx receipt should return a successful status');
    });

    it("should send batch of successful parallel TXs", async () => {
        const weiBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));
        const weiAmountSentBN = convertPhtToWeiBN("0.1");
        const iterations = 1000;
        const gasLimit = 21000;
        const sentFundTxReceipt = Array();

        // It runs every txs in parallel
        console.debug(`Sending ${iterations} TXs in parallel...`);
        for ( let i = 0; i < iterations; i++ ) {
            web3.eth.sendTransaction({
                from: ROOT_ACCOUNT,
                to: NEW_ACCOUNT_ADDR,
                value: weiAmountSentBN,
                gas: gasLimit
            }).on('receipt', function(txReceipt) {
                sentFundTxReceipt.push(txReceipt);
                assert.equal(txReceipt.status, "0x1", "successful TX status expected");
            }).on('error', (error) => {
                console.error(error);
                sentFundTxReceipt.push(error);
                assert.equal(true, false, error)
            });
        }

        let maxWaitTime = Date.now() + 20 * 1000; // Max wait of 20 seconds
        do {
            await waitFor(1);
        } while ( sentFundTxReceipt.length < iterations && maxWaitTime > Date.now() );

        if (maxWaitTime <= Date.now()) {
            assert.fail("Wait time exceeded");
            return;
        }

        const gasCostBN = await calculateGasCostBN(gasLimit);

        // Send back funds to faucet account less gas spent in send tx
        const txReceipt = await web3.eth.sendTransaction({
            from: NEW_ACCOUNT_ADDR,
            to: ROOT_ACCOUNT,
            value: weiAmountSentBN.mul(web3.utils.toBN(iterations)).sub(gasCostBN),
            gas: gasLimit
        });

        assert.equal(txReceipt.status, "0x1", "successful TX status expected");

        const usedGasCost = await calculateGasCostBN(gasLimit * (iterations + 1));
        const weiBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));
        assert.equal(weiBalancePostTxBN.toString(), weiBalancePreTxBN.sub(usedGasCost).toString(), 'from account balance is incorrect')
    });

    it("should send two TXs in parallel with inverse nonce order", async () => {
        const nonce = await web3.eth.getTransactionCount(NEW_ACCOUNT_ADDR);
        const amountToSendBN = convertPhtToWeiBN("0.1");

        const signedTx1 = await web3.eth.signTransaction({
            nonce: nonce + 1,
            gasPrice: '500000000000',
            gasLimit: '21000',
            from: NEW_ACCOUNT_ADDR,
            to: ROOT_ACCOUNT,
            value: amountToSendBN.toString(),
            data: ''
        }, NEW_ACCOUNT_PASS);

        const signedTx2 = await web3.eth.signTransaction({
            nonce: nonce + 2,
            gasPrice: '500000000000',
            gasLimit: '21000',
            from: NEW_ACCOUNT_ADDR,
            to: ROOT_ACCOUNT,
            value: amountToSendBN.toString(),
            data: ''
        }, NEW_ACCOUNT_PASS);

        
        const sendSignedTx2 = new Promise((resolve, reject) => {
            web3.eth.sendSignedTransaction(signedTx2.raw, async (err, hash) => {
                if (err != null) {
                    reject(err);
                    return;
                }
                fetchTxReceipt(hash, txReceiptTimeout).then(function(receipt) {
                    resolve(receipt);
                }).catch(function(err) {
                    reject(err);
                });
            });
        });
        
        await waitFor(0.3);

        const sendSignedTx1 = new Promise((resolve, reject) => {
            web3.eth.sendSignedTransaction(signedTx1.raw, async (err, hash) => {
                if (err != null) {
                    reject(err);
                    return;
                }
                fetchTxReceipt(hash, txReceiptTimeout).then((receipt) => {
                    resolve(receipt);
                }).catch(function(err) {
                    reject(err);
                });
            });
        });

        let successfulReceiptCount = 0;
        let failedReceiptCount = 0;

        try {
            const receipt = await sendSignedTx1;
            if (receipt.status === true) successfulReceiptCount++;
            else failedReceiptCount++;
        } catch ( err ) {
            failedReceiptCount++;
        }

        try {
            const receipt = await sendSignedTx2;
            if (receipt.status === true) successfulReceiptCount++;
            else failedReceiptCount++;
        } catch ( err ) {
            failedReceiptCount++;
        }

        assert.equal(successfulReceiptCount, 2, "only 1 tx should have succeed");
        assert.equal(failedReceiptCount, 0, "not tx should have failed");
    });
});
