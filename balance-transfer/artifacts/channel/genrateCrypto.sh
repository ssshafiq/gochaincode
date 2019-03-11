#!/bin/sh

set -e
# Grab the current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"


echo "Start genrating crypto "

rm -rf "${DIR}"/crypto-config/*
chmod 777 -R  "${DIR}"/crypto-config/

echo "cryptogen"
cryptogen generate --config="${DIR}"/cryptogen.yaml

echo "configtxgen"
configtxgen -profile TwoOrgsOrdererGenesis -outputBlock "${DIR}"/genesis.block

export CHANNEL_NAME=mychannel  && configtxgen -profile TwoOrgsChannel -outputCreateChannelTx "${DIR}"/mychannel.tx -channelID $CHANNEL_NAME

configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate "${DIR}"/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate "${DIR}"/Org2MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org2MSP


echo "Removing old  network-config.yaml"

rm -f ../network-config.yaml

echo "Setting network-config.yaml"

cp ../network-config-template.yaml ../network-config.yaml

AdminOrg1Key="$(ls ${DIR}/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/)"
echo "${AdminOrg1Key}"
sed -i "s/AdminOrg1Key/${AdminOrg1Key}/g"  ../network-config.yaml

AdminOrg2Key="$(ls ${DIR}/crypto-config/peerOrganizations/org2.example.com/users/Admin@org2.example.com/msp/keystore/)"
sed -i "s/AdminOrg2Key/${AdminOrg2Key}/g"  ../network-config.yaml

echo "${AdminOrg2Key}"


echo "setting CA"
CAorg1Key="$(ls ${DIR}/crypto-config/peerOrganizations/org1.example.com/ca | grep sk)"
CAorg2Key="$(ls ${DIR}/crypto-config/peerOrganizations/org2.example.com/ca | grep sk)"

rm -f  ../docker-compose.yaml

cp ../docker-compose-template.yaml ../docker-compose.yaml

sed -i "s/CAorg1Key/${CAorg1Key}/g"  ../docker-compose.yaml
sed -i "s/CAorg2Key/${CAorg2Key}/g"  ../docker-compose.yaml
 echo "End of script"
