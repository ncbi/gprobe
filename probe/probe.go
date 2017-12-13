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

package probe

import (
	"context"
	"google.golang.org/grpc"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
	"time"
)

var (
	// DialTimeout timeout for establishing connection to a gRPC server
	DialTimeout time.Duration = 1 * time.Second
	// OpTimeout timeout for performing gRPC calls
	OpTimeout time.Duration = 1 * time.Second

	ctx        = context.Background()
	connection *grpc.ClientConn
	client     hpb.HealthClient
)

// Connect to specified server
func Connect(serverAddresss string) (err error) {
	lctx, cancel := context.WithTimeout(ctx, DialTimeout)
	defer cancel()
	connection, err = grpc.DialContext(lctx, serverAddresss, grpc.WithInsecure())
	client = hpb.NewHealthClient(connection)
	return err
}

// Disconnect from server
func Disconnect() {
	connection.Close()
	connection = nil
}

// CheckServer checks health of the server overall
func CheckServer() (hpb.HealthCheckResponse_ServingStatus, error) {
	return doCheck(&hpb.HealthCheckRequest{})
}

// CheckService checks health of a specific service on the server
// If nil is passed as argument, the effect is the same as if CheckServer() was called
func CheckService(serviceName string) (hpb.HealthCheckResponse_ServingStatus, error) {
	return doCheck(&hpb.HealthCheckRequest{
		Service: serviceName,
	})
}

func doCheck(request *hpb.HealthCheckRequest) (status hpb.HealthCheckResponse_ServingStatus, err error) {
	lctx, cancel := context.WithTimeout(ctx, OpTimeout)
	defer cancel()

	response, err := client.Check(lctx, request)
	if response != nil {
		status = response.Status
	}
	return
}
