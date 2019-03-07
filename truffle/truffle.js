/*
 * NB: since truffle-hdwallet-provider 0.0.5 you must wrap HDWallet providers in a 
 * function when declaring them. Failure to do so will cause commands to hang. ex:
 * ```
 * mainnet: {
 *     provider: function() { 
 *       return new HDWalletProvider(mnemonic, 'https://mainnet.infura.io/<infura-key>') 
 *     },
 *     network_id: '1',
 *     gas: 4500000,
 *     gasPrice: 10000000000,
 *   },
 */

module.exports = {
    networks: {
        sirius: {
            host: "127.0.0.1",
            port: 8545,
            network_id: "162",
            from: "0xd119b8b038d3a67d34ca1d46e1898881626a082b",
            gasPrice: "500000000000",
        },
        mainnet: {
            host: "127.0.0.1",
            port: 8545,
            network_id: "162",
            from: "0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e",
            gasPrice: "500000000000",
        },
        standalone: {
            host: "127.0.0.1",
            port: 8545,
            network_id: "161",
            from: "0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e",
            gasPrice: "500000000000",
        }
    },
    mocha: { 
        enableTimeouts: false
    }
};
