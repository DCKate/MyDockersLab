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

# Organize the setting of database
## Replica Set
### Shard Replica Sets
1. host.yml
The host file for ansible, include all of the server deployed with.  The alias of server need to be consisted with config.yml.  
2. go.yml
The play used for ansible, and the role of replica need to be repeated as many time as the number of replica set.  It is need to assign which dictionary the info of replica set as variable.
3. config.yml
All of the information need for replica set.

> There are two replica set (replica_sa, replica_sb) deployed on two server (atom1, atom2) in this configuration. 

### Config Server
Deploye as the role of replica set.
1. host.yml
This setting has one replica set for config server (group : configs).
2. go.yml
Since the role is deploye as replica, it need to be repeated as many time as the number of replica set.  However, this setting is assume only one replica set for config server, and once if need multi replica set, it could be modify as the format of shard replica sets.
3. config.yml
The setting on the config server (variable : replica_ca).

> There is a replica set (replica_ca) deployed on two server (atom1, atom2) in this configuration. 
## Router server
Add the setting of server on the group of router_set on config.yml. 

# Note
## Local test with playbook 
Run the test role:
```
$ansible-playbook -i roles/localtest/files/hosts.yml go.yml -e@roles/localtest/files/config.yml -l localserver
```

## Simple Test
It could use the tkmongo on the folder mongomongo to insert 10000 data as the setting of database on config.yml.