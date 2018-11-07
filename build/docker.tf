# create an instance, using the network, security groups, key pair and cloud config defined below
resource "openstack_compute_instance_v2" "docker" {
  count       = "${var.docker_replicas}"
  name        = "docker${count.index+1}"
  image_name  = "${var.image_name}"
  flavor_name = "${var.flavor_name}"
  key_pair    = "${openstack_compute_keypair_v2.docker.name}"

  security_groups = [
    "${openstack_compute_secgroup_v2.ssh.name}",
    "${openstack_compute_secgroup_v2.http.name}",
    "${openstack_compute_secgroup_v2.swarm.name}",
  ]

  network = {
    uuid = "${openstack_networking_network_v2.docker.id}"
  }

  # terraform is not smart enough to realize we need a subnet first
  depends_on = ["openstack_networking_subnet_v2.docker"]
}

# create the internal management network
resource "openstack_networking_network_v2" "docker" {
  name           = "docker-net"
  admin_state_up = "true"
}

# create a subnet on the internal management network
resource "openstack_networking_subnet_v2" "docker" {
  name            = "docker-subnet"
  network_id      = "${openstack_networking_network_v2.docker.id}"
  cidr            = "${var.docker_cidr}"
  ip_version      = 4
  enable_dhcp     = "true"
  dns_nameservers = "${var.nameservers}"
}

# associate the docker subnet with the main router
resource "openstack_networking_router_interface_v2" "docker" {
  router_id = "${openstack_networking_router_v2.router.id}"
  subnet_id = "${openstack_networking_subnet_v2.docker.id}"
}

# appropriate a floating IP for the docker instance
resource "openstack_networking_floatingip_v2" "docker" {
  count = "${var.docker_replicas}"
  pool  = "${var.os_floating_ip_pool}"
}

resource "openstack_compute_floatingip_associate_v2" "docker" {
  count       = "${var.docker_replicas}"
  floating_ip = "${element(openstack_networking_floatingip_v2.docker.*.address, count.index)}"
  instance_id = "${element(openstack_compute_instance_v2.docker.*.id, count.index)}"
}

# create a key pair used for remote access to the manager instance
resource "openstack_compute_keypair_v2" "docker" {
  name       = "docker"
  public_key = "${var.os_keypair}"
}

# create security group for swarm communication
resource "openstack_compute_secgroup_v2" "swarm" {
  # https://docs.docker.com/engine/swarm/swarm-tutorial/#open-protocols-and-ports-between-the-hosts
  name        = "swarm"
  description = "Docker Swarm communication"

  rule {
    from_port   = 2377
    to_port     = 2377
    ip_protocol = "tcp"
    cidr        = "${var.docker_cidr}"
  }

  rule {
    from_port   = 7946
    to_port     = 7946
    ip_protocol = "tcp"
    cidr        = "${var.docker_cidr}"
  }

  rule {
    from_port   = 7946
    to_port     = 7946
    ip_protocol = "udp"
    cidr        = "${var.docker_cidr}"
  }

  rule {
    from_port   = 4789
    to_port     = 4789
    ip_protocol = "udp"
    cidr        = "${var.docker_cidr}"
  }
}
