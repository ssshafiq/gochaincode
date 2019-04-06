#!/bin/bash
#
# Copyright IBM Corp. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

jq --version > /dev/null 2>&1
if [ $? -ne 0 ]; then
	echo "Please Install 'jq' https://stedolan.github.io/jq/ to execute this script"
	echo
	exit 1
fi

starttime=$(date +%s)

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  ./testAPIs.sh -l golang|node"
  echo "    -l <language> - chaincode language (defaults to \"golang\")"
}
# Language defaults to "golang"
LANGUAGE="golang"

# Parse commandline args
while getopts "h?l:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    l)  LANGUAGE=$OPTARG
    ;;
  esac
done



##set chaincode path
function setChaincodePath(){
	LANGUAGE=`echo "$LANGUAGE" | tr '[:upper:]' '[:lower:]'`
	case "$LANGUAGE" in
		"golang")
		CC_SRC_PATH="github.com/example_cc/go"
		;;
		"node")
		CC_SRC_PATH="$PWD/artifacts/src/github.com/example_cc/node"
		;;
		*) printf "\n ------ Language $LANGUAGE is not supported yet ------\n"$
		exit 1
	esac
}

setChaincodePath

#create ORG1_TOKENPatient.txt if not exist

if [[ ! -e ORG1_TOKENPatient.txt ]]; then
    echo "creating ORG1_TOKENPatient.txt"
    touch ORG1_TOKENPatient.txt
fi

#create ORG1_TOKEN.txt if not exist

if [[ ! -e ORG1_TOKEN.txt ]]; then
    echo "creating ORG1_TOKEN.txt"
    touch ORG1_TOKEN.txt
fi

#create ORG2_TOKEN.txt if not exist

if [[ ! -e ORG2_TOKEN.txt ]]; then
    echo "creating ORG2_TOKEN.txt"
    touch ORG2_TOKEN.txt
fi


echo "POST request Enroll patient on Org1  ..."
echo
ORG1_TOKENPatient=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=PatientOrgMTBC&orgName=org-mtbc&userrole=Patient&id=pat001&mspRole=client')
echo $ORG1_TOKENPatient
ORG1_TOKENPatient=$(echo $ORG1_TOKENPatient | jq ".token" | sed "s/\"//g")
echo
echo $ORG1_TOKENPatient >ORG1_TOKENPatient.txt


echo "POST request Enroll on Org1  ..."
echo
ORG1_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=ProviderOrgMTBC&orgName=org-mtbc&userrole=Provider&id=pro001&mspRole=client')
echo $ORG1_TOKEN
ORG1_TOKEN=$(echo $ORG1_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo $ORG1_TOKEN >org1token.txt
echo "ORG1 token is $ORG1_TOKEN"
echo
echo "POST request Enroll on Org2 ..."
echo
ORG2_TOKEN=$(curl -s -X POST \
  http://localhost:4000/users \
  -H "content-type: application/x-www-form-urlencoded" \
  -d 'username=ProviderOrgUNI&orgName=org-uni&userrole=Provider&id=pro002&mspRole=client')
echo $ORG2_TOKEN

ORG2_TOKEN=$(echo $ORG2_TOKEN | jq ".token" | sed "s/\"//g")
echo
echo $ORG2_TOKEN >org2token.txt
echo "ORG2 token is $ORG2_TOKEN"

# echo
# echo
# echo "POST request Create channel  ..."
# echo
# echo "*************************************"
# echo "ORG1 TOKEN IS ********* $(cat org1token.txt)"
# echo "ORG2 TOKEN IS ********* $(cat org2token.txt)"
# echo "*************************************"
# curl -s -X POST \
#   http://localhost:4000/channels \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
# 	"channelName":"mychannel",
# 	"channelConfigPath":"../artifacts/channel/mychannel.tx"
# }'
# echo
# echo
# sleep 5
# echo "POST request Join channel on Org1"
# echo
# curl -s -X POST \
#   http://localhost:4000/channels/mychannel/peers \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
# 	"peers": ["peer0.org1.example.com","peer1.org1.example.com"]
# }'
# echo
# echo

# echo "POST request Join channel on Org2"
# echo
# curl -s -X POST \
#   http://localhost:4000/channels/mychannel/peers \
#   -H "authorization: Bearer $ORG2_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
# 	"peers": ["peer0.org2.example.com","peer1.org2.example.com"]
# }'
# echo
# echo

# echo "POST request Update anchor peers on Org1"
# echo
# curl -s -X POST \
#   http://localhost:4000/channels/mychannel/anchorpeers \
#   -H "authorization: Bearer $ORG1_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
# 	"configUpdatePath":"../artifacts/channel/Org1MSPanchors.tx"
# }'
# echo
# echo

# echo "POST request Update anchor peers on Org2"
# echo
# curl -s -X POST \
#   http://localhost:4000/channels/mychannel/anchorpeers \
#   -H "authorization: Bearer $ORG2_TOKEN" \
#   -H "content-type: application/json" \
#   -d '{
# 	"configUpdatePath":"../artifacts/channel/Org2MSPanchors.tx"
# }'
# echo
# echo
