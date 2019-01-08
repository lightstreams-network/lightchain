# Tests

## Running all tests

Copy `.env.sample` and fill up the values
```
cp .env.sample .env
```

Recommended values:
```
ROOT_ACCOUNT="0x4f5adedca6d869e9f5f7dcf4b7a9dfa8231a095f"
PASSPHRASE="ggarri86"
```

That account corresponds to one of the ones defined
on the genesis block `/setup/genesis.json`.

Once environment is ready we can run the test
```
npm run test
```

### Notes

It is important to mention that accounts[] corresponds to the first accounts
created during the genesis block which can be found in 
the root folder of the project.

```
accounts[0] -> `0x4eaaad8ea38d5ef75ebdeb3d1be59d56f86c4ca9` 
accounts[1] -> `0x4f5adedca6d869e9f5f7dcf4b7a9dfa8231a095f`
```
