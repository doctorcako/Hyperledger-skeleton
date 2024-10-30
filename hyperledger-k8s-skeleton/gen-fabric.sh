source ./scripts/.env.sh
source ./scripts/hyperledger-config-generator.sh
source ./scripts/utils.sh

function main(){
    if [ ! -z $NAMESPACE ]; then
        DOMAIN=$NAMESPACE
    fi
    generateNamespaces
    genVolumes
    generateRCA
    if [ $ENABLE_INTERMEDIATE_CA = true ]; then
        generateICA
    fi
    generateOrderers
    generatePeers
    generateCli
    
}

main