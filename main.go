package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sinakhalili/abigenz/abis"
)

type EtherscanResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  string `json:"result"`
}

func main() {
	log.Println("Hello, sailor")

	etherscan_api_key := os.Getenv("ETHERSCAN_API_KEY")
	if etherscan_api_key == "" {
		log.Printf("Error: ")
		log.Println("Please set the environment variable ETHERSCAN_API_KEY")
		return
	}

	linkpool_free_rpc := "wss://main-light.eth.linkpool.io/ws"
	log.Println("Websocket url: ", linkpool_free_rpc)

	// Connect to the RPC.
	rpcClient, err := rpc.Dial(linkpool_free_rpc)
	if err != nil {
		log.Fatalln("Unable to connect to the RPC!")
	}
	ethClient := ethclient.NewClient(rpcClient)

	if len(os.Args) < 2 {
		log.Println("Usage: abigenz [CONTRACT ADDRESS]")
		return
	}

	contract_address := os.Args[1]
	log.Printf("Looking for contract %s\n", contract_address)
	contract_hex := common.HexToAddress(contract_address)

	contract, err := abis.NewERC20(contract_hex, ethClient)
	if err != nil {
		log.Fatalln(err)
	}

	tokenName, err := contract.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Contract name: ", tokenName)

	/*
	   etherscan calls are made in the following format
	   (https://docs.etherscan.io/api-endpoints/contracts)
	*/
	base, err := url.Parse("https://api.etherscan.io/api")
	if err != nil {
		return
	}

	params := url.Values{}
	params.Add("module", "contract")
	params.Add("action", "getabi")
	params.Add("address", contract_address)
	params.Add("apiKey", etherscan_api_key)
	base.RawQuery = params.Encode()

	resp, err := http.Get(base.String())
	if err != nil {
		log.Println("Oh no! Got: ")

		body, err := io.ReadAll(resp.Body)
		log.Println(string(body))
		log.Fatalln(err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var etherscan_response EtherscanResponse
	err = json.Unmarshal(body, &etherscan_response)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Etherscan says: %s\n", etherscan_response.Message)
	if etherscan_response.Message == "NOTOK" {
		log.Fatal("Error: ", etherscan_response.Result)
	}

	filename := tokenName + ".json"
	if err := os.WriteFile(filename, []byte(etherscan_response.Result), 0666); err != nil {
		log.Fatal(err)
	}
	log.Println("Abi in ", filename)

	// abigen --pkg main --abi abi.json  --out out.go
	go_filename := tokenName + ".go"
	cmd := exec.Command("abigen", "--pkg", "abis", "--abi", filename, "--out", go_filename)
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Go file in ", go_filename)
}
