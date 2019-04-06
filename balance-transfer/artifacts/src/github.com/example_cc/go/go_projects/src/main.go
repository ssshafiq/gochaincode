package main

import (
	implementation "Implementation"
	"fmt"
	inf "Interfaces"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
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

// ========================================
// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)
        u := &implementation.User{}
	// Handle different functions
	if function == "RegisterPatient" {
		return inf.InterfacePatient.RegisterPatient(u, stub, args)
	} else if function == "GetPatientBySSN" {
		return inf.InterfacePatient.GetPatientBySSN(u, stub, args)
	} else if function == "GetPatientByInformation" {
		return inf.InterfacePatient.GetPatientByInformation(u, stub, args)
	} else if function == "RegisterProvider" {
		return inf.InterfaceProvider.RegisterProvider(u, stub, args)
	} else if function == "GetProviderById" {
		return inf.InterfaceProvider.GetProviderById(u, stub, args)
	} else if function == "UpdateProviderAccess" {
		return inf.InterfaceProvider.UpdateProviderAccess(u, stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}
