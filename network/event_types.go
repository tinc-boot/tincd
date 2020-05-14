package network

//go:generate events-gen -p network -E Events -s -P -o events.go

//event:"Stopped"
type NetworkID struct {
	Name string `json:"name"`
}

//event:"PeerDiscovered"
//event:"PeerJoined"
//event:"PeerLeft"
type PeerID struct {
	Network string `json:"network"`
	Node    string `json:"node"`
	Address string `json:"address,omitempty"`
}
