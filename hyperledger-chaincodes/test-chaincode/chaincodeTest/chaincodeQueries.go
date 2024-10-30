package chaincodeTest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type HistoryQueryResult struct {
	Record    *chaindodeContract `json:"record"`
	TxId      string             `json:"txId"`
	Timestamp time.Time          `json:"timestamp"`
	IsDelete  bool               `json:"isDelete"`
}

func (f *chaincodeContract) GetById(ctx contractapi.TransactionContextInterface, id string) (Operation, error) {
	operationBytes, err := ctx.GetStub().GetState(id)
	if err != nil {
		return Operation{}, fmt.Errorf("error getting operation: %w", err)
	}

	var operation chaincodeContract
	if operation != nil {
		err = json.Unmarshal(operationBytes, &operation)
		if err != nil {
			return Operation{}, fmt.Errorf("error deserializing equipment: %v", err)
		}
	}

	return operation, nil
}

func (f *chaincodeContract) GetHistoricalOperationById(ctx contractapi.TransactionContextInterface, operationId string) ([]HistoryQueryResult, error) {
	// Get the history of operation by macAddress
	resultsIterator, err := ctx.GetStub().GetHistoryForKey(operationId)
	if err != nil {
		return nil, fmt.Errorf("failed to get history of operation: %w", err)
	}
	defer resultsIterator.Close()

	var operationHistory []HistoryQueryResult

	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over operation history: %w", err)
		}
		var operation Operation

		if len(result.Value) > 0 {
			err = json.Unmarshal(result.Value, &operation)
			if err != nil {
				return nil, err
			}
		} else {
			operation = Operation{
				Id: operationId,
			}
		}

		timestamp, err := ptypes.Timestamp(result.Timestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to get history timestamps: %w", err)
		}

		record := HistoryQueryResult{
			TxId:      result.TxId,
			Timestamp: timestamp,
			Record:    &operation,
			IsDelete:  result.IsDelete,
		}
		operationHistory = append(operationHistory, record)
	}

	return operationHistory, nil
}

func (c *chaincodeContract) FilterOperationsByParams(ctx contractapi.TransactionContextInterface, params map[string]string) ([]chaincodeContract, error) {
	// If no filters are provided, call getAlloperations
	if len(params) == 0 {
		return getAllOperations(ctx)
	}

	// Create the base selector map
	selector := map[string]interface{}{"selector": map[string]interface{}{}}

	// Add filters based on the provided parameters
	// for key, value := range params {
	// 	// Special handling for integer fields like countermeasureId
	// 	if key == "" {
	// 		// Convert the string value to an integer
	// 		intValue, err := strconv.Atoi(value)
	// 		if err != nil {
	// 			return nil, fmt.Errorf("invalid value for integer field 'countermeasure.countermeasureId': %s", value)
	// 		}
	// 		selector["selector"].(map[string]interface{})[key] = intValue
	// 	} else {
	// 		// Add all other filters as string values
	// 		selector["selector"].(map[string]interface{})[key] = value
	// 	}
	// }

	// Convert the selector to a JSON string
	query, err := json.Marshal(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query selector: %w", err)
	}
	queryString := string(queryBytes)

	// Execute the rich query
	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get Operations: %w", err)
	}
	defer resultsIterator.Close()

	// Create a list to store operations
	var operationList []Operation

	// Iterate over the results and add them to the list
	for resultsIterator.HasNext() {
		result, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate over operation results: %w", err)
		}
		var operation Operation
		err = json.Unmarshal(result.Value, &operation)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal operation: %w", err)
		}
		operationList = append(operationList, operation)
	}

	return operationList, nil
}

func getAlloperations(ctx contractapi.TransactionContextInterface) ([]Operation, error) {
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, fmt.Errorf("error getting operations: %v", err)
	}

	defer func(resultsIterator shim.StateQueryIteratorInterface) {
		err := resultsIterator.Close()
		if err != nil {
			fmt.Printf("error closing iterrator: %v", err)
			return
		}
	}(resultsIterator)

	var operations []Operation
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("error iterating over operations: %w", err)
		}

		var operation operation
		err = json.Unmarshal(queryResponse.Value, &operation)
		if err != nil {
			return nil, fmt.Errorf("error unmarshal: %v", err)
		}

	}

	return operations, nil
}
