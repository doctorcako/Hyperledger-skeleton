package chaincodeTest

import (
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Define the smart contract struct
type chaindodeContract struct {
	contractapi.Contract
}

type Operation struct {
	Id            string          `json:"id"`
	Token         string          `json:"protocol"`
	Action        AssociatedEvent `json:"associatedEvent"`
	ExecutionUser string          `json:"executionUser"`
	CreationTime  time.Time       `json:"creationTime"`
	LastUpdate    time.Time       `json:"lastUpdate"`
}

type AssociatedEvent struct {
	Message     string    `json:"message"`
	EventName   string    `json:"eventName"`
	Description string    `json:"description"`
	Timestamp   time.Time `json:"timestamp"`
}
