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
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	hv1 "google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"flag"
)

var (
	port int
	listenAddr string
	bin string
)

func startServer() net.Listener {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	server := health.NewServer()
	hv1.RegisterHealthServer(grpcServer, server)

	server.SetServingStatus("foo", hv1.HealthCheckResponse_SERVING)
	server.SetServingStatus("bar", hv1.HealthCheckResponse_NOT_SERVING)
	go grpcServer.Serve(listener)
	return listener
}

func init() {
	flag.IntVar(&port, "stub-port", 54321, "port for the stub server")
	flag.StringVar(&bin, "gprobe", "", "path to the gprobe binary")
}

func TestMain(m *testing.M) {
	flag.Parse()
	listenAddr = fmt.Sprintf("%s:%d", "localhost", port)

	lis := startServer()
	result := 0
	defer func() {
		lis.Close()
		os.Exit(result)
	}()

	result = m.Run()
}

func TestShouldReturnServingForRunningServer(t *testing.T) {
	stdout, stderr, exitcode := runBin(t, listenAddr)

	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "SERVING\n", stdout)
	assert.Empty(t, stderr)
}

func TestShouldFailIfServerIsNotListening(t *testing.T) {
	stdout, stderr, exitcode := runBin(t, "nosuchhost:1234")

	assert.Equal(t, 127, exitcode)
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "error", "should print status to STDOUT")
}

func TestShouldReturnServingForHealthyService(t *testing.T) {
	stdout, stderr, exitcode := runBin(t, listenAddr, "foo")

	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "SERVING\n", stdout)
	assert.Empty(t, stderr)
}

func TestShouldReturnNotServingForUnhealthyService(t *testing.T) {
	stdout, stderr, exitcode := runBin(t, listenAddr, "bar")

	assert.Equal(t, 2, exitcode)
	assert.Equal(t, "NOT_SERVING\n", stdout)
	assert.Contains(t, stderr, "health-check failed")
}

func TestShouldNotFailForUnhealthyServiceIfNoFailIsSet(t *testing.T) {
	stdout, stderr, exitcode := runBin(t, "--no-fail", listenAddr, "bar")

	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "NOT_SERVING\n", stdout)
	assert.Empty(t, stderr)
}

func TestShouldFailIfServiceHealthCheckIsNotRegistered(t *testing.T) {
	stdout, stderr, exitcode := runBin(t, listenAddr, "non_registered_service")

	assert.Equal(t, 127, exitcode)
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "NotFound")
}

func runBin(t *testing.T, args ...string) (stdout string, stderr string, exitcode int) {
	gprobe := exec.Command(bin, args...)
	stdoutPipe, _ := gprobe.StdoutPipe()
	stderrPipe, _ := gprobe.StderrPipe()

	err := gprobe.Start()
	if err != nil {
		t.Error(err)
	}

	stdout = readPipe(t, stdoutPipe)
	stderr = readPipe(t, stderrPipe)
	exitcode = waitForExitCode(t, gprobe)

	return
}

func readPipe(t *testing.T, reader io.Reader) string {
	buf := new(bytes.Buffer)
	_, err := io.Copy(buf, reader)
	if err != nil {
		t.Error(err)
	}
	return buf.String()
}

func waitForExitCode(t *testing.T, cmd *exec.Cmd) (exitcode int) {
	err := cmd.Wait()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitcode = status.ExitStatus()
			}
		} else {
			exitcode = -1
			t.Error(err)
		}
	}
	return
}


