#!/bin/bash

echo "Check if inventory.yaml exists..."
if [ -f "inventory.yaml" ]; then
    echo "Starting Ansible..."
    export ANSIBLE_HOST_KEY_CHECKING=False
    export ANSIBLE_LOG_PATH=/home/user/ansible-logs/ansible.log
    ansible-playbook -i inventory.yaml crc-cloud/ansible/playbooks/start.yaml
else
    echo "Could not find inventory file. Exit"
    exit 1
fi
