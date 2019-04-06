
#!/bin/bash

echo "Starting ....."

set -e
# Grab the current directory
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

echo "Removing fabric-client-kv-org-mtbc"

rm -rf "${DIR}"/fabric-client-kv-org-mtbc/*

echo "Removing fabric-client-kv-org-uni"

rm -rf "${DIR}"/fabric-client-kv-org-uni/*

echo "Removing network-config"

OrgMTBCAdminKey="$(ls /tmp/fabric/crypto-config/peerOrganizations/org-mtbc/users/Admin@org-mtbc/msp/keystore/)"

OrgUNIAdminKey="$(ls /tmp/fabric/crypto-config/peerOrganizations/org-uni/users/Admin@org-uni/msp/keystore/)"

echo $OrgMTBCAdminKey
echo $OrgUNIAdminKey


rm -rf "${DIR}"/artifacts/network-config.yaml

cp  "${DIR}"/artifacts/network-config-template.yaml "${DIR}"/artifacts/network-config.yaml

sed -i "s/OrgMTBCAdminKey/${OrgMTBCAdminKey}/g" "${DIR}"/artifacts/network-config.yaml

sed -i "s/OrgUNIAdminKey/${OrgUNIAdminKey}/g" "${DIR}"/artifacts/network-config.yaml


echo "Successfully completed ......"

