/**
 * - Deploy a new smart contract
 * - Cannot call private method
 * - Payable methods receive funds correctly
 * - Only owner access to protected methods
 * - Test latest protection to popular attacks
 */
const chai = require('chai');
chai.use(require('chai-as-promised'));
const assert = chai.assert;

const { extractEnvAccountAndPwd } = require('./utils');

const ValidatorSet = artifacts.require("ValidatorSet");

describe('Governance', () => {
  let ROOT_ACCOUNT = extractEnvAccountAndPwd(process.env.NETWORK).from;
  
  let VALIDATOR1_KEY = "0x012C7DB9A70AA4940014A0CC279BFD18D8E1E224".toLowerCase();
  let VALIDATOR2_KEY = "0x012C7DB9A70AA4940014A0CC279BFD18D8E1E225".toLowerCase();
  let VALIDATOR3_KEY = "0x012C7DB9A70AA4940014A0CC279BFD18D8E1E226".toLowerCase();
  
  let VALIDATOR1_ADDR;
  let VALIDATOR2_ADDR;
  let VALIDATOR3_ADDR;
  
  const COMMON_ACCOUNT_PASS = "password";

  it("should create validators account and top them up, not assertion", async() => {
    VALIDATOR1_ADDR = await web3.eth.personal.newAccount(COMMON_ACCOUNT_PASS);
    VALIDATOR2_ADDR = await web3.eth.personal.newAccount(COMMON_ACCOUNT_PASS);
    VALIDATOR3_ADDR = await web3.eth.personal.newAccount(COMMON_ACCOUNT_PASS);
  });

  it("should add a validator ", async () => {
    const instance = await ValidatorSet.deployed();

    const tx = await instance.addValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR);
    assert.equal(tx.receipt.status, true, "successful TX status expected");
    
    const validatorSetLengthBN = await instance.validatorSetSize.call();
    const validatorSetLength = parseInt(validatorSetLengthBN.toString());
    const validatorPubKey = await instance.validatorPubKey(validatorSetLength-1);
    const validatorAddress = await instance.validatorAddress(VALIDATOR1_KEY);
    
    assert.equal(VALIDATOR1_KEY, validatorPubKey);
    assert.equal(VALIDATOR1_ADDR, validatorAddress);
  });
  
  it("should add a another validator ", async () => {
    const instance = await ValidatorSet.deployed();

    const tx = await instance.addValidator(VALIDATOR2_KEY, VALIDATOR2_ADDR);
    assert.equal(tx.receipt.status, true, "successful TX status expected");

    const validatorSetLengthBN = await instance.validatorSetSize.call();
    const validatorSetLength = parseInt(validatorSetLengthBN.toString());
    const validatorPubKey = await instance.validatorPubKey(validatorSetLength-1);
    const validatorAddress = await instance.validatorAddress(VALIDATOR2_KEY);
    
    assert.equal(VALIDATOR2_KEY, validatorPubKey);
    assert.equal(VALIDATOR2_ADDR, validatorAddress);
  });
  
  it("should retrieve the two added validator pubkeys", async () => {
    const instance = await ValidatorSet.deployed();

    const validatorPubKeys = [];
    const validatorSetLengthBN = await instance.validatorSetSize.call();
    const validatorSetLength = parseInt(validatorSetLengthBN.toString());

    for(let i=0; i < validatorSetLength; i++) {
      const pubKey = await instance.validatorPubKey.call(i, {});
      validatorPubKeys.push(pubKey);
    }

    assert.deepEqual(validatorPubKeys, [VALIDATOR1_KEY, VALIDATOR2_KEY]);
  });

  it("should remove the first validator", async () => {
    const instance = await ValidatorSet.deployed();

    await instance.removeValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR);

    const validatorPubKeys = [];
    const validatorAmountBN = await instance.validatorSetSize.call();
    const validatorAmount = parseInt(validatorAmountBN.toString());

    for(let i=0; i < validatorAmount; i++) {
      const pubKey = await instance.validatorPubKey.call(i, {});
      validatorPubKeys.push(pubKey);
    }

    assert.deepEqual(validatorPubKeys, [VALIDATOR2_KEY]);
  });
  
  it("should fail to and same duplicated validator pubkey", async () => {
    const instance = await ValidatorSet.deployed();
    return assert.isRejected(instance.addValidator(VALIDATOR2_KEY, VALIDATOR3_ADDR));
  });
  
  it("should fail to and delete a no included validator pubkey", async () => {
    const instance = await ValidatorSet.deployed();
    return assert.isRejected(instance.removeValidator(VALIDATOR1_KEY, VALIDATOR3_ADDR));
  });
  
  it("should add a another validator", async () => {
    const instance = await ValidatorSet.deployed();

    const tx = await instance.addValidator(VALIDATOR3_KEY, VALIDATOR3_ADDR);
    assert.equal(tx.receipt.status, true, "successful TX status expected");

    const validatorPubKeys = [];
    const validatorSetLengthBN = await instance.validatorSetSize.call();
    const validatorSetLength = parseInt(validatorSetLengthBN.toString());
    for(let i=0; i < validatorSetLength; i++) {
      const pubKey = await instance.validatorPubKey.call(i, {});
      validatorPubKeys.push(pubKey);
    }

    assert.deepEqual(validatorPubKeys, [VALIDATOR2_KEY, VALIDATOR3_KEY]);
  });
  
  it("should fail to modify validator after freeze", async() => {
    const instance = await ValidatorSet.deployed();

    const tx = await instance.setFreezeStatus(true);
    assert.equal(tx.receipt.status, true, "successful TX status expected");
    return assert.isRejected(instance.addValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR));
  });
});
