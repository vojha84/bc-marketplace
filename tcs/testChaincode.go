package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SampleChaincode struct {
}

//custom data models
type PurchaseOrder struct {
	ID               string `json:"id"`
	ItemID           string `json:"itemId"`
	ProductID        string `json:"productId"`
	LastModifiedDate string `json:"lastModifiedDate"`
	Quantity         int    `json:"quantity"`
	NetValue         int    `json:"netValue"`
}

func GetPurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetPurchaseOrder")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing purchase order ID")
	}

	var purchaseOrderId = args[0]
	bytes, err := stub.GetState(purchaseOrderId)
	if err != nil {
		fmt.Println("Could not fetch purchase order with id "+purchaseOrderId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreatePurchaseOrder(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreatePurchaseOrder")

	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for purchase order creation")
	}

	var purchaseOrderId = args[0]
	var purchaseOrderInput = args[1]

	err := stub.PutState(purchaseOrderId, []byte(purchaseOrderInput))
	if err != nil {
		fmt.Println("Could not save purchase order to ledger", err)
		return nil, err
	}

	var customEvent = "{objectType: 'purchaseOrder', eventType: 'create', payload:'" + purchaseOrderInput + "'}"
	err = stub.SetEvent("eventHub", []byte(customEvent))
	if err != nil {
		return nil, err
	}

	fmt.Println("Successfully saved purchase order")
	return nil, nil

}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Inside INIT for test chaincode")
	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetPurchaseOrder" {
		return GetPurchaseOrder(stub, args)
	}
	return nil, nil
}

func GetCertAttribute(stub shim.ChaincodeStubInterface, attributeName string) (string, error) {
	fmt.Println("Entering GetCertAttribute")
	attr, err := stub.ReadCertAttribute(attributeName)
	if err != nil {
		return "", errors.New("Couldn't get attribute " + attributeName + ". Error: " + err.Error())
	}
	attrString := string(attr)
	return attrString, nil
}

func (t *SampleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "CreatePurchaseOrder" {
		return CreatePurchaseOrder(stub, args)
	} else {
		return nil, errors.New("Invalid function name " + function)
	}
	return nil, nil
}

func main() {
	err := shim.Start(new(SampleChaincode))
	if err != nil {
		fmt.Println("Could not start SampleChaincode")
	} else {
		fmt.Println("SampleChaincode successfully started")
	}

}
