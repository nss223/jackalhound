/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example02.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	//"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//AccountMV is the minimum usable type for save account data
type AccountMV struct {
	Type   string `json:"Type"`
	Owner  string `json:"Owner"`
	Issuer string `json:"Issuer"`
	Other  string `json:"Other"`
}

//AccountCL is the mapping of credit line
type AccountCL struct {
	AccountMV `json:"AccountMV"`
	Balance   int      `json:"Balance"`
	Parent    string   `json:"Parent"`
	Children  []string `json:"Children"`
}

type Account struct {
	AccountMV `json:"AccountMV"`
}

type User struct {
	Accounts map[string]interface{} `json:"Accounts"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("map Init")
	return shim.Success(nil)
}

func (t *SimpleChaincode) init(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("map init")
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("map Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "trade" {
		// Make payment of X units from A to B
		return t.trade(stub, args)
	} else if function == "deleteUser" {
		// Deletes an entity from its state
		return t.userdelete(stub, args)
	} else if function == "deleteAccount" {
		// Deletes an entity from its state
		return t.accountdelete(stub, args)
	} else if function == "createUser" {
		// Deletes an entity from its state
		return t.createUser(stub, args)
	} else if function == "createAccount" {
		// Deletes an entity from its state
		return t.createAccount(stub, args)
	} else if function == "queryUser" {
		// query detail of account
		return t.userquery(stub, args)
	} else if function == "queryAccount" {
		// query detail of account
		return t.accountquery(stub, args)
	} else if function == "queryUserHistory" {
		//query user history by id
		return t.queryUserHistory(stub, args)
	} else if function == "queryAccountHistory" {
		//query account history by id
		return t.queryAccountHistory(stub, args)
	} else if function == "queryuserall" {
		return t.queryuserall(stub, args)
	} else if function == "queryaccountall" {
		return t.queryaccountall(stub, args)
	} else if function == "createAccountCL" {
		return t.createAccountCL(stub, args)
	} else if function == "splitAccountCL" {
		return t.splitAccountCL(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"trade\" \"deleteUser\" \"deleteAccount\" \"queryUser\" \"queryAccount\" \"createUser\"\"createAccount\" \"queryHistory\"")
}

//Test function
func getCertificate(stub shim.ChaincodeStubInterface) interface{} {
	creatorByte, _ := stub.GetCreator()
	certStart := bytes.IndexAny(creatorByte, "-----BEGIN")
	if certStart == -1 {
		return -1 //("No certificate found")
	}
	certText := creatorByte[certStart:]
	bl, _ := pem.Decode(certText)
	if bl == nil {
		return -2 //("Could not decode the PEM structure")
	}

	cert, err := x509.ParseCertificate(bl.Bytes)
	if err != nil {
		return -3 //("ParseCertificate failed")
	}
	uname := cert.Subject.CommonName
	fmt.Println("Name:" + uname)
	return uname //shim.Success([]byte("Called testCertificate " + uname))
}

func (t *SimpleChaincode) userquery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountID string
	var data User
	if len(args) != 1 {
		return shim.Error("Expected 1 parament as users id(in map userquery, id starts with character\"U\")")
	}
	accountID = args[0]
	cert, ok := getCertificate(stub).(string)
	if !ok {
		return shim.Error("Can not get certificate.")
	}
	if (cert != accountID) && !((cert == "Admin@org1.example.com") || (cert == "Admin@org2.example.com") || (cert == "Admin@org3.example.com")) {
		return shim.Error("Operator don't have authority.")
	}
	accountID = "User" + accountID
	raw, err := stub.GetState(accountID)
	if (err != nil) || (raw == nil) {
		return shim.Error("Failed to get state of users account in map userquery")
	}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed in transing to json in map userquery")
	}
	list := data.Accounts
	output := "User " + accountID + " have follow maping account: "
	for aid := range list {
		output = output + " " + aid + "\n"
	}
	fmt.Printf("Query Response: " + output)
	return shim.Success([]byte(output))
}

func (t *SimpleChaincode) accountquery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountID string
	var data Account
	if len(args) == 2 {
		if args[1] == "CL" {
			return t.accountCLquery(stub, args)
		}
		return shim.Error("Expected 1 parament as mapping account id(in map accountquery)")
	}
	if len(args) != 1 {
		return shim.Error("Expected 1 parament as mapping account id(in map accountquery)")
	}
	accountID = args[0]
	accountID = "Account" + accountID
	raw, err := stub.GetState(accountID)
	if (err != nil) || (raw == nil) {
		return shim.Error("Failed to get state of users account in map accountquery")
	}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed in transing to json in map accountquery")
	}
	cert, ok := getCertificate(stub).(string)
	if !ok {
		return shim.Error("Can not get certificate.")
	}
	if (cert != data.Owner) && !((cert == "Admin@org1.example.com") || (cert == "Admin@org2.example.com") || (cert == "Admin@org3.example.com")) {
		return shim.Error("Operator don't have authority.")
	}
	jsonResp := "{\"UserID\":\"" + data.Owner + "\",\"accountID\":\"" + accountID + "\",\"Type\":\"" + data.Type + "\",\"Issuer\":\"" + data.Issuer + "\",\"Other\":\"" + data.Other + "\"}"

	fmt.Printf("Query Response: " + jsonResp)
	return shim.Success([]byte(jsonResp))
}

func (t *SimpleChaincode) accountCLquery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountID string
	var data AccountCL
	accountID = args[0]
	accountID = "Account" + accountID
	raw, err := stub.GetState(accountID)
	if (err != nil) || (raw == nil) {
		return shim.Error("Failed to get state of users account in map accountquery")
	}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed in transing to json in map accountquery")
	}
	cert, ok := getCertificate(stub).(string)
	if !ok {
		return shim.Error("Can not get certificate.")
	}
	if (cert != data.Owner) && !((cert == "Admin@org1.example.com") || (cert == "Admin@org2.example.com") || (cert == "Admin@org3.example.com")) {
		return shim.Error("Operator don't have authority.")
	}
	jsonResp := "{\"UserID\":\"" + data.Owner + "\",\"accountID\":\"" + accountID + "\",\"Type\":\"" + data.Type + "\",\"Issuer\":\"" + data.Issuer + "\",\"Other\":\"" + data.Other + "\",\"Balance\":\"" + strconv.Itoa(data.Balance) + "\",\"Parent\":\"" + data.Parent
	if len(data.Children) != 0 {
		jsonResp = jsonResp + "\",\"ChildrenList\":\""
	}
	for i := 0; i < len(data.Children); i++ {
		jsonResp = jsonResp + data.Children[i] + ","
	}
	jsonResp = jsonResp + "\"}"
	fmt.Printf("Query Response: " + jsonResp)
	return shim.Success([]byte(jsonResp))
}

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("mapping create user")
	var userid string
	var data User
	var raw []byte
	var err error
	if len(args) != 1 {
		return shim.Error("Expected 1 parament as mapping user id")
	}
	userid = "User" + args[0]

	QueryParameters := [][]byte{[]byte("queryAccount"), []byte(args[0]), []byte("mapchannel"), []byte("all")}
	response := stub.InvokeChaincode("regcc", QueryParameters, "regchannel")
	if response.Status != 200 {
		return response
	}
	raw, _ = stub.GetState(userid)
	if raw != nil {
		return shim.Error("User is already exists.")
	}
	data.Accounts = make(map[string]interface{})
	raw, err = json.Marshal(data)
	if err != nil {
		return shim.Error("Failed to trans json")
	}
	err = stub.PutState(userid, raw)
	if err != nil {
		return shim.Error("Failed to save state")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("mapping create asset")
	if len(args) != 5 {
		return shim.Error("Expected 5 parament in mapping createAccount")
	}
	var userid string
	var assetid string
	var data Account
	var raw []byte
	var err error

	userid = "User" + args[0]
	assetid = "Account" + args[1]

	raw, err = stub.GetState(userid)
	if err != nil {
		return shim.Error("User is not exists")
	}

	raw, err = stub.GetState(assetid)
	if raw != nil {
		return shim.Error("asset is already exist.")
	}
	data.Type = args[2]
	data.Issuer = args[3]
	data.Owner = args[0]
	data.Other = args[4]
	raw, err = json.Marshal(data)
	if err != nil {
		return shim.Error("Failed to trans to json in createAccount")
	}
	err = stub.PutState(assetid, raw)
	if err != nil {
		return shim.Error("Failed to put state in createAccount")
	}
	Uraw, _ := stub.GetState(userid)
	var list User
	err = json.Unmarshal(Uraw, &list)
	if err != nil {
		return shim.Error("Failed to trans json in createAccount.")
	}
	list.Accounts[args[1]] = nil
	Uraw, err = json.Marshal(list)
	if err != nil {
		return shim.Error("Failed to save state")
	}
	err = stub.PutState(userid, Uraw)
	if err != nil {
		return shim.Error("Failed to save state.")
	}
	return shim.Success(nil)
}

//trade asset from A to B
func (t *SimpleChaincode) trade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("mapping invoke")
	var A, B string
	var asset string
	var Avalbyte, Bvalbyte, assetbyte []byte
	var err error
	if len(args) != 3 {
		return shim.Error("Expected 3 parament(in map accountquery)")
	}
	A = args[0]
	B = args[1]
	asset = args[2]
	A = "User" + A
	B = "User" + B
	Avalbyte, err = stub.GetState(A)
	if err != nil {
		return shim.Error("User A is not exist.")
	}
	if Avalbyte == nil {
		return shim.Error("User A is not exist.")
	}
	Bvalbyte, err = stub.GetState(B)
	if err != nil {
		return shim.Error("User B is not exist.")
	}
	if Bvalbyte == nil {
		return shim.Error("User B is not exist.")
	}
	assetbyte, err = stub.GetState("Account" + asset)
	if err != nil {
		return shim.Error("asset is not exist.")
	}
	if assetbyte == nil {
		return shim.Error("asset is not exist.")
	}
	var assetdata Account
	err = json.Unmarshal(assetbyte, &assetdata)
	if err != nil {
		return shim.Error("Read data error in mapping invoke.1")
	}
	if assetdata.Owner != args[0] {
		return shim.Error(A + " don't have this asset")
	}
	var Adata, Bdata User
	err = json.Unmarshal(Avalbyte, &Adata)
	if err != nil {
		return shim.Error("Read data error in mapping invoke.2")
	}
	err = json.Unmarshal(Bvalbyte, &Bdata)
	if err != nil {
		return shim.Error("Read data error in mapping invoke.3")
	}
	delete(Adata.Accounts, asset)
	assetdata.Owner = args[1]
	Bdata.Accounts[asset] = nil
	assetbyte, err = json.Marshal(assetdata)
	if err != nil {
		return shim.Error("Save data error in mapping invoke.1")
	}
	err = stub.PutState("Account"+asset, assetbyte)
	if err != nil {
		return shim.Error("Save data error in mapping invoke.2")
	}
	Avalbyte, err = json.Marshal(Adata)
	if err != nil {
		return shim.Error("Save data error in mapping invoke.3")
	}
	err = stub.PutState(A, Avalbyte)
	if err != nil {
		return shim.Error("Save data error in mapping invoke.4")
	}
	Bvalbyte, err = json.Marshal(Bdata)
	if err != nil {
		return shim.Error("Save data error in mapping invoke.5")
	}
	err = stub.PutState(B, Bvalbyte)
	if err != nil {
		return shim.Error("Save data error in mapping invoke.6")
	}
	return shim.Success(nil)
}

func getHistoryListResult(resultsIterator shim.HistoryQueryIteratorInterface) ([]byte, error) {

	defer resultsIterator.Close()
	// buffer is a JSON array containing QueryRecords
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		item, _ := json.Marshal(queryResponse)
		buffer.Write(item)
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")
	fmt.Printf("queryResult:\n%s\n", buffer.String())
	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) queryAccountHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("queryHistory in mapcc")
	var id string
	if len(args) != 1 {
		return shim.Error("Expected 1 parament in mapping queryHistory.")
	}
	id = "Account" + args[0]
	raw, err := stub.GetState(id)
	if (err != nil) || (raw == nil) {
		return shim.Error("User " + id + " is not exists, or getstate error.")
	}
	it, _ := stub.GetHistoryForKey(id)
	result, _ := getHistoryListResult(it)
	return shim.Success(result)
}

func (t *SimpleChaincode) queryUserHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("queryHistory in mapcc")
	var id string
	if len(args) != 1 {
		return shim.Error("Expected 1 parament in mapping queryHistory.")
	}
	id = "User" + args[0]
	raw, err := stub.GetState(id)
	if (err != nil) || (raw == nil) {
		return shim.Error("User " + id + " is not exists, or getstate error.")
	}
	it, _ := stub.GetHistoryForKey(id)
	result, _ := getHistoryListResult(it)
	return shim.Success(result)
}

func (t *SimpleChaincode) userdelete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountID string
	var err error
	var temp []string
	accountID = "User" + args[0]
	var raw []byte
	var data User
	raw, err = stub.GetState(accountID)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed to trans json")
	}
	for asset := range data.Accounts {
		temp[0] = asset
		t.accountdelete(stub, temp)
	}
	err = stub.DelState(accountID)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) accountdelete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountID string
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	accountID = "Account" + args[0]
	err = stub.DelState(accountID)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) queryuserall(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Key string
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}
	Key = "User" + args[0]
	raw, err := stub.GetState(Key)
	if err != nil {
		return shim.Error("Failed to getstate")
	}
	if raw == nil {
		return shim.Error("User is not exist")
	}
	var data User
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed to trans to json")
	}
	return shim.Success(raw)
}

func (t *SimpleChaincode) queryaccountall(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Key string
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1.")
	}
	Key = "Account" + args[0]
	raw, err := stub.GetState(Key)
	if err != nil {
		return shim.Error("Failed to getstate")
	}
	if raw == nil {
		return shim.Error("User is not exist")
	}
	var data Account
	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed to trans to json")
	}
	return shim.Success(raw)
}

func (t *SimpleChaincode) createAccountCL(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("mapping create asset")
	if len(args) != 7 {
		return shim.Error("Expected 7 parament in mapping createAccountCL")
	}
	var userid string
	var assetid string
	var data AccountCL
	var raw []byte
	var err error

	userid = "User" + args[0]
	assetid = "Account" + args[1]

	raw, err = stub.GetState(userid)
	if err != nil {
		return shim.Error("User is not exists")
	}

	raw, err = stub.GetState(assetid)
	if raw != nil {
		return shim.Error("asset is already exist.")
	}
	data.Type = args[2]
	data.Issuer = args[3]
	data.Owner = args[0]
	data.Other = args[4]
	data.Balance, _ = strconv.Atoi(args[5])
	data.Parent = args[6]
	data.Children = make([]string, 0)
	raw, err = json.Marshal(data)
	if err != nil {
		return shim.Error("Failed to trans to json in createAccount")
	}
	err = stub.PutState(assetid, raw)
	if err != nil {
		return shim.Error("Failed to put state in createAccount")
	}
	Uraw, _ := stub.GetState(userid)
	var list User
	err = json.Unmarshal(Uraw, &list)
	if err != nil {
		return shim.Error("Failed to trans json in createAccount.")
	}
	list.Accounts[args[1]] = nil
	Uraw, err = json.Marshal(list)
	if err != nil {
		return shim.Error("Failed to save state")
	}
	err = stub.PutState(userid, Uraw)
	if err != nil {
		return shim.Error("Failed to save state.")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) splitAccountCL(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("mapping split asset")
	if len(args) != 4 {
		return shim.Error("Expected 4 parament in mapping splitAccountCL")
	}
	var vendorid string
	var buyerid string
	var newassetid string
	var assetid string
	var X int //how much creit lines in this tradition
	var data AccountCL
	var raw []byte
	var err error

	buyerid = "User" + args[1]
	assetid = "Account" + args[0]
	newassetid = "Account" + args[2]
	X, _ = strconv.Atoi(args[3])

	raw, err = stub.GetState(buyerid)
	if err != nil || raw == nil {
		return shim.Error("User is not exists")
	}

	raw, err = stub.GetState(newassetid)
	if err != nil && raw != nil {
		return shim.Error("buyer's assetid is already exist.")
	}

	raw, err = stub.GetState(assetid)
	if err != nil || raw == nil {
		return shim.Error("vendor's assetid is not exist.")
	}

	err = json.Unmarshal(raw, &data)
	if err != nil {
		return shim.Error("Failed to trans json in splitAccountCL.")
	}

	vendorid = data.Owner
	cert, ok := getCertificate(stub).(string)
	if !ok {
		return shim.Error("Can not get certificate.")
	}
	if (cert != vendorid) && !((cert == "Admin@org1.example.com") || (cert == "Admin@org2.example.com") || (cert == "Admin@org3.example.com")) {
		return shim.Error("Operator don't have authority.")
	}
	if X <= 0 {
		return shim.Error("non-positive number is inlegal.")
	}
	if X > data.Balance {
		return shim.Error("This credit account balance is not enough.")
	}
	data.Balance = data.Balance - X
	data.Children = append(data.Children, newassetid+" "+args[3])
	raw, err = json.Marshal(data)
	err = stub.PutState(assetid, raw)
	if err != nil {
		return shim.Error("Error in putstate.")
	}
	Args := make([]string, 7)
	Args[0] = args[1]     //userid
	Args[1] = args[2]     //accountid
	Args[2] = data.Type   //accountType
	Args[3] = data.Issuer //accountIssuer
	Args[4] = data.Other  //accountOther
	Args[5] = args[3]     //accountBalance
	Args[6] = args[0]     //accountParent

	response := t.createAccountCL(stub, Args)
	return response
}

// query callback representing the query of a chaincode

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
