# Abigen for zoomers

Just a simple wrapper to fetch abis from etherscan and run [abigen](https://geth.ethereum.org/docs/dapp/native-bindings) on it.

Uses the `name` of a contract if possible.

## Usage
First put your etherscan api key in the environment variable ETHERSCAN_API_KEY.

Then simply run
```
➜ abigenz 0xD63751B0fBef5F0153c69D6429Ed1429B6F79247 # or any contract address
```

Which will output
```
2021/11/17 04:39:41 Hello, sailor
2021/11/17 04:39:41 Websocket url:  wss://main-light.eth.linkpool.io/ws
2021/11/17 04:39:41 Looking for contract 0xD63751B0fBef5F0153c69D6429Ed1429B6F79247
2021/11/17 04:39:41 Contract name:  Uniswap V2
2021/11/17 04:39:42 Etherscan says: OK
2021/11/17 04:39:42 Abi in  Uniswap V2.json
2021/11/17 04:39:42 Go file in  Uniswap V2.go
```

As well as create the files `Uniswap V2.go` and `Uniswap V2.json` in the current directory under the go package "abis".

## Installation

```
go get github.com/sinakhalili/abigenz
```
# abigenz
