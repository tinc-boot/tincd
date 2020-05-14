package apiclient

import (
	"context"
	client "github.com/reddec/jsonrpc2/client"
	network "github.com/tinc-boot/tincd/network"
	"sync/atomic"
)

func Default() *APIClient {
	return &APIClient{BaseURL: "https://example.com/api"}
}

type APIClient struct {
	BaseURL  string
	sequence uint64
}

// Send self description and get known nodes
func (impl *APIClient) Exchange(ctx context.Context, self network.Node) (reply []network.Node, err error) {
	err = client.CallHTTP(ctx, impl.BaseURL, "API.Exchange", atomic.AddUint64(&impl.sequence, 1), &reply, self)
	return
}
