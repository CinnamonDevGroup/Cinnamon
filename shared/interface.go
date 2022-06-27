package shared

import (
	"context"
	"net/rpc"

	"google.golang.org/grpc"

	"github.com/AngelFluffyOokami/Cinnamon/proto"
	"github.com/hashicorp/go-plugin"
)

var cookey string
var cookieVal string

var Handshake = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   cookey,
	MagicCookieValue: cookieVal,
}

func ShareKeys(cookey string, cookieVal string) {
	cookey = cookey
	cookieVal = cookieVal
}

var PluginMap = map[string]plugin.Plugin{
	"kv_grpc": &KVGRPCPlugin{},
	"kv":      &KVPlugin{},
}

type KV interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
}

type KVPlugin struct {
	Impl KV
}

func (p *KVPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*KVPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

type KVGRPCPlugin struct {
	plugin.Plugin

	Impl KV
}

func (p *KVGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterKVServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *KVGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewKVClient(c)}, nil
}
