package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Chaincode init function
func main() {
	// Create a new instance of the chaincode
	chaincode, err := contractapi.NewChaincode(new(chaincodeTest.Contract))
	if err != nil {
		panic(err)
	}

	// Start the chaincode
	if err := chaincode.Start(); err != nil {
		panic(err)
	}

}
