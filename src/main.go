package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/osrg/gobgp/server"
	"google.golang.org/grpc/peer"
)

var (
	grpcAddr = flag.String("grpcAddr", ":50051", "gRPC listen address")
	certFile = flag.String("certFile", "cert.pem", "TLS certificate file")
	keyFile  = flag.String("keyFile", "key.pem", "TLS private key file")
)

// BGPService is the gRPC service that implements the BGP API
type BGPService struct {
	s *server.BgpServer
}

// AddPath is a unary RPC that implements the AddPath API
func (s *BGPService) AddPath(ctx context.Context, path *Path) (*Empty, error) {
	// Get the peer address from the context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}
	// Parse the peer address
	addr, err := net.ResolveTCPAddr("tcp", p.Addr.String())
	if err != nil {
		return nil, err
	}
	// Convert the path to a BGPPath
	bgpPath := path.toBGPPath(addr.IP)
	// Add the path to the BGP server
	_, err = s.s.AddPath(ctx, &server.AddPathRequest{Path: bgpPath})
	if err != nil {
		return nil, err
	}
	return &Empty{}, nil
}

// DeletePath is a unary RPC that implements the DeletePath API
func (s *BGPService) DeletePath(ctx context.Context, path *Path) (*Empty, error) {
	// Get the peer address from the context
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}
	// Parse the peer address
	addr, err := net.ResolveTCPAddr("tcp", p.Addr.String())
	if err != nil {
		return nil, err
	}
	// Convert the path to a BGP
	bgpPath := path.toBGPPath(addr.IP)
	// Delete the path from the BGP server
	_, err = s.s.DeletePath(ctx, &server.DeletePathRequest{Path: bgpPath})
	if err != nil {
		return nil, err
	}
	return &Empty{}, nil
}
