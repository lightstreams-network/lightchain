package database

// IMPORTANT: In order to prevent a replay-attack protection, ensure a UNIQUE CHAIN ID
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md

const Genesis = `
{
    "config": {
        "chainId": 163,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "ByzantiumBlock": 0
    },
	"nonce": "1",
    "difficulty": "1024",
    "gasLimit": "100000000",
    "alloc": {
        "0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e": {
            "balance": "300000000000000000000000000"
        }
    }
}
`
