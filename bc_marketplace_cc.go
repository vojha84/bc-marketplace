/*
Copyright 2016 IBM

Licensed under the Apache License, Version 2.0 (the "License")
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Licensed Materials - Property of IBM
© Copyright IBM Corp. 2016
*/
package main

/*import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
    "strings"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    "github.com/vojha84/uuid"

	
)*/


import (
	
	"encoding/json"
	"errors"
	"fmt"
	"time"
	"crypto/x509"
	"encoding/pem"
	"net/http"
	"net/url"
    "io/ioutil"
    "strconv"
	"strings"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

//Key names for array holding all the keys belonging to a particular type
var landKeysName = "landKeys"
var propertyKeysName = "propertyKeys"
var propertyAdKeysName = "propertyAdKeys"
var buyerKeysName = "buyerKeys"
var sellerKeysName = "sellerKeys"
var bankKeysName = "bankKeys"
var appraiserKeysName = "appraiserKeys"
var auditorKeysName = "appraiserKeys"
var maKeysName = "maKeys"
var scKeysName = "scKeys"
var aaKeysName = "aaKeys"
var maLogKeysName = "maLogKeys"

//Blockchain Log Key 
var bcLogsKey = "bcLogsKey"



//Prefixes for keys inside state
var typeLand = "land:"
var typePermit = "permit:"
var typeMortgageApplication = "ma:"
var typeSalesContract = "sc:"
var typeAppraiserApplication = "aa:"
var typeProperty = "prop:"
var typePropertyAd = "propad:"
var typeBuyer = "buyer:"
var typeSeller = "seller:"
var typeBank = "bank:"
var typeAppraiser = "appraiser:"
var typeUser = "user:"
var typeAuditor = "auditor:"
var typeMALog = "malog:"

//==============================================================================================================================
//	 Object types - Each object type is mapped to an integer which we use to compare types
//==============================================================================================================================
const   BUYER int      =  1
const   SELLER int  =  2
const   BANK int =  3
const   APPRAISER int =  4
const   AUDITOR int =  5
const   USER int =  6
const   LAND int =  7
const   PROPERTY int =  8
const   PROPERTYAD int =  9
const   MORTGAGEAPPLICATION int =  10
const   SALESCONTRACT int =  11
const   APPRAISERAPPLICATION int =  12
const   MALOG int =  13

//==============================================================================================================================
//	 Affiliation types - Each object type is mapped to an integer which we use to compare affiliations
//==============================================================================================================================
const   BUYER_A   int   =  1
const   SELLER_A int  =  2
const   BANK_A int =  3
const   APPRAISER_A  int =  4
const   AUDITOR_A int =  5




// MarketplaceChaincode implementation
type MarketplaceChaincode struct {
}

/**
Data structures have been denormalized for the sake of simiplicity keeping performance in mind.
**/
type Land struct {
	ID string `json:"id"`
	Description string `json:"description"`
	Address string `json:"address"`
	OwnerId string `json:"ownerId"`
	LastModifiedDate string `json:"lastModifiedDate"`
	
}

type Property struct {
	ID string `json:"id"`
	LandID string `json:"landId"`
	PermitID string `json:"permitId"`
	Description string `json:"description"`
	Address string `json:"address"`
	OwnerId string `json:"ownerId"`
	RegisteredPrice int `json:"registeredPrice"`
	LastModifiedDate string `json:"lastModifiedDate"`
}

type PropertyAd struct {
	ID string `json:"id"`
	LandID string `json:"landId"`
	PermitID string `json:"permitId"`
	PropertyID string `json:"propertyId"`
	Description string `json:"description"`
	Address string `json:"address"`
	SellerID string `json:"sellerId"`
	BankID string `json:"bankId"`
	ListedPrice int `json:"listedPrice"`
	LastModifiedDate string `json:"lastModifiedDate"`
}

type FinancialInfo struct {
	MonthlySalary int `json:"monthlySalary"`
	OtherIncome int `json:"otherIncome"`
	OtherExpenditure int `json:"otherExpenditure"`
	MonthlyRent int  `json:"monthlyRent"`
	MonthlyLoanPayment int `json:"monthlyLoanPayment"`

}

type PersonalInfo struct {
	Firstname string `json:"firstname"`
	Lastname string `json:"lastname"`
	DOB string `json:"dob"`
	Phone string `json:"phone"`
	Mobile string  `json:"mobile"`
	Email string `json:"email"`

}


type MortgageApplication struct {
	ID string `json:"id"`
	PropertyId string `json:"propertyId"`
	LandId string `json:"landId"`
	PermitId string `json:"permitId"`
	BuyerId string `json:"buyerId"`
	AppraisalApplicationId string `json:"appraiserApplicationId"`
	SalesContractId string `json:"salesContractId"`
	PersonalInfo  PersonalInfo `json:"personalInfo"`
	FinancialInfo  FinancialInfo `json:"financialInfo"`
	Status  string `json:"status"`
	RequestedAmount  int `json:"requestedAmount"`
	FairMarketValue  int `json:"fairMarketValue"`
	ApprovedAmount  int `json:"approvedAmount"`
	ReviewerId  string `json:"reviewerId"`
	LastModifiedDate string `json:"lastModifiedDate"`
}

type SalesContract struct {
	ID string `json:"id"`
	PropertyId string `json:"propertyId"`
	BuyerId string `json:"buyerId"`
	SellerId string `json:"sellerId"`
	ReviewerId string `json:"reviewerId"`
	BuyerSignature string `json:"buyerSignature"`
	SellerSignature string `json:"sellerSignature"`
	Status string `json:"status"`
	Price int `json:"price"`
	LastModifiedDate string `json:"lastModifiedDate"`
}

type AppraiserApplication struct {
	ID string `json:"id"`
	MortgageApplicationId string `json:"mortgageApplicationId"`
	AppraiserId string `json:"appraiserId"`
	ReviewerId string `json:"reviewerId"`
	PropertyId string `json:"propertyId"`
	Status string `json:"status"`
	FairMarketValue int `json:"fairMarketValue"`
	LastModifiedDate string `json:"lastModifiedDate"`

}

//Parent type that buyer, seller, auditor, appraiser 'inherit from'
//Hack to acheive polymorphism in GO. Probably better way. Needs investigating
type User struct {
        ID string `json:"id"`
        Affiliation int `json:"affiliation"`
}

type Buyer struct {
	ID string `json:"id"`
	Affiliation int `json:"affiliation"`
	MortgageApplications[] string `json:"mortgageApplications"`
	SalesContracts[] string `json:"salesContracts"`

}

type Seller struct {
	ID string `json:"id"`
	Affiliation int `json:"affiliation"`
	SalesContracts[] string `json:"salesContracts"`

}

type Bank struct {
	ID string `json:"id"`
	Affiliation int `json:"affiliation"`
	MortgageApplications[] string `json:"mortgageApplications"`
	SalesContracts[] string `json:"salesContracts"`

}

type Auditor struct {
	ID string `json:"id"`
	Affiliation int `json:"affiliation"`
}

type Appraiser struct {
	ID string `json:"id"`
	Affiliation int `json:"affiliation"`
	AppraiserApplications[] string `json:"appraiserApplications"`
	

}

type ECertResponse struct {
	OK string `json:"OK"`
}	

type MAUpdateSchema struct {
	Status string `json:"status"`
	SalesContractId string `json:"salesContractId"`
	FairMarketValue int `json:"fairMarketValue"`
	ApprovedAmount int `json:"approvedAmount"`
}

type AAUpdateSchema struct{
	Status string `json:"status"`
	FairMarketValue int `json:"fairMarketValue"`
}

type SCUpdateSchema struct{
	Status string `json:"status"`
	BuyerSignature string `json:"buyerSignature"`
	SellerSignature string `json:"sellerSignature"`
	Price int `json:"price"`
}


type MALog struct {
	MortgageApplicationId string `json:"mortgageApplicationId"`
	BuyerId string `json:"buyerId"`
	ReviewerId string `json:"reviewerId"`
	Status string `json:"status"`
	Action string `json:"action"`
	Text string `json:"text"`
	Timestamp string `json:"timestamp"`
}

type MALogHolder struct {
	MALogs [] MALog `json:"MALogs"`
}

/**
Generate initial set of land records
**/
func generateLandRecords(stub *shim.ChaincodeStub)([4]Land, error){
	fmt.Println("Entering generateLandRecords")
	nowTime := time.Now()

	var landRecords [4] Land

	land1 := Land{"land1", "Residential area", "Madison Ave, New York, Ny", "jack24", nowTime.Format("2006-01-02 15:04:05")}
	land2 := Land{"land2", "Residential area", "Fremont, California, CA", "mark14", nowTime.Format("2006-01-02 15:04:05")}
	land3 := Land{"land3", "Residential area", "San Francisco, California, CA", "jane24", nowTime.Format("2006-01-02 15:04:05")}
	land4 := Land{"land4", "Residential area", "Los Angeles, California, CA", "bill24", nowTime.Format("2006-01-02 15:04:05")}

	landRecords[0] = land1
	landRecords[1] = land2
	landRecords[2] = land3
	landRecords[3] = land4
	
	var landKeys[4]string

	for j := 0; j < len(landRecords); j++ {
      fmt.Println(landRecords[j] )

      lBytes, _ := json.Marshal(&landRecords[j])

      err := stub.PutState(typeLand+landRecords[j].ID, lBytes )
      if err!=nil{
      	fmt.Println("generateLandRecords: Could not save land record")
      	return landRecords, err
      }
      landKeys[j] = typeLand+landRecords[j].ID
   }
	
	landKeyBytes, _ := json.Marshal(&landKeys)

    err:= stub.PutState(landKeysName, landKeyBytes )
    if err!=nil{
      	fmt.Println("generateLandRecords: Could not save land records")
      	return landRecords, err
      }

	return landRecords, nil

}


/**
Generate list of registered properties
**/
func generatePropertyList(stub *shim.ChaincodeStub)([4]Property, error){
	fmt.Println("Entering generatePropertyList")
	nowTime := time.Now()

	var propertyList [4] Property

	property1 := Property{"property1", "land1", "permit1",  "Residential House", "4305 22nd street, Flushing, New York, Ny", "jack24", 500000, nowTime.Format("2006-01-02 15:04:05")}
	property2 := Property{"property2", "land2", "permit2",  "Residential House", "2156 Madison Ave, New York, Ny", "mark14", 500000, nowTime.Format("2006-01-02 15:04:05")}
	property3 := Property{"property3", "land3", "permit3",  "Residential House", "660 Madison Ave, New York, Ny", "jane24",  500000, nowTime.Format("2006-01-02 15:04:05")}
	property4 := Property{"property4", "land4", "permit4",  "Residential House", "200 Madison Ave, New York, Ny", "bill24",  500000, nowTime.Format("2006-01-02 15:04:05")}


	propertyList[0] = property1
	propertyList[1] = property2
	propertyList[2] = property3
	propertyList[3] = property4
	
	var pKeys[4]string

	for j := 0; j < len(propertyList); j++ {
      fmt.Println(propertyList[j] )
       pBytes, _ := json.Marshal(&propertyList[j])

      err := stub.PutState(typeProperty+propertyList[j].ID, pBytes )
      if err!=nil{
      	fmt.Println("generatePropertyList: Could not save property record")
      	return propertyList, err
      }
      pKeys[j] = typeProperty+propertyList[j].ID
   }
	
	pKeyBytes, _ := json.Marshal(&pKeys)

    err := stub.PutState(propertyKeysName, pKeyBytes )
    if err!=nil{
      	fmt.Println("generatePropertyList: Could not save property list")
      	return propertyList, err
    }

	return propertyList, nil

}



/**
Generate list of Properties for sale
**/
func generatePropertyAdsList(stub *shim.ChaincodeStub)([4]PropertyAd, error){
	fmt.Println("Entering generatePropertyAdsList")
	nowTime := time.Now()

	var propertyAds [4] PropertyAd


	propertyAd1 := PropertyAd{"propertyAd1", "land1", "permit1", "property1", "description", "	", "jack24", "Bank Of America", 1000000, nowTime.Format("2006-01-02 15:04:05")}
	propertyAd2 := PropertyAd{"propertyAd2", "land2", "permit2", "property2", "description", "2156 Madison Ave, Apartment no: 202, New York, Ny", "mark14", "Wells Fargo Mortgage", 1500000, nowTime.Format("2006-01-02 15:04:05")}
	propertyAd3 := PropertyAd{"propertyAd3", "land3", "permit3", "property3", "description","660 Madison Ave, Apartment no: 302, New York, Ny", "jane24", "CitiMortgage", 2000000, nowTime.Format("2006-01-02 15:04:05")}
	propertyAd4 := PropertyAd{"propertyAd4", "land4", "permit4", "property4", "description","200 Madison Ave, Apartment no: 402, New York, Ny", "bill24", "JP Morgan", 2500000, nowTime.Format("2006-01-02 15:04:05")}


	propertyAds[0] = propertyAd1
	propertyAds[1] = propertyAd2
	propertyAds[2] = propertyAd3
	propertyAds[3] = propertyAd4
	
	var paKeys[4]string

	for j := 0; j < len(propertyAds); j++ {
      fmt.Println(propertyAds[j] )
      paBytes, _ := json.Marshal(&propertyAds[j])

      err := stub.PutState(typePropertyAd+propertyAds[j].ID, paBytes )
       if err!=nil{
      	fmt.Println("generatePropertyAdsList: Could not save property ad %s", err)
      	return propertyAds, err
    	}
      paKeys[j] = typePropertyAd+propertyAds[j].ID;
   }
	
	paKeyBytes, _ := json.Marshal(&paKeys)

    err := stub.PutState(propertyAdKeysName, paKeyBytes )
	 if err!=nil{
      	fmt.Println("generatePropertyAdsList: Could not save property ads list %s", err)
      	return propertyAds, err
    }

	return propertyAds, nil

}


func InitKeys(stub *shim.ChaincodeStub, keyType string)([]byte, error){
	fmt.Println("Entering InitKeys")

	var keys[]string

	keysBytes, _ := json.Marshal(&keys)
	err := stub.PutState(keyType, keysBytes)
    if err != nil {
        fmt.Println("Failed to initialize key collection "+keyType)
        return nil, nil
    }

	fmt.Println("Initialization complete")
	return nil, nil
}



//==============================================================================================================================
//	 get_ecert - Takes the name passed and calls out to the REST API for HyperLedger to retrieve the ecert
//				 for that user. Returns the ecert as retrived including html encoding.
//==============================================================================================================================
func GetEcert(stub *shim.ChaincodeStub, name string) ([]byte, error) {
	
	var cert ECertResponse
	
	peer_address, err := stub.GetState("Peer_Address")
															if err != nil { return nil, errors.New("Error retrieving peer address") }

	response, err := http.Get("http://"+string(peer_address)+"/registrar/"+name+"/ecert") 	// Calls out to the HyperLedger REST API to get the ecert of the user with that name
    
															if err != nil { return nil, errors.New("Error calling ecert API") }
	
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)					// Read the response from the http callout into the variable contents
	
															if err != nil { return nil, errors.New("Could not read body") }
	
	err = json.Unmarshal(contents, &cert)
	
															if err != nil { return nil, errors.New("Could not retrieve ecert for user: "+name) }
															
	return []byte(string(cert.OK)), nil
}

//==============================================================================================================================
//	 get_caller - Retrieves the username of the user who invoked the chaincode.
//				  Returns the username as a string.
//==============================================================================================================================


func GetUsername(stub *shim.ChaincodeStub) (string, error) {

	bytes, err := stub.GetCallerCertificate();
															if err != nil { return "", errors.New("Couldn't retrieve caller certificate") }
	x509Cert, err := x509.ParseCertificate(bytes);				// Extract Certificate from result of GetCallerCertificate						
															if err != nil { return "", errors.New("Couldn't parse certificate")	}
															
	return x509Cert.Subject.CommonName, nil
}

//==============================================================================================================================
//	 check_affiliation - Takes an ecert as a string, decodes it to remove html encoding then parses it and checks the
// 				  		certificates common name. The affiliation is stored as part of the common name.
//==============================================================================================================================

func CheckAffiliation(stub *shim.ChaincodeStub, cert string) (int, error) {																																																					
	
	decodedCert, err := url.QueryUnescape(cert);    				// make % etc normal //
	
															if err != nil { return -1, errors.New("Could not decode certificate") }
	
	pem, _ := pem.Decode([]byte(decodedCert))           				// Make Plain text   //

	x509Cert, err := x509.ParseCertificate(pem.Bytes);				// Extract Certificate from argument //
														
															if err != nil { return -1, errors.New("Couldn't parse certificate")	}

	cn := x509Cert.Subject.CommonName
	
	res := strings.Split(cn,"\\")
	
	affiliation, _ := strconv.Atoi(res[2])
	
	return affiliation, nil
}

//==============================================================================================================================
//	 get_caller_data - Calls the get_ecert and check_role functions and returns the ecert and role for the
//					 name passed.
//==============================================================================================================================

func GetCallerMetadata(stub *shim.ChaincodeStub) (string, int, error){

	fmt.Println("Entering GetCallerMetadata")

	username, err := GetUsername(stub)
		if err != nil {
			fmt.Println("GetCallerMetadata: Could not get username %s", err); 
			return "", -1, err 
		}

		fmt.Println("USER: ")
		fmt.Println(username)


	user, err := GetUser(stub, username)
	if err !=nil {
		fmt.Println("GetCallerMetadata: Could not get user with ID: %s %s", username, err); 
		return "", -1, err
	}

	affiliation := user.Affiliation

	return username, affiliation, nil
}


/**
Fetch list of property Ads
**/
func GetPropertyAds(stub *shim.ChaincodeStub) ([]PropertyAd, []byte, error){
	
	var PropertyAds []PropertyAd
	
	// Get list of all the keys
	keysBytes, err := stub.GetState(propertyAdKeysName)
	if err != nil {
		fmt.Println("Error retrieving property ad keys")
		return PropertyAds, nil, err
	}

	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling property ad keys")
		return PropertyAds, nil, err
	}

	// Get all the keys
	for _, value := range keys {
		paBytes, err := stub.GetState(value)
		
		var pa PropertyAd
		err = json.Unmarshal(paBytes, &pa)
		if err != nil {
			fmt.Println("Error retrieving property ad " + value)
			return PropertyAds, nil, err
		}
		
		fmt.Println("Appending property ad " + value)
		PropertyAds = append(PropertyAds, pa)
	}	
	
	bytes, err := json.Marshal(&PropertyAds)
	if err != nil {
		fmt.Println("Error marshalling property ads ",err)
		return PropertyAds, nil, err
	}

	return PropertyAds, bytes,  nil
}


/**
Get property ad by id
**/
func GetPropertyAd(stub *shim.ChaincodeStub, id string) (PropertyAd, []byte, error){
	var pa PropertyAd

	pid, err := GetStateKey(id, PROPERTYAD)
	if err != nil {
		fmt.Println("Error key for property ad ", err)
		return pa, nil, err
	}

	paBytes, err := stub.GetState(pid)
	if err != nil {
		fmt.Println("Error retrieving property ad " ,err)
		return pa, nil, err
	}
		
	err = json.Unmarshal(paBytes, &pa)
	if err != nil {
		fmt.Println("Error unmarshalling property ad " ,err)
		return pa, nil, err
	}
		
	return pa, paBytes, nil
}


/**
Fetch list of all mortgage applications for a user
**/

func GetMortgageApplications(stub *shim.ChaincodeStub, callerId string, callerAffiliation int, args[]string) ([]byte, error){
	fmt.Println("Entering GetMortgageApplications")

	if callerAffiliation == BUYER_A || callerAffiliation == BANK_A {
		key, err := GetStateKey(callerId, USER)
		var mas[] string
		var mortgageApplications[] MortgageApplication

		if callerAffiliation == BUYER_A {
			
			var user Buyer
			bytes, err := stub.GetState(key)
			if err != nil{
				fmt.Println("GetMortgageApplications: Could not get bytes for buyer ",err)
				return nil, err
			}
			err = json.Unmarshal(bytes, &user)
			if err != nil{
				fmt.Println("GetMortgageApplications: Could not unmarshal buyer ",err)
				return nil, err
			}
			mas = user.MortgageApplications

		}else if callerAffiliation == BANK_A {
			
			var user Bank
			bytes, err := stub.GetState(key)
			if err != nil{
				fmt.Println("GetMortgageApplications: Could not get bytes for bank ",err)
				return nil, err
			}
			err = json.Unmarshal(bytes, &user)
			if err != nil{
				fmt.Println("GetMortgageApplications: Could not unmarshal bank ",err)
				return nil, err
			}
			mas = user.MortgageApplications

		}

		for i:=0 ; i<len(mas) ; i++ {
			ma, _, err := GetMortgageApplication(stub, callerId, callerAffiliation, []string{mas[i]})
			if err != nil{
				fmt.Println("GetMortgageApplications: Could not get mortgageApplication for id: "+mas[i]+" ",err)
				return nil, err
			}
			mortgageApplications = append(mortgageApplications, ma)
		}

		masBytes, err := json.Marshal(&mortgageApplications)
		if err != nil{
				fmt.Println("GetMortgageApplications: Could not marshal mas bytes ",err)
				return nil, err
		}

		return masBytes, nil

	}

	return nil, errors.New("GetMortgageApplications: callerId "+callerId+ " cannot access mortgage applications")
}

/**
Fetch list of all appraiser applications for a user
**/

func GetAppraiserApplications(stub *shim.ChaincodeStub, callerId string, callerAffiliation int, args[]string) ([]byte, error){
	fmt.Println("Entering GetAppraiserApplications")

	if callerAffiliation == APPRAISER_A  {
		key, err := GetStateKey(callerId, USER)
		var mas[] string
		var appraiserApplications[] AppraiserApplication
			
			var user Appraiser
			bytes, err := stub.GetState(key)
			if err != nil{
				fmt.Println("GetAppraiserApplications: Could not get bytes for buyer ",err)
				return nil, err
			}
			err = json.Unmarshal(bytes, &user)
			if err != nil{
				fmt.Println("GetAppraiserApplications: Could not unmarshal buyer ",err)
				return nil, err
			}
			mas = user.AppraiserApplications

		for i:=0 ; i<len(mas) ; i++ {
			ma, _, err := GetAppraiserApplication(stub, callerId, callerAffiliation, []string{mas[i]})
			if err != nil{
				fmt.Println("GetAppraiserApplications: Could not get appraiserApplication for id: "+mas[i]+" ",err)
				return nil, err
			}
			appraiserApplications = append(appraiserApplications, ma)
		}

		masBytes, err := json.Marshal(&appraiserApplications)
		if err != nil{
				fmt.Println("GetAppraiserApplications: Could not marshal mas bytes ",err)
				return nil, err
		}

		return masBytes, nil

	}

	return nil, errors.New("GetAppraiserApplications: callerId "+callerId+ " cannot access appraiser applications")
}

/**
Fetch list of sales contracts for a user
**/
func GetSalesContracts(stub *shim.ChaincodeStub, callerId string, callerAffiliation int, args[]string) ([]byte, error){
	fmt.Println("Entering GetSalesContracts")

	if callerAffiliation == BUYER_A || callerAffiliation == BANK_A || callerAffiliation == SELLER_A {
		key, err := GetStateKey(callerId, USER)
		var mas[] string
		var salesContracts[] SalesContract

		if callerAffiliation == BUYER_A {
			
			var user Buyer
			bytes, err := stub.GetState(key)
			if err != nil{
				fmt.Println("GetSalesContracts: Could not get bytes for buyer ",err)
				return nil, err
			}
			err = json.Unmarshal(bytes, &user)
			if err != nil{
				fmt.Println("GetSalesContracts: Could not unmarshal buyer ",err)
				return nil, err
			}
			mas = user.SalesContracts

		}else if callerAffiliation == BANK_A {
			
			var user Bank
			bytes, err := stub.GetState(key)
			if err != nil{
				fmt.Println("GetSalesContracts: Could not get bytes for bank ",err)
				return nil, err
			}
			err = json.Unmarshal(bytes, &user)
			if err != nil{
				fmt.Println("GetSalesContracts: Could not unmarshal bank ",err)
				return nil, err
			}
			mas = user.SalesContracts

		}else if callerAffiliation == SELLER_A {
			
			var user Seller
			bytes, err := stub.GetState(key)
			if err != nil{
				fmt.Println("GetSalesContracts: Could not get bytes for seller ",err)
				return nil, err
			}
			err = json.Unmarshal(bytes, &user)
			if err != nil{
				fmt.Println("GetSalesContracts: Could not unmarshal seller ",err)
				return nil, err
			}
			mas = user.SalesContracts

		}

		for i:=0 ; i<len(mas) ; i++ {
			ma, _, err := GetSalesContract(stub, callerId, callerAffiliation, []string{mas[i]})
			if err != nil{
				fmt.Println("GetSalesContracts: Could not get sales contract for id: "+mas[i]+" ",err)
				return nil, err
			}
			salesContracts = append(salesContracts, ma)
		}

		masBytes, err := json.Marshal(&salesContracts)
		if err != nil{
				fmt.Println("GetSalesContracts: Could not marshal mas bytes ",err)
				return nil, err
		}

		return masBytes, nil

	}

	return nil, errors.New("GetSalesContracts: callerId "+callerId+ " cannot access sales contracts")
}

/**
Fetch list of all sales contracts for a user
**/
/*func (t *MarketplaceChaincode) GetSalesContracts(stub *shim.ChaincodeStub, id string, affiliation int) ([]byte, error){
	
	fmt.Println("Entering GetSalesContracts")

	var salesContracts []string
	var key = "";


	if(affiliation == BUYER_A){
		key = buyerKeysName
	}else if(affiliation == SELLER_A){
		key = sellerKeysName
	}else if(affiliation == BANK_A){
		key = bankKeysName
	}else if(affiliation == AUDITOR_A){
		key = auditorKeysName
	}else{
		fmt.Println("GetSalesContracts: Invalid affiliation")
		return nil, errors.New("Invalid affiliation")
	}
		
	// Get list of all the keys representing user
	keysBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving user keys")
		return nil, errors.New("Error retrieving user keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling user keys")
		return nil, errors.New("Error unmarshalling user keys")
	}

	// Get all the users
	for _, value := range keys {
		uBytes, err := stub.GetState(value)
		
		if(affiliation == BUYER_A){
			var user Buyer
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling buyer " + value)
				return nil, errors.New("Error unmarshalling buyer " + value)
			}

			if user.ID == id {
				salesContracts = user.SalesContracts
				scBytes, _ := json.Marshal(&salesContracts)
				return scBytes, nil
			}
			
		}else if(affiliation == SELLER_A){
			var user Bank
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling bank " + value)
				return nil, errors.New("Error unmarshalling bank " + value)
			}

			if user.ID == id {
				salesContracts = user.SalesContracts
				scBytes, _ := json.Marshal(&salesContracts)

				return scBytes, nil
			}

		}else if(affiliation == BANK_A){
			var user Bank
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling bank " + value)
				return nil, errors.New("Error unmarshalling bank " + value)
			}

			if user.ID == id {
				salesContracts = user.SalesContracts
				scBytes, _ := json.Marshal(&salesContracts)

				return scBytes, nil
			}

		}else if(affiliation == AUDITOR_A){
			var user Auditor
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling auditor " + value)
				return nil, errors.New("Error unmarshalling auditor " + value)
			}

			if user.ID == id {
				salesContracts = user.SalesContracts
				scBytes, _ := json.Marshal(&salesContracts)
				return scBytes, nil
			}

			
		}

		
		
	}	
	
	fmt.Println("GetSalesContracts: User with id "+id+ "not found")
	return nil, nil
}*/


/**
Get mortgage application for a particular user by id
**/
/*func (t *MarketplaceChaincode) GetSalesContract(stub *shim.ChaincodeStub, userId string, affiliation int, contractId string) ([]byte, error){
	
	fmt.Println("Entering GetSalesContract")

	var mortgageApplication MortgageApplication
	var key = "";


	if(affiliation == BUYER_A){
		key = buyerKeysName
	}else if(affiliation == BANK_A){
		key = bankKeysName
	}else if(affiliation == AUDITOR_A){
		key = auditorKeysName
	}else{
		fmt.Println("GetMortgageApplications: Invalid affiliation")
		return nil, errors.New("Invalid affiliation")
	}
		
	// Get list of all the keys representing user
	keysBytes, err := stub.GetState(key)
	if err != nil {
		fmt.Println("Error retrieving user keys")
		return nil, errors.New("Error retrieving user keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling user keys")
		return nil, errors.New("Error unmarshalling user keys")
	}

	// Get all the users
	for _, value := range keys {
		uBytes, err := stub.GetState(value)
		
		if(affiliation == BUYER_A){
			var user Buyer
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling buyer " + value)
				return nil, errors.New("Error unmarshalling buyer " + value)
			}

			if user.ID == userId {
				mortgageApplications := user.MortgageApplications
				for _, ma := range mortgageApplications{
					if ma == mortgageId {
						maBytes, _ := json.Marshal(&ma)
						return maBytes, nil
					}
				}
				
			}
			
		}else if(affiliation == BANK_A){
			var user Bank
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling bank " + value)
				return nil, errors.New("Error unmarshalling bank " + value)
			}

			if user.ID == userId {
				mortgageApplications := user.MortgageApplications
				for _, ma := range mortgageApplications{
					if ma.ID == mortgageId {
						maBytes, _ := json.Marshal(&ma)
						return maBytes, nil
					}
				}
			}

		}else if(affiliation == AUDITOR_A){
			var user Auditor
			err = json.Unmarshal(uBytes, &user)
			if err != nil {
				fmt.Println("Error unmarshalling auditor " + value)
				return nil, errors.New("Error unmarshalling auditor " + value)
			}

			if user.ID == id {
				mortgageApplications := user.mortgageApplications
				for _, ma := range mortgageApplications{
					if ma.ID == mortgageId {
						maBytes, _ := json.Marshal(&ma)
						return maBytes, nil
					}
				}
			}

			
		}

		
		
	}	
	
	fmt.Println("GetMortgageApplications: Mortgage Application with id "+mortgageId+ " not found for user with id "+userId)
	return nil, nil
}*/

/**
Create a new mortgage application
Add the key to the list of mortgage application keys
Add the key to the list of mortgage applications keys associated with a user
**/


/*func CreateMortgageApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering CreateMortgageApplication")

	mortgageApplicationInput := args[0]
	var mortgageApplication MortgageApplication

	mortgageApplicationId, _ := GetUUID()
	fmt.Println(mortgageApplicationId)
	
	err := json.Unmarshal([]byte(mortgageApplicationInput), &mortgageApplication)
	if err != nil {
		fmt.Println("Error unmarshalling mortgageApplicationInput string " + mortgageApplicationInput+" %s",err)
		return nil, errors.New("Error unmarshalling mortgageApplicationInput string " + mortgageApplicationInput)
	}

	fmt.Println("HEEEEEEEEEEEEEEEEEEE ")
	mortgageApplication.ID = mortgageApplicationId

	fmt.Println("CreateMortgageApplication: Created new Mortgage Application with ID "+mortgageApplicationId)
	fmt.Println(mortgageApplication)

	maBytes, _ := json.Marshal(&mortgageApplication)

	maKey, err  := GetStateKey(mortgageApplication.ID, MORTGAGEAPPLICATION)

	fmt.Println("Generated mortgageApplication key "+maKey)

	err = stub.PutState(maKey, maBytes)
	if err != nil {
		fmt.Println("Error saving mortgageApplication "+mortgageApplication.ID +" to state")
		return nil, errors.New("Error saving mortgageApplication "+mortgageApplication.ID +" to state")
	}

	ok, err := AddMortgageApplicationKey(stub, maKey)

	fmt.Println(ok)

	if err != nil {
		return nil, err
	}

	userKey, err := GetStateKey(callerId, BUYER)
	
	user, err := GetBuyer(stub, userKey)		

	mas := user.MortgageApplications
	user.MortgageApplications = append(mas, maKey)

	err = SaveBuyer(stub, user, userKey)

	if err != nil {	
		fmt.Printf("CreateMortgageApplication: Failed to store updated user with id"+ userKey + ": %s", err)
		return nil, errors.New("CreateMortgageApplication: Failed to store updated user with id"+ userKey ) 
	}

	fmt.Println("Here...............")
	
	return maBytes, nil
}
*/

func CreateMortgageApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering CreateMortgageApplication")

	if len(args) < 2 {
		fmt.Println("CreateMortgageApplication: expected two arguments")
		return nil, errors.New("Could not create MortgageApplication. Invalid input")
	}

	mortgageApplicationId := args[0]
	mortgageApplicationInput := args[1]

	maKey, err  := GetStateKey(mortgageApplicationId, MORTGAGEAPPLICATION)

	fmt.Println("Generated mortgageApplication key "+maKey)

	var ma MortgageApplication
	err = json.Unmarshal([]byte(mortgageApplicationInput), &ma)
	if err !=nil {
		fmt.Println("CreateMortgageApplication: Could not unmarshal mortgageApplicationInput", err)
		return nil, err
	}

	bankId := ma.ReviewerId

	err = stub.PutState(maKey, []byte(mortgageApplicationInput))
	if err != nil {
		fmt.Println("Error saving mortgageApplication "+mortgageApplicationId +" to state", err)
		return nil, err
	}

	ok, err := AddKey(stub, maKey, maKeysName)

	fmt.Println(ok)

	if err != nil {
		return nil, err
	}

	userKey, err := GetStateKey(callerId, USER)
	
	user, err := GetBuyer(stub, userKey)		

	mas := user.MortgageApplications
	//Store the external mortgage application id generated by front end as foreign key in user
	user.MortgageApplications = append(mas, mortgageApplicationId)

	err = SaveBuyer(stub, user, userKey)

	if err != nil {	
		fmt.Printf("CreateMortgageApplication: Failed to store updated user with id"+ userKey + ": ", err)
		return nil, err
	}

	bankKey, err := GetStateKey(bankId, USER)

	bank, err := GetBank(stub, bankKey)		

	bmas := bank.MortgageApplications
	//Store the external mortgage application id generated by front end as foreign key in user
	bank.MortgageApplications = append(bmas, mortgageApplicationId)

	err = SaveBank(stub, bank, bankKey)

	if err != nil {	
		fmt.Printf("CreateMortgageApplication: Failed to store updated bank with id"+ bankKey + ": ", err)
		return nil, err
	}


	fmt.Println("CreateMortgageApplication: Successfully created and stored mortgageApplication with ID: "+mortgageApplicationId)

	AppendMALog(stub, "CreateMortgageApplication", callerId+" Submitted new MortgageApplication", "Submitted", mortgageApplicationId)
	
	return nil, nil
}

/**
Return a Mortgage application based on access rights
**/
func GetMortgageApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) (MortgageApplication, []byte, error){
	fmt.Println("Entering GetMortgageApplication")

	var ma MortgageApplication

	if(len(args) < 1 ){
		fmt.Println("CreateMortgageApplication: expected 1 argument")
		return ma, nil, errors.New("Could not create MortgageApplication. Invalid input")
	}

	maId := args[0]

	maKey, err  := GetStateKey(maId, MORTGAGEAPPLICATION)

	fmt.Println("Generated mortgageApplication key "+maKey)

	bytes, err := stub.GetState(maKey)
	if err != nil {
		fmt.Println("GetMortgageApplication: Could not fetch mortgageApplication with ID : "+maId)
		return ma, nil, err
	}

	
	err = json.Unmarshal(bytes, &ma)
	if err != nil {
		fmt.Println("GetMortgageApplication: Could not unmarshal mortgageApplication with ID : "+maId)
		return ma, nil, err
	}

	if callerId == ma.BuyerId || callerId == ma.ReviewerId || callerAffiliation == AUDITOR_A {
		//Caller is permitted to access mortgage application
		return ma, bytes, nil
	}else{
		fmt.Println("GetMortgageApplication: Caller with ID "+callerId+ " and affiliation "+string(callerAffiliation)+" does not have rights to access mortgageApplication")
		return ma, nil, errors.New("User "+callerId+ "does not have rights to access mortgageApplication with id "+maId)
	}

	
}

/**
Updates Mortgage application based on access rights
**/
func UpdateMortgageApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering UpdateMortgageApplication")

	if(len(args) < 2){
		fmt.Println("UpdateMortgageApplication: No parameters provided for update")
		return nil, errors.New("Could not update mortgageApplication. No parameters provided for update ")
	}

	id := args[0]
	var currentStatus string
	var updates MAUpdateSchema
	var statusChanged bool = false
	var scIdChanged bool = false
	var amChanged bool = false

	ma, _, err := GetMortgageApplication(stub, callerId, AUDITOR_A, []string{id})
	if err != nil {	
		return nil, err
	}

	err = json.Unmarshal([]byte(args[1]), &updates )
		if err != nil {
			fmt.Println("UpdateMortgageApplication: Could not unmarshal updates ", err)
			return nil, err
		}

	var msg string

	if callerId == ma.ReviewerId {
		//Valid user to update the application
	
		status := strings.TrimSpace(updates.Status)
		if len(status) > 0{
			currentStatus = ma.Status
			ma.Status = status
			statusChanged = true
			msg += callerId+ " changed status from "+currentStatus+" to "+status
		}

		salesContractId :=  strings.TrimSpace(updates.SalesContractId)
		if len(salesContractId) > 0 {
			ma.SalesContractId = salesContractId
			if statusChanged == true {
				msg += "and updated sales contract Id to "+salesContractId+"."
			}else {
				msg += callerId+" updated sales contract Id to "+salesContractId+"."
			}
			scIdChanged = true
			
		}

		approvedAmount :=  updates.ApprovedAmount
		
		if approvedAmount != 0 {
			ma.ApprovedAmount = approvedAmount
			if statusChanged == true || scIdChanged == true{
				msg += "and updated approved amount to "+strconv.Itoa(approvedAmount)+"."
			}else {
				msg += callerId+" updated approved amount to "+strconv.Itoa(approvedAmount)+"."
			}
			amChanged = true

		}

		if statusChanged == true || scIdChanged == true || amChanged == true{
			bytes, err := SaveMortgageApplication(stub, ma, id)
			if err != nil {
				fmt.Println("SaveMortgageApplication: Could not save mortgageApplication ",err)
				return nil, err
			}
			AppendMALog(stub, "UpdateMortgageApplication", msg, ma.Status, id)
			return bytes, nil
		}else{
			fmt.Println("SaveMortgageApplication: Nothing to update")
			return nil, nil
		}
		
		
		/*if statusChanged == true && scIdChanged == true{
			msg = callerId+ " changed status from "+currentStatus+" to "+status+" and updated sales contract Id: "+salesContractId
		}else if statusChanged == true && scIdChanged == false{
			msg = callerId+ " changed status from "+currentStatus+" to "+status
		}else if statusChanged == false && scIdChanged == true{
			msg = callerId+" updated sales contract Id: "+salesContractId
		}*/

		

	}else if callerAffiliation == APPRAISER_A{
		fairMarketValue :=  updates.FairMarketValue
		
		if fairMarketValue != 0 {
			ma.FairMarketValue = fairMarketValue

			bytes, err := SaveMortgageApplication(stub, ma, id)
			if err != nil {
				fmt.Println("SaveMortgageApplication: Could not save mortgageApplication ",err)
				return nil, err
			}
			AppendMALog(stub, "UpdateMortgageApplication", msg, ma.Status, id)
			return bytes, nil
		}else{
			fmt.Println("SaveMortgageApplication: Nothing to update")
			return nil, nil
		}
	}else{
		fmt.Println("UpdateMortgageApplication: User with id "+callerId+ "does not have rights to update the mortgage application")
		return nil, errors.New("User with id "+callerId+ "does not have rights to update the mortgage application")
	}


}

/**
Create a new Appraiser Application
**/
func CreateAppraiserApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering CreateAppraiserApplication")

	if len(args) < 2 {
		fmt.Println("CreateAppraiserApplication: expected two arguments")
		return nil, errors.New("Could not create CreateAppraiserApplication. Invalid input")
	}

	if callerAffiliation != BANK_A {
		//Caller is not allowed to create an appraiser application
		fmt.Println("CreateAppraiserApplication: "+callerId+" is not allowed to create appraiser application")
		return nil, errors.New(callerId+" is not allowed to create appraiser application")
	}

	appraiserApplicationId := args[0]
	appraiserApplicationInput := args[1]

	maKey, err  := GetStateKey(appraiserApplicationId, APPRAISERAPPLICATION)

	fmt.Println("Generated appraiserApplication key "+maKey)

	err = stub.PutState(maKey, []byte(appraiserApplicationInput))
	if err != nil {
		fmt.Println("Error saving CreateAppraiserApplication "+appraiserApplicationId +" to state")
		return nil, errors.New("Error saving CreateAppraiserApplication "+appraiserApplicationId +" to state")
	}

	var aa AppraiserApplication
	err = json.Unmarshal([]byte(appraiserApplicationInput), &aa)
	if err !=nil {
		fmt.Println("CreateAppraiserApplication: Could not unmarshal appraiserApplicationInput", err)
		return nil, err
	}

	ok, err := AddKey(stub, maKey, aaKeysName)

	fmt.Println(ok)

	if err != nil {
		return nil, err
	}

	userKey, err := GetStateKey(aa.AppraiserId, USER)
	
	user, err := GetAppraiser(stub, userKey)		

	mas := user.AppraiserApplications
	user.AppraiserApplications = append(mas, appraiserApplicationId)

	err = SaveAppraiser(stub, user, userKey)

	if err != nil {	
		fmt.Printf("CreateAppraiserApplication: Failed to store updated user with id"+ userKey + ": %s", err)
		return nil, errors.New("CreateAppraiserApplication: Failed to store updated user with id"+ userKey ) 
	}


	fmt.Println("CreateAppraiserApplication: Successfully created and stored appraiserApplication with ID: "+appraiserApplicationId)

	AppendMALog(stub, "CreateAppraiserApplication", callerId+" Submitted new AppraiserApplication", "Submitted", appraiserApplicationId)
	
	return nil, nil
}

/**
Return a Appraiser application based on access rights
**/
func GetAppraiserApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) (AppraiserApplication, []byte, error){
	fmt.Println("Entering GetAppraiserApplication")

	var ma AppraiserApplication

	if(len(args) < 1 ){
		fmt.Println("GetAppraiserApplication: expected 1 argument")
		return ma, nil, errors.New("Could not GetAppraiserApplication. Invalid input")
	}

	maId := args[0]

	maKey, err  := GetStateKey(maId, APPRAISERAPPLICATION)

	fmt.Println("Generated appraiserApplication key "+maKey)

	bytes, err := stub.GetState(maKey)
	if err != nil {
		fmt.Println("GetAppraiserApplication: Could not fetch appraiserApplication with ID : "+maId)
		return ma, nil, err
	}

	
	err = json.Unmarshal(bytes, &ma)
	if err != nil {
		fmt.Println("GetAppraiserApplication: Could not unmarshal appraiserApplication with ID : "+maId)
		return ma, nil, err
	}

	if callerId == ma.AppraiserId || callerId == ma.ReviewerId || callerAffiliation == AUDITOR_A {
		//Caller is permitted to access mortgage application
		return ma, bytes, nil
	}else{
		fmt.Println("GetAppraiserApplication: Caller with ID "+callerId+ " and affiliation "+string(callerAffiliation)+" does not have rights to access mortgageApplication")
		return ma, nil, errors.New("User "+callerId+ "does not have rights to access appraiserApplication with id "+maId)
	}

	
}

/**
Updates Appraiser application based on access rights
**/
func UpdateAppraiserApplication(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering UpdateAppraiserApplication")

	if(len(args) < 2){
		fmt.Println("UpdateAppraiserApplication: No parameters provided for update")
		return nil, errors.New("Could not update appraiserApplication. No parameters provided for update ")
	}

	id := args[0]
	var currentStatus string
	var updates AAUpdateSchema
	var statusChanged bool = false
	var mvChanged bool = false
	
	ma, _, err := GetAppraiserApplication(stub, callerId, callerAffiliation, []string{id})
	if err != nil {	
		return nil, err
	}

	if callerId == ma.AppraiserId {
		//Valid user to update the application
		err = json.Unmarshal([]byte(args[1]), &updates )
		if err != nil {
			fmt.Println("UpdateAppraiserApplication: Could not unmarshal updates %s", err)
			return nil, err
		}


		status := strings.TrimSpace(updates.Status)
		if len(status) > 0{
			currentStatus = ma.Status
			ma.Status = status
			statusChanged = true
		}

		fairMarketValue :=  updates.FairMarketValue
		
		if fairMarketValue != 0 {
			ma.FairMarketValue = fairMarketValue
			mvChanged = true
		}

		bytes, err := SaveAppraiserApplication(stub, ma, id)
		if err != nil {
			fmt.Println("SaveAppraiserApplication: Could not save appraiser application ",err)
			return nil, err
		}

		bytes, err = UpdateMortgageApplication(stub, callerId, callerAffiliation, []string{ma.MortgageApplicationId,`{"fairMarketValue":`+strconv.Itoa(fairMarketValue)+`}`})
		if err != nil {
			fmt.Println("SaveAppraiserApplication: Could not update mortgage application ",err)
			return nil, err
		}

		var msg string
		var fmvStr string
		if mvChanged == true {
			fmvStr = strconv.Itoa(fairMarketValue)
		}

		if statusChanged == true && mvChanged == true{
			msg = callerId+ " changed status from "+currentStatus+" to "+status+" and updated fair market value: "+fmvStr
		}else if statusChanged == true && mvChanged == false{
			msg = callerId+ " changed status from "+currentStatus+" to "+status
		}else if statusChanged == false && mvChanged == true{
			msg = callerId+" updated fair market value: "+fmvStr
		}

		AppendMALog(stub, "UpdateAppraiserApplication", msg, status, id)
		return bytes, nil

	}else{
		fmt.Println("UpdateAppraiserApplication: User with id "+callerId+ "does not have rights to update the appraiser application")
		return nil, errors.New("User with id "+callerId+ "does not have rights to update the appraiser application")
	}
}


func CreateSalesContract(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering CreateSalesContract")

	if len(args) < 2 {
		fmt.Println("CreateSalesContract: expected two arguments")
		return nil, errors.New("Could not create CreateSalesContract. Invalid input")
	}

	if callerAffiliation != BUYER_A {
		//Caller is not allowed to create an sales contract
		fmt.Println("CreateSalesContract: "+callerId+" is not allowed to create seller contract")
		return nil, errors.New(callerId+" is not allowed to create seller contract")
	}

	salesContractId := args[0]
	salesContractInput := args[1]

	maKey, err  := GetStateKey(salesContractId, SALESCONTRACT)

	fmt.Println("Generated salesContract key "+maKey)

	var sc SalesContract
	err = json.Unmarshal([]byte(salesContractInput), &sc)
	if err !=nil {
		fmt.Println("CreateSalesContract: Could not unmarshal salesContractInput", err)
		return nil, err
	}
	
	sellerId := sc.SellerId
	bankId := sc.ReviewerId

	err = stub.PutState(maKey, []byte(salesContractInput))
	if err != nil {
		fmt.Println("Error saving CreateSalesContract "+salesContractId +" to state")
		return nil, errors.New("Error saving CreateSalesContract "+salesContractId +" to state")
	}

	ok, err := AddKey(stub, maKey, scKeysName)

	fmt.Println(ok)

	if err != nil {
		return nil, err
	}

	userKey, err := GetStateKey(sellerId, USER)
	
	user, err := GetSeller(stub, userKey)		

	mas := user.SalesContracts
	user.SalesContracts = append(mas, salesContractId)

	err = SaveSeller(stub, user, userKey)

	if err != nil {	
		fmt.Printf("CreateSalesContract: Failed to store updated user with id"+ userKey + ": %s", err)
		return nil, errors.New("CreateSalesContract: Failed to store updated user with id"+ userKey ) 
	}

	buyerKey, err := GetStateKey(callerId, USER)
	
	buyer, err := GetBuyer(stub, buyerKey)		

	bmas := buyer.SalesContracts
	buyer.SalesContracts = append(bmas, salesContractId)

	err = SaveBuyer(stub, buyer, buyerKey)

	if err != nil {	
		fmt.Printf("CreateSalesContract: Failed to store updated user with id"+ buyerKey + ": %s", err)
		return nil, errors.New("CreateSalesContract: Failed to store updated user with id"+ buyerKey ) 
	}

	bankKey, err := GetStateKey(bankId, USER)
	
	bank, err := GetBank(stub, bankKey)		

	bas := bank.SalesContracts
	bank.SalesContracts = append(bas, salesContractId)

	err = SaveBank(stub, bank, userKey)

	if err != nil {	
		fmt.Printf("CreateSalesContract: Failed to store updated user with id"+ bankKey + ": %s", err)
		return nil, errors.New("CreateSalesContract: Failed to store updated user with id"+ bankKey ) 
	}


	fmt.Println("CreateSalesContract: Successfully created and stored salesContract with ID: "+salesContractId)

	AppendMALog(stub, "CreateSalesContract", callerId+" Submitted new SalesContract", "Submitted", salesContractId)
	
	return nil, nil
}

/**
Return a Seller application based on access rights
**/
func GetSalesContract(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) (SalesContract, []byte, error){
	fmt.Println("Entering GetSalesContract")

	var ma SalesContract

	if(len(args) < 1 ){
		fmt.Println("GetSalesContract: expected 1 argument")
		return ma, nil, errors.New("Could not GetSalesContract. Invalid input")
	}

	maId := args[0]

	maKey, err  := GetStateKey(maId, SALESCONTRACT)

	fmt.Println("Generated salesContract key "+maKey)

	bytes, err := stub.GetState(maKey)
	if err != nil {
		fmt.Println("GetSalesContract: Could not fetch salesContract with ID : "+maId)
		return ma, nil, err
	}

	
	err = json.Unmarshal(bytes, &ma)
	if err != nil {
		fmt.Println("GetSalesContract: Could not unmarshal salesContract with ID : "+maId)
		return ma, nil, err
	}

	if callerId == ma.SellerId || callerId == ma.BuyerId || callerAffiliation == AUDITOR_A || callerAffiliation == BANK_A {
		//Caller is permitted to access sales contract
		return ma, bytes, nil
	}else{
		fmt.Println("GetSalesContract: Caller with ID "+callerId+ " and affiliation "+strconv.Itoa(callerAffiliation)+" does not have rights to access mortgageContract")
		return ma, nil, errors.New("User "+callerId+ "does not have rights to access salesContract with id "+maId)
	}

	
}

/**
Updates Seller application based on access rights
**/
func UpdateSalesContract(stub *shim.ChaincodeStub, callerId string, callerAffiliation int , args[]string) ([]byte, error){
	fmt.Println("Entering GetSalesContract")

	if(len(args) < 2){
		fmt.Println("UpdateSalesContract: No parameters provided for update")
		return nil, errors.New("Could not update salesContract. No parameters provided for update ")
	}

	id := args[0]
	var currentStatus string
	var updates SCUpdateSchema
	
	
	ma, _, err := GetSalesContract(stub, callerId, callerAffiliation, []string{id})
	if err != nil {	
		return nil, err
	}

	if callerId == ma.SellerId || callerId == ma.BuyerId {
		//Valid user to update the contract
		err = json.Unmarshal([]byte(args[1]), &updates )
		if err != nil {
			fmt.Println("UpdateSalesContract: Could not unmarshal updates %s", err)
			return nil, err
		}

		var logs[] string

		status := strings.TrimSpace(updates.Status)
		if len(status) > 0{
			currentStatus = ma.Status
			ma.Status = status
			logs = append(logs, "changed status from "+currentStatus+" to "+status+"")
		}

		bs := strings.TrimSpace(updates.BuyerSignature)
		if len(bs) > 0{
			ma.BuyerSignature = bs
			logs = append(logs, "Buyer: "+ma.BuyerId+" Signed")
		}

		ss := strings.TrimSpace(updates.SellerSignature)
		if len(ss) > 0{
			ma.BuyerSignature = ss
			logs = append(logs, "Seller: "+ma.SellerId+" Signed")
		}

		price := updates.Price
		if price != 0{
			ma.Price = price
			logs = append(logs, "Price updated to: "+strconv.Itoa(price))
		}


		bytes, err := SaveSalesContract(stub, ma, id)
		if err != nil {
			return nil, err
		}

		var msg string
		for _, log := range logs{
			msg+= " "+log
		}

		AppendMALog(stub, "UpdateSalesContract", msg, status, id)
		return bytes, nil

	}else{
		fmt.Println("UpdateSalesContract: User with id "+callerId+ "does not have rights to update the seller application")
		return nil, errors.New("User with id "+callerId+ "does not have rights to update the seller application")
	}
}


/**
Save Mortgage Application to the ledger
**/
func SaveMortgageApplication(stub *shim.ChaincodeStub, ma MortgageApplication, id string)([]byte, error){
	fmt.Println("Entering SaveMortgageApplication")
	if &ma != nil {
		bytes, _ := json.Marshal(&ma)
		maKey, err  := GetStateKey(id, MORTGAGEAPPLICATION)
		err = stub.PutState(maKey, bytes)
		if err != nil {
			fmt.Println("SaveMortgageApplication: Could not save mortgage application ", err)
			return nil, err
		}
		return bytes, nil
	}else{
		return nil, errors.New("Invalid mortgageApplication input")
	}

}


/**
Gets the Buyer from the state if it exists or creates a new one
**/
func GetBuyer(stub *shim.ChaincodeStub, id string)(Buyer, error){
	fmt.Println("Entering Buyer")

	var buyer Buyer
	bytes, err := stub.GetState(id)						
				
		if err != nil {	
			fmt.Printf("GetBuyer: Could not get user with id "+ id + ": %s", err)
			return buyer, errors.New("GetBuyer: Failed to get user with id "+ id) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetBuyer: buyer with id does not exist: "+ id + ": %s", err)
			fmt.Println("GetBuyer: creating a buyer with id: "+id)
			
			mas := []string {}
			sc := []string {}
			buyer = Buyer{id, BUYER_A, mas, sc}
			fmt.Println(buyer)

			bytes, err := json.Marshal(&buyer)
			if err != nil {	
				fmt.Printf("GetBuyer: Could not marshal buyer : %s", err)
				return buyer, errors.New("GetBuyer: Could not marshal buyer with id "+id) 
			}

			err = stub.PutState(id, bytes)
			if err != nil {	
				fmt.Printf("GetBuyer: Could not save buyer : %s", err)
				return buyer, errors.New("GetBuyer: Could not save buyer with id "+id) 
			}

			return buyer, nil
		}
	
			err = json.Unmarshal(bytes, &buyer)
			if err != nil {	
				fmt.Printf("GetBuyer: Could not unmarshal buyer : %s", err)
				return buyer, errors.New("GetBuyer: Could not unmarshal buyer with id "+id) 
			}
		

		return buyer, nil
}

func SaveBuyer(stub *shim.ChaincodeStub, buyer Buyer, id string)(error){
	fmt.Println("Entering SaveBuyer")
	if &buyer != nil {
		bytes, _ := json.Marshal(buyer)
		err := stub.PutState(id, bytes)
		if err != nil {
			fmt.Println("SaveBuyer: Could not save buyer %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid buyer input")
	}
}

/**
Gets the Bank from the state if it exists or creates a new one
**/
func GetBank(stub *shim.ChaincodeStub, id string)(Bank, error){
	fmt.Println("Entering GetBank")

	var bank Bank
	bytes, err := stub.GetState(id)						
				
		if err != nil {	
			fmt.Printf("GetBank: Could not get user with id "+ id + ": %s", err)
			return bank, errors.New("GetBank: Failed to get user with id "+ id) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetBank: bank with id does not exist: "+ id + ": %s", err)
			fmt.Println("GetBank: creating a bank with id: "+id)
			
			var mas = []string {}
			var sc = [] string {}
			bank = Bank{id, BANK_A, mas, sc}
			fmt.Println(bank)

			bytes, err := json.Marshal(&bank)
			if err != nil {	
				fmt.Printf("GetBank: Could not marshal bank : %s", err)
				return bank, errors.New("GetBank: Could not marshal bank with id "+id) 
			}

			err = stub.PutState(id, bytes)
			if err != nil {	
				fmt.Printf("GetBank: Could not save bank : %s", err)
				return bank, errors.New("GetBank: Could not save bank with id "+id) 
			}

			return bank, nil
		}
	
			err = json.Unmarshal(bytes, &bank)
			if err != nil {	
				fmt.Printf("GetBank: Could not unmarshal bank : %s", err)
				return bank, errors.New("GetBank: Could not unmarshal bank with id "+id) 
			}
		

		return bank, nil
}

func SaveBank(stub *shim.ChaincodeStub, bank Bank, id string)(error){
	fmt.Println("Entering SaveBank")
	if &bank != nil {
		bytes, _ := json.Marshal(&bank)
		err := stub.PutState(id, bytes)
		if err != nil {
			fmt.Println("SaveBank: Could not save bank %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid bank input")
	}
}

/**
Save Appraiser Application to the ledger
**/
func SaveAppraiserApplication(stub *shim.ChaincodeStub, ma AppraiserApplication, id string)([]byte, error){
	fmt.Println("Entering SaveAppraiserApplication")
	if &ma != nil {
		bytes, _ := json.Marshal(&ma)
		aaKey, err  := GetStateKey(id, APPRAISERAPPLICATION)
		err = stub.PutState(aaKey, bytes)
		if err != nil {
			fmt.Println("SaveAppraiserApplication: Could not save appraiser application %s", err)
			return nil, err
		}
		return bytes, nil
	}else{
		return nil, errors.New("Invalid appraiserApplication input")
	}

}


/**
Gets the Appraiser from the state if it exists or creates a new one
**/
func GetAppraiser(stub *shim.ChaincodeStub, id string)(Appraiser, error){
	fmt.Println("Entering Appraiser")

	var appraiser Appraiser
	bytes, err := stub.GetState(id)						
				
		if err != nil {	
			fmt.Printf("GetAppraiser: Could not get user with id "+ id + ": %s", err)
			return appraiser, errors.New("GetAppraiser: Failed to get user with id "+ id) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetAppraiser: appraiser with id does not exist: "+ id + ": %s", err)
			fmt.Println("GetAppraiser: creating a appraiser with id: "+id)
			
			aa := []string {}
			
			appraiser = Appraiser{id, APPRAISER_A, aa}
			fmt.Println(appraiser)

			bytes, err := json.Marshal(&appraiser)
			if err != nil {	
				fmt.Printf("GetAppraiser: Could not marshal appraiser : %s", err)
				return appraiser, errors.New("GetAppraiser: Could not marshal appraiser with id "+id) 
			}

			err = stub.PutState(id, bytes)
			if err != nil {	
				fmt.Printf("GetAppraiser: Could not save appraiser : %s", err)
				return appraiser, errors.New("GetAppraiser: Could not save appraiser with id "+id) 
			}

			return appraiser, nil
		}
	
			err = json.Unmarshal(bytes, &appraiser)
			if err != nil {	
				fmt.Printf("GetAppraiser: Could not unmarshal appraiser : %s", err)
				return appraiser, errors.New("GetAppraiser: Could not unmarshal appraiser with id "+id) 
			}
		

		return appraiser, nil
}

func SaveAppraiser(stub *shim.ChaincodeStub, appraiser Appraiser, id string)(error){
	fmt.Println("Entering SaveAppraiser")
	if &appraiser != nil {
		bytes, _ := json.Marshal(&appraiser)
		err := stub.PutState(id, bytes)
		if err != nil {
			fmt.Println("SaveAppraiser: Could not save appraiser %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid appraiser input")
	}
}

/**
Gets the Seller from the state if it exists or creates a new one
**/
func GetSeller(stub *shim.ChaincodeStub, id string)(Seller, error){
	fmt.Println("Entering Seller")

	var seller Seller
	bytes, err := stub.GetState(id)						
				
		if err != nil {	
			fmt.Printf("GetSeller: Could not get user with id "+ id + ": %s", err)
			return seller, errors.New("GetSeller: Failed to get user with id "+ id) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetSeller: seller with id does not exist: "+ id + ": %s", err)
			fmt.Println("GetSeller: creating a seller with id: "+id)
			
			sc := []string {}
			
			seller = Seller{id, SELLER_A, sc}
			fmt.Println(seller)

			bytes, err := json.Marshal(&seller)
			if err != nil {	
				fmt.Printf("GetSeller: Could not marshal seller : %s", err)
				return seller, errors.New("GetSeller: Could not marshal seller with id "+id) 
			}

			err = stub.PutState(id, bytes)
			if err != nil {	
				fmt.Printf("GetSeller: Could not save seller : %s", err)
				return seller, errors.New("GetSeller: Could not save seller with id "+id) 
			}

			return seller, nil
		}
	
			err = json.Unmarshal(bytes, &seller)
			if err != nil {	
				fmt.Printf("GetSeller: Could not unmarshal seller : %s", err)
				return seller, errors.New("GetSeller: Could not unmarshal seller with id "+id) 
			}
		

		return seller, nil
}

/**
Saves seller state to the ledger
**/
func SaveSeller(stub *shim.ChaincodeStub, seller Seller, id string)(error){
	fmt.Println("Entering SaveSeller")
	if &seller != nil {
		bytes, _ := json.Marshal(&seller)
		err := stub.PutState(id, bytes)
		if err != nil {
			fmt.Println("SaveSeller: Could not save seller %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid seller input")
	}
}

/**
Save Seller Application to the ledger
**/
func SaveSalesContract(stub *shim.ChaincodeStub, ma SalesContract, id string)([]byte, error){
	fmt.Println("Entering SaveSalesContract")
	if &ma != nil {
		bytes, _ := json.Marshal(&ma)
		scKey, err  := GetStateKey(id, SALESCONTRACT)
		err = stub.PutState(scKey, bytes)
		if err != nil {
			fmt.Println("SaveSalesContract: Could not save seller application %s", err)
			return nil, err
		}
		return bytes, nil
	}else{
		return nil, errors.New("Invalid sellerApplication input")
	}

}

/**
Gets the Auditor from the state if it exists or creates a new one
**/
func GetAuditor(stub *shim.ChaincodeStub, id string)(Auditor, error){
	fmt.Println("Entering GetAuditor")

	var auditor Auditor
	bytes, err := stub.GetState(id)						
				
		if err != nil {	
			fmt.Printf("GetAuditor: Could not get user with id "+ id + ": %s", err)
			return auditor, errors.New("GetAuditor: Failed to get user with id "+ id) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetAuditor: auditor with id does not exist: "+ id + ": %s", err)
			fmt.Println("GetAuditor: creating a auditor with id: "+id)
			
			
			auditor = Auditor{id, AUDITOR_A}
			fmt.Println(auditor)

			bytes, err := json.Marshal(&auditor)
			if err != nil {	
				fmt.Printf("GetAuditor: Could not marshal auditor : %s", err)
				return auditor, errors.New("GetAuditor: Could not marshal auditor with id "+id) 
			}

			err = stub.PutState(id, bytes)
			if err != nil {	
				fmt.Printf("GetAuditor: Could not save auditor : %s", err)
				return auditor, errors.New("GetAuditor: Could not save auditor with id "+id) 
			}

			return auditor, nil
		}
	
			err = json.Unmarshal(bytes, &auditor)
			if err != nil {	
				fmt.Printf("GetAuditor: Could not unmarshal auditor : %s", err)
				return auditor, errors.New("GetAuditor: Could not unmarshal auditor with id "+id) 
			}
		

		return auditor, nil
}

func SaveAuditor(stub *shim.ChaincodeStub, auditor Auditor, id string)(error){
	fmt.Println("Entering SaveAuditor")
	if &auditor != nil {
		bytes, _ := json.Marshal(&auditor)
		err := stub.PutState(id, bytes)
		if err != nil {
			fmt.Println("SaveAuditor: Could not save auditor %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid auditor input")
	}
}

/**
Get a the parent User type from state. Will contain only ID and Affiliation
//Hack for polymorphism
**/
func GetUser(stub *shim.ChaincodeStub, id string)(User, error){
	fmt.Println("Entering GetUser")

	var user User

	key, err := GetStateKey(id, USER)
	if err !=nil {
		fmt.Println("GetUser: Could not get key for user %s", err); 
		return user, err
	}

	bytes , err :=stub.GetState(key)
	if err !=nil {
		fmt.Println("GetUser: Could not get user bytes for user from state %s", err); 
		return user, err
	}


	err = json.Unmarshal(bytes, &user)
	if err !=nil {
		fmt.Println("GetUser: Could not unmarshal user %s", err); 
		return user, err
	}

	return user, nil

}
/**
Add the new id to the list of of keys
**/
func AddKey(stub *shim.ChaincodeStub, id string, keysName string)(bool, error){
	fmt.Println("Entering AddKey")

	var maKeys[] string

	maKeysBytes, err := stub.GetState(keysName)
	if err != nil || len(maKeysBytes) == 0 {
		fmt.Println("AddKey: keys not found for "+keysName+". Creating now...")
		maKeys = []string{}
	}else{
		err = json.Unmarshal(maKeysBytes, &maKeys)
		if err != nil {
			fmt.Println("AddKey: Error unmarshalling  keys %s ", err)
			return false, err
		}
	}

	maKeys = append(maKeys, id)

	bytes, _ := json.Marshal(&maKeys)
	
	err = stub.PutState(maKeysName, bytes)
	if err != nil{
		fmt.Printf("AddKey: Error storing key: %s", err);
		return false, err
	}

	return true, nil

}


/**
Key used for storing object of type buyer
**/
func  GetStateKey(id string, otype int)(string, error){

	if otype == MORTGAGEAPPLICATION {
		return typeMortgageApplication+id, nil
	}else if otype == SALESCONTRACT {
		return typeSalesContract+id, nil
	}else if otype == APPRAISERAPPLICATION {
		return typeAppraiserApplication+id, nil
	}else if otype == USER {
		return typeUser+id, nil
	}else if otype == BUYER {
		return typeBuyer+id, nil
	}else if otype == SELLER {
		return typeSeller+id, nil
	}else if otype == BANK {
		return typeBank+id, nil
	}else if otype == APPRAISER {
		return typeAppraiser+id, nil
	}else if otype == AUDITOR {
		return typeAuditor+id, nil
	}else if otype == LAND {
		return typeLand+id, nil
	}else if otype == PROPERTY {
		return typeProperty+id, nil
	}else if otype == PROPERTYAD {
		return typePropertyAd+id, nil
	}else if otype == MALOG {
		return typeMALog+id, nil
	}else{
		fmt.Println("GetStateKey: Invalid type "+string(otype))
		return "", errors.New("Invalid type")
	}
}

/**
Adds Log for Mortgage Application changes
**/
func AppendMALog(stub *shim.ChaincodeStub, action string, text string, status string, id string)(error){
	fmt.Println("Entering AppendMALog")


	nowTime := time.Now()
	key, _ := GetStateKey(id, MALOG)

	lh, err := GetMALogHolder(stub, key)

	var log MALog
	log.MortgageApplicationId = id
	log.BuyerId = ""
	log.ReviewerId = ""
	log.Text = text
	log.Action = action
	log.Status = status
	log.Timestamp = nowTime.Format("2006-01-02 15:04:05")

	lh.MALogs = append(lh.MALogs, log)

	err = SaveMALogHolder(stub, lh, key)
	if err != nil {
		return err
	}

	keys , err := GetMALogKeys(stub)
	if err != nil {
		return err
	}

	keys = append(keys, key)
	SaveMALogKeys(stub, keys)

	bcLogs, err := GetBCLogs(stub)
	if err != nil {
		return err
	}

	bcLogs = append(bcLogs, log)
	SaveBCLogs(stub, bcLogs)

	return nil
}

/**
Gets the Buyer from the state if it exists or creates a new one
**/
func GetMALogHolder(stub *shim.ChaincodeStub, id string)(MALogHolder, error){
	fmt.Println("Entering GetMALogHolder")

	var lh MALogHolder
	bytes, err := stub.GetState(id)						
				
		if err != nil {	
			fmt.Printf("GetMALogHolder: Could not get logHolder with id %s"+ id , err)
			return lh, errors.New("GetMALogHolder: Failed to get logHolder with id "+ id) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetMALogHolder: logHolder with id does not exist: %s"+ id , err)
			fmt.Println("GetMALogHolder: creating a logHolder with id: "+id)
			
			logs := []MALog {}
			
			lh = MALogHolder{logs}
			fmt.Println(lh)

			bytes, err := json.Marshal(&lh)
			if err != nil {	
				fmt.Printf("GetMALogHolder: Could not marshal logHolder : %s", err)
				return lh, errors.New("GetMALogHolder: Could not marshal logHolder with id "+id) 
			}

			err = stub.PutState(id, bytes)
			if err != nil {	
				fmt.Printf("GetMALogHolder: Could not save logHolder : %s", err)
				return lh, errors.New("GetMALogHolder: Could not save logHolder with id "+id) 
			}

			return lh, nil
		}
	
			err = json.Unmarshal(bytes, &lh)
			if err != nil {	
				fmt.Printf("GetBuyer: Could not unmarshal buyer : %s", err)
				return lh, errors.New("GetBuyer: Could not unmarshal buyer with id "+id) 
			}
		

		return lh, nil
}

func SaveMALogHolder(stub *shim.ChaincodeStub, lh MALogHolder, id string)(error){
	fmt.Println("Entering SaveMALogHolder")
	if &lh != nil {
		bytes, _ := json.Marshal(lh)
		err := stub.PutState(id, bytes)
		if err != nil {
			fmt.Println("SaveMALogHolder: Could not save logHolder %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid logHolder input")
	}
}


/**
Gets the BCLogHolder from the state if it exists or creates a new one
**/
func GetBCLogs(stub *shim.ChaincodeStub)([]MALog, error){
	fmt.Println("Entering GetBCLogs")

	var logs[] MALog
	bytes, err := stub.GetState(bcLogsKey)						
				
		if err != nil {	
			fmt.Printf("GetBCLogs: Could not get logs with id %s"+ bcLogsKey , err)
			return logs, errors.New("GetBCLogs: Failed to get logs with id "+ bcLogsKey) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetBCLogs: logs with id does not exist: %s"+ bcLogsKey , err)
			fmt.Println("GetBCLogs: creating logs with id: "+bcLogsKey)
			
			logs := []MALog {}

			return logs, nil
		}
	
			err = json.Unmarshal(bytes, &logs)
			if err != nil {	
				fmt.Printf("GetBCLogs: Could not unmarshal logs : %s", err)
				return logs, errors.New("GetBCLogs: Could not unmarshal logs with id "+bcLogsKey) 
			}
		

		return logs, nil
}

func SaveBCLogs(stub *shim.ChaincodeStub, logs []MALog)(error){
	fmt.Println("Entering SaveBCLogs")
	if &logs != nil {
		bytes, _ := json.Marshal(&logs)
		err := stub.PutState(bcLogsKey, bytes)
		if err != nil {
			fmt.Println("SaveBCLogs: Could not save logs %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid logs input")
	}
}


/**
Gets the BCLogHolder from the state if it exists or creates a new one
**/
func GetMALogKeys(stub *shim.ChaincodeStub)([]string, error){
	fmt.Println("Entering GetMALogKeys")

	var keys[] string
	bytes, err := stub.GetState(maLogKeysName)						
				
		if err != nil {	
			fmt.Printf("GetMALogKeys: Could not get logs with id %s"+ maLogKeysName , err)
			return keys, errors.New("GetMALogKeys: Failed to get logs with id "+ maLogKeysName) 
		}

		if(len(bytes) == 0){
			fmt.Printf("GetMALogKeys: logs with id does not exist: %s"+ maLogKeysName , err)
			fmt.Println("GetMALogKeys: creating logs with id: "+maLogKeysName)
			
			keys := []string {}

			return keys, nil
		}
	
			err = json.Unmarshal(bytes, &keys)
			if err != nil {	
				fmt.Printf("GetMALogKeys: Could not unmarshal logs : %s", err)
				return keys, errors.New("GetMALogKeys: Could not unmarshal logs with id "+maLogKeysName) 
			}
		

		return keys, nil
}

func SaveMALogKeys(stub *shim.ChaincodeStub, keys []string)(error){
	fmt.Println("Entering SaveMALogKeys")
	if &keys != nil {
		bytes, _ := json.Marshal(&keys)
		err := stub.PutState(maLogKeysName, bytes)
		if err != nil {
			fmt.Println("SaveMALogKeys: Could not save logs %s", err)
			return err
		}
		return nil
	}else{
		return errors.New("Invalid logs input")
	}
}

/**
Create a user and store all related data and metadata
**/

func CreateUser(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Println("Entering CreateUser")
	if len(args) < 2 {
		fmt.Println("CreateUser: Did not recieve enough parameters for creating a user")
		return nil, errors.New("Did not recieve enough parameters for creating a user")
	}

	id := args[0]
	if len(strings.TrimSpace(id)) == 0{
		return nil, errors.New("Invalid user Id")
	}
	affiliationStr := args[1]
	if len(strings.TrimSpace(affiliationStr)) == 0{
		return nil, errors.New("Invalid affiliation")
	}
	affiliation, err := strconv.Atoi(affiliationStr)
	if affiliation == 0 || err !=nil {
		return nil, errors.New("Invalid affiliation")
	}

	key, err := GetStateKey(id, USER)
	if err!=nil {
			fmt.Println("CreateUser: Could not get key for user ",err)
			return nil, err
	}

	if affiliation == BUYER_A{
		
		_, err := GetBuyer(stub, key)
		if err!=nil {
			fmt.Println("CreateUser: Could not create user  ",err)
			return nil, err
		}

	}else if affiliation == SELLER_A{
		_, err := GetSeller(stub, key)
		if err!=nil {
			fmt.Println("CreateUser: Could not create user  ",err)
			return nil, err
		}
		
	}else if affiliation == BANK_A{
		_, err := GetBank(stub, key)
		if err!=nil {
			fmt.Println("CreateUser: Could not create user  ",err)
			return nil, err
		}
		
		
	}else if affiliation == APPRAISER_A{
		_, err := GetAppraiser(stub, key)
		if err!=nil {
			fmt.Println("CreateUser: Could not create user  ",err)
			return nil, err
		}
		
		
	}else if affiliation == AUDITOR_A{
		_, err := GetAuditor(stub, key)
		if err!=nil {
			fmt.Println("CreateUser: Could not create user %s ",err)
			return nil, err
		}
		
		
	}else{
		return nil, errors.New("Invalid user type")
	}
		
	fmt.Println("CreateUser: Successfully created user with ID: "+id)
	return []byte(id), nil
	

}

/**
Returns all transaction records for a mortgage application
**/
func GetAuditorMALogs(stub *shim.ChaincodeStub, callerId string, callerAffiliation int, args []string) ([]byte, error){
	fmt.Println("GetAuditorMALogs")
	
	if len(args) < 1{
		fmt.Println("GetAuditorMALogs: Mortgage Application ID missing")
		return nil, errors.New("Mortgage Application ID missing")
	}

	if callerAffiliation != AUDITOR_A{
		fmt.Println("GetAuditorMALogs: caller "+callerId+" does not have rights to access auditor logs")
		return nil, errors.New("caller "+callerId+" does not have rights to access auditor logs")
	}

	key, _ := GetStateKey(args[0], MALOG)

	lh, err := GetMALogHolder(stub, key)
	if err !=nil {
		fmt.Println("GetAuditorMALogs: Could not fetch MALogHolder for key "+key+" ",err)
		return nil, err
	}

	maLogs := lh.MALogs
	bytes, err := json.Marshal(&maLogs)
	if err !=nil {
		fmt.Println("GetAuditorMALogs: Could not marshal maLogs ",err)
		return nil, err
	}

	return bytes, nil

}


/**
Returns all transaction records for this blockchain network
**/
func GetAuditorBCLogs(stub *shim.ChaincodeStub, callerId string, callerAffiliation int, args []string) ([]byte, error){
	fmt.Println("GetAuditorBCLogs")
	
	if len(args) < 1{
		fmt.Println("GetAuditorBCLogs: Mortgage Application ID missing")
		return nil, errors.New("Mortgage Application ID missing")
	}

	if callerAffiliation != AUDITOR_A{
		fmt.Println("GetAuditorBCLogs: caller "+callerId+" does not have rights to access auditor logs")
		return nil, errors.New("caller "+callerId+" does not have rights to access auditor logs")
	}

	bcLogs, err := GetBCLogs(stub)
	if err != nil {
		fmt.Println("GetAuditorBCLogs: Could not fetch bc logs ",err)
		return nil,err
	}

	bytes, err := json.Marshal(&bcLogs)
	if err !=nil {
		fmt.Println("GetAuditorBCLogs: Could not marshal bcLogs ",err)
		return nil, err
	}

	return bytes, nil

}


// Generates UUID 
/*func GetUUID()(string, error){
	 u4 := uuid.NewV4()
	 var str string = u4.String()
	 fmt.Println(str)
	 return str, nil
}*/

/**
Initialize all dependencies and setup the state
**/
func Setup(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	fmt.Println("Entering Setup")    

	lrec, err := generateLandRecords(stub)
	if err !=nil {
		fmt.Println("Could not generateLandRecords  ", err)
		return nil, err
	}
	fmt.Println(lrec)

	prec, err := generatePropertyList(stub)
	if err !=nil {
		fmt.Println("Could not generateLandRecords  ", err)
		return nil, err
	}
	fmt.Println(prec)

	parec, err := generatePropertyAdsList(stub)
	if err !=nil {
		fmt.Println("Could not generateLandRecords  ", err)
		return nil, err
	}
	fmt.Println(parec)
	
	fmt.Println("Setup complete")
	return nil, nil
}

func (t *MarketplaceChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if function == "Setup" {
        fmt.Println("Firing setup")
        return Setup(stub, args)
    }
	return nil, nil
}

func (t *MarketplaceChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	//need one arg
	/*if len(args) < 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting ......")
	}*/

	username, affiliation, err := GetCallerMetadata(stub)
	if err !=nil {
		return nil, err
	}

	if &username != nil && len(strings.TrimSpace(username)) == 0 {
		return nil, errors.New("Invoke: Could not get username")
	}

	if affiliation <= 0 {
		return nil, errors.New("Invoke: Could not get affiliation")
	} 

	fmt.Println("Caller Metadata: ",username, affiliation);

	


	if function == "GetMortgageApplication" {
		fmt.Println("Getting MortgageApplication")
		_, bytes, err := GetMortgageApplication(stub, username, affiliation, args)
		if err != nil {
			fmt.Println("Error from GetMortgageApplication")
			return nil, err
		} else {
			fmt.Println("All success, returning ma")
			return bytes, nil		 
		}
	}else if function == "GetAppraiserApplication" {
		fmt.Println("Getting AppraiserApplication")
		_, bytes, err := GetAppraiserApplication(stub, username, affiliation, args)
		if err != nil {
			fmt.Println("Error from GetAppraiserApplication")
			return nil, err
		} else {
			fmt.Println("All success, returning ma")
			return bytes, nil		 
		}
	}else if function == "GetSalesContract" {
		fmt.Println("Getting GetSalesContract")
		_, bytes, err := GetSalesContract(stub, username, affiliation, args)
		if err != nil {
			fmt.Println("Error from GetSalesContract")
			return nil, err
		} else {
			fmt.Println("All success, returning sales contract")
			return bytes, nil		 
		}
	}else if function == "GetPropertyAds" {
		fmt.Println("Getting GetPropertyAds")
		_, bytes, err := GetPropertyAds(stub)
		if err != nil {
			fmt.Println("Error from GetPropertyAds")
			return nil, err
		} else {
			fmt.Println("All success, returning property ads")
			return bytes, nil		 
		}
	}else if function == "GetPropertyAd" {
		fmt.Println("Getting GetPropertyAd")
		_, bytes, err := GetPropertyAd(stub, args[0])
		if err != nil {
			fmt.Println("Error from GetPropertyAd")
			return nil, err
		} else {
			fmt.Println("All success, returning property ad")
			return bytes, nil		 
		}
	}else if function == "GetMortgageApplications" {
		fmt.Println("Getting GetMortgageApplications")
		return GetMortgageApplications(stub, username, affiliation, args)
	}else if function == "GetAppraiserApplications" {
		fmt.Println("Getting GetAppraiserApplications")
		return GetAppraiserApplications(stub, username, affiliation, args)
	}else if function == "GetSalesContracts" {
		fmt.Println("Getting GetSalesContracts")
		return GetSalesContracts(stub, username, affiliation, args)
	}else if function == "GetAuditorMALogs" {
		fmt.Println("Getting GetAuditorMALogs")
		return GetAuditorMALogs(stub, username, affiliation, args)
	}else if function == "GetAuditorBCLogs" {
		fmt.Println("Getting GetAuditorBCLogs")
		return GetAuditorBCLogs(stub, username, affiliation, args)
	}

	return nil, errors.New("Invalid function name")
	
}

func (t *MarketplaceChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("Entering Invoke")
	fmt.Println("run is running " + function)
	
	if function == "CreateUser"{
		fmt.Println("Firing CreateUser")
        return CreateUser(stub, args)
    }
    if function == "Setup" {
        fmt.Println("Firing Setup")
        return Setup(stub, args)
    }

	username, affiliation, err := GetCallerMetadata(stub)
	if err !=nil {
		return nil, err
	}

	if &username != nil && len(strings.TrimSpace(username)) == 0 {
		return nil, errors.New("Invoke: Could not get username")
	}

	if affiliation <= 0 {
		return nil, errors.New("Invoke: Could not get affiliation")
	} 

	fmt.Println("Caller Metadata: ",username, affiliation);


	if function == "CreateMortgageApplication" {
		fmt.Println("Firing CreateMortgageApplication")
		return CreateMortgageApplication(stub, username, affiliation, args)
	}else if function == "UpdateMortgageApplication" {
		fmt.Println("Firing UpdateMortgageApplication")
		return UpdateMortgageApplication(stub, username, affiliation, args)
	}else if function == "CreateAppraiserApplication" {
		fmt.Println("Firing CreateAppraiserApplication")
		return CreateAppraiserApplication(stub, username, affiliation, args)
	}else if function == "UpdateAppraiserApplication" {
		fmt.Println("Firing UpdateAppraiserApplication")
		return UpdateAppraiserApplication(stub, username, affiliation, args)
	}else if function == "CreateSalesContract" {
		fmt.Println("Firing CreateSalesContract")
		return CreateSalesContract(stub, username, affiliation, args)
	}else if function == "UpdateSalesContract" {
		fmt.Println("Firing UpdateSalesContract")
		return UpdateSalesContract(stub, username, affiliation, args)
	}else if function == "CreateUser" {
        fmt.Println("Firing CreateUser")
        return CreateUser(stub, args)
    }else if function == "Setup" {
        fmt.Println("Firing Setup")
        return Setup(stub, args)
    }

	return nil, errors.New("Received unknown function invocation")
}


func main(){
	
	err := shim.Start(new(MarketplaceChaincode))
	if err != nil {
		fmt.Println("Error starting Marketplace chaincode: %s", err)
	}

	fmt.Println("MarketplaceChaincode Successfully started")
	
}