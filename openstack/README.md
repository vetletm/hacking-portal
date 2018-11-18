# Openstack package
This package allows us to check status, reboot, and rebuild any machine present in our infrastructure. We use the github.com/gophercloud/gophercloud library as our entrypoint into SkyHiGh (NTNU's internal openstack implementation).


Difficulties:
- The library itself is huge and somewhat weirdly documented. Each package has several sub-packages each with their own godoc page, making it rather challenging to find the correct methods and functions.

What we'd like to do, but can't due to time:
- Implement an admin interface where we can create, delete, reboot any server belonging to any group.
- Add new keypairs to a server. The challenge was that Openstack's API does not allow appending a new keypair to a server without rebuilding. Which would have resulted in us going far beyond our original scope.
