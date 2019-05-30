# Change Log

## [Unreleased](https://github.com/lightstreams-network/lightchain/tree/HEAD)

[Full Changelog](https://github.com/lightstreams-network/lightchain/compare/v1.1.2...HEAD)

**Implemented enhancements:**

- Replace RPCApi communication between Ethereum and Consensus [\#166](https://github.com/lightstreams-network/lightchain/issues/166)
- Upgrade to latest Tendermint version v0.31.5 [\#161](https://github.com/lightstreams-network/lightchain/issues/161)
- Adjust tendermint timeout settings [\#136](https://github.com/lightstreams-network/lightchain/issues/136)
- Force flag should raise an interactive confirmation before deleting the datadir [\#98](https://github.com/lightstreams-network/lightchain/issues/98)
- \#161 Upgrades Tendermint to version v0.31.5. [\#167](https://github.com/lightstreams-network/lightchain/pull/167) ([EnchanterIO](https://github.com/EnchanterIO))

**Fixed bugs:**

- Comited block logs are attempting to reference a rootHash as blockHash by mistake [\#152](https://github.com/lightstreams-network/lightchain/issues/152)
- \[URGENT\] Investigate client node corrupt state  [\#146](https://github.com/lightstreams-network/lightchain/issues/146)
- lightchain init/run tracelog default value path [\#137](https://github.com/lightstreams-network/lightchain/issues/137)
- Ethereum `debug` API not loaded [\#133](https://github.com/lightstreams-network/lightchain/issues/133)
- Lightchain doesn't shutdown gracefuly since last Tendermint update 0.30 [\#122](https://github.com/lightstreams-network/lightchain/issues/122)
- Prevent unexpected load of the default tendermint config  [\#108](https://github.com/lightstreams-network/lightchain/issues/108)
- Invalid TX corrupts the local node's blockchain state [\#70](https://github.com/lightstreams-network/lightchain/issues/70)
- Fixed dockerfile [\#147](https://github.com/lightstreams-network/lightchain/pull/147) ([ggarri](https://github.com/ggarri))

**Closed issues:**

- Make Install Fails on Ubuntu 18.10 [\#129](https://github.com/lightstreams-network/lightchain/issues/129)
- Switch to LocalClient for ABCI connection instead of a current Socket [\#110](https://github.com/lightstreams-network/lightchain/issues/110)

**Merged pull requests:**

- \#166 Replacing usage of consensus RPC calls by direct method calls [\#168](https://github.com/lightstreams-network/lightchain/pull/168) ([ggarri](https://github.com/ggarri))
- \#136 Correlate timeouts to network empty block intervals [\#165](https://github.com/lightstreams-network/lightchain/pull/165) ([ggarri](https://github.com/ggarri))
- \#98 Adding confirmation before removing folder [\#164](https://github.com/lightstreams-network/lightchain/pull/164) ([ggarri](https://github.com/ggarri))
- \#135 catching and handling consensus config loading error [\#163](https://github.com/lightstreams-network/lightchain/pull/163) ([ggarri](https://github.com/ggarri))
- \#133 Enable ethereum debug api when it is requested [\#159](https://github.com/lightstreams-network/lightchain/pull/159) ([ggarri](https://github.com/ggarri))
- \#134 Adding test to validate the correctness of eth\_estimateGas funct… [\#158](https://github.com/lightstreams-network/lightchain/pull/158) ([ggarri](https://github.com/ggarri))
- 137 bug tracer [\#157](https://github.com/lightstreams-network/lightchain/pull/157) ([ggarri](https://github.com/ggarri))
- \#137 Moving tracer.log file into local node datadir [\#154](https://github.com/lightstreams-network/lightchain/pull/154) ([ggarri](https://github.com/ggarri))
- \#152 Fixes `log.BlockHash = bs.header.Root`. [\#153](https://github.com/lightstreams-network/lightchain/pull/153) ([EnchanterIO](https://github.com/EnchanterIO))
- \#146 Panic in case it cannot persist local state to prevent corruptin… [\#151](https://github.com/lightstreams-network/lightchain/pull/151) ([ggarri](https://github.com/ggarri))
- \#110 Migrates Tendermint ABCI communication from socket to local client. [\#150](https://github.com/lightstreams-network/lightchain/pull/150) ([EnchanterIO](https://github.com/EnchanterIO))

## [v1.1.2](https://github.com/lightstreams-network/lightchain/tree/v1.1.2) (2019-05-16)
[Full Changelog](https://github.com/lightstreams-network/lightchain/compare/v1.0.0...v1.1.2)

**Fixed bugs:**

- Websockets + Web3: `Error: Connection not open` [\#141](https://github.com/lightstreams-network/lightchain/issues/141)
- Unable to start node with Go 1.12 [\#128](https://github.com/lightstreams-network/lightchain/issues/128)

**Closed issues:**

- Tokens distribution script + Smart Contract [\#117](https://github.com/lightstreams-network/lightchain/issues/117)
- Make the --datadir flag optional  [\#49](https://github.com/lightstreams-network/lightchain/issues/49)

**Merged pull requests:**

- Version update: v1.1.2 [\#144](https://github.com/lightstreams-network/lightchain/pull/144) ([ggarri](https://github.com/ggarri))
- \#141 Implementing support for --wsapi and --wsorigins [\#143](https://github.com/lightstreams-network/lightchain/pull/143) ([ggarri](https://github.com/ggarri))
- upgrade geth [\#140](https://github.com/lightstreams-network/lightchain/pull/140) ([0x13a](https://github.com/0x13a))
- \#117 Removes distribution pkg [\#132](https://github.com/lightstreams-network/lightchain/pull/132) ([EnchanterIO](https://github.com/EnchanterIO))
- feat: makes makefile self documenting [\#131](https://github.com/lightstreams-network/lightchain/pull/131) ([0x13a](https://github.com/0x13a))
- Include help section [\#130](https://github.com/lightstreams-network/lightchain/pull/130) ([ggarri](https://github.com/ggarri))
- feat: makes datadir flag optional and default to ~/.lightchain [\#127](https://github.com/lightstreams-network/lightchain/pull/127) ([0x13a](https://github.com/0x13a))
- \[NOTASK\] Reduce total cost of running truffle test suite [\#126](https://github.com/lightstreams-network/lightchain/pull/126) ([ggarri](https://github.com/ggarri))
- \[NOTASK\] Improving dev scripting to integrate mainnet usage [\#125](https://github.com/lightstreams-network/lightchain/pull/125) ([ggarri](https://github.com/ggarri))
- \[NOTASK\] Improve build script [\#124](https://github.com/lightstreams-network/lightchain/pull/124) ([ggarri](https://github.com/ggarri))

## [v1.0.0](https://github.com/lightstreams-network/lightchain/tree/v1.0.0) (2019-03-27)
[Full Changelog](https://github.com/lightstreams-network/lightchain/compare/v0.12.0...v1.0.0)

**Closed issues:**

- Update readme with mainnet network instructions [\#120](https://github.com/lightstreams-network/lightchain/issues/120)
- Update tests to latest Web3 API 1.0.x [\#114](https://github.com/lightstreams-network/lightchain/issues/114)
- Write a test for "double spend attack" from user, balance perspective [\#113](https://github.com/lightstreams-network/lightchain/issues/113)
- Optimize consensus configuration [\#112](https://github.com/lightstreams-network/lightchain/issues/112)
- Update README to introduce Mainnet [\#109](https://github.com/lightstreams-network/lightchain/issues/109)
- Configure database genesis file [\#106](https://github.com/lightstreams-network/lightchain/issues/106)
- Configure  consensus configs before launch [\#105](https://github.com/lightstreams-network/lightchain/issues/105)
- Generate genesis alloc map [\#104](https://github.com/lightstreams-network/lightchain/issues/104)
- Generate 4 priv + pub keys for mainnet in isolated, secured env and configure consensus genesis [\#103](https://github.com/lightstreams-network/lightchain/issues/103)
- Add mainet network flow [\#101](https://github.com/lightstreams-network/lightchain/issues/101)
- Integration/unit tests proving fundamental functionalities \(TX handling / State / DB / Syncing... [\#76](https://github.com/lightstreams-network/lightchain/issues/76)
- Improve cli-doc in `docs.lightstreams.network` [\#62](https://github.com/lightstreams-network/lightchain/issues/62)
- Implement logic to allow the approval/rejection of Validators [\#16](https://github.com/lightstreams-network/lightchain/issues/16)
- Define Reward strategy for Validators [\#15](https://github.com/lightstreams-network/lightchain/issues/15)
- Implement socket communication with Tendermint [\#1](https://github.com/lightstreams-network/lightchain/issues/1)

**Merged pull requests:**

- \#106 Adding token allocation to mainnet genesis file [\#123](https://github.com/lightstreams-network/lightchain/pull/123) ([ggarri](https://github.com/ggarri))
- Mainnet [\#121](https://github.com/lightstreams-network/lightchain/pull/121) ([ggarri](https://github.com/ggarri))
- 117 distribution sc script [\#119](https://github.com/lightstreams-network/lightchain/pull/119) ([EnchanterIO](https://github.com/EnchanterIO))
- \#113 Adds 06\_double\_spend.js simplified test. [\#116](https://github.com/lightstreams-network/lightchain/pull/116) ([EnchanterIO](https://github.com/EnchanterIO))
- Updates Truffle to 5.0 and Web3 to 1.0, migrates tests [\#115](https://github.com/lightstreams-network/lightchain/pull/115) ([EnchanterIO](https://github.com/EnchanterIO))
- Docs update [\#107](https://github.com/lightstreams-network/lightchain/pull/107) ([mikesmo](https://github.com/mikesmo))

## [v0.12.0](https://github.com/lightstreams-network/lightchain/tree/v0.12.0) (2019-03-06)
[Full Changelog](https://github.com/lightstreams-network/lightchain/compare/v0.10.0...v0.12.0)

**Fixed bugs:**

- Fix initialization process to consider an empty `datadir` a valid folder [\#93](https://github.com/lightstreams-network/lightchain/issues/93)
- Fix issues with not closing gracefully after an interrupt signal [\#92](https://github.com/lightstreams-network/lightchain/issues/92)
- Init cmd should exit if directory already exists [\#79](https://github.com/lightstreams-network/lightchain/issues/79)
- 92 fix unreachable trap signall [\#100](https://github.com/lightstreams-network/lightchain/pull/100) ([EnchanterIO](https://github.com/EnchanterIO))
- \#93 Allows `lightchain init` to be executed in empty data dir. [\#99](https://github.com/lightstreams-network/lightchain/pull/99) ([EnchanterIO](https://github.com/EnchanterIO))
- \#79 Prevents lightchain node dir to be overwritten accidentally on init. [\#88](https://github.com/lightstreams-network/lightchain/pull/88) ([EnchanterIO](https://github.com/EnchanterIO))

**Closed issues:**

- Upgrade `tendermint` and `go-ethereum` deps to latest stable version [\#84](https://github.com/lightstreams-network/lightchain/issues/84)
- Assert blockchain state on submited transaction [\#77](https://github.com/lightstreams-network/lightchain/issues/77)
- Running `lightchain` without any flags should display usage [\#60](https://github.com/lightstreams-network/lightchain/issues/60)
- Protection against DOS tx spam [\#46](https://github.com/lightstreams-network/lightchain/issues/46)
- Review and decide Gas Price for network [\#13](https://github.com/lightstreams-network/lightchain/issues/13)

**Merged pull requests:**

- \#46 Increases min gas price to 500000000000 and fixes ABCI sad path. [\#95](https://github.com/lightstreams-network/lightchain/pull/95) ([EnchanterIO](https://github.com/EnchanterIO))
- \#70 Adding documentation to work around corrupted local state [\#91](https://github.com/lightstreams-network/lightchain/pull/91) ([ggarri](https://github.com/ggarri))
- \#60 Running `lightchain` now shows usage instructions in terminal. [\#89](https://github.com/lightstreams-network/lightchain/pull/89) ([EnchanterIO](https://github.com/EnchanterIO))
- Feature/77 trace state on submitted tx [\#86](https://github.com/lightstreams-network/lightchain/pull/86) ([EnchanterIO](https://github.com/EnchanterIO))

## [v0.10.0](https://github.com/lightstreams-network/lightchain/tree/v0.10.0) (2019-02-26)
[Full Changelog](https://github.com/lightstreams-network/lightchain/compare/v0.9.1...v0.10.0)

**Fixed bugs:**

- \[Bug\] net\_version and chainId do not match [\#65](https://github.com/lightstreams-network/lightchain/issues/65)
- Research possible memory leak [\#64](https://github.com/lightstreams-network/lightchain/issues/64)
- Close process gracefully [\#36](https://github.com/lightstreams-network/lightchain/issues/36)
- \#84 Upgraded go-ethereum and tendermint version | Implemented fix pro… [\#85](https://github.com/lightstreams-network/lightchain/pull/85) ([ggarri](https://github.com/ggarri))
- \[\#65\] Set ethCfg.NetworkId as Genesis.Config.ChainId [\#66](https://github.com/lightstreams-network/lightchain/pull/66) ([ggarri](https://github.com/ggarri))

**Closed issues:**

- Assert genesis blockchain state [\#80](https://github.com/lightstreams-network/lightchain/issues/80)
- Integrate Prometheus exporter [\#78](https://github.com/lightstreams-network/lightchain/issues/78)
- Add instructions to launch lightchain with RPC open [\#74](https://github.com/lightstreams-network/lightchain/issues/74)
- Create a dedicated Ethereum Wallet for testing proposes in Sirius [\#68](https://github.com/lightstreams-network/lightchain/issues/68)
- Optimize log outputs [\#61](https://github.com/lightstreams-network/lightchain/issues/61)
- Update README with latest instructions [\#47](https://github.com/lightstreams-network/lightchain/issues/47)

**Merged pull requests:**

- 78 prometheus intergration [\#83](https://github.com/lightstreams-network/lightchain/pull/83) ([ggarri](https://github.com/ggarri))
- \#80 Creates `tracer pkg` with first trace "AssertPersistedGenesisBlock". [\#82](https://github.com/lightstreams-network/lightchain/pull/82) ([EnchanterIO](https://github.com/EnchanterIO))
- \#74 Replacing last usages of ethereum logger by engine logger [\#81](https://github.com/lightstreams-network/lightchain/pull/81) ([ggarri](https://github.com/ggarri))
- \#74 Adds docs how to run Lightchain with RPC open. [\#75](https://github.com/lightstreams-network/lightchain/pull/75) ([EnchanterIO](https://github.com/EnchanterIO))
- Feature/70 consensus db pkgs cleanup no breaking changes [\#73](https://github.com/lightstreams-network/lightchain/pull/73) ([EnchanterIO](https://github.com/EnchanterIO))
- \#70 Cleans up consensus and database pkg. Adds docs and better logging. [\#72](https://github.com/lightstreams-network/lightchain/pull/72) ([EnchanterIO](https://github.com/EnchanterIO))
- Fixed typo [\#71](https://github.com/lightstreams-network/lightchain/pull/71) ([azappella](https://github.com/azappella))
- \[\#68\] Adding new ethereum account for testing proposes [\#69](https://github.com/lightstreams-network/lightchain/pull/69) ([ggarri](https://github.com/ggarri))
- \#61 Adds context to all logging across pkgs and corrects wording. [\#67](https://github.com/lightstreams-network/lightchain/pull/67) ([EnchanterIO](https://github.com/EnchanterIO))

## [v0.9.1](https://github.com/lightstreams-network/lightchain/tree/v0.9.1) (2019-02-01)
**Fixed bugs:**

- Explorer is displaying bad block time due to issue in internal calculation on JS side [\#56](https://github.com/lightstreams-network/lightchain/issues/56)
- Tests are not passing on sirius ntw [\#45](https://github.com/lightstreams-network/lightchain/issues/45)
- Fix wrong `chaindata` generated [\#39](https://github.com/lightstreams-network/lightchain/issues/39)
- Each network should have a unique chain ID [\#37](https://github.com/lightstreams-network/lightchain/issues/37)
- Fix issue with Tendermint not excluding empty blocks [\#29](https://github.com/lightstreams-network/lightchain/issues/29)
- Improves wording in tests and math calculation should be done on BN [\#22](https://github.com/lightstreams-network/lightchain/issues/22)
- Fix `eth\_syncing` issue [\#19](https://github.com/lightstreams-network/lightchain/issues/19)
- Close Lightchain gracefully [\#14](https://github.com/lightstreams-network/lightchain/issues/14)
- \[Bug\] Restore correctly Ethereum.Blockchain head state [\#12](https://github.com/lightstreams-network/lightchain/issues/12)
- \#19 Considers highest block to be current block +1 if not synced. [\#55](https://github.com/lightstreams-network/lightchain/pull/55) ([EnchanterIO](https://github.com/EnchanterIO))
- \#37 Fixes unique chain ID per network, EIP-155. [\#44](https://github.com/lightstreams-network/lightchain/pull/44) ([EnchanterIO](https://github.com/EnchanterIO))
- Bug/19 eth web3 syncing real numbers and tests [\#42](https://github.com/lightstreams-network/lightchain/pull/42) ([EnchanterIO](https://github.com/EnchanterIO))
- \#29 Fixed tendermint detection of empty blocks [\#30](https://github.com/lightstreams-network/lightchain/pull/30) ([ggarri](https://github.com/ggarri))
- Performs math calculations on BN instead of strings. [\#23](https://github.com/lightstreams-network/lightchain/pull/23) ([EnchanterIO](https://github.com/EnchanterIO))

**Closed issues:**

- Create PULL\_REQUEST\_TEMPLATE [\#54](https://github.com/lightstreams-network/lightchain/issues/54)
- Create ISSUE\_TEMPLATE [\#53](https://github.com/lightstreams-network/lightchain/issues/53)
- Upgrade Tendermint to v0.29.1 [\#52](https://github.com/lightstreams-network/lightchain/issues/52)
- RPC should not be exposed on Sirius validators [\#51](https://github.com/lightstreams-network/lightchain/issues/51)
- Include lightchain cli-docs in docs.lightstreams.network  [\#48](https://github.com/lightstreams-network/lightchain/issues/48)
- Adjust truffle test to support both networks without any code changes [\#41](https://github.com/lightstreams-network/lightchain/issues/41)
- Remove sirius keystore files from repository [\#38](https://github.com/lightstreams-network/lightchain/issues/38)
- Investigate a block never contains more than 1 trx [\#35](https://github.com/lightstreams-network/lightchain/issues/35)
- Remove usage of filesytem setup files in favor of go constants  [\#32](https://github.com/lightstreams-network/lightchain/issues/32)
- Implement `stand-alone` network flag [\#27](https://github.com/lightstreams-network/lightchain/issues/27)
- Implement a full integration of Tendermint   [\#18](https://github.com/lightstreams-network/lightchain/issues/18)
- Replace if possible usages of ETH by PHT token within blockchain [\#17](https://github.com/lightstreams-network/lightchain/issues/17)
- Create docker-compose.yml to easily launch lightchain + tendermint [\#11](https://github.com/lightstreams-network/lightchain/issues/11)
- Create a Dockerfile for Tendermint [\#10](https://github.com/lightstreams-network/lightchain/issues/10)
- Create a Dockerfile for lightchain  [\#9](https://github.com/lightstreams-network/lightchain/issues/9)
- Implement `--networkId` flag reading [\#6](https://github.com/lightstreams-network/lightchain/issues/6)
- End-to-end web3 \(truffle\) integration tests [\#4](https://github.com/lightstreams-network/lightchain/issues/4)
- Codebase skeleton [\#3](https://github.com/lightstreams-network/lightchain/issues/3)
- Basic Ethereum blockchain functionalities \(transfers, contracts\) [\#2](https://github.com/lightstreams-network/lightchain/issues/2)

**Merged pull requests:**

- \[\#36\] Implemented Consensus engine closing event [\#58](https://github.com/lightstreams-network/lightchain/pull/58) ([ggarri](https://github.com/ggarri))
- \#37 web3.eth.syncing now returns real-time sync data status. [\#50](https://github.com/lightstreams-network/lightchain/pull/50) ([EnchanterIO](https://github.com/EnchanterIO))
- \#39 Provisionally preventing duplicated chaindb folder by correlating… [\#40](https://github.com/lightstreams-network/lightchain/pull/40) ([ggarri](https://github.com/ggarri))
- \#32 Removed usage of filesystem from initialization process [\#34](https://github.com/lightstreams-network/lightchain/pull/34) ([ggarri](https://github.com/ggarri))
- Adds support for web3.eth.syncing feature. [\#33](https://github.com/lightstreams-network/lightchain/pull/33) ([EnchanterIO](https://github.com/EnchanterIO))
- \#11 Including dockerfile and updated README after skeleton refactor [\#31](https://github.com/lightstreams-network/lightchain/pull/31) ([ggarri](https://github.com/ggarri))
- \#27 Implemented stand alone node flag [\#28](https://github.com/lightstreams-network/lightchain/pull/28) ([ggarri](https://github.com/ggarri))
- \#3 Project skeleton refactor [\#26](https://github.com/lightstreams-network/lightchain/pull/26) ([ggarri](https://github.com/ggarri))
- Restructure of the project [\#25](https://github.com/lightstreams-network/lightchain/pull/25) ([EnchanterIO](https://github.com/EnchanterIO))
- \#3 First baby steps for skeleton refactor [\#20](https://github.com/lightstreams-network/lightchain/pull/20) ([ggarri](https://github.com/ggarri))
- \#4 Implemented first batch of Blockchain tests [\#8](https://github.com/lightstreams-network/lightchain/pull/8) ([ggarri](https://github.com/ggarri))



\* *This Change Log was automatically generated by [github_changelog_generator](https://github.com/skywinder/Github-Changelog-Generator)*