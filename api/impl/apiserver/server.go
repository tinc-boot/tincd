// Code generated by jsonrpc2. DO NOT EDIT.
//go:generate jsonrpc2-gen -i ../../interface.go -I API -o ./server.go --package apiserver --go ../apiclient/client.go --go-package apiclient --go-linked
package apiserver

import (
	"context"
	"encoding/json"
	jsonrpc2 "github.com/reddec/jsonrpc2"
	api "github.com/tinc-boot/tincd/api"
	network "github.com/tinc-boot/tincd/network"
)

func RegisterAPI(router *jsonrpc2.Router, wrap api.API) []string {
	router.RegisterFunc("API.Exchange", func(ctx context.Context, params json.RawMessage, positional bool) (interface{}, error) {
		var args struct {
			Arg0 network.Node `json:"self"`
		}
		var err error
		if positional {
			err = jsonrpc2.UnmarshalArray(params, &args.Arg0)
		} else {
			err = json.Unmarshal(params, &args)
		}
		if err != nil {
			return nil, err
		}
		return wrap.Exchange(ctx, args.Arg0)
	})

	return []string{"API.Exchange"}
}
