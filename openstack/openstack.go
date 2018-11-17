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

	result := servers.Reboot(client, machine.UUID, servers.RebootOpts{Type: servers.SoftReboot})

	// shit's trippin balls
	return result.ErrResult.Result.Err
}

// Status takes server UUID and checks its status
func Status(uuid string) error {
	// TODO: Check if UUID is in database
	// TODO: Check if server is ACTIVE
	return nil
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
