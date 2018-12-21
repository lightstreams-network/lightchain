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
            network_id: "*",
            from: "0x4f5adedca6d869e9f5f7dcf4b7a9dfa8231a095f"
        }
    }
};
