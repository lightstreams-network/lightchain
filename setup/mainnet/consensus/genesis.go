package consensus

const Genesis = `
{
    "genesis_time": "2019-03-31T00:00:00.000000000Z",
    "chain_id": "lightstreams-mainnet",
    "consensus_params": {
        "block_size": {
            "max_bytes": "22020096",
            "max_gas": "-1"
        },
        "evidence": {
            "max_age": "100000"
        },
        "validator": {
            "pub_key_types": [
                "ed25519"
            ]
        }
    },
    "validators": [
        {
            "address": "34A193CE63F3E51BA875AA42F89E14ACB713D5C1",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "/axfQYNcqjYjJOTb6jMPNGsbuLD65/JfrEBstfvQXOY="
            },
            "power": "10",
            "name": "mainnet-validator-1"
        },
    ],
    "app_hash": ""
}
`