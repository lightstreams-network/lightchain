package database

// IMPORTANT: In order to prevent a replay-attack protection, ensure a UNIQUE CHAIN ID
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md

const Genesis = `
{
    "config": {
        "chainId": 162,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "ByzantiumBlock": 0,
        "ConstantinopleBlock": 558000,
        "PetersburgBlock": 558001
    },
    "nonce": "1",
    "difficulty": "1024",
    "gasLimit": "100000000",
    "alloc": {
        "0x61f3b85929c20980176d4a09771284625685c40e": {
            "balance": "300000000000000000000000000"
        }
    }
}
`
