//go:build darwin && arm64

// Package vmnet is the Go binding for Apple's vmnet framework, as
// Virtualization.framework's VZVmnetNetworkDeviceAttachment consumes it
// (macOS 26), plus PerVM, one such network created for one virtual machine.
package vmnet

import (
	"errors"
	"fmt"
	"net/netip"
)

// ErrNetwork indicates the per-VM vmnet network could not be created or
// interrogated.
var ErrNetwork = errors.New("vm network failure")

// PerVM is a vmnet network created for ONE virtual machine, plus the
// addressing a guest needs to use it.
//
// It replaces Virtualization.framework's NAT attachment
// (VZNATNetworkDeviceAttachment), and the reason is not performance, it is
// ownership. That attachment silently joins the host-wide vmnet SHARED
// network: one bridge, one subnet, one DHCP server, shared by every VM
// product on the machine. Its configuration is system state
// (/Library/Preferences/SystemConfiguration/com.apple.vmnet.plist,
// /etc/bootpd.plist) that any root-privileged peer rewrites for everyone —
// observed 2026-08-01, a peer client renumbered the network from
// 192.168.64.0/24 to 192.168.138.0/23 and left bootpd with dhcp_enabled
// false, after which every guest on it failed to boot while the bridge and
// NAT engine kept working for its owner. There is no defensive fix for that
// from inside a tenant.
//
// A per-VM network removes the shared state entirely: the framework
// allocates a distinct subnet per network (measured: two concurrent VMs got
// 192.168.64.0/24 and 192.168.65.0/24), arbitrates collisions itself
// (creating a network on a subnet already in use fails at create time
// rather than corrupting routing), and reaps the bridge with the process.
//
// It also makes the guest address known host-side before the VM starts, so a
// guest that can be configured from the host (over a control channel) needs
// no DHCP: about a millisecond, against the 85-300ms a lease acquisition
// costs.
//
// Requires macOS 26 (vmnet_network_configuration). It needs neither the
// com.apple.vm.networking entitlement nor root: the
// com.apple.security.virtualization entitlement is enough (verified).
type PerVM struct {
	// Gateway is the network's host-side address: the guest's default route,
	// its NAT inside address, and where the framework answers DNS.
	Gateway netip.Addr
	// Guest is the address set aside for the VM: the first address after the
	// gateway. A guest configured from the host takes it; a guest that asks
	// DHCP instead is handed an address by the framework, which may differ.
	Guest netip.Addr
	// Prefix is the subnet the framework allocated.
	Prefix netip.Prefix

	network *Network
}

// NewPerVM creates a per-VM vmnet network in shared mode and derives its
// addressing.
//
// The subnet is deliberately NOT requested: letting the framework allocate
// makes it the single arbiter across every VM on the host (ours and other
// products'), which is exactly the conflict-avoidance a hand-rolled
// allocator would have to reimplement badly.
func NewPerVM() (*PerVM, error) {
	config, err := NewNetworkConfiguration(SharedMode)
	if err != nil {
		return nil, fmt.Errorf("%w: vmnet configuration (macOS 26+ required): %w", ErrNetwork, err)
	}

	network, err := NewNetwork(config)
	if err != nil {
		return nil, fmt.Errorf("%w: creating vmnet network: %w", ErrNetwork, err)
	}

	prefix, err := network.IPv4Subnet()
	if err != nil {
		network.Release()

		return nil, fmt.Errorf("%w: reading allocated subnet: %w", ErrNetwork, err)
	}

	gateway, guest, err := addressing(prefix)
	if err != nil {
		network.Release()

		return nil, err
	}

	return &PerVM{
		Gateway: gateway,
		Guest:   guest,
		Prefix:  prefix,
		network: network,
	}, nil
}

// Network is the underlying vmnet network, to build the VM's network device
// attachment from.
func (n *PerVM) Network() *Network {
	return n.network
}

// GuestCIDR is the guest address in "addr/bits" form.
func (n *PerVM) GuestCIDR() string {
	return fmt.Sprintf("%s/%d", n.Guest, n.Prefix.Bits())
}

// addressing derives the gateway and guest addresses from the prefix the
// framework reported.
//
// vmnet_network_get_ipv4_subnet reports the address the NETWORK ITSELF holds
// — the bridge's, i.e. the gateway (observed: "192.168.64.1/24", not
// .0/24). So the gateway is taken from that value directly rather than
// assumed to be the first host address; only when the framework reports the
// network base (a form this API has not been seen to use, but which the type
// permits) is the first host address derived instead. The guest then takes
// the next address after the gateway.
func addressing(prefix netip.Prefix) (gateway, guest netip.Addr, err error) {
	if !prefix.Addr().Is4() {
		return netip.Addr{}, netip.Addr{}, fmt.Errorf("%w: subnet %s is not IPv4", ErrNetwork, prefix)
	}

	base := prefix.Masked().Addr()

	gateway = prefix.Addr()
	if gateway == base {
		gateway = base.Next()
	}

	guest = gateway.Next()
	if !prefix.Contains(guest) {
		return netip.Addr{}, netip.Addr{}, fmt.Errorf(
			"%w: subnet %s has no room for a guest address beside gateway %s", ErrNetwork, prefix, gateway,
		)
	}

	return gateway, guest, nil
}

// Close releases the network's reservation with the framework. The subnet the
// daemon set aside for this network stays booked, host-wide, until the
// reference is released: a process that exits without this leaves it behind
// until the daemon restarts, and 64 such exits leave no subnet for anyone.
// Idempotent; the VM must be stopped first, since its device still attaches
// here.
func (n *PerVM) Close() {
	if n == nil || n.network == nil {
		return
	}

	n.network.Release()
	n.network = nil
}
