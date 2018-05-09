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

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//Account is the type for save account data
type Account struct {
	Balance int    `json:"Balance"`
	Owner   string `json:"Owner"`
	Issuer  string `json:"Issuer"`
	Other   string `json:"Other"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("pointcc Init")
	return shim.Success(nil)
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("pointcc createAccount")
	//_, args := stub.GetFunctionAndParameters()
	var accountID string // accountID
	var emptyaccount Account
	var err error
	var data []byte

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 as [\"Owner\",\"accountID\",\"Balance\",\"Issuer\",\"Other\"]")
	}

	accountID = args[1]
	emptyaccount.Balance, _ = strconv.Atoi(args[2])
	if emptyaccount.Balance < 0 {
		return shim.Error("Account balance cannot be set as negative number")
	}

	emptyaccount.Owner = args[0]
	certsname, ok := getCertificate(stub).(string)
	if !ok {
		return shim.Error("Read certificate error")
	}
	if emptyaccount.Owner == certsname {
		emptyaccount.Balance = 0
	} else if (certsname != "admin" && certsname != "Admin@org1.example.com") {
		return shim.Error(certsname + ": you don't have authority")
	} else {
		fmt.Printf("currently operatator is admin")
	}

	emptyaccount.Issuer = args[3]
	emptyaccount.Other = args[4]
	QueryParameters := [][]byte{[]byte("queryAccount"), []byte(emptyaccount.Owner), []byte(accountID), []byte("all")}
	response := stub.InvokeChaincode("regcc", QueryParameters, "mychannel")
	if response.Status != 200 {
		return response
	}

	data, err = json.Marshal(emptyaccount)
	if err != nil {
		return shim.Error("A wrong input\n")
	}

	if data == nil {
		return shim.Error("A wrong output\n")
	}

	// Write the state to the ledger
	err = stub.PutState(accountID, data)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) setAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("pointcc setAccount")
	//_, args := stub.GetFunctionAndParameters()
	var accountID string // accountID
	var emptyaccount Account
	var err error
	var data []byte

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 as [\"Owner\",\"accountID\",\"Balance\",\"Issuer\",\"Other\"]")
	}

	accountID = args[1]
	emptyaccount.Balance, _ = strconv.Atoi(args[2])
	if emptyaccount.Balance < 0 {
		return shim.Error("Account balance cannot be set as negative number")
	}
	emptyaccount.Owner = args[0]
	certsname, ok := getCertificate(stub).(string)
	if !ok {
		return shim.Error("Read certificate error")
	}
	if (certsname != "admin" && certsname != "Admin@org1.example.com") {
		return shim.Error("You don't have authority")
	} else {
        fmt.Printf("Crruntely operater is admin")
    }

	emptyaccount.Issuer = args[3]
	emptyaccount.Other = args[4]
	QueryParameters := [][]byte{[]byte("queryAccount"), []byte(emptyaccount.Owner), []byte(accountID), []byte("all")}
	response := stub.InvokeChaincode("regcc", QueryParameters, "mychannel")
	if response.Status != 200 {
		return response
	}

	data, err = json.Marshal(emptyaccount)
	if err != nil {
		return shim.Error("A wrong input\n")
	}

	if data == nil {
		return shim.Error("A wrong output\n")
	}

	// Write the state to the ledger
	err = stub.PutState(accountID, data)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("pointcc Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "trade" {
		// Make payment of X units from A to B
		return t.trade(stub, args)
	} else if function == "deleteAccount" {
		// Deletes an entity from its state
		return t.deleteAccount(stub, args)
	} else if function == "createAccount" {
		// my test chaincode createAccount function
		return t.createAccount(stub, args)
	} else if function == "queryAccount" {
		// query detail of account
		return t.queryAccount(stub, args)
	} else if function == "queryHistory" {
		// query history of account
		return t.historyquery(stub, args)
	} else if function == "setAccount" {
		// Set the balance of Account
		return t.setAccount(stub, args)
	} else if function == "extrade" {
		// multi asset traditon
		return t.extrade(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"trade\" \"deleteAccount\" \"queryAccount\" \"setAccount\" \"createAccount\"")
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

//test query method function
func (t *SimpleChaincode) queryAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var accountID string
	var Key string
	var account Account
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the acountID to query and 1 parament")
	}

	accountID = args[0]
	Key = args[1]
	//return shim.Error(accountID + " " + Key)

	// Get the state from the ledger
	data, err := stub.GetState(accountID)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + accountID + "\"}"
		return shim.Error(jsonResp)
	}

	if data == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + accountID + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal(data, &account)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to read data of " + accountID + "\"}"
		return shim.Error(jsonResp)
	}

	// expercting avlible json key name
	if !((Key == "all") || (Key == "Balance") || (Key == "Owner") || (Key == "Issuer") || (Key == "Other")) {
		return shim.Error("Invalid json key name. Expecting \"all\" \"Balance\" \"Issuer\" \"Owner\" \"Other\" ")
	} else if Key == "all" {
		jsonResp := "{\"accountID\":\"" + accountID + "\",\"Issuer\":\"" + account.Issuer + "\",\"Owner\":\"" + account.Owner + "\",\"Balance\":\"" + strconv.Itoa(account.Balance) + "\",\"Other\":\"" + account.Other + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return shim.Success([]byte(jsonResp))
	} else if Key == "Balance" {
		jsonResp := "{\"accountID\":\"" + accountID + "\",\"Balance\":\"" + strconv.Itoa(account.Balance) + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	} else if Key == "Owner" {
		jsonResp := "{\"accountID\":\"" + accountID + "\",\"Owner\":\"" + account.Owner + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	} else if Key == "Issuer" {
		jsonResp := "{\"accountID\":\"" + accountID + "\",\"Issuer\":\"" + account.Issuer + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	} else if Key == "Other" {
		jsonResp := "{\"accountID\":\"" + accountID + "\",\"Other\":\"" + account.Other + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	} else {
        return shim.Error("What do you wanna me do?")
    }
}

func (t *SimpleChaincode) trade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A, B string //two accountID, A trans X token to B
	var X int
	var Aact, Bact Account //account of A and B
	var Aval, Bval int     //account balance of A and B
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	A = args[0]
	B = args[1]

	// Get A's state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	err = json.Unmarshal(Avalbytes, &Aact)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Aval = Aact.Balance

	//Verificate if the transaction was authorized
	Currentuser := getCertificate(stub)
	if (Currentuser == "Admin@org1.example.com") || (Currentuser == "Admin@org2.example.com") || (Currentuser == "Admin@org3.example.com") {
		fmt.Printf("Current operator is Administrator: ")
		fmt.Println(Currentuser)
		fmt.Printf(".\n")
	} else if Currentuser == -1 {
		return shim.Error("No certificate found")
	} else if Currentuser == -2 {
		return shim.Error("Could not decode the PEM structure")
	} else if Currentuser == -3 {
		return shim.Error("ParseCertificate failed")
	} else if Currentuser != Aact.Owner {
		return shim.Error("Current operator don't have authority!")
	}

	// Get B's state from the ledger
	Bvalbytes, err := stub.GetState(B)
	if err != nil {
		return shim.Error("Failed to get B state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity B not found")
	}
	err = json.Unmarshal(Bvalbytes, &Bact)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Bval = Bact.Balance

	if Aact.Issuer != Bact.Issuer {
		return shim.Error("Different types of Points, can not be traded")
	}

	// Perform the execution
	X, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if X <= 0 {
		return shim.Error("Invalid transaction amount, expecting a postive integer value")
	}
	Aval = Aval - X
	Bval = Bval + X
	if (Aval < 0) || (Bval < 0) {
		fmt.Printf("Insufficient balance")
		return shim.Error("Insufficient balance")
	}
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	Aact.Balance = Aval
	Bact.Balance = Bval
	Avalbytes, _ = json.Marshal(Aact)
	Bvalbytes, _ = json.Marshal(Bact)
	err = stub.PutState(A, Avalbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(B, Bvalbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte("Done!"))
}

func (t *SimpleChaincode) extrade(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Aa, Ab, Ba, Bb string //four accountID, Aa trans Xa token to Ba, Bb trans Xb token to Ab.
	var Xa, Xb int
	var Aacta, Aactb, Bacta, Bactb Account //account of A and B
	var Avala, Avalb, Bvala, Bvalb int     //account balance of A and B
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Aa = args[0]
	Ba = args[1]

	Ab = args[3]
	Bb = args[4]

	// Get Aa's state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := stub.GetState(Aa)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	err = json.Unmarshal(Avalbytes, &Aacta)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Avala = Aacta.Balance

	// Get Ab's state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err = stub.GetState(Ab)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Avalbytes == nil {
		return shim.Error("Entity not found")
	}
	err = json.Unmarshal(Avalbytes, &Aactb)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Avalb = Aactb.Balance

	// Get Ba's state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Bvalbytes, err := stub.GetState(Ba)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	err = json.Unmarshal(Bvalbytes, &Bacta)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Bvala = Bacta.Balance

	// Get Bb's state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Bvalbytes, err = stub.GetState(Bb)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if Bvalbytes == nil {
		return shim.Error("Entity not found")
	}
	err = json.Unmarshal(Bvalbytes, &Bactb)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	Bvalb = Bactb.Balance

	//Verificate if the transaction was authorized
	Currentuser := getCertificate(stub)
	if (Currentuser != "admin" && Currentuser != "Admin@org1.example.com") {
		return shim.Error(Currentuser + ": you don't have authority")
	} else {
		fmt.Printf("currently operatator is admin")
	}

	if (Aacta.Issuer != Bacta.Issuer) || (Aactb.Issuer != Bactb.Issuer) {
		return shim.Error("Different types of Points, can not be traded")
	}

	// Perform the execution
	Xa, err = strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if Xa <= 0 {
		return shim.Error("Invalid transaction amount, expecting a postive integer value")
	}
	Xb, err = strconv.Atoi(args[5])
	if err != nil {
		return shim.Error("Invalid transaction amount, expecting a integer value")
	}
	if Xb <= 0 {
		return shim.Error("Invalid transaction amount, expecting a postive integer value")
	}
	Avala = Avala - Xa
	Bvala = Bvala + Xa
	Avalb = Avalb + Xb
	Bvalb = Bvalb - Xb
	if (Avala < 0) || (Bvala < 0) || (Avalb < 0) || (Bvalb < 0) {
		fmt.Printf("Insufficient balance")
		return shim.Error("Insufficient balance")
	}
	fmt.Printf("After tradition, Avala = %d, Bvala = %d\n, Avalb = %d, Bvalb = %d\n", Avala, Bvala, Avalb, Bvalb)

	// Write the state back to the ledger
	Aacta.Balance = Avala
	Bacta.Balance = Bvala
	Aactb.Balance = Avalb
	Bactb.Balance = Bvalb
	Avalabytes, _ := json.Marshal(Aacta)
	Bvalabytes, _ := json.Marshal(Bacta)
	Avalbbytes, _ := json.Marshal(Aactb)
	Bvalbbytes, _ := json.Marshal(Bactb)

	err = stub.PutState(Aa, Avalabytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(Ba, Bvalabytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(Ab, Avalbbytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(Bb, Bvalbbytes)
	if err != nil {
		return shim.Error(err.Error())
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

func (t *SimpleChaincode) historyquery(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Printf("historyquery in mapcc")
	var id string
	if len(args) != 1 {
		return shim.Error("Expected 1 parament in mapping historyquery.")
	}
	id = args[0]
	raw, err := stub.GetState(id)
	if (err != nil) || (raw == nil) {
		return shim.Error("User " + id + " is not exists, or getstate error.")
	}
	it, _ := stub.GetHistoryForKey(id)
	result, _ := getHistoryListResult(it)
	return shim.Success(result)
}

// Deletes an entity from state
func (t *SimpleChaincode) deleteAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]
	Avalbyte, err := stub.GetState(A)
	var Aact Account
	err = json.Unmarshal(Avalbyte, &Aact)
	if err != nil {
		return shim.Error("Failed to get state.")
	}

	// Delete the key from the state in ledger
	err = stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}
	/*
		QueryParameters := [][]byte{[]byte("invoke"), []byte(Aact.Owner), []byte("delete"), []byte(A), []byte("pointchannel"), []byte("points"), []byte(Aact.Issuer)}
		response := stub.InvokeChaincode("regcc", QueryParameters, "regchannel")
		if response.Status == 500 {
			return response
		}
	*/

	return shim.Success(nil)
}

// query callback representing the query of a chaincode

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
