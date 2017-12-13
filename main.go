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

package main

import (
	"fmt"
	"github.com/ncbi/gprobe/probe"
	"github.com/urfave/cli"
	"google.golang.org/grpc/health/grpc_health_v1"
	"os"
	"time"
)

var version string

const (
	// ExitCodeUsage is returned if application used incorrectly
	ExitCodeUsage = 1
	// ExitCodeHealthCheckNegative is returned if health status is not SERVING
	ExitCodeHealthCheckNegative = 2
	// ExitCodeUnexpected is returned if any other error happens
	ExitCodeUnexpected = 127
)

// appInput holds all parsed CLI flags and arguments
type appInput struct {
	timeout       time.Duration
	noFail        bool
	serverAddress string
	serviceName   string
}

// mainFn holds main application business logic
type mainFn func(appInput *appInput) *cli.ExitError

func createApp(mainFn mainFn) *cli.App {
	app := cli.NewApp()
	appInput := &appInput{}

	app.Name = "gprobe"
	app.Usage = "universal gRPC health-checker. See https://github.com/grpc/grpc/blob/master/doc/health-checking.md"
	app.UsageText = "gprobe [options] server_address [service_name]"
	app.Version = version
	app.HideHelp = true
	app.OnUsageError = func(context *cli.Context, err error, isSubcommand bool) error {
		cli.ShowAppHelp(context)
		return cli.NewExitError(err.Error(), ExitCodeUsage)
	}
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Fprintf(c.App.Writer, "%s\n", c.App.Version)
	}
	app.Flags = []cli.Flag{
		cli.DurationFlag{
			Name:        "timeout, t",
			Usage:       "Operation timeout",
			Destination: &appInput.timeout,
			Value:       1 * time.Second,
		},
		cli.BoolFlag{
			Name:        "no-fail, n",
			Usage:       "Do not fail if service status is other than SERVING. Note: this has no effect on server check",
			Destination: &appInput.noFail,
		},
	}
	app.Action = func(c *cli.Context) error {
		switch c.NArg() {
		case 2:
			appInput.serviceName = c.Args().Get(1)
			appInput.serverAddress = c.Args().Get(0)
			break
		case 1:
			appInput.serverAddress = c.Args().Get(0)
			break
		default:
			return c.App.OnUsageError(c, fmt.Errorf("exactly 1 to 2 arguments are required"), false)
		}
		// Pass all input to mainFn
		return mainFn(appInput)
	}

	return app
}

func main() {
	createApp(appMain).Run(os.Args)
}

func appMain(appInput *appInput) *cli.ExitError {
	probe.OpTimeout = appInput.timeout
	err := probe.Connect(appInput.serverAddress)
	if err != nil {
		return cli.NewExitError(err.Error(), ExitCodeUnexpected)
	}
	defer probe.Disconnect()

	status, err := probe.CheckService(appInput.serviceName)
	if err != nil {
		return cli.NewExitError(err.Error(), ExitCodeUnexpected)
	}

	fmt.Fprintln(os.Stdout, status.String())
	if !(appInput.noFail || status == grpc_health_v1.HealthCheckResponse_SERVING) {
		return cli.NewExitError("health-check failed", ExitCodeHealthCheckNegative)
	}

	// for some reason returning nil here causes err == nil to be false in urfave/cli/errors.go:79
	return cli.NewExitError("", 0)
}
