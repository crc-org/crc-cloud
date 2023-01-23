#!/bin/bash
#CONST
export KUBECONFIG="/opt/kubeconfig"
LOG_PATH="/tmp"
LOG_FILE="$LOG_PATH/_RANDOM_SUFFIX_.log"
DNSMASQ_CONF="/var/srv/dnsmasq.conf"
CLUSTER_HEALTH_SLEEP=8
CLUSTER_HEALTH_RETRIES=500
STEPS_SLEEP_TIME=10
#REPLACED VARS
IIP="${IIP:-"_IIP_"}"
EIP="${EIP:-"_EIP_"}"
PULL_SECRET="${PULL_SECRET:-"_PULL_SECRET_"}"
PASS_DEVELOPER="${PASS_DEVELOPER:-"_PASS_DEVELOPER_"}"
PASS_KUBEADMIN="${PASS_KUBEADMIN:-"_PASS_KUBEADMIN_"}"
PASS_REDHAT="${PASS_REDHAT:-"_PASS_REDHAT_"}"
MAXIMUM_LOGIN_RETRY=500

pr_info() {
    echo "[INF] $1" | tee -a $LOG_FILE > /dev/null
}

pr_error() {
    echo "[ERR] $1" | tee -a $LOG_FILE > /dev/null
}

pr_end() {
    echo "[END] $1" | tee -a $LOG_FILE > /dev/null
}

stop_if_failed(){
	EXIT_CODE=$1
	MESSAGE=$2
	if [[ $EXIT_CODE != 0 ]]
	then
		pr_error "$MESSAGE" 
		exit $EXIT_CODE
	fi
}

replace_default_ca() {
    USER="system:admin"
    GROUP="system:masters"
    USER_SUBJ="/O=${GROUP}/CN=${USER}"
    NAME="custom"
    CA_SUBJ="/OU=openshift/CN=admin-kubeconfig-signer-custom"
    VALIDITY=3650
    pr_info "replacing the default cluster CA and invalidating default kubeconfig"
    openssl genrsa -out $NAME-ca.key 4096
    stop_if_failed $? "failed to generate CA private key"
    openssl req -x509 -new -nodes -key $NAME-ca.key -sha256 -days $VALIDITY -out $NAME-ca.crt -subj "$CA_SUBJ"
    stop_if_failed $? "failed to generate CA certificate"
    openssl req -nodes -newkey rsa:2048 -keyout $USER.key -subj "$USER_SUBJ" -out $USER.csr
    stop_if_failed $? "failed to issue the CSR"
    openssl x509 -extfile <(printf "extendedKeyUsage = clientAuth") -req -in $USER.csr \
       -CA $NAME-ca.crt -CAkey $NAME-ca.key -CAcreateserial -out $USER.crt -days $VALIDITY -sha256
    stop_if_failed $? "failed to generate new admin certificate"
    oc create configmap client-ca-custom -n openshift-config --from-file=ca-bundle.crt=$NAME-ca.crt
    stop_if_failed $? "failed to create user certficate ConfigMap"
    oc patch apiserver cluster --type=merge -p '{"spec": {"clientCA": {"name": "client-ca-custom"}}}'
    stop_if_failed $? "failed to patch API server with newly created certificate"
    oc create configmap admin-kubeconfig-client-ca -n openshift-config --from-file=ca-bundle.crt=$NAME-ca.crt \
    --dry-run -o yaml | oc replace -f -
    stop_if_failed $? "failed to replace OpenShift CA"
}

login () {
    pr_info "logging in again to update $KUBECONFIG"
    COUNTER=0
    until `oc login --insecure-skip-tls-verify=true -u kubeadmin -p "$PASS_KUBEADMIN" https://api.crc.testing:6443 > /dev/null 2>&1`
    do 
        [ $COUNTER == $MAXIMUM_LOGIN_RETRY ] && stop_if_failed 1 "impossible to login on OpenShift, installation failed."
        pr_info "logging into OpenShift with updated credentials try $COUNTER, hang on...."
        sleep 5
        ((COUNTER++))
    done
}

wait_for_resource() {
    local resource=$1
    local retry=0
    local max_retry=20
    until `oc get $resource > /dev/null 2>&1`
    do
        [ $retry == $max_retry ] && stop_if_failed 1 "impossible to get resource ${resource}"
        pr_info "waiting for ${resource} to become available try $retry, hang on...."
        sleep 5
        ((retry++))
    done
}

#Replaces the default pubkey with the new one just generated to avoid the mysterious service to replace it later on :-\
replace_default_pubkey() {
    pr_info "Updating the public key resource for machine config operator"
    local pub_key=$(tr -d '\n\r' < /home/core/id_rsa.pub)
    wait_for_resource machineconfig
    oc patch machineconfig 99-master-ssh -p "{\"spec\": {\"config\": {\"passwd\": {\"users\": [{\"name\": \"core\", \"sshAuthorizedKeys\": [\"${pub_key}\"]}]}}}}" --type merge
    stop_if_failed $? "failed to update public key to machine config operator"
}

setup_dsnmasq(){
    pr_info "writing Dnsmasq conf on $DNSMASQ_CONF"
         cat << EOF > /var/srv/dnsmasq.conf
user=root
port= 53
bind-interfaces
expand-hosts
log-queries
local=/crc.testing/
domain=crc.testing
address=/apps-crc.testing/$IIP
address=/api.crc.testing/$IIP
address=/api-int.crc.testing/$IIP
address=/crc-wz8dw-master-0.crc.testing/192.168.126.11
EOF

    stop_if_failed  $? "failed to write Dnsmasq configuration in $DNSMASQ_CONF"
    pr_info  "adding Dnsmasq as primary DNS"
    sleep 2
    nmcli connection modify Wired\ connection\ 1 ipv4.dns "10.88.0.8,169.254.169.254"
    stop_if_failed  $? "failed to modify NetworkManager settings"
    pr_info  "restarting NetworkManager"
    sleep 2
    systemctl restart NetworkManager 
    stop_if_failed $? "failed to restart NetworkManager"
    pr_info  "enabling & starting Dnsmasq service"
    systemctl enable crc-dnsmasq.service
    systemctl start crc-dnsmasq.service
    sleep 2
    stop_if_failed $? "failed to start Dnsmasq service"
}

enable_and_start_kubelet() {
    pr_info  "enabling & starting Kubelet service"
    systemctl enable kubelet
    systemctl start kubelet
    stop_if_failed $? "failed to start Kubelet service"
}

check_cluster_unhealthy() {
    WAIT="authentication|console|etcd|ingress|openshift-apiserver"
    [ ! -z $1 ] && WAIT=$1

    until `oc get co > /dev/null 2>&1` 
    do
        pr_info "waiting Openshift API to become healthy, hang on...."
        sleep 2
    done

    for i in $(oc get co | grep -P "$WAIT" | awk '{ print $3 }')
    do
        if [[ $i == "False" ]] 
        then
            return 0
        fi
    done
    return 1
}

wait_cluster_become_healthy () {
    COUNTER=0
    W="[ALL]"
    [ ! -z $1 ] && W="[$1]" 
    while $(check_cluster_unhealthy $1)
    do
        sleep $CLUSTER_HEALTH_SLEEP
        if [[ $COUNTER == $CLUSTER_HEALTH_RETRIES ]]
        then
            return 1
        fi
        pr_info "checking for the $COUNTER time if the OpenShift Cluster has become healthy, hang on....$W"
	    ((COUNTER++))
    done
    pr_info "cluster has become ready in $(expr $COUNTER \* $CLUSTER_HEALTH_SLEEP) seconds"
    return 0
}


patch_pull_secret() {
    pr_info  "patching OpenShift pull secret"
    oc patch secret pull-secret -p "{\"data\":{\".dockerconfigjson\":\"$PULL_SECRET\"}}" -n openshift-config --type merge
    stop_if_failed $? "failed patch OpenShift pull secret"
    sleep $STEPS_SLEEP_TIME
}

create_certificate_and_patch_secret() {
    pr_info  "creating OpenShift secrets"
    openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout nip.key -out nip.crt -subj "/CN=$EIP.nip.io" -addext "subjectAltName=DNS:apps.$EIP.nip.io,DNS:*.apps.$EIP.nip.io,DNS:api.$EIP.nip.io"
    oc create secret tls nip-secret --cert=nip.crt --key=nip.key -n openshift-config
    stop_if_failed $? "failed patch OpenShift pull secret"
    sleep $STEPS_SLEEP_TIME
}

patch_ingress_config() {
    pr_info  "patching Cluster Ingress config"
    cat <<EOF > ingress-patch.yaml
spec:
  appsDomain: apps.$EIP.nip.io
  componentRoutes:
  - hostname: console-openshift-console.apps.$EIP.nip.io
    name: console
    namespace: openshift-console
    servingCertKeyPairSecret:
      name: nip-secret
  - hostname: oauth-openshift.apps.$EIP.nip.io
    name: oauth-openshift
    namespace: openshift-authentication
    servingCertKeyPairSecret:
      name: nip-secret
EOF
    oc patch ingresses.config.openshift.io cluster --type=merge --patch-file=ingress-patch.yaml
    stop_if_failed $? "failed patch Cluster Ingress config"
    #sleep $STEPS_SLEEP_TIME
}

patch_api_server() {
    pr_info  "patching API server"
    oc patch apiserver cluster --type=merge -p '{"spec":{"servingCerts": {"namedCertificates":[{"names":["api.'$EIP'.nip.io"],"servingCertificate": {"name": "nip-secret"}}]}}}'
    stop_if_failed $? "failed patch API server"
    #sleep $STEPS_SLEEP_TIME
}

patch_default_route() {
    pr_info  "patching default route"
    oc patch -p '{"spec": {"host": "default-route-openshift-image-registry.'$EIP'.nip.io"}}' route default-route -n openshift-image-registry --type=merge
    stop_if_failed $? "failed patch default route"
    #sleep $STEPS_SLEEP_TIME
}

set_credentials() {
    pr_info  "setting cluster credentials"
    podman run --rm -ti xmartlabs/htpasswd developer $PASS_DEVELOPER > htpasswd.developer
    stop_if_failed $? "failed to set developer password"
    podman run --rm -ti xmartlabs/htpasswd kubeadmin $PASS_KUBEADMIN > htpasswd.kubeadmin
    stop_if_failed $? "failed to set kubeadmin password"
    podman run --rm -ti xmartlabs/htpasswd redhat $PASS_REDHAT > htpasswd.redhat
    stop_if_failed $? "failed to set redhat password"

    cat htpasswd.developer > htpasswd.txt
    cat htpasswd.kubeadmin >> htpasswd.txt
    cat htpasswd.redhat >> htpasswd.txt
    sed -i '/^\s*$/d' htpasswd.txt

    oc create secret generic htpass-secret  --from-file=htpasswd=htpasswd.txt -n openshift-config --dry-run=client -o yaml > /tmp/htpass-secret.yaml
    stop_if_failed $? "failed to create Cluster secret"
    oc replace -f /tmp/htpass-secret.yaml
    stop_if_failed $? "failed to replace Cluster secret"
}

setup_dsnmasq

enable_and_start_kubelet
replace_default_pubkey
wait_cluster_become_healthy "etcd|openshift-apiserver"
stop_if_failed $? "failed to recover Cluster after $(expr $CLUSTER_HEALTH_RETRIES \* $CLUSTER_HEALTH_SLEEP) seconds"

set_credentials
replace_default_ca
login
stop_if_failed $? "failed to recover Cluster after $(expr $CLUSTER_HEALTH_RETRIES \* $CLUSTER_HEALTH_SLEEP) seconds"




patch_pull_secret
wait_cluster_become_healthy "etcd|openshift-apiserver"
stop_if_failed $? "failed to recover Cluster after $(expr $CLUSTER_HEALTH_RETRIES \* $CLUSTER_HEALTH_SLEEP) seconds"

create_certificate_and_patch_secret
wait_cluster_become_healthy "etcd|openshift-apiserver"
stop_if_failed $? "failed to recover Cluster after $(expr $CLUSTER_HEALTH_RETRIES \* $CLUSTER_HEALTH_SLEEP) seconds"

#PATCHES authentication|console|etcd|kube-apiserver
patch_ingress_config
#wait_cluster_become_healthy "etcd|openshift-apiserver"
patch_api_server
#wait_cluster_become_healthy "etcd|openshift-apiserver|ingress|network|dns"
patch_default_route
#wait_cluster_become_healthy "etcd|openshift-apiserver|authentication"

wait_cluster_become_healthy "authentication|console|etcd|ingress|openshift-apiserver"

until `oc get route console-custom -n openshift-console > /dev/null 2>&1` 
do
    pr_info "waiting for console route to become ready, hang on...."
    sleep 2
done 

CONSOLE_ROUTE=`oc get route console-custom -n openshift-console -o json | jq -r '.spec.host'`
pr_end $CONSOLE_ROUTE