package implementation

import (
	entity "Model"
	"encoding/json"
	"fmt"
	"strings"
	"strconv"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/logging"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// ============================================================
// RegisterPatient - create a new Provider, store into chaincode state
// ============================================================
func (u *User) RegisterProvider(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	providerId := strings.ToLower(args[0])
	providerEHR := strings.ToLower(args[1])
	providerEHRUrl := strings.ToLower(args[2])
	firstname := strings.ToLower(args[3])
	lastname := strings.ToLower(args[4])
	speciality := strings.ToLower(args[5])

	// ==== Check if marble already exists ====
	/*patientData, err := stub.GetState(patientId)
	if err != nil {
		return shim.Error("Failed to get marble: " + err.Error())
	} else if patientData != nil {
		fmt.Println("This marble already exists: " + patientData)
		return shim.Error("This marble already exists: " + patientData)
	}*/

	//==== Create Provider object and marshal to JSON ====
	objectType := "Provider"
	provider := &entity.Provider{objectType, providerId, providerEHR, providerEHRUrl, firstname, lastname, speciality}
	//fmt.Println(Provider.firstname)

	providerJSONasBytes, err := json.Marshal(provider)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//marbleJSONasString := `{"docType":"Marble",  "name": "` + marbleName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//marbleJSONasBytes := []byte(str)

	//Alternatively, build the marble json string manually if you don't want to use struct marshalling
	//patientJSONasString := `{"docType":"Patient",  "patientId": "` + patientId + `", "patientSSN": "` + patientSSN + `", "patientUrl": ` + patientUrl + `, "firstname": "` + firstname + `, "DOB": "` + DOB + `, "email": "` + email + `, "mobile": "` + mobile + `"}`
	//patientJSONasBytes := []byte(patientJSONasString)

	// === Save Provider to state ===

	//err = stub.PutPrivateData("patientDetails", providerId, providerJSONasBytes)
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

func (u *User) GetProviderById(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

func (u *User) UpdateProviderAccess(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}


	logger := logging.NewLogger("log")

	logger.Debugf("log is debugging: %s", args[0])

	var patientDetails entity.PatientDetailsUnmarshal
	var val []byte = []byte("`" + args[0] + "`")
	
	s, err1 := strconv.Unquote(string(val))
	

	if err1 != nil{
		return shim.Error("Error in unquote -->\n" +s+"-->\n"+ args[0] +err1.Error())
	}

	err := json.Unmarshal([]byte(s), &patientDetails)
	
	
	if err != nil{
		return shim.Error("Error in unmarshal input json -->\n" +s+"-->\n"+ args[0] +err.Error())
	}

	patientId, err := getAttribute(stub, "id")
	if err != nil {
		return shim.Error("Fail to get Attribute from private DB " + err.Error())
	} 

	patientDetailsAsBytes, err := stub.GetPrivateData("patientDetails", "123")


	if err != nil {
		return shim.Error("Fail to get patint from private DB " + err.Error())
	}
	
	

	var patientDetailsDB entity.PatientDetails

	err = json.Unmarshal(patientDetailsAsBytes, &patientDetailsDB) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error("ID"+err.Error())
	}

	patientDetailsDB.Allergies.ProviderConsent =  append( patientDetailsDB.Allergies.ProviderConsent , patientDetails.Allergies.ProviderConsent[0])
	patientDetailsDB.Immunization.ProviderConsent =  append( patientDetailsDB.Immunization.ProviderConsent , patientDetails.Immunization.ProviderConsent[0])
	patientDetailsDB.Medications.ProviderConsent =  append( patientDetailsDB.Medications.ProviderConsent , patientDetails.Medications.ProviderConsent[0])
	patientDetailsDB.PastMedicalHx.ProviderConsent =  append( patientDetailsDB.PastMedicalHx.ProviderConsent , patientDetails.PastMedicalHx.ProviderConsent[0])
	
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
		return shim.Error("Error in put private data in two org "+err.Error())
	}

	return shim.Success([]byte("Success"))
}

