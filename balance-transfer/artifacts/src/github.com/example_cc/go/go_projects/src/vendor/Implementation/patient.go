package implementation
import (
	entity "Model"
	"encoding/json"
	"fmt"
	"strings"
	"time"
	inf "Interfaces"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type User inf.User
// ============================================================
// RegisterPatient - create a new Patient, store into chaincode state
// ============================================================
func (u *User) RegisterPatient(stub shim.ChaincodeStubInterface, args []string) pb.Response {
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

	mspRole, err := getAttribute(stub, "mspRole")
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
		patient := &entity.Patient{objectType, patientId, patientSSN, patientUrl, firstname, lastname, DOB}
		patientJSONasBytes, err := json.Marshal(patient)

		//==== Create patientMedications object and marshal to JSON ====
		providerId, err := getAttribute(stub, "id")
		if err != nil {
			return shim.Error("Fails to get id " + err.Error())
		}

		providerAsByte, err := stub.GetState(providerId)
		if err != nil {
			return shim.Error("Fails to get provider: " + err.Error())
		}

		var provider entity.Provider
		err = json.Unmarshal(providerAsByte, &provider)

		if err != nil {
			return shim.Error("Fails to unmarshal provider " + err.Error())
		}

		var patientdetails entity.PatientDetails
		patientdetails.Medications.ObjectType = "Medications"
		patientdetails.Medications.Patient = *patient
		var defaultConsent entity.Consent
		defaultConsent.Provider = provider
		defaultConsent.StartTime = time.Now().Format("01-02-2006")
		defaultConsent.EndTime = time.Now().AddDate(1, 0, 0).Format("01-02-2006")
		patientdetails.Medications.ProviderConsent = []entity.Consent{}
		patientdetails.Medications.ProviderConsent = append(patientdetails.Medications.ProviderConsent, defaultConsent)

		//==== Create patientAllergies object and marshal to JSON ====
		patientdetails.Allergies.ObjectType = "Allergies"
		patientdetails.Allergies.Patient = *patient
		patientdetails.Allergies.ProviderConsent = []entity.Consent{}
		patientdetails.Allergies.ProviderConsent = append(patientdetails.Allergies.ProviderConsent, defaultConsent)

		//==== Create patientImmunizations object and marshal to JSON ====
		patientdetails.Immunization.ObjectType = "Immunizations"
		patientdetails.Immunization.Patient = *patient
		patientdetails.Immunization.ProviderConsent = []entity.Consent{}
		patientdetails.Immunization.ProviderConsent = append(patientdetails.Immunization.ProviderConsent, defaultConsent)

		//==== Create patientPastMedicalHx object and marshal to JSON ====
		patientdetails.PastMedicalHx.ObjectType = "PastMedicalHx"
		patientdetails.PastMedicalHx.Patient = *patient
		patientdetails.PastMedicalHx.ProviderConsent = []entity.Consent{}
		patientdetails.PastMedicalHx.ProviderConsent = append(patientdetails.PastMedicalHx.ProviderConsent, defaultConsent)

		//==== Create patientFamilyHx object and marshal to JSON ====
		patientdetails.FamilyHx.ObjectType = "FamilyHx"
		patientdetails.FamilyHx.Patient = *patient
		patientdetails.FamilyHx.ProviderConsent = []entity.Consent{}
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

func inTimeSpan(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func (u *User) GetPatientBySSN(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	fmt.Println("In seachpatient by sssn")

	//   0
	// "bob"
	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	role, err := getAttribute(stub, "userrole")

	userId, err := getAttribute(stub, "id")

	if err != nil {
		return shim.Error(err.Error())
	}

	ssn := strings.ToLower(args[0])

	queryString := fmt.Sprintf("{\"selector\":{\"patientssn\":\"%s\"}}", ssn)

	queryResults, err := getQueryResultForQueryString(stub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}

	var tempArray []entity.PatientUnmarshal
	err = json.Unmarshal(queryResults, &tempArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	var key string
	for _, patient := range tempArray {

		key = patient.Key

	}

	fmt.Println("=======Role==============")
	fmt.Println(role)

	if strings.HasPrefix(role, "Patient") {

		if strings.Contains(role, key) {

			patientDetailsBytes, err := stub.GetPrivateData("patientDetails", key)
			if err != nil {
				return shim.Error("Patient not found " + key + "role " + role + "patient details " + string(patientDetailsBytes))
			}
			return shim.Success(patientDetailsBytes)
		} else {
			return shim.Error("unAuthorized role: " + role + "key: " + key)
		}

	} else if strings.HasPrefix(role, "Provider") {

		patientDetailsBytes, err := stub.GetPrivateData("patientDetailsIn2Orgs", key)
		if err != nil {
			return shim.Error("Patient not found " + key + "role " + role + "patient details " + string(patientDetailsBytes))
		}

		var patientDetailsDB entity.PatientDetails

		err = json.Unmarshal(patientDetailsBytes, &patientDetailsDB) //unmarshal it aka JSON.parse()
		if err != nil {
			return shim.Error(err.Error())
		}

		current, _ := time.Parse("01-02-2006", time.Now().Format("01-02-2006"))
		patientDetailsDB = checkMedicationConsent(patientDetailsDB, userId, current)
		patientDetailsDB = checkAllergiesConsent(patientDetailsDB, userId, current)
		patientDetailsDB = checkImmunizationConsent(patientDetailsDB, userId, current)
		patientDetailsDB = checkPastMedicalHxConsent(patientDetailsDB, userId, current)
		patientDetailsDB = checkFamilyHxConsent(patientDetailsDB, userId, current)

		patientDetailsIn2OrgsBytes, err := json.Marshal(&patientDetailsDB)

		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(patientDetailsIn2OrgsBytes)

	} else {
		// When other
		return shim.Error("Only patients, doctors and pharmacies can access medical details")
	}
}

func checkMedicationConsent(patientDetails entity.PatientDetails, userId string, current time.Time) entity.PatientDetails {

	found := false
	for _, consents := range patientDetails.Medications.ProviderConsent {

		if strings.Compare(consents.Provider.ProviderId, userId) == 0 {

			start, _ := time.Parse("01-02-2006", "01-02-2006")

			end, _ := time.Parse("01-02-2006", consents.EndTime)

			if !inTimeSpan(start, end, current) {

				patientDetails.Medications = entity.Medications{}

			} else {
				found = true
			}

		}
	}

	if found == false {
		patientDetails.Medications = entity.Medications{}
	}

	return patientDetails

}

func checkAllergiesConsent(patientDetails entity.PatientDetails, userId string, current time.Time) entity.PatientDetails {

	found := false
	for _, consents := range patientDetails.Allergies.ProviderConsent {

		if strings.Compare(consents.Provider.ProviderId, userId) == 0 {

			start, _ := time.Parse("01-02-2006", "01-02-2006")

			end, _ := time.Parse("01-02-2006", consents.EndTime)

			if !inTimeSpan(start, end, current) {

				patientDetails.Allergies = entity.Allergies{}

			} else {
				found = true
			}

		}
	}

	if found == false {
		patientDetails.Allergies = entity.Allergies{}
	}

	return patientDetails

}

func checkImmunizationConsent(patientDetails entity.PatientDetails, userId string, current time.Time) entity.PatientDetails {

	found := false
	for _, consents := range patientDetails.Immunization.ProviderConsent {

		if strings.Compare(consents.Provider.ProviderId, userId) == 0 {

			start, _ := time.Parse("01-02-2006", "01-02-2006")

			end, _ := time.Parse("01-02-2006", consents.EndTime)

			if !inTimeSpan(start, end, current) {

				patientDetails.Immunization = entity.Immunization{}

			} else {
				found = true
			}

		}
	}

	if found == false {
		patientDetails.Immunization = entity.Immunization{}
	}

	return patientDetails

}

func checkPastMedicalHxConsent(patientDetails entity.PatientDetails, userId string, current time.Time) entity.PatientDetails {

	found := false
	for _, consents := range patientDetails.PastMedicalHx.ProviderConsent {

		if strings.Compare(consents.Provider.ProviderId, userId) == 0 {

			start, _ := time.Parse("01-02-2006", "01-02-2006")

			end, _ := time.Parse("01-02-2006", consents.EndTime)

			if !inTimeSpan(start, end, current) {

				patientDetails.PastMedicalHx = entity.PastMedicalHx{}

			} else {
				found = true
			}

		}
	}

	if found == false {
		patientDetails.PastMedicalHx = entity.PastMedicalHx{}
	}

	return patientDetails

}

func checkFamilyHxConsent(patientDetails entity.PatientDetails, userId string, current time.Time) entity.PatientDetails {

	found := false
	for _, consents := range patientDetails.FamilyHx.ProviderConsent {

		if strings.Compare(consents.Provider.ProviderId, userId) == 0 {

			start, _ := time.Parse("01-02-2006", "01-02-2006")

			end, _ := time.Parse("01-02-2006", consents.EndTime)

			if !inTimeSpan(start, end, current) {

				patientDetails.FamilyHx = entity.FamilyHx{}

			} else {
				found = true
			}

		}
	}

	if found == false {
		patientDetails.FamilyHx = entity.FamilyHx{}
	}

	return patientDetails

}

func (u *User) GetPatientByInformation(stub shim.ChaincodeStubInterface, args []string) pb.Response {

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

