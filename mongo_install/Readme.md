System requires:
Ansible use 2.5.3
VirtualBox use 5.2.12
Vagrant use 2.1.1

How to run: 
1. start vm with vagrant
```
$vagrant up
```
2. run ansible
```
$ansible-playbook  -i hosts.yml go.yml -e@config.yml -l mongo
```

*Run test role use:
```
$ansible-playbook  -i roles/localtest/files/hosts.yml go.yml -e@roles/localtest/files/config.yml -l localserver
```