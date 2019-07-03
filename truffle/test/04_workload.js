/**
 * - Execute 100 transaction in parallel
 */

const { convertPhtToWeiBN, fetchTxReceipt, extractEnvAccountAndPwd, calculateGasCostBN, waitFor } = require('./utils');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

const sendAsyncTx = (from, to, nonce, amountToSendBN) => {
    return new Promise(async (resolve, reject) => {
        web3.eth.sendTransaction({
            nonce: web3.utils.toHex(nonce),
            gasLimit: web3.utils.toHex('21000'),
            from: from,
            to: to,
            value: amountToSendBN.toString(),
        }, async (err, hash) => {
            if (err != null) {
                reject(err);
                return;
            }

            try {
                const receipt = await fetchTxReceipt(hash, 5);
                resolve(receipt)
            } catch ( err ) {
                reject(err)
            }
        });
    })
};

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
            value: convertPhtToWeiBN("1")
        });

        assert.equal(txReceipt.status, "0x1", 'tx receipt should return a successful status');
    });

    it("should send batch of successful parallel TXs", async () => {
        const weiBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));
        const weiAmountSentBN = convertPhtToWeiBN("0.1");
        const iterations = 100;
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

        let maxWaitTime = Date.now() + (iterations/15) * 1000; // Max wait of 20 seconds
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
        let successfulReceiptCount = 0;
        let failedReceiptCount = 0;

        const sendTx1 = sendAsyncTx(NEW_ACCOUNT_ADDR, ROOT_ACCOUNT, nonce + 1, amountToSendBN);
        const sendTx2 = sendAsyncTx(NEW_ACCOUNT_ADDR, ROOT_ACCOUNT, nonce, amountToSendBN);

        try {
            const receipt = await sendTx1;
            if (receipt.status === true) successfulReceiptCount++;
            else failedReceiptCount++;
        } catch ( err ) {
            console.error(err);
            failedReceiptCount++;
        }

        try {
            const receipt = await sendTx2;
            if (receipt.status === true) successfulReceiptCount++;
            else failedReceiptCount++;
        } catch ( err ) {
            console.error(err);
            failedReceiptCount++;
        }

        assert.equal(successfulReceiptCount, 2, "both txs should have succeed");
        assert.equal(failedReceiptCount, 0, "not tx should have failed");
    });
    
    it("should send two TXs in inverse nonce order", async () => {
        const nonce = await web3.eth.getTransactionCount(NEW_ACCOUNT_ADDR);
        const amountToSendBN = convertPhtToWeiBN("0.1");
        let successfulReceiptCount = 0;
        let failedReceiptCount = 0;

        const sendTx1 = sendAsyncTx(NEW_ACCOUNT_ADDR, ROOT_ACCOUNT, nonce + 1, amountToSendBN);

        try {
            const receipt = await sendTx1;
            if (receipt.status === true) successfulReceiptCount++;
            else failedReceiptCount++;
        } catch ( err ) {
            console.error(err);
            failedReceiptCount++;
        }

        const sendTx2 = sendAsyncTx(NEW_ACCOUNT_ADDR, ROOT_ACCOUNT, nonce, amountToSendBN);

        try {
            const receipt = await sendTx2;
            if (receipt.status === true) successfulReceiptCount++;
            else failedReceiptCount++;
        } catch ( err ) {
            console.error(err);
            failedReceiptCount++;
        }

        assert.equal(successfulReceiptCount, 1, "both txs should have succeed");
        assert.equal(failedReceiptCount, 1, "not tx should have failed");
    });
});
