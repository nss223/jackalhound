#!/bin/sh
#
# Copyright IBM Corp All Rights Reserved
#
# SPDX-License-Identifier: Apache-2.0
#
export PATH=$GOPATH/src/github.com/hyperledger/fabric/build/bin:${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}
CHANNEL_NAME=mychannel

# remember sk, by Jack
sk_old=$(find crypto-config/peerOrganizations/org1.example.com/ca -name "*_sk" -exec basename {} \;)

# remove previous crypto material and config transactions
rm -fr config/*
rm -fr crypto-config/*

# generate crypto material
cryptogen generate --config=./crypto-config.yaml
if [ "$?" -ne 0 ]; then
  echo "Failed to generate crypto material..."
  exit 1
fi

# update sk, by Jack
sk_new=$(find crypto-config/peerOrganizations/org1.example.com/ca -name "*_sk" -exec basename {} \;)
sed -i "s/"$sk_old"/"$sk_new"/" ./docker-compose.yml

# generate genesis block for orderer
configtxgen -profile OneOrgOrdererGenesis -outputBlock ./config/genesis.block
if [ "$?" -ne 0 ]; then
  echo "Failed to generate orderer genesis block..."
  exit 1
fi

gen_crypto_4_channel () {
    CHANNEL_NAME=$1

    # generate channel configuration transaction
    configtxgen -profile OneOrgChannel -outputCreateChannelTx ./config/channel_$CHANNEL_NAME.tx -channelID $CHANNEL_NAME
    if [ "$?" -ne 0 ]; then
        echo "Failed to generate channel configuration transaction..."
        exit 1
    fi

    # generate anchor peer transaction
    configtxgen -profile OneOrgChannel -outputAnchorPeersUpdate ./config/Org1MSPanchors_$CHANNEL_NAME.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
    if [ "$?" -ne 0 ]; then
        echo "Failed to generate anchor peer update for Org1MSP..."
        exit 1
    fi
}

. ./fabric.conf
for chan in $CHANNEL
do
    gen_crypto_4_channel $chan
done
