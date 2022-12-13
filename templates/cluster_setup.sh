#!/bin/bash
export KUBECONFIG="/opt/kubeconfig"
LOG_PATH="/tmp"
LOG_FILE="$LOG_PATH/_RANDOM_SUFFIX_.log"
IIP="_IIP_"
DNSMASQ_CONF="/var/srv/dnsmasq.conf"
CLUSTER_HEALTH_SLEEP=2
CLUSTER_HEALTH_RETRIES=500


pr_info() {
    echo "_INF: $1" | tee -a $LOG_FILE
}

pr_error() {
    echo "_ERR: $1" | tee -a $LOG_FILE
}

pr_end() {
    echo "_END: " | tee -a $LOG_FILE
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

setup_dsnmasq(){
    pr_info "Writing Dnsmasq conf on $DNSMASQ_CONF"
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

    stop_if_failed  $? "Failed to write Dnsmasq configuration in $DNSMASQ_CONF"
    pr_info  "Adding Dnsmasq as primary DNS"
    sleep 2
    nmcli connection modify Wired\ connection\ 1 ipv4.dns "10.88.0.8,169.254.169.254"
    stop_if_failed  $? "Failed to modify NetworkManager settings"
    pr_info  "Restarting NetworkManager"
    sleep 2
    systemctl restart NetworkManager 
    stop_if_failed $? "Failed to restart NetworkManager"
    pr_info  "Enabling & starting Dnsmasq service"
    systemctl enable crc-dnsmasq.service
    systemctl start crc-dnsmasq.service
    sleep 2
    stop_if_failed $? "Failed to start Dnsmasq service"
}

enable_and_start_kubelet() {
    pr_info  "Enabling & starting Kubelet service"
    systemctl enable kubelet
    systemctl start kubelet
    stop_if_failed $? "Failed to start Kubelet service"
}

check_cluster_unhealthy() {
    RES=1
    while [[ $RES != 0 ]]
    do
        sleep 2
        pr_info "$RES waiting Openshift API to become ready, hang on...."
        oc get co > /dev/null 2>&1
        RES=$?
    done
    for i in $(oc get co | grep -P "authentication|console|etcd|kube-apiserver"| awk '{ print $3 }')
    do
        if [[ $i == "False" ]] 
        then
            return 0
        fi
    done
    return 1
}

wait_cluster_become_healthy () {
    counter=0
    while check_cluster_unhealthy
    do
        sleep $CLUSTER_HEALTH_SLEEP
        if [[ $counter == $CLUSTER_HEALTH_RETRIES ]]
        then
            return 1
        fi
	((counter++))
    done
    return 1

}





setup_dsnmasq
enable_and_start_kubelet
pr_info "waiting cluster to become healthy"
wait_cluster_become_healthy
stop_if_failed $? "Failed to recover Cluster after $(expr $CLUSTER_HEALTH_RETRIES \* $CLUSTER_HEALTH_SLEEP) seconds"