// PUBLIC DOMAIN NOTICE
// National Center for Biotechnology Information
//
// This software/database is a "United States Government Work" under the
// terms of the United States Copyright Act.  It was written as part of
// the author's official duties as a United States Government employee and
// thus cannot be copyrighted.  This software/database is freely available
// to the public for use. The National Library of Medicine and the U.S.
// Government have not placed any restriction on its use or reproduction.
//
// Although all reasonable efforts have been taken to ensure the accuracy
// and reliability of the software and data, the NLM and the U.S.
// Government do not and cannot warrant the performance or results that
// may be obtained by using this software or data. The NLM and the U.S.
// Government disclaim all warranties, express or implied, including
// warranties of performance, merchantability or fitness for any particular
// purpose.
//
// Please cite the author in any work or product based on this material.

package acctest

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	hv1 "google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

// StartServer starts new gRPC application with simple health service.
// It is callers responsibility to Stop the server
func StartServer(port int, certFile string, keyFile string) (*grpc.Server, *health.Server, error) {
	transportCredentials, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}
	return doStart(port, grpc.Creds(transportCredentials))
}

// StartInsecureServer starts new gRPC application with simple health service.
// It is callers responsibility to Stop the server
func StartInsecureServer(port int) (*grpc.Server, *health.Server, error) {
	return doStart(port)
}

func doStart(port int, options ...grpc.ServerOption) (server *grpc.Server, service *health.Server, err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	server = grpc.NewServer(options...)
	service = health.NewServer()
	hv1.RegisterHealthServer(server, service)

	go server.Serve(listener)
	return server, service, nil
}

// StartEmptyServer starts gRPC server application with no services
func StartEmptyServer(port int) (server *grpc.Server, err error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	server = grpc.NewServer()

	go server.Serve(listener)
	return server, nil
}
