package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	cid "github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/pkg/errors"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type marble struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	Name       string `json:"name"`    //the fieldtags are needed to keep case from bouncing around
	Color      string `json:"color"`
	Size       int    `json:"size"`
	Owner      string `json:"owner"`
}

// Information of Patient

type Patient struct {
	ObjectType       string `json:docType"`
	PatientId        string `json:"patientId"`
	PatientSSN       string `json:"patientssn"`
	PatientUrl       string `json:"patienturl"`
	PatientFirstname string `json:"firstname"` //docType is used to distinguish the various types of objects in state database
	PatientLastname  string `json:"lastname"`  //the fieldtags are needed to keep case from bouncing around
	DOB              string `json:"dob"`
}

// Information of Provider
type Provider struct {
	ObjectType        string `json:docType"`
	ProviderId        string `json:"providerId"`
	ProviderEHR       string `json:"providerehr"`
	ProviderEHRURL    string `json:"providerehrurl"`
	ProviderFirstname string `json:"firstname"` //docType is used to distinguish the various types of objects in state database
	ProviderLastname  string `json:"lastname"`  //the fieldtags are needed to keep case from bouncing around
	Speciality        string `json:"speciality"`
}

type Consent struct {
	ObjectType string   `json:docType"`
	Provider   Provider `json:"provider"`
	StartTime  string   `json:"starttime"`
	EndTime    string   `json:"endtime"`
}

type PatientDetails struct {
	Medications   Medications   `json:"medications"`
	Allergies     Allergies     `json:"allergies"`
	Immunization  Immunization  `json:"immunization"`
	PastMedicalHx PastMedicalHx `json:"pastMedicalHx"`
	FamilyHx      FamilyHx      `json:"familyHx"`
}

type Medications struct {
	ObjectType      string    `json:docType"`
	Patient         Patient   `json:"patient"`
	ProviderConsent []Consent `json:"providerconsent"`
}
type Allergies struct {
	ObjectType      string    `json:docType"`
	Patient         Patient   `json:"patient"`
	ProviderConsent []Consent `json:"providerconsent"`
}
type Immunization struct {
	ObjectType      string    `json:docType"`
	Patient         Patient   `json:"patient"`
	ProviderConsent []Consent `json:"providerconsent"`
}
type PastMedicalHx struct {
	ObjectType      string    `json:docType"`
	Patient         Patient   `json:"patient"`
	ProviderConsent []Consent `json:"providerconsent"`
}
type FamilyHx struct {
	ObjectType      string    `json:docType"`
	Patient         Patient   `json:"patient"`
	ProviderConsent []Consent `json:"providerconsent"`
}

type PatientUnmarshal struct {
	Key    string
	Record Patient `json:"Patient"`
}

type PatientDetailsUnmarshal struct {
	_id           string
	_rev          string
	Allergies     Allergies     `json:"allergies"`
	FamilyHx      FamilyHx      `json:"familyHx"`
	Immunization  Immunization  `json:"immunization"`
	Medications   Medications   `json:"medications"`
	PastMedicalHx PastMedicalHx `json:"pastMedicalHx"`
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
	if function == "transferMarble" { //change owner of a specific marble
		return t.transferMarble(stub, args)
	} else if function == "transferMarblesBasedOnColor" { //transfer all marbles of a certain color
		return t.transferMarblesBasedOnColor(stub, args)
	} else if function == "delete" { //delete a marble
		return t.delete(stub, args)
	} else if function == "readMarble" { //read a marble
		return t.readMarble(stub, args)
	} else if function == "queryMarblesByOwner" { //find marbles for owner X using rich query
		return t.queryMarblesByOwner(stub, args)
	} else if function == "queryMarbles" { //find marbles based on an ad hoc rich query
		return t.queryMarbles(stub, args)
	} else if function == "getHistoryForMarble" { //get history of values for a marble
		return t.getHistoryForMarble(stub, args)
	} else if function == "getMarblesByRange" { //get marbles based on range query
		return t.getMarblesByRange(stub, args)
	} else if function == "getMarblesByRangeWithPagination" {
		return t.getMarblesByRangeWithPagination(stub, args)
	} else if function == "queryMarblesWithPagination" {
		return t.queryMarblesWithPagination(stub, args)
	} else if function == "RegisterPatient" {
		return t.RegisterPatient(stub, args)
	} else if function == "GetPatientBySSN" {
		return t.GetPatientBySSN(stub, args)
	} else if function == "GetPatientByInformation" {
		return t.GetPatientByInformation(stub, args)
	} else if function == "RegisterProvider" {
		return t.RegisterProvider(stub, args)
	} else if function == "GetProviderById" {
		return t.GetProviderById(stub, args)
	} else if function == "UpdateProviderAccess" {
		return t.UpdateProviderAccess(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// RegisterPatient - create a new Patient, store into chaincode state
// ============================================================
func (t *SimpleChaincode) RegisterPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	fmt.Println("- start register patient")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	mspRole, err := t.getAttribute(stub, "mspRole")
	if err != nil {
		return shim.Error("Fails to get mspRole " + err.Error())
	}

	if mspRole == "client" || mspRole == "admin" {

		patientId := strings.ToLower(args[0])
		patientSSN := strings.ToLower(args[1])
		patientUrl := strings.ToLower(args[2])
		firstname := strings.ToLower(args[3])
		lastname := strings.ToLower(args[4])
		DOB := strings.ToLower(args[5])

		// ==== Check if marble already exists ====
		/*patientData, err := stub.GetState(patientId)
		if err != nil {
			return shim.Error("Failed to get marble: " + err.Error())
		} else if patientData != nil {
			fmt.Println("This marble already exists: " + patientData)
			return shim.Error("This marble already exists: " + patientData)
		}*/

		//==== Create Patient object and marshal to JSON ====
		objectType := "Patient"
		patient := &Patient{objectType, patientId, patientSSN, patientUrl, firstname, lastname, DOB}
		patientJSONasBytes, err := json.Marshal(patient)

		//==== Create patientMedications object and marshal to JSON ====
		providerId, err := t.getAttribute(stub, "id")
		if err != nil {
			return shim.Error("Fails to get id " + err.Error())
		}

		providerAsByte, err := stub.GetState(providerId)
		if err != nil {
			return shim.Error("Fails to get provider: " + err.Error())
		}

		var provider Provider
		err = json.Unmarshal(providerAsByte, &provider)

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		var patientdetails PatientDetails
		patientdetails.Medications.ObjectType = "Medications"
		patientdetails.Medications.Patient = *patient
		var defaultConsent Consent
		defaultConsent.Provider = provider
		defaultConsent.StartTime = time.Now().Format("01-02-2006")
		defaultConsent.EndTime = time.Now().Format("01-02-2006")
		patientdetails.Medications.ProviderConsent = []Consent{}
		patientdetails.Medications.ProviderConsent = append(patientdetails.Medications.ProviderConsent, defaultConsent)

		//==== Create patientAllergies object and marshal to JSON ====
		patientdetails.Allergies.ObjectType = "Allergies"
		patientdetails.Allergies.Patient = *patient
		patientdetails.Allergies.ProviderConsent = []Consent{}
		patientdetails.Allergies.ProviderConsent = append(patientdetails.Allergies.ProviderConsent, defaultConsent)

		//==== Create patientImmunizations object and marshal to JSON ====
		patientdetails.Immunization.ObjectType = "Immunizations"
		patientdetails.Immunization.Patient = *patient
		patientdetails.Immunization.ProviderConsent = []Consent{}
		patientdetails.Immunization.ProviderConsent = append(patientdetails.Immunization.ProviderConsent, defaultConsent)

		//==== Create patientPastMedicalHx object and marshal to JSON ====
		patientdetails.PastMedicalHx.ObjectType = "PastMedicalHx"
		patientdetails.PastMedicalHx.Patient = *patient
		patientdetails.PastMedicalHx.ProviderConsent = []Consent{}
		patientdetails.PastMedicalHx.ProviderConsent = append(patientdetails.PastMedicalHx.ProviderConsent, defaultConsent)

		//==== Create patientFamilyHx object and marshal to JSON ====
		patientdetails.FamilyHx.ObjectType = "FamilyHx"
		patientdetails.FamilyHx.Patient = *patient
		patientdetails.FamilyHx.ProviderConsent = []Consent{}
		patientdetails.FamilyHx.ProviderConsent = append(patientdetails.FamilyHx.ProviderConsent, defaultConsent)

		PatientDetailsJSONasBytes, err := json.Marshal(&patientdetails)

		if err != nil {
			return shim.Error(err.Error())
		}

		// === Save patientDetails to state ===
		err = stub.PutPrivateData("patientDetails", patientId, PatientDetailsJSONasBytes)
		//err = stub.PutState(patientId, PatientDetailsJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//=== Save Patient to state ===
		err = stub.PutState(patientId, patientJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//  ==== Index the Patient to enable name-based range queries, e.g. return all Patients ====
		//  An 'index' is a normal key/value entry in state.
		//  The key is a composite key, with the elements that you want to range query on listed first.
		//  In our case, the composite key is based on indexName~color~name.
		//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
		indexName := "fname~lname"
		fnameLnameIndexKey, err := stub.CreateCompositeKey(indexName, []string{patient.PatientFirstname, patient.PatientLastname})
		if err != nil {
			return shim.Error(err.Error())
		}
		//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		stub.PutState(fnameLnameIndexKey, value)

		// ==== Marble saved and indexed. Return success ====
		//fmt.Println("- end register patient")
		return shim.Success(nil)
	}
	return shim.Error("Unauthorized!")

}

// ============================================================
// RegisterPatient - create a new Provider, store into chaincode state
// ============================================================
func (t *SimpleChaincode) RegisterProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	fmt.Println("- start register patient")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("4th argument must be a non-empty string")
	}

	mspRole, err := t.getAttribute(stub, "mspRole")
	if err != nil {
		return shim.Error("Fails to get mspRole " + err.Error())
	}

	if mspRole == "admin" {
		providerId := strings.ToLower(args[0])
		providerEHR := strings.ToLower(args[1])
		providerEHRUrl := strings.ToLower(args[2])
		firstname := strings.ToLower(args[3])
		lastname := strings.ToLower(args[4])
		speciality := strings.ToLower(args[5])

		//==== Create Provider object and marshal to JSON ====
		objectType := "Provider"
		provider := &Provider{objectType, providerId, providerEHR, providerEHRUrl, firstname, lastname, speciality}

		providerJSONasBytes, err := json.Marshal(provider)
		if err != nil {
			return shim.Error(err.Error())
		}

		//Registering user to CouchDB
		err = stub.PutState(providerId, providerJSONasBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		//  ==== Index the Provider to enable name-based range queries, e.g. return all Patients ====
		//  An 'index' is a normal key/value entry in state.
		//  The key is a composite key, with the elements that you want to range query on listed first.
		//  In our case, the composite key is based on indexName~color~name.
		//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
		indexName := "fname~lname"
		fnameLnameIndexKey, err := stub.CreateCompositeKey(indexName, []string{provider.ProviderFirstname, provider.ProviderLastname})
		if err != nil {
			return shim.Error(err.Error())
		}
		//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
		//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
		value := []byte{0x00}
		stub.PutState(fnameLnameIndexKey, value)

		// ==== Marble saved and indexed. Return success ====
		//fmt.Println("- end register patient")
		return shim.Success(nil)
	}
	return shim.Error("Unauthorized!")

}

// ==============================================
// Search Patient using its SSN
// ==============================================
func (t *SimpleChaincode) GetPatientBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("In seachpatient by sssn")

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	role, err := t.getAttribute(stub, "role")
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("=======Role==============")
	fmt.Println(role)

	if strings.HasPrefix(role, "Patient") {
		ssn := strings.ToLower(args[0])

		queryString := fmt.Sprintf("{\"selector\":{\"patientssn\":\"%s\"}}", ssn)

		queryResults, err := getQueryResultForQueryString(stub, queryString)
		if err != nil {
			return shim.Error(err.Error())
		}

		var tempArray []PatientUnmarshal
		err = json.Unmarshal(queryResults, &tempArray)
		if err != nil {
			return shim.Error(err.Error())
		}

		var key string
		for _, patient := range tempArray {

			key = patient.Key

		}

		if strings.Contains(role, key) {

			patientDetailsBytes, err := stub.GetPrivateData("patientDetails", key)
			if err != nil {
				return shim.Error("Patient not found " + key + "role " + role + "patient details " + string(patientDetailsBytes))
			}
			return shim.Success(patientDetailsBytes)
		} else {
			return shim.Error("unAuthorized role: " + role + "key: " + key)
		}

	} else {
		// When other
		return shim.Error("Only patients, doctors and pharmacies can access medical details")
	}
}

func (t *SimpleChaincode) GetProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	id := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"_id\":\"%s\"}}", id)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// ==============================================
// Search Patient using its Info
// ==============================================
func (t *SimpleChaincode) GetPatientByInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	fname := strings.ToLower(args[0])
	lname := strings.ToLower(args[1])
	dob := strings.ToLower(args[2])

	queryString := fmt.Sprintf("{\"selector\":{\"firstname\":\"%s\",\"lastname\":\"%s\",\"dob\":\"%s\"}}", fname, lname, dob)

	// queryString := fmt.Sprintf("{\"selector\":{\"ObjectType\":\"Patient\",\"_id\":\"%s\"}}", id)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

func (t *SimpleChaincode) UpdateProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	logger := logging.NewLogger("log")

	logger.Debugf("log is debugging: %s", args[0])

	var patientDetails PatientDetailsUnmarshal
	var val []byte = []byte("`" + args[0] + "`")

	s, err1 := strconv.Unquote(string(val))

	if err1 != nil {
		return shim.Error("Error in unquote -->\n" + s + "-->\n" + args[0] + err1.Error())
	}

	err := json.Unmarshal([]byte(s), &patientDetails)

	if err != nil {
		return shim.Error("Error in unmarshal input json -->\n" + s + "-->\n" + args[0] + err.Error())
	}

	patientId, err := t.getAttribute(stub, "id")
	if err != nil {
		return shim.Error("Fail to get Attribute from private DB " + err.Error())
	}

	patientDetailsAsBytes, err := stub.GetPrivateData("patientDetails", "123")

	if err != nil {
		return shim.Error("Fail to get patint from private DB " + patientDetails._id + err.Error())
	}

	var patientDetailsDB PatientDetails

	err = json.Unmarshal(patientDetailsAsBytes, &patientDetailsDB) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error("ID" + patientDetails._id + err.Error())
	}

	patientDetailsDB.Allergies.ProviderConsent = append(patientDetailsDB.Allergies.ProviderConsent, patientDetails.Allergies.ProviderConsent[0])
	patientDetailsDB.Immunization.ProviderConsent = append(patientDetailsDB.Immunization.ProviderConsent, patientDetails.Immunization.ProviderConsent[0])
	patientDetailsDB.Medications.ProviderConsent = append(patientDetailsDB.Medications.ProviderConsent, patientDetails.Medications.ProviderConsent[0])
	patientDetailsDB.PastMedicalHx.ProviderConsent = append(patientDetailsDB.PastMedicalHx.ProviderConsent, patientDetails.PastMedicalHx.ProviderConsent[0])

	PatientDetailsJSONasBytes, err := json.Marshal(&patientDetailsDB)

	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save patientDetails to state ===
	// err = stub.PutPrivateData("patientDetails", patientId, PatientDetailsJSONasBytes)
	// if err != nil {
	// 	return shim.Error( "Error in put private data in one org" +err.Error())
	// }

	err = stub.PutPrivateData("patientDetailsIn2Orgs", patientId, PatientDetailsJSONasBytes)
	if err != nil {
		return shim.Error("Error in put private data in two org " + err.Error())
	}

	return shim.Success([]byte("Success"))
}

func (s *SimpleChaincode) getAttribute(stub shim.ChaincodeStubInterface, key string) (string, error) {

	role, ok, err := cid.GetAttributeValue(stub, key)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("role attribute is missing")
	}

	return role, nil
}

// ===============================================
// readMarble - read a marble from chaincode state
// ===============================================
func (t *SimpleChaincode) readMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the marble to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + name + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}

// ==================================================
// delete - remove a marble key/value pair from state
// ==================================================
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var jsonResp string
	var marbleJSON marble
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	marbleName := args[0]

	// to maintain the color~name index, we need to read the marble first and get its color
	valAsbytes, err := stub.GetState(marbleName) //get the marble from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + marbleName + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + marbleName + "\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal([]byte(valAsbytes), &marbleJSON)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to decode JSON of: " + marbleName + "\"}"
		return shim.Error(jsonResp)
	}

	err = stub.DelState(marbleName) //remove the marble from chaincode state
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}

	// maintain the index
	indexName := "color~name"
	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{marbleJSON.Color, marbleJSON.Name})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  Delete index entry to state.
	err = stub.DelState(colorNameIndexKey)
	if err != nil {
		return shim.Error("Failed to delete state:" + err.Error())
	}
	return shim.Success(nil)
}

// ===========================================================
// transfer a marble by setting a new owner name on the marble
// ===========================================================
func (t *SimpleChaincode) transferMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	marbleName := args[0]
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferMarble ", marbleName, newOwner)

	marbleAsBytes, err := stub.GetState(marbleName)
	if err != nil {
		return shim.Error("Failed to get marble:" + err.Error())
	} else if marbleAsBytes == nil {
		return shim.Error("Marble does not exist")
	}

	marbleToTransfer := marble{}
	err = json.Unmarshal(marbleAsBytes, &marbleToTransfer) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	marbleToTransfer.Owner = newOwner //change the owner

	marbleJSONasBytes, _ := json.Marshal(marbleToTransfer)
	err = stub.PutState(marbleName, marbleJSONasBytes) //rewrite the marble
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end transferMarble (success)")
	return shim.Success(nil)
}

// ===========================================================================================
// constructQueryResponseFromIterator constructs a JSON array containing query results from
// a given result iterator
// ===========================================================================================
func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
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
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return &buffer, nil
}

// ===========================================================================================
// addPaginationMetadataToQueryResults adds QueryResponseMetadata, which contains pagination
// info, to the constructed query results
// ===========================================================================================
func addPaginationMetadataToQueryResults(buffer *bytes.Buffer, responseMetadata *pb.QueryResponseMetadata) *bytes.Buffer {

	buffer.WriteString("[{\"ResponseMetadata\":{\"RecordsCount\":")
	buffer.WriteString("\"")
	buffer.WriteString(fmt.Sprintf("%v", responseMetadata.FetchedRecordsCount))
	buffer.WriteString("\"")
	buffer.WriteString(", \"Bookmark\":")
	buffer.WriteString("\"")
	buffer.WriteString(responseMetadata.Bookmark)
	buffer.WriteString("\"}}]")

	return buffer
}

// ===========================================================================================
// getMarblesByRange performs a range query based on the start and end keys provided.

// Read-only function results are not typically submitted to ordering. If the read-only
// results are submitted to ordering, or if the query is used in an update transaction
// and submitted to ordering, then the committing peers will re-execute to guarantee that
// result sets are stable between endorsement time and commit time. The transaction is
// invalidated by the committing peers if the result set has changed between endorsement
// time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (t *SimpleChaincode) getMarblesByRange(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	startKey := args[0]
	endKey := args[1]

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// ==== Example: GetStateByPartialCompositeKey/RangeQuery =========================================
// transferMarblesBasedOnColor will transfer marbles of a given color to a certain new owner.
// Uses a GetStateByPartialCompositeKey (range query) against color~name 'index'.
// Committing peers will re-execute range queries to guarantee that result sets are stable
// between endorsement time and commit time. The transaction is invalidated by the
// committing peers if the result set has changed between endorsement time and commit time.
// Therefore, range queries are a safe option for performing update transactions based on query results.
// ===========================================================================================
func (t *SimpleChaincode) transferMarblesBasedOnColor(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0       1
	// "color", "bob"
	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	color := args[0]
	newOwner := strings.ToLower(args[1])
	fmt.Println("- start transferMarblesBasedOnColor ", color, newOwner)

	// Query the color~name index by color
	// This will execute a key range query on all keys starting with 'color'
	coloredMarbleResultsIterator, err := stub.GetStateByPartialCompositeKey("color~name", []string{color})
	if err != nil {
		return shim.Error(err.Error())
	}
	defer coloredMarbleResultsIterator.Close()

	// Iterate through result set and for each marble found, transfer to newOwner
	var i int
	for i = 0; coloredMarbleResultsIterator.HasNext(); i++ {
		// Note that we don't get the value (2nd return variable), we'll just get the marble name from the composite key
		responseRange, err := coloredMarbleResultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// get the color and name from color~name composite key
		objectType, compositeKeyParts, err := stub.SplitCompositeKey(responseRange.Key)
		if err != nil {
			return shim.Error(err.Error())
		}
		returnedColor := compositeKeyParts[0]
		returnedMarbleName := compositeKeyParts[1]
		fmt.Printf("- found a marble from index:%s color:%s name:%s\n", objectType, returnedColor, returnedMarbleName)

		// Now call the transfer function for the found marble.
		// Re-use the same function that is used to transfer individual marbles
		response := t.transferMarble(stub, []string{returnedMarbleName, newOwner})
		// if the transfer failed break out of loop and return error
		if response.Status != shim.OK {
			return shim.Error("Transfer failed: " + response.Message)
		}
	}

	responsePayload := fmt.Sprintf("Transferred %d %s marbles to %s", i, color, newOwner)
	fmt.Println("- end transferMarblesBasedOnColor: " + responsePayload)
	return shim.Success([]byte(responsePayload))
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
func (t *SimpleChaincode) queryMarblesByOwner(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	owner := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"docType\":\"marble\",\"owner\":\"%s\"}}", owner)

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
func (t *SimpleChaincode) queryMarbles(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

// ====== Pagination =========================================================================
// Pagination provides a method to retrieve records with a defined pagesize and
// start point (bookmark).  An empty string bookmark defines the first "page" of a query
// result.  Paginated queries return a bookmark that can be used in
// the next query to retrieve the next page of results.  Paginated queries extend
// rich queries and range queries to include a pagesize and bookmark.
//
// Two examples are provided in this example.  The first is getMarblesByRangeWithPagination
// which executes a paginated range query.
// The second example is a paginated query for rich ad-hoc queries.
// =========================================================================================

// ====== Example: Pagination with Range Query ===============================================
// getMarblesByRangeWithPagination performs a range query based on the start & end key,
// page size and a bookmark.

// The number of fetched records will be equal to or lesser than the page size.
// Paginated range queries are only valid for read only transactions.
// ===========================================================================================
func (t *SimpleChaincode) getMarblesByRangeWithPagination(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	startKey := args[0]
	endKey := args[1]
	//return type of ParseInt is int64
	pageSize, err := strconv.ParseInt(args[2], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	bookmark := args[3]

	resultsIterator, responseMetadata, err := stub.GetStateByRangeWithPagination(startKey, endKey, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)

	fmt.Printf("- getMarblesByRange queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return shim.Success(buffer.Bytes())
}

// ===== Example: Pagination with Ad hoc Rich Query ========================================================
// queryMarblesWithPagination uses a query string, page size and a bookmark to perform a query
// for marbles. Query string matching state database syntax is passed in and executed as is.
// The number of fetched records would be equal to or lesser than the specified page size.
// Supports ad hoc queries that can be defined at runtime by the client.
// If this is not desired, follow the queryMarblesForOwner example for parameterized queries.
// Only available on state databases that support rich query (e.g. CouchDB)
// Paginated queries are only valid for read only transactions.
// =========================================================================================
func (t *SimpleChaincode) queryMarblesWithPagination(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//   0
	// "queryString"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	queryString := args[0]
	//return type of ParseInt is int64
	pageSize, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return shim.Error(err.Error())
	}
	bookmark := args[2]

	queryResults, err := getQueryResultForQueryStringWithPagination(stub, queryString, int32(pageSize), bookmark)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryStringWithPagination executes the passed in query string with
// pagination info. Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryStringWithPagination(stub shim.ChaincodeStubInterface, queryString string, pageSize int32, bookmark string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, responseMetadata, err := stub.GetQueryResultWithPagination(queryString, pageSize, bookmark)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	bufferWithPaginationInfo := addPaginationMetadataToQueryResults(buffer, responseMetadata)

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", bufferWithPaginationInfo.String())

	return buffer.Bytes(), nil
}

func (t *SimpleChaincode) getHistoryForMarble(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	marbleName := args[0]

	fmt.Printf("- start getHistoryForMarble: %s\n", marbleName)

	resultsIterator, err := stub.GetHistoryForKey(marbleName)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForMarble returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}
