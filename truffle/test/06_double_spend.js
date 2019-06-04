
const { convertPhtToWeiBN, calculateGasCostBN, extractEnvAccountAndPwd, fetchTxReceipt } = require('./utils');

describe('TestSimplifiedUserBalanceDoubleSpend', () => {
    let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
    let DOUBLE_SPEND_ACC_PWD;
    let DOUBLE_SPEND_ACC;

    it('should transfer 5 PHTs to an account that will be used for double spend tests', async function() {
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

        const signedTx1 = await web3.eth.signTransaction({
            nonce: '0x0',
            gasPrice: '0x746A528800',
            gasLimit: '0x5208',
            from: DOUBLE_SPEND_ACC,
            to: ROOT_ACCOUNT,
            value: amountToSendBN.toString(),
            data: ''
        }, DOUBLE_SPEND_ACC_PWD);

        const signedTx2 = await web3.eth.signTransaction({
            nonce: '0x1',
            gasPrice: '0x746A528800',
            gasLimit: '0x5208',
            from: DOUBLE_SPEND_ACC,
            to: ROOT_ACCOUNT,
            value: amountToSendBN.toString(),
            data: ''
        }, DOUBLE_SPEND_ACC_PWD);

        const sendSignedTx1 = new Promise(function(resolve, reject) {
            web3.eth.sendSignedTransaction(signedTx1.raw, function (err, hash) {
                if (err != null) {
                    resolve();
                } else {
                    fetchTxReceipt(hash, 10).then(function (receipt) {
                        resolve(receipt);
                    }).catch(function (err) {
                        reject();
                    });
                }
            });
        });
        const sendSignedTx2 = new Promise(function(resolve, reject) {
            web3.eth.sendSignedTransaction(signedTx1.raw, function (err, hash) {
                if (err != null) {
                    resolve();
                } else {
                    fetchTxReceipt(hash, 10).then(function (receipt) {
                        resolve(receipt);
                    }).catch(function (err) {
                        reject();
                    });
                }
            });
        });

        const txReceipts = await Promise.all([sendSignedTx1, sendSignedTx2]);

        let receiptsFoundCount = 0;
        txReceipts.forEach(function (receipt) {
            if (receipt === undefined) {
            } else {
                receiptsFoundCount++;
                assert.equal(receipt.status, "0x1", 'tx receipt should return a successful status');
            }
        });

        assert.equal(receiptsFoundCount, 1, "only 1 tx should succeed and have a valid receipt");

        const rootBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(ROOT_ACCOUNT));

        assert.equal(
            rootBalancePostTxBN.toString(),
            rootBalancePreTxBN.add(amountToSendBN).toString(),
            "root account should have received exactly 4 PHTs ignoring malicious second double spend TX"
        );

        const gasCostBN = await calculateGasCostBN(21000);
        const doubleSpendAccBalancePostTxBN = web3.utils.toBN(await web3.eth.getBalance(DOUBLE_SPEND_ACC));
        const expectedDoubleSpendAccBalancePostTxBN = doubleSpendAccBalancePreTxBN.sub(gasCostBN).sub(amountToSendBN);

        assert.equal(doubleSpendAccBalancePostTxBN.toString(), expectedDoubleSpendAccBalancePostTxBN.toString(), "funds from 1 valid TX should be deducted");
    });
});
