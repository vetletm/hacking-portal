# use http backend for tfstate storage, config provided through
# `terraform init` to keep it confidential
# https://www.terraform.io/docs/backends/config.html#partial-configuration
terraform {
  backend "http" {}
}

# create a router that the internal network will use as a gateway
resource "openstack_networking_router_v2" "router" {
  name                = "router1"
  admin_state_up      = "true"
  external_network_id = "${var.os_external_network}"
}

# create a SSH key data type containing locally resourced keys
data "tls_public_key" "imt2681_terraform" {
  private_key_pem = "${file("imt2681_terraform")}"
}

# create security group for internal access
resource "openstack_compute_secgroup_v2" "ssh" {
  name        = "ssh"
  description = "SSH Access"

  rule {
    from_port   = 22
    to_port     = 22
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = -1
    to_port     = -1
    ip_protocol = "icmp"
    cidr        = "0.0.0.0/0"
  }
}

resource "openstack_compute_secgroup_v2" "http" {
  name        = "http"
  description = "HTTP access"

  rule {
    from_port   = 80
    to_port     = 80
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }

  rule {
    from_port   = 443
    to_port     = 443
    ip_protocol = "tcp"
    cidr        = "0.0.0.0/0"
  }
}
