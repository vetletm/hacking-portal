package openstack

import (
	"log"
	"os"

	"hacking-portal/db"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

var machines *db.MachineDatabase
var provider *gophercloud.ProviderClient

// Reboot takes server UUID and attempts to reboot it
func Reboot(uuid string) error {
	machine, err := machines.FindByID(uuid)
	if err != nil {
		return err
	}

	opts := gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")}

	client, err := openstack.NewComputeV2(provider, opts)
	if err != nil {
		return err
	}

	// Attempt to reboot the server and return the error
	err = servers.Reboot(client, machine.UUID, servers.RebootOpts{Type: servers.SoftReboot}).ExtractErr()
	if err != nil {
		log.Println("Server: ", machine.UUID, " failed to reboot")
		return err
	}

	log.Println("Server: ", machine.UUID, " was rebooted")
	return err
}

// Status takes server UUID and checks its status
func Status(uuid string) (string, error) {
	// Check if uuid is in database
	machine, err := machines.FindByID(uuid)
	if err != nil {
		return "", err
	}

	// Set openstack region name
	opts := gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")}

	// Configure client for connecting to openstack
	client, err := openstack.NewComputeV2(provider, opts)
	if err != nil {
		return "", err
	}

	// Get the server object
	server, err := servers.Get(client, machine.UUID).Extract()
	if err != nil {
		return "", err
	}

	// Print the status and return
	log.Println("Server ", server.ID, " is ", server.Status)
	return server.Status, err
}

// Init attempts to setup a connection
func Init() {
	AuthOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatal("Attempted to set authoptions, error: ", err)
		return
	}
	AuthOpts.DomainName = os.Getenv("OS_USER_DOMAIN_NAME")

	provider, err = openstack.AuthenticatedClient(AuthOpts)
	if err != nil {
		log.Fatal("Attempted to set provider, error: ", err)
		return
	}
}