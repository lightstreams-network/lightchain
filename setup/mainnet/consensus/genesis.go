package consensus

const Genesis = `
{
  "genesis_time": "2019-03-05T15:45:37.820518684Z",
  "chain_id": "test-chain-7ovHO6",
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
      "address": "DF161D1DAFFABF814842B3800A8A150F03D85774",
      "pub_key": {
        "type": "tendermint/PubKeyEd25519",
        "value": "EZg9QOMa05g/cSSjeXLfuwBORT/9ETfIjFxvACxOyMk="
      },
      "power": "10",
      "name": ""
    }
  ],
  "app_hash": ""
}
`
