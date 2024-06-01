terraform {
  required_providers {
    openvswitch = {
      source = "example.com/local/openvswitch"
      version = "1.0.0"
    }
  }
}

provider "openvswitch" {}

resource "openvswitch_bridge" "sample_bridge" {
  name = "testbr0"
  ip_address = "192.168.100.1/24"
}


