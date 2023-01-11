#!/bin/bash

_bash_ec2_init() {
    # TO AVOID AWS CLI PAGINATION
    export AWS_PAGER=""

    # CHECK CONTAINER VARS 
    if [ $CONTAINER ] 
    then
        [[ -z $AWS_ACCESS_KEY_ID ]] &&\
        stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
        [[ -z $AWS_SECRET_ACCESS_KEY ]] &&\
        stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
        [[ -z $AWS_DEFAULT_REGION ]] &&\
        stop_if_failed 1 "AWS_ACCESS_KEY_ID must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
        [[ -z $TARGET_INFRASTRUCTURE ]] &&\
        stop_if_failed 1 "TARGET_INFRASTRUCTURE must be set, please refer to AWS CLI documentation https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html"
    fi

    # CONSTS
    INSTANCE_DESCRIPTION="instance_description.json"

    # CAN BE OVERRIDE BY CONTAINER ENV VARS 
    [ -z $INSTANCE_TYPE ] && INSTANCE_TYPE="c6in.2xlarge"
    [ -z $AMI_ID ] && AMI_ID="ami-0569ce8a44f2351be"

    # DEPENDENCIES CHECKS
    AWS=`which aws 2>/dev/null`
    [[ $? != 0 ]] && stop_if_failed 1 "[DEPENDENCY MISSING]: aws cli, please install it and try again"
}


# PRIVATE AWS METHODS
_bash_get_ec2_instance_id() {
    $JQ -r '.Instances[0].InstanceId' $1
}

_bash_get_ec2_instance_public_ip(){
    INSTANCE_IP=""
    while [ -z $INSTANCE_IP ]
    do
        INSTANCE_IP=`$AWS ec2 describe-instances --instance-ids $1 --query 'Reservations[*].Instances[*].PublicIpAddress' --output text`
    done
    echo "$INSTANCE_IP"
}

_bash_get_ec2_instance_private_ip(){
    $JQ -r '.Instances[0].PrivateIpAddress' $1
}

_bash_create_ec2_resources() {
    pr_info "creating EC2 resources"
    RESOURCES_NAME="crc-cloud-$RANDOM_SUFFIX"
    GROUPID=`aws ec2 create-security-group --group-name $RESOURCES_NAME --description "crc-cloud security group run timestamp: $RUN_TIMESTAMP" --no-paginate | $JQ -r '.GroupId'`
    stop_if_failed $? "failed to create EC2 security group"
    # KEYPAIR (Created just because mandatory, will be swapped manually fore core user later on)
    $AWS ec2 create-key-pair --key-name $RESOURCES_NAME --no-paginate
    stop_if_failed $? "failed to create EC2 keypair"
    # SG SETUP
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 22 --cidr 0.0.0.0/0 --no-paginate
    stop_if_failed $? "failed to create SSH rule for security group"
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 6443 --cidr 0.0.0.0/0 --no-paginate
    stop_if_failed $? "failed to create API rule for security group"
    $AWS ec2 authorize-security-group-ingress --group-name $RESOURCES_NAME --protocol tcp --port 443 --cidr 0.0.0.0/0 --no-paginate
    stop_if_failed $? "failed to create HTTPS rule for security group"
    # CREATE INSTANCE
    $AWS ec2 run-instances --no-paginate --tag-specifications "ResourceType=instance,Tags=[{Key=Name,Value=$RESOURCES_NAME}]" \
    --image-id $AMI_ID --instance-type $INSTANCE_TYPE --security-group-ids $GROUPID --key-name $RESOURCES_NAME > $WORKDIR/$INSTANCE_DESCRIPTION

    stop_if_failed $? "failed to launch EC2 instance"
}

_bash_destroy_ec2_resources() {
  ID=$WORKDIR/$INSTANCE_DESCRIPTION
  [ ! -f  $ID ] && stop_if_failed 1 "Missing EC2 resource descriptor $INSTANCE_DESCRIPTION in $WORKDIR"
  INSTANCE_ID=`$JQ -r '.Instances[0].InstanceId' $ID`
  stop_if_failed $? "Failed to parse InstanceId from $ID"
  KEY_NAME=`$JQ -r '.Instances[0].KeyName' $ID`
  stop_if_failed $? "Failed to parse KeyName from $ID"
  SG_ID=`$JQ -r '.Instances[0].SecurityGroups[0].GroupId' $ID`
  stop_if_failed $? "Failed to parse SecurityGroupName from $ID"
  AZ=`$JQ -r '.Instances[0].Placement.AvailabilityZone' $ID`
  stop_if_failed $? "Failed to parse AvailabilityZone from $ID"

  # CHECK IF INSTANCE IS FOUND 
  $AWS ec2 describe-instance-status --instance-id $INSTANCE_ID  > /dev/null 2>&1
  stop_if_failed $? "instance $INSTANCE_ID not found, is your AWS_PROFILE exported and working?"
  # KILL INSTANCE
  pr_info "terminating instance $INSTANCE_ID"
  $AWS ec2 terminate-instances --instance-ids $INSTANCE_ID
  stop_if_failed $? "failed to terminate instance $INSTANCE_ID"
  # WAIT FOR INSTANCE TO TERMINATE
  while [[ `$AWS ec2 describe-instance-status --instance-id $INSTANCE_ID | $JQ -r ".InstanceStatuses[0].SystemStatus.Status"` != null ]]
  do 
    pr_info "waiting instance $INSTANCE_ID to shutdown, hang on...."
    sleep 1
  done

  pr_info "removing security group $SG_ID"
  TRY=0
  until `$AWS ec2 delete-security-group --group-id $SG_ID`
  do
    pr_info "waiting the security group to be removable try $TRY"
    ((TRY++))
    [[ $TRY == $TEARDOWN_MAX_RETRIES ]] && stop_if_failed 1 "failed to remove $SG_ID, are you sure that this sg exists?!"
    sleep 6
  done
  pr_info "removed security group $SG_ID"
  pr_info "removing keypair $KEY_NAME"
  $AWS --region ${AZ::-1} ec2 delete-key-pair --key-name $KEY_NAME
  stop_if_failed $? "failed to remove keypair $KEY_NAME"
  pr_end "everything has been cleaned up!"
  exit 0
}