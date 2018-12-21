# Tests

## Running all tests


```
geth attach ipc:${HOME}/.lightchain/geth.ipc
```

Before you can run the tests, the main accounts: must be unlocked. For simplicity of local environment, we can unlock
them for a long period of time using password `ggarri86`.

```
web3.personal.unlockAccount(web3.personal.listAccounts[0], `ggarri86`, 600000)
web3.personal.unlockAccount(web3.personal.listAccounts[1], `ggarri86`, 600000)
```

Where: 
* accounts[0] -> `0x4eaaad8ea38d5ef75ebdeb3d1be59d56f86c4ca9` 
* accounts[1] -> `0x4f5adedca6d869e9f5f7dcf4b7a9dfa8231a095f`

Once environment is ready we can run the test
```
truffle test --network=sirius
```
