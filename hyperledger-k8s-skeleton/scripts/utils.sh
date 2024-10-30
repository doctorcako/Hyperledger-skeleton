function getDomain {
   if [ $# -ne 1 ]; then
      echo "Usage: getDomain <ORG>"
      exit 1
   fi
   orgsarr=($ORGS)
   for i in "${!orgsarr[@]}"; do
      if [[ "${orgsarr[$i]}" = "${1}" ]]; then
           DOMAIN=${orgsarr[$i]}
           return
      fi
   done
}

function checkNamespaceDeclared(){
   if [ -z $NAMESPACE ]; then
      echo "Declare HLF_ENV: export HLF_ENV=dev" 
      exit 0
   fi
}