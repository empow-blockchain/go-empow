# EMPOW BLOCKCHAIN - Social Network on Blockchain

### How to run full node

- Machine requirements

	- 	CPU: 4 cores or more (8 cores recommended)
	- 	Memory: 8GB or more (16GB recommended)
	- 	Disk: 1TB or more (5TB HDD recommended)
	- 	Network: access to Internet with port tcp: 30000 opened (If you want to enable rpc for node, please open port 30001, 30002)
	
- Run the boot script to start a full node:

	`
curl https://raw.githubusercontent.com/empow-blockchain/go-empow/master/script/boot.sh | INET=testnet bash
`

	*INET : mainnet, testnet (default: mainnet)*
