const Web3 = require("web3");
const async = require("async");
const BN = require('bn.js');

const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");
const port = 8545;
const web3 = new Web3(new Web3.providers.HttpProvider("http://localhost:" + port));

contract('HelloBlockchainWorld', function(accounts) {
    it("should return the msg.sender when reading the owner attribute", function() {
        return HelloBlockchainWorld.deployed().then(function(instance) {
            return instance.owner.call();
        }).then(function(owner) {
            assert.equal(owner, accounts[0], "Owner doesn't match the msg.sender");
        });
    });

    it('should query account balance using web3', async function () {
        const balance = await web3.eth.getBalance(accounts[0]);

        assert(balance > 0, "Sirius blockchain coinbase account should have balance greater than 0");
    });

    it('should transfer funds between accounts', async function () {
        const amount = new BN(web3.utils.toWei("1", "ether"));
        const toBalancePreTx = await web3.eth.getBalance(accounts[1]);
        const receipt = await web3.eth.sendTransaction({from: accounts[0],to: accounts[1], value: amount});
        const toBalancePostTx = await web3.eth.getBalance(accounts[1]);
        const expectedBalance = new BN(toBalancePreTx).add(amount);

        assert.equal(toBalancePostTx, expectedBalance, "recipient account balance is incorrect");
        assert.equal(receipt.gasUsed, 21000, "transfer should consume fixed amount of gas for security purposes");
        assert.equal(receipt.status, true, "tx receipt should return a successful status");
    });

    it('should fail transfer on insufficient funds', async function () {
        const balance = await web3.eth.getBalance(accounts[0]);
        const overflowAmount = new BN(balance).add(new BN(web3.utils.toWei("1", "ether")));

        try {
            const receipt = await web3.eth.sendTransaction({from: accounts[0],to: accounts[1], value: overflowAmount});
        } catch (e) {
            assert.equal("Returned error: insufficient funds for gas * price + value", e.message);
        }
    });

    it("should fail transaction and cost it gas if contract transaction gets rejected", async function() {
        const from = accounts[1];
        const balancePreTx = await web3.eth.getBalance(from);
        let txReceipt = null;

        await HelloBlockchainWorld.deployed()
        .then(
            function(instance) {
                return instance.incrementHelloCount({ from: from });
            }
        )
        .then(
            function(response) {
            },
            function(error) {
                txReceipt = error.receipt;
            }
        );

        const gasPrice = new BN(await web3.eth.getGasPrice());
        const gasUsed = new BN(txReceipt.gasUsed);
        const txCost = gasUsed.mul(gasPrice);
        const expectedBalance = new BN(balancePreTx).sub(txCost);
        const balancePostTx = await web3.eth.getBalance(from);

        assert.equal(balancePostTx, expectedBalance);
        assert.equal(txReceipt.status, "0x0", "failed TX status expected");
    });
});