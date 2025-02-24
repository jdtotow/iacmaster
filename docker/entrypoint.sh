#!/bin/bash
tfenv install $TERRAFORM_VERSION
tfenv use $TERRAFORM_VERSION
terraform -chdir=$WORKING_DIR init 
terraform -chdir=$WORKING_DIR plan 
terraform -chdir=$WORKING_DIR apply 