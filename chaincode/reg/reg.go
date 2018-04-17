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

//Account is the type for save account  data
type Account struct {
	ChannelID   string `json:"ChannelID"`
	AccountType string `json:"AccountType"`
	Issuer      string `json:"Issuer"`
}

//User is the typr for save account  for user
type User struct {
	Accountnum int                `json:Accountnum`
	Accounts   map[string]Account `json:Accounts`
}

type MapUser struct {
	Accounts map[string]interface{} `json:Accounts`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("reg Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "createAccount" {
		// give user a new account
		return t.createAccount(stub, args)
	} else if function == "deleteUser" {
		// Deletes an entity from its state
		return t.deleteUser(stub, args)
	} else if function == "createUser" {
		// my test chaincode init function
		return t.createUser(stub, args)
	} else if function == "queryUser" {
		// query user's all account id
		return t.query(stub, []string{args[0], "all", "all"})
	} else if function == "queryAccount" {
		// query detail of account
		return t.query(stub, args)
	} else if function == "deleteAccount" {
		// Deletes an entity from its state
		return t.deleteAccount(stub, args)
	} else if function == "setAssetByAccount" {
		// set some message to account
		return t.setAssetByAccount(stub, args)
	} else if function == "queryHistory" {
		// set some message to account
		return t.queryHistory(stub, args)
	} else if function == "queryall" {
		return t.queryall(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"createAccount\" \"deleteAccount\" \"deleteUser\" \"createUser\" \"queryall\" \"queryUser\" \"queryAccount\" \"setAssetByAccount\"")
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
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var accountID string
	var Key string
	var alists string
	var err error
	var list User
	var k string

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query and 2 parament userid/accountID/Key")
	}

	A = args[0]
	accountID = args[1]
	Key = args[2]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal(Avalbytes, &list)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to read data of " + A + "\"}"
		return shim.Error(jsonResp)
	}

	//if accountID = "all", then we return all account's ID.
	//Note: we can't register account  named "all"!!!!!!!!!!!!!!!!!!!!!!!!!!!!
	if accountID == "all" {
		alists = A + " have " + strconv.Itoa(list.Accountnum) + " Account: "
		fmt.Printf(A+" have %d Accounts.\n", list.Accountnum)
		i := 0
		for k = range list.Accounts {
			fmt.Printf("Account%d: %s\n", i, k)
			i++
			alists = alists + k + " "
		}
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(alists),
		}
	}
	//Load the account  data.
	var data Account
	data, ok := list.Accounts[accountID]
	if ok == false {
		fmt.Printf("lost key")
		for k = range list.Accounts {
			fmt.Printf(k)
		}
		// expercting avlible accountID
		return shim.Error("{\"Error\":\"This account " + accountID + " is not in Accounts list of " + A + ", you can set accountID = \"all\" to query Account list of " + A + "\"}")
	}

	// expercting avlible json key name
	if !((Key == "all") || (Key == "ChannelID") || (Key == "AccountType") || (Key == "Issuer")) {
		return shim.Error("Invalid json key name. Expecting \"all\" \"ChannelID\" \"AccountType\" \"Issuer\" ")
	} else if Key == "all" {
		jsonResp := "{\"UserID\":\"" + A + "\",\"accountID\":\"" + accountID + "\",\"ChannelID\":\"" + data.ChannelID + "\",\"AccountType\":\"" + data.AccountType + "\",\"Issuer\":\"" + data.Issuer + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return shim.Success([]byte(jsonResp))
	} else if Key == "ChannelID" {
		jsonResp := "{\"UserID\":\"" + A + "\",\"accountID\":\"" + accountID + "\",\"ChannelID\":\"" + data.ChannelID + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	} else if Key == "AccountType" {
		jsonResp := "{\"UserID\":\"" + A + "\",\"accountID\":\"" + accountID + "\",\"AccountType\":\"" + data.AccountType + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	} else if Key == "Issuer" {
		jsonResp := "{\"UserID\":\"" + A + "\",\"accountID\":\"" + accountID + "\",\"Issuer\":\"" + data.Issuer + "\"}"
		fmt.Printf("Query Response:%s\n", jsonResp)
		return pb.Response{
			Status:  200,
			Message: "OK",
            Payload: []byte(jsonResp),
		}
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) createUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("regcc Init")
	//_, args := stub.GetFunctionAndParameters()
	var A string // Entities
	var emptyuser User
	var err error
	var data []byte

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	emptyuser.Accountnum = 0
	emptyuser.Accounts = make(map[string]Account)

	data, err = json.Marshal(emptyuser)
	if err != nil {
		return shim.Error("a wrong way\n")
	}

	// Initialize the chaincode
	A = args[0]

	// Write the state to the ledger
	err = stub.PutState(A, data)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userid string  //entity
	var account string //accountID, if account  don't exists then error except event = createAccount, else putstate raw to it
	var raw Account    //what we want to change
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 userid/accountid/raw//ChannelID/AccountType/Issuer")
	}

	userid = args[0]
	account = args[1]
	raw.ChannelID = args[2]
	raw.AccountType = args[3]
	raw.Issuer = args[4]

	Currentuser := getCertificate(stub)
	if (Currentuser == "Admin@org1.example.com" || Currentuser == "admin") {
		fmt.Printf("Current operator is Administrator: ")
		fmt.Println(Currentuser)
	} else if Currentuser == -1 {
		return shim.Error("No certificate found")
	} else if Currentuser == -2 {
		return shim.Error("Could not decode the PEM structure")
	} else if Currentuser == -3 {
		return shim.Error("ParseCertificate failed")
	} else if Currentuser != userid {
		return shim.Error("Current operator don't have authority!")
	}

	uservalbytes, err := stub.GetState(userid)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if uservalbytes == nil {
		return shim.Error("User is not exist, you can initalizate it")
	}

	var userdata User
	err = json.Unmarshal(uservalbytes, &userdata)
	if err != nil {
		return shim.Error("Failed to trans json")
	}

	_, ok := userdata.Accounts[account]
	if ok {
		return shim.Error("account already exists.\n")
	}
	if account == "all" {
		return shim.Error("\"all\"is reserved name, you can't use it.\n")
	}
	userdata.Accounts[account] = raw
	userdata.Accountnum++
	uservalbytes, _ = json.Marshal(userdata)
	err = stub.PutState(userid, uservalbytes)
	if err != nil {
		return shim.Error("Failed to putstate.\n")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) setAssetByAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userid string  //entity
	var account string //accountID, if account  don't exists then error except event = add, else putstate raw to it
	var raw Account    //what we want to change
	var err error

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5 userid/accountid/raw//ChannelID/AccountType/Issuer")
	}

	userid = args[0]
	account = args[1]
	raw.ChannelID = args[2]
	raw.AccountType = args[3]
	raw.Issuer = args[4]

	Currentuser := getCertificate(stub)
	if (Currentuser == "Admin@org1.example.com") || (Currentuser == "admin") {
		fmt.Printf("Current operator is Administrator: ")
		fmt.Println(Currentuser)
	} else if Currentuser == -1 {
		return shim.Error("No certificate found")
	} else if Currentuser == -2 {
		return shim.Error("Could not decode the PEM structure")
	} else if Currentuser == -3 {
		return shim.Error("ParseCertificate failed")
	} else if Currentuser != userid {
		return shim.Error("Current operator don't have authority!")
	}

	uservalbytes, err := stub.GetState(userid)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if uservalbytes == nil {
		return shim.Error("User is not exist, you can initalizate it")
	}

	var userdata User
	err = json.Unmarshal(uservalbytes, &userdata)
	if err != nil {
		return shim.Error("Failed to trans json")
	}

	_, ok := userdata.Accounts[account]
	if !ok {
		return shim.Error("account  don't exists.\n")
	}
	userdata.Accounts[account] = raw
	uservalbytes, _ = json.Marshal(userdata)
	err = stub.PutState(userid, uservalbytes)
	if err != nil {
		return shim.Error("Failed to putstate.\n")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) deleteAccount(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var userid string  //entity
	var account string //accountID, if account  don't exists then error except event = add, else putstate raw to it
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 userid/accountid")
	}

	userid = args[0]
	account = args[1]

	Currentuser := getCertificate(stub)
	if (Currentuser == "Admin@org1.example.com") || (Currentuser == "admin") {
		fmt.Printf("Current operator is Administrator: ")
		fmt.Println(Currentuser)
	} else if Currentuser == -1 {
		return shim.Error("No certificate found")
	} else if Currentuser == -2 {
		return shim.Error("Could not decode the PEM structure")
	} else if Currentuser == -3 {
		return shim.Error("ParseCertificate failed")
	} else if Currentuser != userid {
		return shim.Error("Current operator don't have authority!")
	}

	uservalbytes, err := stub.GetState(userid)
	if err != nil {
		return shim.Error("Failed to get state")
	}
	if uservalbytes == nil {
		return shim.Error("User is not exist, you can initalizate it")
	}

	var userdata User
	err = json.Unmarshal(uservalbytes, &userdata)
	if err != nil {
		return shim.Error("Failed to trans json")
	}

	_, ok := userdata.Accounts[account]
	if !ok {
		return shim.Success([]byte("Account don't exists.\n"))
	}
	delete(userdata.Accounts, account)
	userdata.Accountnum--
	uservalbytes, _ = json.Marshal(userdata)
	err = stub.PutState(userid, uservalbytes)
	if err != nil {
		return shim.Error("Failed to putstate.\n")
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

func (t *SimpleChaincode) queryHistory(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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
func (t *SimpleChaincode) deleteUser(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	A := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(A)
	if err != nil {
		return shim.Error("Failed to delete state")
	}

	return shim.Success(nil)
}

func (t *SimpleChaincode) queryall(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Key string
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments, expecting 1.")
	}
	Key = args[0]
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
	output := "Resbonse: "
	var QueryParameters [][]byte
	var response pb.Response
	for k, v := range data.Accounts {
		if v.ChannelID == "mychannel" {
			QueryParameters = [][]byte{[]byte("queryAccount"), []byte(k), []byte("all")}
			response = stub.InvokeChaincode("pointcc", QueryParameters, "mychannel")
			if response.Status == 200 {
				output = output + "pointasset-Payload:" + string(response.Payload) + "-message:" + response.Message + "\n"
			}
		} else if v.ChannelID == "mychannel" {
			QueryParameters = [][]byte{[]byte("queryall"), []byte(Key)}
			response = stub.InvokeChaincode("mapcc", QueryParameters, "mychannel")
			if response.Status == 200 {
				raw := response.Payload
				var data2 MapUser
				err = json.Unmarshal(raw, &data2)
				if err != nil {
					return shim.Error("read mapchannel error.")
				}
				for k2 := range data2.Accounts {
					QueryParameters = [][]byte{[]byte("queryAccount"), []byte(k2)}
					response = stub.InvokeChaincode("mapcc", QueryParameters, "mychannel")
					if response.Status == 200 {
						output = output + "mapasset-Payload:" + string(response.Payload) + "-message:" + response.Message + "\n"
					}
				}
			}
		}
	}
	return shim.Success([]byte(output))
}

// query callback representing the query of a chaincode

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
