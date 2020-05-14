package api

import (
	"context"
	"tinc-boot/tincd/network"
)

type API interface {
	// Send self description and get known nodes
	Exchange(ctx context.Context, self network.Node) ([]network.Node, error)
}
