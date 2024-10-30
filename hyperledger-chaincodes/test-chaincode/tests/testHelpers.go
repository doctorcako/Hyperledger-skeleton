package testHelpers

import (
	"smart-contracts/chaincode/tests/mocks"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/transaction.go -fake-name TransactionContext . transactionContext
type transactionContext interface {
	contractapi.TransactionContextInterface
}

// go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/chaincodestub.go -fake-name ChaincodeStub . chaincodeStub
type chaincodeStub interface {
	shim.ChaincodeStubInterface
}

// go run github.com/maxbrunsfeld/counterfeiter/v6 -o mocks/statequeryiterator.go -fake-name StateQueryIterator . stateQueryIterator
type stateQueryIterator interface {
	shim.StateQueryIteratorInterface
}

// SetupTest sets up the mock test environment
func SetupTest() (*mocks.ChaincodeStub, *mocks.TransactionContext) {

}
