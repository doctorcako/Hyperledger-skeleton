package chaincodeTest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const ()

type BaseLine struct {
	IsBaseLine   bool      `json:"isBaseLine"`
	Added        string    `json:"added"`
	BaseLineTime time.Time `json:"baseLineTime"`
}

func (t *chaindodeContract) Create(ctx contractapi.TransactionContextInterface,
	eventName,
	description string,
	aeTime time.Time,
	created time.Time,
	executionUser,
	channel string,
) (string, error) {

	operationBytes, err := ctx.GetStub().GetState(operationId)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %w", err)
	}
	if operationBytes != nil {
		return "", fmt.Errorf("operation already exists")
	}

}

func (t *chaindodeContract) UpdateCountermeasureExecution(ctx contractapi.TransactionContextInterface,
	operationId,
	executionStatus string, // ERROR OR OK
	executionResult string,
	startTime,
	endTime time.Time,
	channel string,
) (string, error) {

	// UPDATE THE operation
	operationBytes, err := ctx.GetStub().GetState(operationId)
	if err != nil {
		return "", fmt.Errorf("error getting operation: %w", err)
	}

	var operation Operation
	if operationBytes != nil {
		err = json.Unmarshal(operationBytes, &operation)
		if err != nil {
			return "", fmt.Errorf("error deserializing operation: %w", err)
		}
	} else {
		return "", fmt.Errorf("the operation doesn't exist: %w", err)
	}

	operationBytes, err = json.Marshal(operation)
	if err != nil {
		return "", fmt.Errorf("error serializing operation: %w", err)
	}

	err = ctx.GetStub().PutState(operation.Id, operationBytes)
	if err != nil {
		return "", fmt.Errorf("error putting state : %w", err)
	}

	return string(operationBytes), nil
}
