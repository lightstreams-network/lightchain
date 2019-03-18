/**
 * - Execute 100 transaction in parallel
 */

describe('Ethereum API', () => {
  it("should assert whether network version does not match Genesis chainId", async () => {
    const networkVersion = await web3.eth.net.getId();
    if (process.env.NETWORK === 'mainnet') {
      assert.equal(networkVersion, "163", "Network version is not expected one");
    } else if (process.env.NETWORK === 'sirius') {
      assert.equal(networkVersion, "162", "Network version is not expected one");
    } else if (process.env.NETWORK === 'standalone') {
      assert.equal(networkVersion, "161", "Network version is not expected one");
    } else {
      assert.equal(true, false, "Invalid selected network");
    }
  });
});
