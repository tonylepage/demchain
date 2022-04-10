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

import (
	"errors"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DEMstore Chaincode implementation
type DEMstore struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Measurement struct {
	Location		string `json:"location"`
	Measuredepoch		string `json:"measuredepoch"`
	Rtt				string `json:"rtt"`
	CDN				string `json:"cdn"`
	Provider		string `json:"provider"`
}

// InitLedger add a base set of performance data 
func (s *DEMstore) InitLedger(ctx contractapi.TransactionContextInterface) error {
	measurements := []Measurement{
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "Stackpath", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "Fastly", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "Akamai", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "Cloudflare", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "CloudFront", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "GMA", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "Aliyun", provider: "Tony-test"},
		{location: "Taipei, Taiwan", measuredepoch: "1649410093", rtt: "3000", cdn: "CDNetworks", provider: "Tony-test"},
	}

	for _, measurement := range measurements {
		measurementJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil { 
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}


// CreateMeasurement issues a new measurement to the world state with given details.
func (s *DEMstore) CreateMeasurement(ctx contractapi.TransactionContextInterface, location string, measuredepoch string, rtt string, cdn string, provider string) error {
	exists, err := s.MeasurementExists(ctx, location, cdn)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the measurement at %s for %s already exists", location, cdn)
	}

	measurement := Measurement{
		Location: location,
		Measuredepoch: measuredepoch,
		Rtt: rtt,
		CDN: cdn,
		Provider: provider,
	}
	measurementJSON, err := json.Marshal(measurement)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(location)
}




func (t *DEMstore) Init(ctx contractapi.TransactionContextInterface, A string, Aval int, B string, Bval int) error {
	fmt.Println("DEMstore Init")
	var err error
	// Initialize the chaincode
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)
	// Write the state to the ledger
	err = ctx.GetStub().PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return err
	}

	return nil
}

// Transaction makes payment of X units from A to B
func (t *ABstore) Invoke(ctx contractapi.TransactionContextInterface, A, B string, X int) error {
	var err error
	var Aval int
	var Bval int
	// Get the state from the ledger
	// TODO: will be nice to have a GetAllState call to ledger
	Avalbytes, err := ctx.GetStub().GetState(A)
	if err != nil {
		return fmt.Errorf("Failed to get state")
	}
	if Avalbytes == nil {
		return fmt.Errorf("Entity not found")
	}
	Aval, _ = strconv.Atoi(string(Avalbytes))

	Bvalbytes, err := ctx.GetStub().GetState(B)
	if err != nil {
		return fmt.Errorf("Failed to get state")
	}
	if Bvalbytes == nil {
		return fmt.Errorf("Entity not found")
	}
	Bval, _ = strconv.Atoi(string(Bvalbytes))

	// Perform the execution
	Aval = Aval - X
	Bval = Bval + X
	fmt.Printf("Aval = %d, Bval = %d\n", Aval, Bval)

	// Write the state back to the ledger
	err = ctx.GetStub().PutState(A, []byte(strconv.Itoa(Aval)))
	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(B, []byte(strconv.Itoa(Bval)))
	if err != nil {
		return err
	}

	return nil
}

// Delete  an entity from state
func (t *DEMstore) Delete(ctx contractapi.TransactionContextInterface, A string) error {

	// Delete the key from the state in ledger
	err := ctx.GetStub().DelState(A)
	if err != nil {
		return fmt.Errorf("Failed to delete state")
	}

	return nil
}

// Query callback representing the query of a chaincode
func (t *DEMstore) Query(ctx contractapi.TransactionContextInterface, A string) (string, error) {
	var err error
	// Get the state from the ledger
	Avalbytes, err := ctx.GetStub().GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return "", errors.New(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return "", errors.New(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return string(Avalbytes), nil
}

func main() {
	cc, err := contractapi.NewChaincode(new(DEMstore))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting DEMstore chaincode: %s", err)
	}
}
