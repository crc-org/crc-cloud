#!/bin/bash


stop_if_failed(){
	EXIT_CODE=$1
	MESSAGE=$2
	if [[ $EXIT_CODE != 0 ]]
	then
		echo "$MESSAGE"
		exit $EXIT_CODE
	fi
}

setup_dsnmasq(){
    echo "Writing Dnsmasq conf on $DNSMASQ_CONF"
        sudo cat << EOF > $DNSMASQ_CONF
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
    echo  "Adding Dnsmasq as primary DNS"
    sleep 2
    sudo nmcli connection modify Wired\ connection\ 1 ipv4.dns "10.88.0.8,169.254.169.254"
    stop_if_failed  $? "Failed to modify NetworkManager settings"
    echo  "Restarting NetworkManager"
    sleep 2
    sudo systemctl restart NetworkManager 
    stop_if_failed $? "Failed to restart NetworkManager"
    echo  "Enabling & starting Dnsmasq service"
    sudo systemctl enable crc-dnsmasq.service
    sudo systemctl start crc-dnsmasq.service
    sleep 2
    stop_if_failed $? "Failed to start Dnsmasq service"
}

enable_and_start_kubelet() {
    echo  "Enabling & starting Kubelet service"
    sudo systemctl enable kubelet
    sudo systemctl start kubelet
    stop_if_failed $? "Failed to start Kubelet service"
}

check_cluster_unhealthy() {
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





#setup_dsnmasq
#enable_and_start_kubelet
echo "waiting cluster to become healthy"
wait_cluster_become_healthy
stop_if_failed $? "Failed to recover Cluster after $(expr $CLUSTER_HEALTH_RETRIES \* $CLUSTER_HEALTH_SLEEP) seconds"