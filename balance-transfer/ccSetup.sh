#!/bin/bash
starttime=$(date +%s)

LANGUAGE="golang"
cc=$1
ORG1_TOKEN=$(cat org1token.txt)
ORG2_TOKEN=$(cat org2token.txt)
CC_SRC_PATH="github.com/example_cc/go/go_projects/src"

echo "POST Install chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org1.example.com\",\"peer1.org1.example.com\"],
	\"chaincodeName\":\"$cc\",
	\"chaincodePath\":\"$CC_SRC_PATH\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"v0\"
}"
echo
echo

echo "POST Install chaincode on Org2"
echo
curl -s -X POST \
  http://localhost:4000/chaincodes \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"peers\": [\"peer0.org2.example.com\",\"peer1.org2.example.com\"],
	\"chaincodeName\":\"$cc\",
	\"chaincodePath\":\"$CC_SRC_PATH\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"chaincodeVersion\":\"v0\"
}"
echo
echo

echo "POST instantiate chaincode on Org1"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d "{
	\"chaincodeName\":\"$cc\",
	\"chaincodeVersion\":\"v0\",
	\"chaincodeType\": \"$LANGUAGE\",
	\"args\":[\"\"]
}"
echo
echo

while true
do

echo "Select one of the options below:"
echo "1) Exit Chaincode"
echo "2) Register Patient"
echo "3) Query Patient by SSN"
echo "4) Query Patient by Information"
echo "5) Register Provider"
echo "6) Query Provider by providerId"
echo "7) Update Provider Access"

read option

case $option in
"1") break;;
"2") echo "Registering Patient Ibrahim"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/$cc \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer0.org2.example.com"],
	"fcn":"RegisterPatient",
	"args":["12345","321","patient.mtbc.com#123","ibrahim","132","123"]
}'
echo
echo ;;

"3") echo "GET query chaincode on peer0 of Org2 for Patient"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/$cc?peer=peer0.org2.example.com&fcn=GetPatientBySSN&args=%5B%22321%22%5D" \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json"
echo
echo ;;

"4") 
echo "GET query chaincode on peer0 of Org2 for Patient"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/$cc?peer=peer0.org2.example.com&fcn=GetPatientByInformation&args=%5B%22ibrahim%22%2C%22smith%22%2C%2212%2F02%2F2019%22%5D" \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json"
echo
echo ;;

"5") 
echo "Registering Provider Saad"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/$cc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com","peer0.org2.example.com"],
	"fcn":"RegisterProvider",
	"args":["1234","TalkEHR","secure.talkehr.com","saad","buth","gyno"]
}'
echo
echo ;;

"6") echo "GET query chaincode on peer0 of Org2 for Patient"
echo
curl -s -X GET \
  "http://localhost:4000/channels/mychannel/chaincodes/$cc?peer=peer0.org2.example.com&fcn=GetProviderById&args=%5B%22789%22%5D" \
  -H "authorization: Bearer $ORG2_TOKEN" \
  -H "content-type: application/json"
echo
echo ;;


"7") 
echo "Update Provider Access"
echo
curl -s -X POST \
  http://localhost:4000/channels/mychannel/chaincodes/$cc \
  -H "authorization: Bearer $ORG1_TOKEN" \
  -H "content-type: application/json" \
  -d '{
	"peers": ["peer0.org1.example.com"],
	"fcn":"UpdateProviderAccess",
	"args":["{\"allergies\": {\"ObjectType\": \"Allergies\",\"providerconsent\": [ { \"ObjectType\": \"\", \"endtime\": \"03-18-2019\", \"provider\": { \"ObjectType\": \"Provider\", \"firstname\": \"faisal\", \"lastname\": \"faisal\", \"providerId\": \"provider001\", \"providerehr\": \"mtbc\", \"providerehrurl\": \"mtbc\", \"speciality\": \"faisal\" }, \"starttime\": \"03-18-2019\" } ] }, \"familyHx\": { \"ObjectType\": \"FamilyHx\", \"providerconsent\": [ { \"ObjectType\": \"\", \"endtime\": \"03-18-2019\", \"provider\": { \"ObjectType\": \"Provider\", \"firstname\": \"faisal\", \"lastname\": \"faisal\", \"providerId\": \"provider001\", \"providerehr\": \"mtbc\", \"providerehrurl\": \"mtbc\", \"speciality\": \"faisal\" }, \"starttime\": \"03-18-2019\" } ] }, \"immunization\": { \"ObjectType\": \"Immunizations\", \"providerconsent\": [ { \"ObjectType\": \"\", \"endtime\": \"03-18-2019\", \"provider\": { \"ObjectType\": \"Provider\", \"firstname\": \"faisal\", \"lastname\": \"faisal\", \"providerId\": \"provider001\", \"providerehr\": \"mtbc\", \"providerehrurl\": \"mtbc\", \"speciality\": \"faisal\" }, \"starttime\": \"03-18-2019\" } ] }, \"medications\": { \"ObjectType\": \"Medications\", \"providerconsent\": [ { \"ObjectType\": \"\", \"endtime\": \"03-18-2019\", \"provider\": { \"ObjectType\": \"Provider\", \"firstname\": \"faisal\", \"lastname\": \"faisal\", \"providerId\": \"provider001\", \"providerehr\": \"mtbc\", \"providerehrurl\": \"mtbc\", \"speciality\": \"faisal\" }, \"starttime\": \"03-18-2019\" } ] }, \"pastMedicalHx\": { \"ObjectType\": \"PastMedicalHx\",  \"providerconsent\": [ { \"ObjectType\": \"\", \"endtime\": \"03-18-2019\", \"provider\": { \"ObjectType\": \"Provider\", \"firstname\": \"faisal\", \"lastname\": \"faisal\", \"providerId\": \"provider001\", \"providerehr\": \"mtbc\", \"providerehrurl\": \"mtbc\", \"speciality\": \"faisal\" }, \"starttime\": \"03-18-2019\" } ] }}"]
}'
echo
echo ;;

esac





done
# echo "GET query Block by blockNumber"
# echo
# BLOCK_INFO=$(curl -s -X GET \
#   "http://localhost:4000/channels/mychannel/blocks/1?peer=peer0.org1.example.com" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json")
# echo $BLOCK_INFO
# # Assign previvious block hash to HASH
# HASH=$(echo $BLOCK_INFO | jq -r ".header.previous_hash")
# echo

# echo "GET query Transaction by TransactionID"
# echo
# curl -s -X GET http://localhost:4000/channels/mychannel/transactions/$TRX_ID?peer=peer0.org1.example.com \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo


# echo "GET query Block by Hash - Hash is $HASH"
# echo
# curl -s -X GET \
#   "http://localhost:4000/channels/mychannel/blocks?hash=$HASH&peer=peer0.org1.example.com" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "cache-control: no-cache" \
#   -H "content-type: application/json" \
#   -H "x-access-token: $ORG1_TOKEN"
# echo
# echo

# echo "GET query ChainInfo"
# echo
# curl -s -X GET \
#   "http://localhost:4000/channels/mychannel?peer=peer0.org1.example.com" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

# echo "GET query Installed chaincodes"
# echo
# curl -s -X GET \
#   "http://localhost:4000/chaincodes?peer=peer0.org1.example.com" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

# echo "GET query Instantiated chaincodes"
# echo
# curl -s -X GET \
#   "http://localhost:4000/channels/mychannel/chaincodes?peer=peer0.org1.example.com" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo

# echo "GET query Channels"
# echo
# curl -s -X GET \
#   "http://localhost:4000/channels?peer=peer0.org1.example.com" \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json"
# echo
# echo


echo "Total execution time : $(($(date +%s)-starttime)) secs ..."
