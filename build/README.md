### Terraform

Used to create the VMs, their networks, floating IPs etc on OpenStack.  
Usage:
```bash
# source OpenStack RC v3 file
source openstack.rc

# deploy with Terraform
terraform init
terraform apply
```

### Ansible

Used to install the software on all the machines, including Docker, MongoDB and the Go application, latter two inside a Docker Swarm.

Before running, make sure to populate the credentials in `group_vars/docker1.yml`.  
Usage:
```bash
ansible-playbook -i interface playbook.yml
```
