#!/bin/bash
function generateNamespaces() {
    for ORG in $ORGS; do
        if [ ! -d $CONFIG/config.state ]; then
            mkdir -p $CONFIG/config.state
        fi
        if [ ! -d $CONFIG/$ORG ]; then
            mkdir -p $CONFIG/$ORG
        fi
        sed -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%ORG%/${ORG}/g" ${TEMPLATES}/00-fabric-namespace.yaml > ${CONFIG}/$ORG/fabric-namespace-$ORG-$HLF_ENV.yaml
    done
}

function genVolumes {
    for ORG in $ORGS; do
        if [ ! -d $CONFIG/$ORG/pvc ];then
            mkdir -p $CONFIG/$ORG/pvc
        fi
        sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" ${TEMPLATES}/01-fabric-pvc-rca-scripts.yaml > ${CONFIG}/$ORG/pvc/$ORG-fabric-pvc-rca-scripts.yaml
        sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" ${TEMPLATES}/01-fabric-pvc-rca-data.yaml > ${CONFIG}/$ORG/pvc/$ORG-fabric-pvc-rca-data.yaml
        sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" ${TEMPLATES}/01-fabric-pvc-rca.yaml > ${CONFIG}/$ORG/pvc/$ORG-fabric-pvc-rca.yaml
        sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" ${TEMPLATES}/01-fabric-pvc-ica.yaml > ${CONFIG}/$ORG/pvc/$ORG-fabric-pvc-ica.yaml
    done

    if [ ! -d $CONFIG/orderer/pvc ];then
        mkdir -p $CONFIG/orderer/pvc
    fi
    for ORG in $ORDERER_ORGS; do
        sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%EFSSERVER%/${EFSSERVER}/g" ${TEMPLATES}/01-fabric-pvc-orderer.yaml > ${CONFIG}/orderer/pvc/$ORG-fabric-pvc-orderer.yaml
    done
}

function generateRCA() {
    local newOrg=$1
    local isNewOrg=false

    if [ -f $CONFIG/config.state/rca-ports.sh ]; then
        source $CONFIG/config.state/rca-ports.sh
        for portValue in "${RCA_PORTS_IN_USE[@]}"
        do
            if [[ portValue -gt rcaport ]]; then
                rcaport=$((portValue+0))
            fi
        done
        rcaport=$((rcaport+1))
    else
        RCA_PORTS_IN_USE=()
    fi

    orgsArray=($CA_ORGS)
    for key in ${!orgsArray[@]}; do
        local ORG="${orgsArray[${key}]}"

        if [ ! -d $CONFIG/$ORG/rca/ ];then
            mkdir -p $CONFIG/$ORG/rca/
        fi

        if [ -z "${RCA_PORTS_IN_USE[${key}]}" ]; then
            RCA_PORTS_IN_USE+=( $rcaport )
        elif [ "${RCA_PORTS_IN_USE[${key}]}" -eq "$rcaport" ] ; then
            RCA_PORTS_IN_USE[$key]=$rcaport
        fi
        
        if [ -z $newOrg ]; then
            rcaport=$((rcaport+1))
            sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%FABRICORGS%/${ORGS}/g" ${TEMPLATES}/02-fabric-deployment-rca.yaml > ${CONFIG}/$ORG/rca/$ORG-fabric-deployment-rca.yaml
            sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%PORT%/${rcaport}/g" ${TEMPLATES}/03-fabric-service-rca.yaml > ${CONFIG}/$ORG/rca/$ORG-fabric-service-rca.yaml
        elif [[ "$newOrg" != *"$ORG"* ]]; then 
            isNewOrg=true
        else
            isNewOrg=false
        fi
    done

    if [ $isNewOrg = "true" ]; then
        sed -e "s/%ORG%/${newOrg}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%FABRICORGS%/${ORGS}/g" -e "s/%DIR%/${DIR//\//\\/}/g" ${TEMPLATES}/02-fabric-deployment-rca.yaml > ${CONFIG}/$newOrg/$newOrg-fabric-deployment-rca-$newOrg.yaml
        sed -e "s/%ORG%/${newOrg}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%PORT%/${rcaport}/g" ${TEMPLATES}/03-fabric-service-rca.yaml > ${CONFIG}/fabric-service-rca-$newOrg.yaml
        RCA_PORTS_IN_USE+=( $rcaport )
    fi 

    declare -p RCA_PORTS_IN_USE > $CONFIG/config.state/rca-ports.sh
}

function generateICA() {
    local newOrg=$1
    local isNewOrg=false

    if [ -f $CONFIG/config.state/ica-ports.sh ]; then
        source $CONFIG/config.state/ica-ports.sh
        for portValue in "${ICA_PORTS_IN_USE[@]}"
        do
            if [[ portValue -gt rcaport ]]; then
                icaport=$((portValue+0))
            fi
        done
        icaport=$((icaport+1))
    else
        ICA_PORTS_IN_USE=()
    fi

    orgsArray=($CA_ORGS)
    for key in ${!orgsArray[@]}; do
        local ORG="${orgsArray[${key}]}"
        getDomain $ORG

        if [ ! -d $CONFIG/ica ];then
            mkdir -p $CONFIG/ica
        fi

        if [ -z "${ICA_PORTS_IN_USE[${key}]}" ]; then
            ICA_PORTS_IN_USE+=( $icaport )
        elif [ "${ICA_PORTS_IN_USE[${key}]}" -eq "$icaport" ] ; then
            ICA_PORTS_IN_USE[$key]=$icaport
        fi
        
        if [ -z $newOrg ]; then
            icaport=$((icaport+1))
            sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%FABRICORGS%/${ORGS}/g" -e "s/%DIR%/${DIR//\//\\/}/g" ${TEMPLATES}/02-fabric-deployment-ica.yaml > ${CONFIG}/ica/fabric-deployment-ica-$ORG.yaml
            sed -e "s/%ORG%/${ORG}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%PORT%/${icaport}/g" ${TEMPLATES}/03-fabric-service-ica.yaml > ${CONFIG}/ica/fabric-service-ica-$ORG.yaml
        elif [[ "$newOrg" != *"$ORG"* ]]; then 
            isNewOrg=true
        else
            isNewOrg=false
        fi
    done

    if [ $isNewOrg = "true" ]; then
        getDomain $newOrg
        sed -e "s/%ORG%/${newOrg}/g" -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%FABRICORGS%/${ORGS}/g" -e "s/%DIR%/${DIR//\//\\/}/g" ${TEMPLATES}/03-fabric-deployment-rca.yaml > ${CONFIG}/ica/fabric-deployment-rca-$newOrg.yaml
        sed -e "s/%ORG%/${newOrg}/g" -e "s/%HLF_ENV%/${DOHLF_ENVAIN}/g" -e "s/%PORT%/${icaport}/g" ${TEMPLATES}/03-fabric-service-rca.yaml > ${CONFIG}/ica/fabric-service-rca-$newOrg.yaml
        ICA_PORTS_IN_USE+=( $icaport )
    fi 

    declare -p ICA_PORTS_IN_USE > $CONFIG/config.state/ica-ports.sh
}

function generateOrderers() {
    local newOrdererNum=$1
    local isNewOrderer=false 

    if [ -f $CONFIG/config.state/orderer-ports.sh ]; then
        source $CONFIG/config.state/orderer-ports.sh
        for portValue in "${ORDERER_PORTS_IN_USE[@]}"
        do
          if [[ portValue -gt ordererport ]]; then
            ordererport=$((portValue+0))
            ordererauxport=$((ordererauxport+1000))
          fi
        done
        ordererport=$((ordererport+1))
    else
        ORDERER_PORTS_IN_USE=()
    fi
    
    if [ ! -d $CONFIG/orderer/orderers ];then
        mkdir -p $CONFIG/orderer/orderers
    fi 

    orderersArray=($ORDERER_ORGS)
    for key in ${!orderersArray[@]}; do
        local ORDERER="${orderersArray[${key}]}"
        
        if [ -z "${ORDERER_PORTS_IN_USE[${key}]}" ]; then
            ORDERER_PORTS_IN_USE+=( $ordererport )
        elif [ "${ORDERER_PORTS_IN_USE[${key}]}" -eq "$ordererport" ] ; then
            ORDERER_PORTS_IN_USE[$key]=$ordererport
        fi


        if [ -z $newOrdererNum ]; then
            ordererport=$((ordererport+1))
            sed -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${key}/g" -e "s/%ORG%/${ORDERER}/g"  ${TEMPLATES}/04-fabric-deployment-orderer.yaml > ${CONFIG}/orderer/orderers/$ORDERER-fabric-deployment.yaml
            sed -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${key}/g" -e "s/%ORG%/${ORDERER}/g" -e "s/%PORT%/${ordererport}/g"  ${TEMPLATES}/05-fabric-service-orderer.yaml > ${CONFIG}/orderer/orderers/$ORDERER-fabric-service.yaml
        elif [ "$newOrdererNum" != $key ]; then 
            isNewOrderer=true
        else
            isNewOrderer=false
        fi
    done

    if [ $isNewOrderer = "true" ]; then
        sed -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${newOrdererNum}/g" -e "s/%ORG%/orderer${newOrdererNum}/g"  ${TEMPLATES}/04-fabric-deployment-orderer.yaml > ${CONFIG}/orderer/orderers/orderer$newOrdererNum-fabric-deployment.yaml
        sed -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${newOrdererNum}/g" -e "s/%ORG%/orderer${newOrdererNum}/g" -e "s/%PORT%/${ordererport}/g"  ${TEMPLATES}/04-fabric-service-orderer.yaml > ${CONFIG}/orderer/orderers/orderer$newOrdererNum-fabric-service.yaml
        ORDERER_PORTS_IN_USE+=( $ordererport )
    fi
    declare -p ORDERER_PORTS_IN_USE > $CONFIG/config.state/orderer-ports.sh
}

function generatePeers(){
    local OrgPeer=$1
    local peerId=$2
    local newOrgPeer=true
    
    orgsArray=($ORGS)
    CURRENT_ORGS=(${orgsArray[@]})

    for key in ${!orgsArray[@]}; do
        local ORG="${orgsArray[${key}]}"

        if [ ! -d $CONFIG/$ORG/peers ];then
            mkdir -p $CONFIG/$ORG/peers
        fi 
        
        if [[ "$OrgPeer" == *"$ORG"* ]]; then
            OrgPeer=$ORG
        fi

        if [ -d $CONFIG/config.state/peer-ports/peers-$ORG ]; then
            if [ -f ./$CONFIG/config.state/peer-ports/peers-$ORG/peer-ports.sh ]; then
                source $CONFIG/config.state/peer-ports/peers-$ORG/peer-ports.sh
                for portValue in "${PEER_PORTS_IN_USE[@]}"
                do
                if [[ portValue -gt peerport ]]; then
                    peerport=$((portValue+0))
                fi
                done

                if [ ! -z $OrgPeer ]; then
                    peerport=$((1*$key+peerport+1000+$peerId+5))
                else
                    peerport=$((peerport+1000))
                fi
            else
                PEER_PORTS_IN_USE=()
            fi

            if [ -f $CONFIG/config.state/peer-ports/peers-$ORG/couchdb-ports.sh ]; then
                source $CONFIG/config.state/peer-ports/peers-$ORG/couchdb-ports.sh
                for portValue in "${COUCHDB_PORTS_IN_USE[@]}"
                do
                if [[ portValue -gt couchdbport ]]; then
                    couchdbport=$((portValue+0))
                fi
                done
                couchdbport=$((couchdbport+1))
            else
                COUCHDB_PORTS_IN_USE=()
            fi
        else
            mkdir -p $CONFIG/config.state/peer-ports/peers-$ORG
        fi


        AUX_PEERS_PER_ORG=$PEERS_PER_ORG
        for (( PEER_NUM=0; PEER_NUM < $AUX_PEERS_PER_ORG; PEER_NUM++ )); do
            
            if [ -z "${PEER_PORTS_IN_USE[${PEER_NUM}]}" ]; then
                PEER_PORTS_IN_USE+=( $peerport )
            elif [ "${PEER_PORTS_IN_USE[${PEER_NUM}]}" -eq "$peerport" ] ; then
                PEER_PORTS_IN_USE[$key]=$peerport
            fi

            if [ -z "${COUCHDB_PORTS_IN_USE[${PEER_NUM}]}" ]; then
                COUCHDB_PORTS_IN_USE+=( $couchdbport )
            elif [ "${COUCHDB_PORTS_IN_USE[${PEER_NUM}]}" -eq "$couchdbport" ] ; then
                COUCHDB_PORTS_IN_USE[$key]=$couchdbport
            fi
           
            if [ -z $OrgPeer ] || [ ! -z $OrgPeer -a $peerId -eq $PEER_NUM ] || [ $OrgPeer != $ORG ] ; then
                peerport=$((peerport+5))
                peerport2=$((peerport+1))
                MSP="$(tr '[:lower:]' '[:upper:]' <<< ${ORG:0:1})${ORG:1}"
                sed -e "s/%ORG%/${ORG}/g" -e "s/%MSP%/${MSP}/g"  -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${PEER_NUM}/g" ${TEMPLATES}/06-fabric-deployment-peer.yaml > ${CONFIG}/$ORG/peers/peer$PEER_NUM-$ORG-fabric-deployment.yaml
                sed -e "s/%ORG%/${ORG}/g"  -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${PEER_NUM}/g" -e "s/%PORT%/${couchdbport}/g" -e "s/%OUT_PORT%/${outcouchdb}/g" ${TEMPLATES}/07-fabric-service-couchdb.yaml > ${CONFIG}/$ORG/peers/peer$PEER_NUM-$ORG-couchdb-fabric-service.yaml                    
                sed -e "s/%ORG%/${ORG}/g"  -e "s/%HLF_ENV%/${HLF_ENV}/g" -e "s/%NUM%/${PEER_NUM}/g" -e "s/%PORT_1%/${peerport}/g" -e "s/%PORT_2%/${peerport2}/g" ${TEMPLATES}/08-fabric-service-peer.yaml > ${CONFIG}/$ORG/peers/peer$PEER_NUM-$ORG-fabric-service.yaml
                couchdbport=$((couchdbport+1))
                outcouchdb=$((outcouchdb + 1000))
            else
                AUX_PEERS_PER_ORG=$(($AUX_PEERS_PER_ORG+1))
                outcouchdb=$((outcouchdb + 1001))
            fi
        done
        declare -p PEER_PORTS_IN_USE > $CONFIG/config.state/peer-ports/peers-$ORG/peer-ports.sh
        declare -p COUCHDB_PORTS_IN_USE > $CONFIG/config.state/peer-ports/peers-$ORG/couchdb-ports.sh

    done
}

function generateCli() {
    orgsArray=($ORGS)
    for key in ${!orgsArray[@]}; do
        local ORG="${orgsArray[${key}]}"
        sed -e "s/%ORG%/${ORG}/g"  -e "s/%HLF_ENV%/${HLF_ENV}/g" ${TEMPLATES}/09-fabric-cli.yaml > ${CONFIG}/$ORG/fabric-deployment-cli.yaml
    done
}