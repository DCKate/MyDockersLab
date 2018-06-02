# System requires:
* Ansible  2.5.3
* VirtualBox  5.2.12
* Vagrant  2.1.1

# How to run:

1. Use vagrant to start vm:
```
$vagrant up
```
2. Run playbook to deploy on vm:
```
$ansible-playbook -i hosts.yml go.yml -e@config.yml -l mongo
```

## Local test with playbook 
Run the test role:
```
$ansible-playbook -i roles/localtest/files/hosts.yml go.yml -e@roles/localtest/files/config.yml -l localserver
```

# Organize the setting of database

## Replica Set
1. host.yml
The host file for ansible, include all of the server deployed with.  The alias of server need to be consisted with config.yml.  
2. go.yml
The play used for ansible, and the role of replica need to repeat as many time as the number of replica set.  It is need to assign which dictionary the info of replica set as variable.
3. config.yml
All of the information need for replica set.

> There are two replica set (replica_sa, replica_sb) deployed on two server (atom1, atom2) in this configuration. 

## Shard
todo 