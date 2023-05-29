package main

import (
	"log"
	"supplychain/chaincode"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	supplyChaincode, err := contractapi.NewChaincode(&chaincode.SmartContract{})
	
	if err != nil {
		log.Panicf("Error creating supply chaincode: %v", err)
	}

	if err := supplyChaincode.Start(); err != nil {
		log.Panicf("Error starting supply chaincode: %v", err)
	}
}
