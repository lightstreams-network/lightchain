/**
 * - Deploy a new smart contract
 * - Cannot call private method
 * - Payable methods receive funds correctly
 * - Only owner access to protected methods
 * - Test latest protection to popular attacks
 */

const { convertPhtToWeiBN, extractEnvAccountAndPwd } = require('./utils');

const Validators = artifacts.require("Validators");

contract('Governance', () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  
  let VALIDATOR1_KEY = "012C7DB9A70AA4940014A0CC279BFD18D8E1E224";
  
  let VALIDATOR1_ADDR;
  
  const COMMON_ACCOUNT_PASS = "password";

  it("should create validators account and top them up, not assertion", async() => {
    VALIDATOR1_ADDR = await web3.eth.personal.newAccount(COMMON_ACCOUNT_PASS);
    
    await web3.eth.sendTransaction({
      from: ROOT_ACCOUNT,
      to: VALIDATOR1_ADDR,
      value: convertPhtToWeiBN("5")
    });
  });

  it("should add an new validator ", async () => {
    const instance = await Validators.deployed();

    const estimatedGas = await instance.addValidator.estimateGas(VALIDATOR1_KEY, VALIDATOR1_ADDR);
    const tx = await instance.addValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR, {
      from: ROOT_ACCOUNT,
      gas: estimatedGas
    });
    const txReceipt = tx.receipt;
    
    const validatorAddress = await instance.validatorAddress(VALIDATOR1_KEY);
    
    assert.equal(txReceipt.status, true, "successful TX status expected");
    assert.equal(VALIDATOR1_ADDR, validatorAddress);
  });

  it("should not allow to add a new validator ", async () => {
    const instance = await Validators.deployed();
    await web3.eth.personal.unlockAccount(VALIDATOR1_ADDR, COMMON_ACCOUNT_PASS, 1000);

    let txReceipt;
    try {
      const tx = await instance.addValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR, {
        from: VALIDATOR1_ADDR
      });
      txReceipt = tx.receipt;
    } catch (e) {
      txReceipt = e.receipt;
      if (typeof e.receipt === 'undefined') {
        assert(false, e.message)
      }
    }
    
    assert.equal(txReceipt.status, false, "successful TX status expected");
  });
  
  it("should not allow to remove a validator ", async () => {
    const instance = await Validators.deployed();
    await web3.eth.personal.unlockAccount(VALIDATOR1_ADDR, COMMON_ACCOUNT_PASS, 1000);

    let txReceipt;
    try {
      const tx = await instance.removeValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR, {
        from: VALIDATOR1_ADDR
      });
      txReceipt = tx.receipt;
    } catch (e) {
      txReceipt = e.receipt;
      if (typeof e.receipt === 'undefined') {
        assert(false, e.message)
      }
    }
    
    assert.equal(txReceipt.status, false, "successful TX status expected");
  });
  
  it("should remove validator", async () => {
    const instance = await Validators.deployed();

    const tx = await instance.removeValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR, {
      from: ROOT_ACCOUNT
    });
    const txReceipt = tx.receipt;
    
    const validatorAddress = await instance.validatorAddress(VALIDATOR1_KEY);
    
    assert.equal(txReceipt.status, true, "successful TX status expected");
    assert.equal(validatorAddress, "0x0000000000000000000000000000000000000000");
  });
});
