#!/bin/bash

set -euf -o pipefail
tfenv install $TERRAFORM_VERSION
tfenv use $TERRAFORM_VERSION
terraform -chdir=$WORKING_DIR init 
terraform -chdir=$WORKING_DIR plan -out=object.tfplan 
echo "$WORKING_DIR/variables.tfvars"
terraform -chdir=$WORKING_DIR apply -var-file=$WORKING_DIR/variables.tfvars