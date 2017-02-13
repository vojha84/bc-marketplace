package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type SampleChaincode struct {
}

//custom data models
type WayBill struct {
	ID               string `json:"id"`
	AssetType        string `json:"assetType"`
	LastModifiedDate string `json:"lastModifiedDate"`
	Quantity         int    `json:"quantity"`
}

func GetWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering GetWayBill")

	if len(args) < 1 {
		fmt.Println("Invalid number of arguments")
		return nil, errors.New("Missing way bill ID")
	}

	var wayBillId = args[0]
	bytes, err := stub.GetState(wayBillId)
	if err != nil {
		fmt.Println("Could not fetch loan way bill with id "+wayBillId+" from ledger", err)
		return nil, err
	}
	return bytes, nil
}

func CreateWayBill(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Entering CreateWayBill")

	if len(args) < 2 {
		fmt.Println("Invalid number of args")
		return nil, errors.New("Expected atleast two arguments for way bill creation")
	}

	var wayBillId = args[0]
	var wayBillInput = args[1]

	err := stub.PutState(wayBillId, []byte(wayBillInput))
	if err != nil {
		fmt.Println("Could not save way bill to ledger", err)
		return nil, err
	}

	fmt.Println("Successfully saved way bill")
	return nil, nil

}

func (t *SampleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	return nil, nil
}

func (t *SampleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "GetWayBill" {
		return GetWayBill(stub, args)
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
	if function == "CreateWayBill" {
		return CreateWayBill(stub, args)
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

	// var str = `{"id":"wb1","assetType":"chip","lastModifiedDate":"21/09/2016 2:30pm","quantity":10}`
	// stub := shim.NewMockStub("mock", new(SampleChaincode))

	// args := []string{"wb1", str}
	// _, err := stub.MockInvoke("123", "CreateWayBill", args)
	// fmt.Println(err)

	// bytes, err := stub.MockQuery("GetWayBill", []string{"wb1"})
	// var wb WayBill
	// err = json.Unmarshal(bytes, &wb)

	// fmt.Println(wb)

}
