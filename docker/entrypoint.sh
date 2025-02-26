#!/bin/bash

set -euf -o pipefail
tfenv install $TERRAFORM_VERSION
tfenv use $TERRAFORM_VERSION
printenv
echo "Executing terraform init ..."
terraform -chdir=$WORKING_DIR init 
echo "Executing terraform plan ..."
terraform -chdir=$WORKING_DIR plan -var-file=$WORKING_DIR/variables.tfvars -out=object.tfplan 
echo "Executing terraform applying the infrastructure ..."
terraform -chdir=$WORKING_DIR apply -var-file=$WORKING_DIR/variables.tfvars object.tfplan