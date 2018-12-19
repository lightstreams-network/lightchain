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

        assert.equal(expectedBalance, toBalancePostTx, "recipient account balance is incorrect");
        assert.equal(21000, receipt.gasUsed, "transfer should consume fixed amount of gas for security purposes");
        assert.equal(true, receipt.status, "tx receipt should return a successful status")
    });
});