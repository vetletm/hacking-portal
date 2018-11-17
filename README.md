# Hacking Portal

This project deploys a web interface to manage and assign Kali VMs in a lab environment used for the Ethical Hacking project in the IMT3004 course at NTNU Gjøvik. It allows an admin (lecturer) to assign up to three Kali VMs to each group of students, and it allows students to restart their assigned VMs without having access to the OpenStack environment in where they are hosted. This was originally a part of the PEMA bachelor assignment for the class of '16.

The original scope was to authenticate using LDAP and NTNU's own authentication infrastructure, integrate with OpenStack for VM assignments and management, store intermediate information in a MongoDB database, all of which hosted in a Docker Swarm on VMs in NTNU's own OpenStack environment. The scope also included a tasks/answers interface, but time didn't allow it.

Technologies used:
- Go
- MongoDB
- OpenStack
- Docker (Swarm)
- Terraform (orchestrating the test deployment)
- Ansible (provision the test deployment)

Difficulties:
- We had some issues with MongoDB replias in Docker Swarm, so we ended up running a single instance
	- Documentation on this was not easy to come by
- We also had some issues (initially) with the OpenStack package for Go
	- Documentation was outdated and the package was very complex, but it was the only option

We managed to get a working example with a good enough UI/UX using the Bootstrap CSS framework, enough to impress the IMT3004 lecturers that suggested the project topic.

We estimate having used 30-40 hours _each_ during the development of this project, way beyond the 20 hour minimum.
We learned the value of describing an aptly sized scope beforehand (which we slightly failed), and the value of good documentation (which we struggled with finding).

## Usage

Deploy with Terraform and Ansible (see the `/build` directory), then visit the public IP of any of the docker VMs.
