const Validators = artifacts.require("Validators.sol");

module.exports = (deployer) => {
  console.log("Deploying `Validators.sol`...");
  deployer.deploy(Validators).then(() => {
    console.log("Validators.sol successfully deployed!");
  });
};
