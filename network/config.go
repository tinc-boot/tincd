package network

import (
	"fmt"
	"strconv"
	"strings"
	"tinc-boot/tincd/config"
)

// Main configuration for network (tinc.conf)
type Config struct {
	Name       string   `json:"name"`                 // self node name
	Port       uint16   `json:"port"`                 // listening port
	Interface  string   `json:"interface"`            // interface name (for Darwin should be empty)
	Mode       string   `json:"mode"`                 // mode, should be switch always
	Mask       int      `json:"mask"`                 // subnet mask size
	DeviceType string   `json:"deviceType,omitempty"` // device type (tap for most)
	Device     string   `json:"device,omitempty"`     // device name
	ConnectTo  []string `json:"connectTo,omitempty"`  // list of public nodes (automatically index)
	Broadcast  string   `json:"broadcast"`            // broadcast mode (mst)
}

// Upgrade few parameters of self node. Empty parameters are ignored
type Upgrade struct {
	Port    uint16    `json:"port,omitempty"`    // listening port
	Address []Address `json:"address,omitempty"` // list of public addresses
	Device  string    `json:"device,omitempty"`  // custom device name
}

// Public address
type Address struct {
	Host string `json:"host"`           // domain name or IP
	Port uint16 `json:"port,omitempty"` // optional port
}

func (addr *Address) String() string {
	if addr.Port != 0 {
		return fmt.Sprintf("%s %v", addr.Host, addr.Port)
	}
	return addr.Host
}

func (addr *Address) Scan(value string) error {
	hp := strings.SplitN(strings.TrimSpace(value), " ", 2)
	addr.Host = hp[0]
	if len(hp) == 1 {
		return nil
	}
	v, err := strconv.ParseUint(hp[1], 10, 16)
	addr.Port = uint16(v)
	return err
}

// Node configuration (as in hosts directory)
type Node struct {
	Name      string    `json:"name"`                                 // node name
	Subnet    string    `json:"subnet"`                               // subnet (should same for all nodes in network)
	Port      uint16    `json:"port"`                                 // optional listening port
	IP        string    `json:"ip"`                                   // VPN ip
	Address   []Address `json:"address,omitempty"`                    // list of public addresses
	PublicKey string    `json:"publicKey" tinc:"RSA PUBLIC KEY,blob"` // public RSA key
	Version   int       `json:"version"`                              // version. should be updated only by node-owner
}

func (cfg *Config) Build() (text []byte, err error) {
	return config.Marshal(cfg)
}

func (cfg *Config) Parse(text []byte) error {
	return config.Unmarshal(text, cfg)
}

func (n *Node) Build() (text []byte, err error) {
	return config.Marshal(n)
}

func (n *Node) Parse(data []byte) error {
	return config.Unmarshal(data, n)
}
