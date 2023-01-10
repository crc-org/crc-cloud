terraform apply \
-var "access_key=xxxxx" \
-var "secret_key=xxxxx" \
-var "region=us-west-2" \
-var "ami=ami-0569ce8a44f2351be" \
-var "instance_type=c6in.2xlarge" \
-state=terraform.tfstate \
-auto-approve \
-chdir=./terraform/aws
