#!/bin/bash
AMI_ID="ami-0569ce8a44f2351be"
PUBKEY="id_rsa"
RUN_TIMESTAMP=`date +%s`
BASE_WORKDIR="workdir"
WORKDIR="$BASE_WORKDIR/$RUN_TIMESTAMP"
RANDOM_SUFFIX=`echo $RANDOM | md5sum | head -c 8`
RANDOM_SUFFIX_FILE="$WORKDIR/suffix"
prepare_workdir() {
    mkdir $WORKDIR
    echo $RANDOM_SUFFIX > $RANDOM_SUFFIX_FILE
    rm $BASE_WORKDIR/latest | true /dev/null 2>&1
    ln -s $(pwd)/$WORKDIR $(pwd)/$BASE_WORKDIR/latest
}


prepare_swap_keys() {
    ssh-keygen -m PEM -f $WORKDIR/$PUBKEY -q -N ""
    cp templates/swap_keys.sh $WORKDIR
    sed "s#_PUBKEY_#$(cat $WORKDIR/$PUBKEY.pub)#" templates/swap_keys.sh > $WORKDIR/swap_keys.sh
    chmod +x $WORKDIR/swap_keys.sh
}

create_instances(){
    aws ec2 create-key-pair --key-name openspot-ng-$RANDOM_SUFFIX
    aws ec2 import-key-pair --key-name openspot-ng-$RANDOM_SUFFIX --public-key-material file:///path/to/my/key.pem
    aws ec2 create-security-group --group-name openspot-ng-$RANDOM_SUFFIX --description "openspot-ng security group run timestamp: $RUN_TIMESTAMP"
    aws ec2 authorize-security-group-ingress --group-name openspot-ng-$RANDOM_SUFFIX --protocol tcp --port 22 --cidr 0.0.0.0/0
    aws ec2 authorize-security-group-ingress --group-name openspot-ng-$RANDOM_SUFFIX --protocol tcp --port 6443 --cidr 0.0.0.0/0
    aws ec2 authorize-security-group-ingress --group-name openspot-ng-$RANDOM_SUFFIX --protocol tcp --port 443 --cidr 0.0.0.0/0
    aws ec2 run-instances --image-id ami-xxxxxxxx --instance-type c6i.2xlarge --user-data file://path/to/user-data-script.txt --security-group-ids sg-xxxxxxxx --key-name my-key-pair
}

prepare_workdir
prepare_swap_keys

#
#


