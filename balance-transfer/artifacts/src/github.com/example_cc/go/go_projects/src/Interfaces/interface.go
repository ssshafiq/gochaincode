package Interfaces

import (

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)
type User struct {
}

//Repository repository interface
type InterfacePatient interface {
	RegisterPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetPatientByInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response
	}

type InterfaceProvider interface {
	RegisterProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response
	GetProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response
	UpdateProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response
}
