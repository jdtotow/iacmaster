#!/bin/bash

set -euf -o pipefail
tfenv install $TERRAFORM_VERSION
tfenv use $TERRAFORM_VERSION
az login --service-principal --username $ARM_CLIENT_ID --password $ARM_CLIENT_SECRET --tenant $ARM_TENANT_ID 
az account set --subscription $ARM_SUBSCRIPTION_ID
terraform -chdir=$WORKING_DIR init 
terraform -chdir=$WORKING_DIR plan -out=object.tfplan 
echo "$WORKING_DIR/variables.tfvars"
terraform -chdir=$WORKING_DIR apply object.tfplan 