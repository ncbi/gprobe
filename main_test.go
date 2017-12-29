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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/urfave/cli"
	"time"
)

func Test_createConfig_args_narg1(t *testing.T) {
	// given
	args := cli.Args{"server"}
	flags := &appFlags{}

	// when
	config, err := createConfig(flags, args)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "server", config.serverAddress)
	assert.Empty(t, config.serviceName)
}

func Test_createConfig_args_narg2(t *testing.T) {
	// given
	args := cli.Args{"server", "svc"}
	flags := &appFlags{}

	// when
	config, err := createConfig(flags, args)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "server", config.serverAddress)
	assert.Equal(t, "svc", config.serviceName)
}

func Test_createConfig_args_narg3(t *testing.T) {
	// given
	args := cli.Args{"foo", "bar", "baz"}
	flags := &appFlags{}

	// when
	_, err := createConfig(flags, args)

	// then
	assert.Error(t, err)
}

func Test_createConfig_args_narg0(t *testing.T) {
	// given
	args := cli.Args{}
	flags := &appFlags{}

	// when
	_, err := createConfig(flags, args)

	// then
	assert.Error(t, err)
}

func Test_createConfig_flags_empty(t *testing.T) {
	// given
	args := cli.Args{"foo"}
	flags := &appFlags{}

	// when
	config, err := createConfig(flags, args)

	// then
	assert.NoError(t, err)
	assert.False(t, config.noFail)
}

func Test_createConfig_flags(t *testing.T) {
	// given
	args := cli.Args{"foo"}
	flags := &appFlags{
		noFail:  true,
		timeout: time.Minute,
	}

	// when
	config, err := createConfig(flags, args)

	// then
	assert.NoError(t, err)
	assert.Equal(t, time.Minute, config.timeout)
	assert.True(t, config.noFail)
}
