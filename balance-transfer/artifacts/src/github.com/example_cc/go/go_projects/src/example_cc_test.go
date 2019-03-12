package main

import (
	"fmt"
	"testing"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func checkInit(t *testing.T, stub *shim.MockStub, args [][]byte) {
	res := stub.MockInit("1", args)
	if res.Status != shim.OK {
		fmt.Println("Init failed", string(res.Message))
		t.FailNow()
	}
	if res.Status == shim.OK {
		fmt.Println("Init success")
	}
}

func checkState(t *testing.T, stub *shim.MockStub, name string, value string) {
	bytes := stub.State[name]
	fmt.Println(".............State............")
	if bytes == nil {
		fmt.Println("State", name, "failed to get value")
		t.FailNow()
	}
	if string(bytes) != value {
		fmt.Println("State value", name, "was not", value, "as expected")
		t.FailNow()
	}
	fmt.Println(".............State............")
}

func Test_Init(t *testing.T) {
	simpleCC := new(SimpleChaincode)
	mockStub := shim.NewMockStub("mockstub", simpleCC)
	txId := "mockTxID"

	mockStub.MockTransactionStart(txId)
	response := simpleCC.Init(mockStub)
	mockStub.MockTransactionEnd(txId)
	if s := response.GetStatus(); s != 200 {
		fmt.Println("Init test failed")
		t.FailNow()
	}
}

func Test_initMarble(t *testing.T) {
	simpleCC := new(SimpleChaincode)
	mockStub := shim.NewMockStub("mockstub", simpleCC)
	txId := "mockTxID"

	args := []string{"123123", "asdasd", "asdasd", "asdasd", "123", "234"}
	args1 := []string{"123123"}

	fmt.Println("---------Register---------")
	mockStub.MockTransactionStart(txId)
	response := simpleCC.RegisterPatient(mockStub, args)
	fmt.Println("---------Get Patient---------")
	response = simpleCC.readMarble(mockStub, args1)
	mockStub.MockTransactionEnd(txId)

	fmt.Println("Status: " + fmt.Sprint(response.GetStatus()))
	fmt.Println("Payload: " + string(response.GetPayload()))
	fmt.Println("Message: " + response.GetMessage())

	if s := response.GetStatus(); s != 200 {
		fmt.Println("initMarble test failed")
		t.FailNow()
	}
}
