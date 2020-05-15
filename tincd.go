package tincd

import (
	"context"
	"fmt"
	"github.com/tinc-boot/tincd/internal"
	"github.com/tinc-boot/tincd/network"
	"net"
	"path/filepath"
)

// Base TINCD running instance. All methods should be goroutine safe
type Tincd interface {
	// Events bus
	Events() *network.Events
	// Stop service. Non-blocking, could be called several times
	Stop()
	// Get last service error (if exists)
	Error() error
	// Get wait channel. Will be close after stop
	Done() <-chan struct{}
	// Check that service is still running
	IsRunning() bool
	// Check that node (by name) is connected to the network
	IsActive(node string) bool
	// List of all connected peers
	Peers() []string
	// Get network definition
	Definition() *network.Network
}

// Start tincd (and tinc-web-boot protocol) services. Not blocking after start. If sudo is true it will try to ask
// administrative privileges for each platform (graphically if possible)
func Start(ctx context.Context, nw *network.Network, sudo bool) (*netImpl, error) {
	if !nw.IsDefined() {
		return nil, fmt.Errorf("network %s is not defined", nw.Name())
	}
	tincBin, err := internal.DetectTincBinary()
	if err != nil {
		return nil, fmt.Errorf("detect tinc binary: %w", err)
	}
	impl := &netImpl{
		definition: nw,
		tincBin:    tincBin,
	}
	return impl, impl.initAndStart(ctx, sudo)
}

// Start tincd (and tinc-web-boot protocol) services based on configuration in directory. Not blocking after start.
// If sudo is true it will try to ask
// administrative privileges for each platform (graphically if possible)
func StartFromDir(ctx context.Context, directory string, sudo bool) (*netImpl, error) {
	abs, err := filepath.Abs(directory)
	if err != nil {
		return nil, err
	}
	return Start(ctx, &network.Network{Root: abs}, sudo)
}

// Create (but not start) and configure new network in specified location with pre-parsed subnet. IP will be generated
// randomly.
// Base name of the location will be used as name of network.
func CreateNet(location string, subnet *net.IPNet) (*network.Network, error) {
	abs, err := filepath.Abs(location)
	if err != nil {
		return nil, err
	}
	name := filepath.Base(abs)

	if !network.IsValidName(name) {
		return nil, fmt.Errorf("invalid network name")
	}
	netw := &network.Network{Root: abs}
	return netw, netw.Configure(subnet)
}

// Create (but not start) and configure new network in specified location. Subnet should be defined in CIDR.
// IP will be generated randomly.
// Base name of the location will be used as name of network.
func Create(location string, subnet string) (*network.Network, error) {
	_, ipnet, err := net.ParseCIDR(subnet)
	if err != nil {
		return nil, err
	}
	return CreateNet(location, ipnet)
}
