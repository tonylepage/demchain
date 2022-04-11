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
	Measuredepoch	string `json:"measuredepoch"`
	Rtt				string `json:"rtt"`
	CDN				string `json:"cdn"`
	Provider		string `json:"provider"`
}

// InitLedger add a base set of performance data 
func (s *DEMstore) InitLedger(ctx contractapi.TransactionContextInterface) error {
	measurements := []Measurement{
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "Stackpath", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "Fastly", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "Akamai", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "Cloudflare", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "CloudFront", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "GMA", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "Aliyun", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: "1649410093", Rtt: "3000", CDN: "CDNetworks", Provider: "Tony-test"},
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


func main() {
	cc, err := contractapi.NewChaincode(new(DEMstore))
	if err != nil {
		panic(err.Error())
	}
	if err := cc.Start(); err != nil {
		fmt.Printf("Error starting DEMstore chaincode: %s", err)
	}
}
