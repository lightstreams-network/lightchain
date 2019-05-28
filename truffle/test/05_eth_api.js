/**
 * - Execute 100 transaction in parallel
 */

const Web3 = require('web3');
const { extractEnvAccountAndPwd } = require('./utils');
const HelloBlockchainWorld = artifacts.require("HelloBlockchainWorld");

describe('Ethereum API', () => {
  let _web3;
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  
  it("should initialize library web3", async () => {
    _web3 = new Web3(web3._provider.host, null, {
      defaultAccount: ROOT_ACCOUNT,
      defaultGasPrice: "500000000000"
    });
  });
  
  it("should assert whether network version does not match Genesis chainId", async () => {
    const networkVersion = await _web3.eth.net.getId();
    if (process.env.NETWORK === 'mainnet') {
      assert.equal(networkVersion, "163", "Network version is not expected one");
    } else if (process.env.NETWORK === 'sirius') {
      assert.equal(networkVersion, "162", "Network version is not expected one");
    } else if (process.env.NETWORK === 'standalone') {
      assert.equal(networkVersion, "161", "Network version is not expected one");
    } else {
      assert.equal(true, false, "Invalid selected network");
    }
  });
  
  it("should assert whether deployment does not successfully complete using estimated gas", async () => {
    const myContract = new _web3.eth.Contract(HelloBlockchainWorld.abi);
    const estimatedGas = await myContract.deploy({
      data: HelloBlockchainWorld.bytecode,
      arguments: []
    }).estimateGas();
    
    const deployContract = () => {
      return new Promise((resolve, reject) => {
        myContract.deploy({
          data: HelloBlockchainWorld.bytecode,
          arguments: []
        }).send({ gas: estimatedGas })
            .on('error', (err) => {
              reject(err)
            })
            .on('confirmation', (confirmationNumber, receipt) => {
              resolve(receipt);
            })
      });
    };

    const receipt = await deployContract();
    assert.equal(receipt.status, true, "Gas estimation did not provide correct value.");
  });
});
