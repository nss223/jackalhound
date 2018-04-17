#!/bin/bash
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
# Exit on first error
set -e

# don't rewrite paths for Windows Git Bash users
export MSYS_NO_PATHCONV=1

starttime=$(date +%s)

# launch network; create channel and join peer to channel

# Now launch the CLI container in order to install, instantiate chaincode
# and prime the ledger with our 10 cars

./start.sh
docker-compose -f ./docker-compose.yml up -d cli

VERSION="1.0"

install_instantiate() {
    docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode install -n "$1"cc -v "$VERSION" -p github.com/$1

    ARGS='{"Args":[""]}'
    docker exec -e "CORE_PEER_LOCALMSPID=Org1MSP" -e "CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp" cli peer chaincode instantiate -o orderer.example.com:7050 -C mychannel -n "$1"cc -v $VERSION -c "$ARGS" -P "OR ('Org1MSP.member','Org2MSP.member')"
}

CC_NAME="reg point map"
for i in $CC_NAME
do
    install_instantiate $i
done
