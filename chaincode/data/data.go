/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/util"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type data struct {
	//ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	ID         string   `json:"id"`
	URI        string   `json:"uri"`
	Key        string   `json:"key"`
	ClearHash  string   `json:"clear_hash"`
	CipherHash string   `json:"cipher_hash"`
	Doc        string   `json:"doc"`
	Creater    string   `json:"creater"`
	Owner      []string `json:"owner"`
	Pid        string   `json:"pid"`
	Timestamp  string   `json:"timestamp"`
}

// Check data is valid
func dataIsValid(uri, key, clearhash, cipherhash string) bool {
	// 1. uri.GET() = DEHASH(cipherhash)
	// 2. DEHASH(cipherhash) = key.ENCODE(clearhash)
	return true
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "create" { //create a new data item
		return t.create(stub, args, "")
	} else if function == "checkOut" { //checkout data item
		return t.checkOut(stub, args)
	} else if function == "addOwner" { //append owner of a data item
		return t.addOwner(stub, args)
	} else if function == "modify" { //set a new version of a data item
		return t.modify(stub, args)
	} else if function == "tracelog" { //get a trace of a data item
		return t.tracelog(stub, args)
	} else if function == "queryByOwner" { //find data items for owner X using rich query
		return t.queryByOwner(stub, args)
	} else if function == "hstory" { //get history of values for a data item
		return t.history(stub, args)
	}

	fmt.Println("Invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initMarble - create a new marble, store into chaincode state
// ============================================================
func (t *SimpleChaincode) create(stub shim.ChaincodeStubInterface, args []string, pid string) pb.Response {
	var err error

	// 0	1		2		3		4		5		6
	// id 	uri		key		chash	chash	doc		creater		owner	trace
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	fmt.Println("- create data")
	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(strconv.Itoa(i) + "argument must be a non-empty string")
		}
	}

	creater, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	timestamp, _ := stub.GetTxTimestamp()
	timestr := time.Unix(timestamp.Seconds, int64(timestamp.Nanos)).String()

	id := args[0]
	uri := args[1]
	key := args[2]
	clearhash := args[3]
	cipherhash := args[4]
	doc := args[5]
	owner := []string{creater}

	if !dataIsValid(key, clearhash, cipherhash) {
		return shim.Error("Data is not valid")
	}

	// ==== Check if data item already exists ====
	dataAsBytes, err := stub.GetState(id)
	if err != nil {
		return shim.Error("Failed to get data: " + err.Error())
	} else if dataAsBytes != nil {
		fmt.Println("This data already exists: " + id)
		return shim.Error("This data already exists: " + id)
	}

	// ==== Create data object and marshal to JSON ====
	//objectType := "data"
	data := &data{id, uri, key, clearhash, cipherhash, doc, creater, owner, pid, timestr}
	dataJSONasBytes, err := json.Marshal(data)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save data to state ===
	err = stub.PutState(id, dataJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// leave compositekey if needed
	// //  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	// //  An 'index' is a normal key/value entry in state.
	// //  The key is a composite key, with the elements that you want to range query on listed first.
	// //  In our case, the composite key is based on indexName~color~name.
	// //  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	// indexName := "color~name"
	// colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{marble.Color, marble.Name})
	// if err != nil {
	// 	return shim.Error(err.Error())
	// }
	// //  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	// //  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	// value := []byte{0x00}
	// stub.PutState(colorNameIndexKey, value)

	// ==== data saved and indexed. Return success ====
	fmt.Println("- end create data")
	return shim.Success(nil)
}

func (t *SimpleChaincode) tracelog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting id of the data to query")
	}

	var id, jsonResp string
	var dataJSON data
	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the data from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Data does not exist, id: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &dataJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	currentUser, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	} else if !util.Contains(dataJSON.Owner, currentUser) && !util.IsAdmin(currentUser) {
		return shim.Error("Permission denied")
	}

	log.Println("- trace begin from the leaf: " + id)
	jsonResp = "["
	jsonResp += string(valAsbytes)
	// trace `pid` until root
	var pid = dataJSON.Pid
	for "" != pid {
		valAsbytes, err = stub.GetState(pid) //get the data from chaincode state
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
			return shim.Error(jsonResp)
		} else if valAsbytes == nil {
			jsonResp = "{\"Error\":\"Data does not exist, id: " + id + "\"}"
			return shim.Error(jsonResp)
		}
		err = json.Unmarshal([]byte(valAsbytes), &dataJSON)
		if err != nil {
			jsonResp = "{\"Error\":\"Failed to decode JSON of: " + id + "\"}"
			return shim.Error(jsonResp)
		}
		jsonResp += ", "
		jsonResp += string(valAsbytes)
		pid = dataJSON.Pid
	}

	jsonResp += "]"
	return shim.Success([]byte(jsonResp))
}

func (t *SimpleChaincode) addOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0	1
	// id 	newowner
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	log.Println("- add owner")

	var id, jsonResp string
	var dataJSON data
	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the data from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Data does not exist, id: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &dataJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	currentUser, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	} else if !util.Contains(dataJSON.Owner, currentUser) && !util.IsAdmin(currentUser) {
		return shim.Error("Permission denied")
	}

	newowner := args[1]
	if len(newowner) <= 0 {
		return shim.Error("Argument must be a non-empty string")
	}
	dataJSON.Owner = append(dataJSON.Owner, newowner) // append the owner

	dataJSONasBytes, _ := json.Marshal(dataJSON)
	err = stub.PutState(id, dataJSONasBytes) //rewrite the data
	if err != nil {
		return shim.Error(err.Error())
	}

	log.Println("- end add owner")
	return shim.Success(nil)
}

func (t *SimpleChaincode) modify(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 0	1 	 2    3 	4 	 5		6
	// oid  id 	 uri  key  chash chash doc
	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	log.Println("- modify data item")

	var dataJSON data
	var id, jsonResp string

	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the data from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Data does not exist, id: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &dataJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	currentUser, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	} else if !util.Contains(dataJSON.Owner, currentUser) && !util.IsAdmin(currentUser) {
		return shim.Error("Permission denied")
	}

	return t.create(stub, args[1:], id)
}

// ===============================================
// readMarble - read a marble from chaincode state
// ===============================================
func (t *SimpleChaincode) checkOut(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting id of the data to query")
	}

	var id, jsonResp string
	var dataJSON data
	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the data from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Data does not exist, id: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &dataJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	currentUser, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	} else if !util.Contains(dataJSON.Owner, currentUser) && !util.IsAdmin(currentUser) {
		return shim.Error("Permission denied")
	}

	return shim.Success(valAsbytes)
}

// =======Rich queries =========================================================================
// Two examples of rich queries are provided below (parameterized query and ad hoc query).
// Rich queries pass a query string to the state database.
// Rich queries are only supported by state database implementations
//  that support rich query (e.g. CouchDB).
// The query string is in the syntax of the underlying state database.
// With rich queries there is no guarantee that the result set hasn't changed between
//  endorsement time and commit time, aka 'phantom reads'.
// Therefore, rich queries should not be used in update transactions, unless the
// application handles the possibility of result set changes between endorsement and commit time.
// Rich queries can be used for point-in-time queries against a peer.
// ============================================================================================

// ===== Example: Parameterized rich query =================================================
// queryMarblesByOwner queries for marbles based on a passed in owner.
// This is an example of a parameterized query where the query logic is baked into the chaincode,
// and accepting a single query parameter (owner).
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) queryByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	//if len(args) < 1 {
	//	return shim.Error("Incorrect number of arguments. Expecting 1")
	//}

	//owner := strings.ToLower(args[0])
	owner, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	}

	//queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"marble\",\"owner\":\"%s\"}}", owner)
	queryString := fmt.Sprintf("{\"selector\":{\"owner\":{\"$elemMatch\":{\"$eq\":\"%s\"}}}}", owner)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ===== Example: Ad hoc rich query ========================================================
// queryMarbles uses a query string to perform a query for marbles.
// Query string matching state database syntax is passed in and executed as is.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryMarblesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// =========================================================================================
func (t *SimpleChaincode) queryData(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	queryString := args[0]

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(stub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := stub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
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
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"value\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) history(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	var id, jsonResp string
	var dataJSON data
	id = args[0]
	valAsbytes, err := stub.GetState(id) //get the data from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Data does not exist, id: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	err = json.Unmarshal([]byte(valAsbytes), &dataJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + id + "\"}"
		return shim.Error(jsonResp)
	}
	currentUser, err := util.GetUser(stub)
	if err != nil {
		return shim.Error(err.Error())
	} else if !util.Contains(dataJSON.Owner, currentUser) && !util.IsAdmin(currentUser) {
		return shim.Error("Permission denied")
	}

	fmt.Printf("- start history, id: %s\n", id)

	resultsIterator, err := stub.GetHistoryForKey(id)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	bResult, err := util.HistoryIter2json(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(bResult)
}
