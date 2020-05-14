package tincd

import "time"

const (
	GreetInterval     = 5 * time.Second // interval between attempts to "greet" new nodes
	CommunicationPort = 4655            // default communication port inside VPN
)
