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
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"

	hv1 "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port        int
	cert        string
	key         string
	bin         string
	stubSrvAddr string
)

func init() {
	flag.IntVar(&port, "stub-port", 54321, "port for the stub server")
	flag.StringVar(&bin, "gprobe", "../gprobe", "path to the gprobe binary")
}

func TestMain(m *testing.M) {
	flag.Parse()
	stubSrvAddr = fmt.Sprintf("%s:%d", "localhost", port)
	os.Exit(m.Run())
}

func TestShouldReturnServingForRunningServer(t *testing.T) {
	// given
	srv, _, err := StartInsecureServer(port)
	if err != nil {
		log.Fatalf("can't start stub server: %v", err)
	}
	defer srv.GracefulStop()

	// when
	stdout, stderr, exitcode := runBin(t, stubSrvAddr)

	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "SERVING\n", stdout)
	assert.Empty(t, stderr)
}

func TestShouldFailIfServerIsNotListening(t *testing.T) {
	// given no server

	// when
	stdout, stderr, exitcode := runBin(t, stubSrvAddr)

	// then
	assert.Equal(t, 127, exitcode)
	assert.Empty(t, stdout)
	assert.Contains(t, stderr, "error")
}

func TestShouldReturnServingForHealthyService(t *testing.T) {
	// given
	srv, svc, err := StartInsecureServer(port)
	if err != nil {
		log.Fatalf("can't start stub server: %v", err)
	}
	defer srv.GracefulStop()
	svc.SetServingStatus("foo", hv1.HealthCheckResponse_SERVING)

	// when
	stdout, stderr, exitcode := runBin(t, stubSrvAddr, "foo")

	// then
	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "SERVING\n", stdout)
	assert.Empty(t, stderr)
}

func TestShouldReturnNotServingForUnhealthyService(t *testing.T) {
	// given
	srv, svc, err := StartInsecureServer(port)
	if err != nil {
		log.Fatalf("can't start stub server: %v", err)
	}
	defer srv.GracefulStop()
	svc.SetServingStatus("foo", hv1.HealthCheckResponse_NOT_SERVING)

	// when
	stdout, stderr, exitcode := runBin(t, stubSrvAddr, "foo")

	// then
	assert.Equal(t, 2, exitcode)
	assert.Equal(t, "NOT_SERVING\n", stdout)
	assert.Contains(t, stderr, "health-check failed")
}

func TestShouldNotFailForUnhealthyServiceIfNoFailIsSet(t *testing.T) {
	// given
	srv, svc, err := StartInsecureServer(port)
	if err != nil {
		log.Fatalf("can't start stub server: %v", err)
	}
	defer srv.GracefulStop()
	svc.SetServingStatus("foo", hv1.HealthCheckResponse_NOT_SERVING)

	// when
	stdout, stderr, exitcode := runBin(t, "--no-fail", stubSrvAddr, "foo")

	// then
	assert.Equal(t, 0, exitcode)
	assert.Equal(t, "NOT_SERVING\n", stdout)
	assert.Empty(t, stderr)
}

func TestShouldFailIfServiceHealthCheckIsNotRegistered(t *testing.T) {
	// given
	srv, _, err := StartInsecureServer(port)
	if err != nil {
		log.Fatalf("can't start stub server: %v", err)
	}
	defer srv.GracefulStop()

	// when
	stdout, stderr, exitcode := runBin(t, stubSrvAddr, "foo")

	// then
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
