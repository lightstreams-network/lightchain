# Tests

## Running all tests

Copy `.env.sample` and fill up the values
```
cp .env.sample .env
```

In case you are running over `standalone` network:
```
ROOT_ACCOUNT="0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e"
PASSPHRASE="WelcomeToSirius"
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
accounts[0] -> `0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e`
```
