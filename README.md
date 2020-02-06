# EMPOW BLOCKCHAIN - Social Network on Blockchain

### How to run full node

- Machine requirements

	- 	CPU: 4 cores or more (8 cores recommended)
	- 	Memory: 8GB or more (16GB recommended)
	- 	Disk: 1TB or more (5TB HDD recommended)
	- 	Network: access to Internet with port tcp: 30000 opened (If you want to enable rpc for node, please open port 30001, 30002)
	
- Run the boot script to start a full node:

	`
curl https://raw.githubusercontent.com/empow-blockchain/go-empow/master/script/boot.sh | INET=mainnet bash
`

	*INET : mainnet, testnet (default: mainnet)*

### Build Empow Blockchain
- Install Golang

- Install Git LFS

- Config Environment Variable (GOPATH)

	- Edit file `~/.profile`
	
	- Add these 2 lines to the end of the file
	```shell
	export GOPATH=$(go env GOPATH)
	export PATH=$PATH:$GOPATH/bin
	```

- Pull code from github to golang folder

	```shell
	go get -d github.com/empow-blockchain/go-empow
	```

- Build code

	```shell
	cd $GOPATH/src/github.com/empow-blockchain/go-empow
	git lfs pull
	make build install
	cd vm/v8vm/v8/; make deploy; cd ../../..
	```
- Run blockchain

	```shell
	iserver  -f ./config/iserver.yml
	```
