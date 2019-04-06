package entity

type Patient struct {
	ObjectType       string `json:docType"`
	PatientId        string `json:"patientId"`
	PatientSSN       string `json:"patientssn"`
	PatientUrl       string `json:"patienturl"`
	PatientFirstname string `json:"firstname"` //docType is used to distinguish the various types of objects in state database
	PatientLastname  string `json:"lastname"`  //the fieldtags are needed to keep case from bouncing around
	DOB              string `json:"dob"`
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
type PatientUnmarshal struct 
{
	Key string
	Record Patient `json:"Patient"`
}
type PatientDetailsUnmarshal struct {
	_id string
	_rev string
	Allergies     Allergies     `json:"allergies"`
	FamilyHx      FamilyHx      `json:"familyHx"`
	Immunization  Immunization  `json:"immunization"`
	Medications   Medications   `json:"medications"`
	PastMedicalHx PastMedicalHx `json:"pastMedicalHx"`
}