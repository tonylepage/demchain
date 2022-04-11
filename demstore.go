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
	"encoding/json"
	"crypto/sha1"
	"fmt"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// DEMstore Chaincode implementation
type DEMstore struct {
	contractapi.Contract
}

// Asset describes basic details of what makes up a simple asset
type Measurement struct {
	ID				string 	`json:"ID"`
	Location		string 	`json:"location"`
	Measuredepoch	int 	`json:"measuredepoch"`
	Rtt				int 	`json:"rtt"`
	CDN				string 	`json:"cdn"`
	Provider		string 	`json:"provider"`
}

// InitLedger add a base set of performance data 
func (s *DEMstore) InitLedger(ctx contractapi.TransactionContextInterface) error {
	measurements := []Measurement{
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "Stackpath", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "Fastly", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "Akamai", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "Cloudflare", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "CloudFront", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "GMA", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "Aliyun", Provider: "Tony-test"},
		{Location: "Taipei, Taiwan", Measuredepoch: 1649410093, Rtt: 3000, CDN: "CDNetworks", Provider: "Tony-test"},
	}

	for _, measurement := range measurements {
		measurementJSON, err := json.Marshal(measurement)
		if err != nil {
			return err
		}

		measurementID := s.GetHashID(ctx, measurement.Location, measurement.CDN)

		err = ctx.GetStub().PutState(measurementID, measurementJSON)
		if err != nil { 
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}


// CreateMeasurement issues a new measurement to the world state with given details.
func (s *DEMstore) CreateMeasurement(ctx contractapi.TransactionContextInterface, location string, measuredepoch string, rtt string, cdn string, provider string) error {
	measurementID := s.GetHashID(ctx, location, cdn)
	exists, err := s.MeasurementExists(ctx, measurementID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the measurement at %s for %s already exists", location, cdn)
	}

	measurement := Measurement{
		ID: measurementID,
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

// ReadMeasurement returns the asset stored in the world state with given id.
func (s *DEMstore) ReadMeasurement(ctx contractapi.TransactionContextInterface, id string) (*Measurement, error) {
    measurementJSON, err := ctx.GetStub().GetState(id)
    if err != nil {
      return nil, fmt.Errorf("failed to read from world state: %v", err)
    }
    if assetJSON == nil {
      return nil, fmt.Errorf("the measurement %s does not exist", id)
    }

    var measurement Measurement
    err = json.Unmarshal(measurementJSON, &measurement)
    if err != nil {
      return nil, err
    }

    return &measurement, nil
}

// GetAllMeasurements returns all measurements found in world state
func (s *DEMstore) GetAllMeasurements(ctx contractapi.TransactionContextInterface) ([]*Measurement, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all measurements in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var measurements []*Measurement
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
		return nil, err
		}

		var measurement Measurement
		err = json.Unmarshal(queryResponse.Value, &measurement)
		if err != nil {
		return nil, err
		}
		measurements = append(measurements, &measurement)
	}

	return measurements, nil
}

// UpdateMeasurement updates an existing measurement in the world state with provided parameters.
func (s *DEMstore) UpdateMeasurement(ctx contractapi.TransactionContextInterface, location string, measuredepoch string, rtt string, cdn string, provider string) error {
	measurementID := s.GetHashID(ctx, location, cdn)
	exists, err := s.MeasurementExists(ctx, measurementID)
	if err != nil {
	  return err
	}
	if !exists {
	  return fmt.Errorf("the measurement at %s for %s does not exist", location, cdn)
	}

	// overwriting original asset with new asset
	measurement := Measurement{
		ID: measurementID,
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

	return ctx.GetStub().PutState(measurementID, measurementJSON)
}

// AssetExists returns true when asset with given ID exists in world state
func (s *DEMstore) MeasurementExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	measurementJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
	  return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return measurementJSON != nil, nil
}

// GetHashID returns the hash of city and cdn to be used as a key
func (s *DEMstore) GetHashID(ctx contractapi.TransactionContextInterface, location string, cdn string) (string, error) {

	rawID := location + cdn
	hash := sha1.New()
	hash.Write([]byte(rawID))
	rhash := hash.Sum(nil)

	return rhash, nil
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
