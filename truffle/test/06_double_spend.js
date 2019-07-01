const { convertPhtToWeiBN, calculateGasCostBN, extractEnvAccountAndPwd, fetchTxReceipt } = require('./utils');

describe('TestSimplifiedUserBalanceDoubleSpend', () => {
    let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
    let DOUBLE_SPEND_ACC_PWD;
    let DOUBLE_SPEND_ACC;

    it('should transfer 0.12 PHTs to an account that will be used for double spend tests', async function() {
        DOUBLE_SPEND_ACC_PWD = "doublespend";
        DOUBLE_SPEND_ACC = await web3.eth.personal.newAccount(DOUBLE_SPEND_ACC_PWD);
        await web3.eth.personal.unlockAccount(DOUBLE_SPEND_ACC, DOUBLE_SPEND_ACC_PWD, 1000);

        const txReceipt = await web3.eth.sendTransaction({
            from: ROOT_ACCOUNT,
            to: DOUBLE_SPEND_ACC,
            value: convertPhtToWeiBN("0.12")
        });

        assert.equal(txReceipt.status, "0x1", 'tx receipt should return a successful status');
    });

    it('perform simplified balance based double spend', async function() {
        const amountToSendBN = convertPhtToWeiBN("0.1");
        const rootBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));
        const doubleSpendAccBalancePreTxBN = web3.utils.toBN(await web3.eth.getBalance(DOUBLE_SPEND_ACC));
        const txReceiptTimeout = 10;

        const signedTx1 = await web3.eth.signTransaction({
            nonce: '0',
            gasPrice: '500000000000',
            gasLimit: '21000',
            from: DOUBLE_SPEND_ACC,
            to: ROOT_ACCOUNT,
            value: amountToSendBN.toString(),
            data: ''
        }, DOUBLE_SPEND_ACC_PWD);

        const signedTx2 = await web3.eth.signTransaction({
            nonce: '1',
            gasPrice: '500000000000',
            gasLimit: '21000',
            from: DOUBLE_SPEND_ACC,
            to: ROOT_ACCOUNT,
            value: amountToSendBN.toString(),
            data: ''
        }, DOUBLE_SPEND_ACC_PWD);

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

        assert.equal(successfulReceiptCount, 1, "only 1 tx should have succeed");
        assert.equal(failedReceiptCount, 1, "only 1 tx should have failed");

        const rootBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));

        assert.equal(
            rootBalancePostTxBN.toString(),
            rootBalancePreTxBN.add(amountToSendBN).toString(),
            `root account should have received exactly ${amountToSendBN} PHTs ignoring malicious second double spend TX`
        );

        const gasCostBN = await calculateGasCostBN(21000);
        const doubleSpendAccBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(DOUBLE_SPEND_ACC));
        const expectedDoubleSpendAccBalancePostTxBN = doubleSpendAccBalancePreTxBN.sub(gasCostBN).sub(amountToSendBN);

        assert.equal(doubleSpendAccBalancePostTxBN.toString(), expectedDoubleSpendAccBalancePostTxBN.toString(), "funds from 1 valid TX should be deducted");
    });
});
