package consensus

const Genesis = `
{
    "genesis_time": "2018-12-16T15:20:31.254715867Z",
    "chain_id": "lightstreams-sirius-testnet",
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
            "address": "2DE3B810E4EAC51A10A3740D15AE92142B01DC7B",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "Eg4FsTNkUETQ330nlaCRlbRPPKMArLlnsPwHGURL8Xs="
            },
            "power": "10",
            "name": "sirius-validator-node1"
        },
        {
            "address": "21B62C962F60618DB4DA34CD623C39FC7CC35CE6",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "jkxtFjBaS/cfK76HRlaW65O8woOHgmlC5bguzlQWsD4="
            },
            "power": "10",
            "name": "sirius-validator-node2"
        },
        {
            "address": "B7DFDC424F64EC89C42A502F89EB29F811144CA1",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "j5XWFDrV58Bqz3UdVKBAA4jGps5OTPv0PnaXqCEhZR8="
            },
            "power": "10",
            "name": "sirius-validator-node3"
        },
        {
            "address": "9A171D8706D6E09BE57326402D231B969DF5D872",
            "pub_key": {
                "type": "tendermint/PubKeyEd25519",
                "value": "+PaREFjqJYsU1lpMkFiShhB52agOZyHSzSp4p0TFaak="
            },
            "power": "10",
            "name": "sirius-validator-node4"
        }
    ],
    "app_hash": ""
}
`
