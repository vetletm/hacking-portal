# create an instance, using the network, security groups, key pair and cloud config defined below
resource "openstack_compute_instance_v2" "kali" {
  count       = "${var.kali_replicas}"
  name        = "kali${count.index+1}"
  image_name  = "${var.image_name}"
  flavor_name = "${var.flavor_name}"
  key_pair    = "${openstack_compute_keypair_v2.kali.name}"

  security_groups = [
    "${openstack_compute_secgroup_v2.ssh.name}",
    "${openstack_compute_secgroup_v2.http.name}",
  ]

  network = {
    uuid = "${openstack_networking_network_v2.kali.id}"
  }

  # terraform is not smart enough to realize we need a subnet first
  depends_on = ["openstack_networking_subnet_v2.kali"]
}

# create the internal management network
resource "openstack_networking_network_v2" "kali" {
  name           = "kali-net"
  admin_state_up = "true"
}

# create a subnet on the internal management network
resource "openstack_networking_subnet_v2" "kali" {
  name            = "kali-subnet"
  network_id      = "${openstack_networking_network_v2.kali.id}"
  cidr            = "${var.kali_cidr}"
  ip_version      = 4
  enable_dhcp     = "true"
  dns_nameservers = "${var.nameservers}"
}

# associate the kali subnet with the main router
resource "openstack_networking_router_interface_v2" "kali" {
  router_id = "${openstack_networking_router_v2.router.id}"
  subnet_id = "${openstack_networking_subnet_v2.kali.id}"
}

# appropriate a floating IP for the docker instance
resource "openstack_networking_floatingip_v2" "kali" {
  count = "${var.kali_replicas}"
  pool  = "${var.os_floating_ip_pool}"
}

resource "openstack_compute_floatingip_associate_v2" "kali" {
  count       = "${var.kali_replicas}"
  floating_ip = "${element(openstack_networking_floatingip_v2.kali.*.address, count.index)}"
  instance_id = "${element(openstack_compute_instance_v2.kali.*.id, count.index)}"
}

# create a key pair used for remote access to the manager instance
resource "openstack_compute_keypair_v2" "kali" {
  name       = "kali"
  public_key = "${var.os_keypair}"
}
