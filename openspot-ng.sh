#!/bin/bash
source ./common.sh
export AWS_PROFILE="redhat"

#COMMANDS
CURL=`which curl`
JQ=`which jq`
MD5SUM=`which md5sum`
AWS=`which aws`
HEAD=`which head`
SSHKEYGEN=`which ssh-keygen`
SED=`which sed`
NC=`which nc`
SSH=`which ssh`
SCP=`which scp`


#VARIABLES
AMI_ID="ami-0569ce8a44f2351be"
INSTANCE_TYPE="c6i.2xlarge"
PUBKEY="id_rsa"
RUN_TIMESTAMP=`date +%s`
BASE_WORKDIR="workdir"
WORKDIR="$BASE_WORKDIR/$RUN_TIMESTAMP"
INSTANCE_DESCRIPTION="instance_description.json"
RANDOM_SUFFIX=`echo $RANDOM | $MD5SUM | $HEAD -c 8`
RANDOM_SUFFIX_FILE="$WORKDIR/suffix"


prepare_workdir() {
    mkdir $WORKDIR
    echo $RANDOM_SUFFIX > $RANDOM_SUFFIX_FILE
    rm $BASE_WORKDIR/latest | true /dev/null 2>&1
    ln -s $(pwd)/$WORKDIR $(pwd)/$BASE_WORKDIR/latest
}


prepare_swap_keys() {
    $SSHKEYGEN -m PEM -f $WORKDIR/$PUBKEY -q -N ""
    cp templates/swap_keys.sh $WORKDIR
    $SED "s#_PUBKEY_#$(cat $WORKDIR/$PUBKEY.pub)#" templates/swap_keys.sh > $WORKDIR/swap_keys.sh
    chmod +x $WORKDIR/swap_keys.sh
}

create_instances(){
    RESOURCES_NAME="openspot-ng-$RANDOM_SUFFIX"
    GROUPID=`aws ec2 create-security-group --group-name $RESOURCES_NAME --description "openspot-ng security group run timestamp: $RUN_TIMESTAMP" --no-paginate | $JQ -r '.GroupId'`
    #KEYPAIR (Created just because mandatory, will be swapped manually fore core user later on)
    $AWS ec2 create-key-pair --key-name $RESOURCES_NAME --no-paginate
    #SG SETUP
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 22 --cidr 0.0.0.0/0 --no-paginate
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 6443 --cidr 0.0.0.0/0 --no-paginate
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 443 --cidr 0.0.0.0/0 --no-paginate
    #CREATE INSTANCE
    $AWS ec2 run-instances --no-paginate --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$RESOURCES_NAME}]" --image-id $AMI_ID --instance-type $INSTANCE_TYPE --user-data file://./$WORKDIR/swap_keys.sh --security-group-ids $GROUPID --key-name $RESOURCES_NAME > $WORKDIR/$INSTANCE_DESCRIPTION 
}

swap_ssh_key() {
    $SCP -o StrictHostKeychecking=no -i id_ecdsa_crc $WORKDIR/id_rsa.pub core@$EIP:.
    $SSH -o StrictHostKeychecking=no -i id_ecdsa_crc core@$EIP "cat /home/core/id_rsa.pub > /home/core/.ssh/authorized_keys"
}


#prepare_workdir
#prepare_swap_keys
#create_instances

#INSTANCE_ID=`get_instance_id $WORKDIR/$INSTANCE_DESCRIPTION`
IIP=`get_instance_private_ip $WORKDIR/$INSTANCE_DESCRIPTION`
#EIP=`get_instance_public_ip $INSTANCE_ID`

#wait_instance_readiness $EIP

#swap_ssh_key
echo $IIP


