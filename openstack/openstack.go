package openstack

import (
	"log"
	"os"
	"strings"

	"hacking-portal/db"
	"hacking-portal/models"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// For accessing methods on machines database
var machines *db.MachineDatabase

// For connecting to openstack
var client *gophercloud.ServiceClient

// Reboot takes server UUID and attempts to reboot it
func Reboot(uuid string) error {
	// Check if uuid is in database
	_, err := machines.FindByID(uuid)
	if err != nil {
		return err
	}

	// Attempt to reboot the server and return the error
	err = servers.Reboot(client, uuid, servers.RebootOpts{Type: servers.SoftReboot}).ExtractErr()
	if err != nil {
		log.Println("Server", uuid, "failed to reboot")
		return err
	}

	log.Println("Server", uuid, "was rebooted")
	return err
}

// Status takes server UUID and checks its status
func Status(uuid string) (string, error) {
	// Check if uuid is in database
	_, err := machines.FindByID(uuid)
	if err != nil {
		return "", err
	}

	// Get the server object
	server, err := servers.Get(client, uuid).Extract()
	if err != nil {
		return "", err
	}

	// Print the status and return
	log.Println("Server", server.ID, "is", server.Status)
	return server.Status, err
}

// getFloating find the floating IP and returns it
func getFloating(server servers.Server) string {
	// Iterate through Addresses until floating is found
	for _, networkAddresses := range server.Addresses {
		for _, element := range networkAddresses.([]interface{}) {
			address := element.(map[string]interface{})

			if address["OS-EXT-IPS:type"] == "floating" {
				return address["addr"].(string)
			}
		}
	}
	// If nothing was found
	return ""
}

// Init attempts to setup a connection
func Init() {
	// source options from environment
	authOpts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		log.Fatal("Attempted to set authoptions, error: ", err)
	}
	authOpts.DomainName = os.Getenv("OS_USER_DOMAIN_NAME")

	// authenticate with the OpenStack API
	provider, err := openstack.AuthenticatedClient(authOpts)
	if err != nil {
		log.Fatal("Attempted to set provider, error: ", err)
	}

	// grab a new compute client
	if client, err = openstack.NewComputeV2(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	}); err != nil {
		log.Fatal("Failed to initialize OpenStack client", err)
	}

	// grab a list of servers, which is paginated
	allPages, err := servers.List(client, servers.ListOpts{
		AllTenants: true,
	}).AllPages()
	if err != nil {
		log.Fatal("Failed to get server list from OpenStack", err)
	}

	// get all the servers from the paginated list
	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		log.Fatal("Failed to get all servers from OpenStack", err)
	}

	// iterate through all servers and attempt to put them into the database
	for _, server := range allServers {
		if strings.HasPrefix(strings.ToLower(server.Name), "kali") {
			// machine exists, update in database
			if err := machines.Upsert(models.Machine{
				Name:    server.Name,
				UUID:    server.ID,
				Address: getFloating(server),
			}); err != nil {
				log.Fatal("Attempted to insert new machine into db, error:", err)
				return
			}
		}
	}
}
