package consensus

const Genesis = `
{
    "genesis_time": "2019-03-06T15:20:31.254715867Z",
    "chain_id": "lightstreams-mainnet",
    "consensus_params": {
        "block": {
            "max_bytes": "22020096",
            "max_gas": "-1",
			"time_iota_ms": "1000"
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
            "address": "012C7DB9A70AA4940014A0CC279BFD18D8E1E224",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "hAHyC5KAfQA59DgUL8/JEAy8+mdH1SrflkXas2yiAp8="
            },
            "power": "10",
            "name": "mainnet-validator-node1"
        },
        {
            "address": "DC7F3826B5B38B154D2E1B36F1B6C6E71BF3DCCC",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "kSoT+yG1YSIzMAVy4+lvEIfBzsMqrhUiN5TBTRyg360="
            },
            "power": "10",
            "name": "mainnet-validator-node2"
        },
        {
            "address": "996B7A6C17800857CCBA5CDBE0D5803DCC54BB73",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "Wps2hXJtROjwAiAh0WxONvPOPkKI4uo8ChZymUeO8Nc="
            },
            "power": "10",
            "name": "mainnet-validator-node3"
        },
        {
            "address": "E522C9A549925FBF695B3EF5BE8E90551047DD44",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "2lvmdOFBSV45SDN8mHkJZ2I6I3GTxlc7Ld5LbOP4AU0="
            },
            "power": "10",
            "name": "mainnet-validator-node4"
        }
    ],
    "app_hash": ""
}
`
