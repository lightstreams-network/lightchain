const Validators = artifacts.require("Validators.sol");

const VALIDATOR1_KEY = "AAE6B62F536602EE3E52BA80A8FAB225479BD3B7";
const VALIDATOR1_ADDR = "0xe3af0f7f5bfe871453b765db0b28e79063540c73";

module.exports = (deployer) => {
  console.log("Deploying `Validators.sol`...");
  deployer.deploy(Validators).then((instance) => {
    console.log("Validators.sol successfully deployed!");
    return instance;
  }).then((instance) => {
    return instance.addValidator(VALIDATOR1_KEY, VALIDATOR1_ADDR);
  }).then((tx) => {
    if(tx.receipt.status) {
      console.log(`Added ${VALIDATOR1_KEY}:${VALIDATOR1_ADDR}`);      
    } else {
      console.error(`Error ${JSON.stringify(tx)}`);      
    }
  });
};
