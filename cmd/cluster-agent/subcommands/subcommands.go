// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !windows && kubeapiserver
// +build !windows,kubeapiserver

package subcommands

import (
	"github.com/DataDog/datadog-agent/cmd/cluster-agent/command"
	cmdcheck "github.com/DataDog/datadog-agent/cmd/cluster-agent/subcommands/check"
	cmdconfig "github.com/DataDog/datadog-agent/cmd/cluster-agent/subcommands/config"
	cmdconfigcheck "github.com/DataDog/datadog-agent/cmd/cluster-agent/subcommands/configcheck"
	cmdhealth "github.com/DataDog/datadog-agent/cmd/cluster-agent/subcommands/health"
	cmdstart "github.com/DataDog/datadog-agent/cmd/cluster-agent/subcommands/start"
	cmdversion "github.com/DataDog/datadog-agent/cmd/cluster-agent/subcommands/version"
)

// ClusterAgentSubcommands returns SubcommandFactories for the subcommands
// supported with the current build flags.
func ClusterAgentSubcommands() []command.SubcommandFactory {
	return []command.SubcommandFactory{
		cmdstart.Commands,
		cmdversion.Commands,
		cmdcheck.Commands,
		cmdconfig.Commands,
		cmdhealth.Commands,
		cmdconfigcheck.Commands,
		// config_check.go
		// flare.go

		// clusterchecks.go
		// diagnose.go
		// dummy.go
		// health.go
		// metadata_mapper_digest.go
		// secret_helper.go
		// status.go
		// telemetry.go
	}
}
