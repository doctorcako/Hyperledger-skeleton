#!/bin/bash

APP_PATH=$(pwd)
TEMPLATES=$APP_PATH/templates
CONFIG=$APP_PATH/config/$HLF_ENV

ORDERER_TYPE="kafka"

CA_ORGS="org1 org2 org3 orderer"
ORGS="org1 org2 org3"
ORDERER_ORGS="orderer0 orderer1 orderer2"
PEERS_PER_ORG=1

rcaport=30100
icaport=30200
ordererport=30300
peerport=30400
couchdbport=30500
outcouchdb=5984


ENABLE_INTERMEDIATE_CA=false