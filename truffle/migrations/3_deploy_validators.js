const ValidatorSet = artifacts.require("ValidatorSet.sol");

module.exports = (deployer) => {
  console.log("Deploying `ValidatorSet.sol`...");
  
  deployer.deploy(ValidatorSet).then((instance) => {
    console.log(`Validators.sol successfully!`);
    return instance;
  });
};
