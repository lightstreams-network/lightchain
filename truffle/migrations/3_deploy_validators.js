const ValidatorSet = artifacts.require("ValidatorSet.sol");

const VALIDATOR1_KEY = "AA260068DDA65DF1F57AE2B4C2A5B7E7B30B34B1";
const VALIDATOR1_ADDR = "0xe3af0f7f5bfe871453b765db0b28e79063540c73";

module.exports = (deployer) => {
  console.log("Deploying `ValidatorSet.sol`...");
  deployer.deploy(ValidatorSet).then((instance) => {
    console.log(`Validators.sol successfully deployed at ${instance.address}!`);
    return instance;
  }).then((instance) => {
    return instance.addValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR);
  }).then((tx) => {
    if(tx.receipt.status) {
      console.log(`Added ${VALIDATOR1_KEY}:${VALIDATOR1_ADDR} at `);      
    } else {
      console.error(`Error ${JSON.stringify(tx)}`);      
    }
  });
};
