# Tests

## Running all tests

Before you can run the tests, the main accounts, accounts[0] and [1]
must be unlocked. For simplicity of local environment, we can unlock
them for a long period of time using password "ggarri86".

```
geth attach ipc:/Users/enchanterio/.lightchain/geth.ipc

web3.personal.unlockAccount(web3.personal.listAccounts[0], null, 600000)
web3.personal.unlockAccount(web3.personal.listAccounts[1], null, 600000)
```

```
truffle test --network=sirius
```