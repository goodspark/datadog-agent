// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package processes

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/DataDog/datadog-agent/test/new-e2e/utils/clients"
	"github.com/DataDog/datadog-agent/test/new-e2e/utils/credentials"
	"github.com/DataDog/datadog-agent/test/new-e2e/utils/infra"
	"github.com/DataDog/test-infra-definitions/aws/ec2/ec2"

	"github.com/pulumi/pulumi/sdk/v3/go/auto"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

const initScript = `#!/bin/bash

set -ex

export DEBIAN_FRONTEND=noninteractive

apt -y update && apt -y install docker.io

sudo chmod o+rw /var/run/docker.sock
`

// EC2TestEnv provides an ec2 test environment
type EC2TestEnv struct {
	ddAPIKey   string
	instanceIP string

	sshClient   *ssh.Client
	stackOutput auto.UpResult
}

// NewEC2TestEnv creates an EC2TestEnv in the aws sandbox env
func NewEC2TestEnv(name string) (*EC2TestEnv, error) {
	ec2TestEnv := &EC2TestEnv{}
	credentialsManager := credentials.NewManager()

	// Retrieving necessary secrets
	sshKey, err := credentialsManager.GetCredential(credentials.AWSSSMStore, "agent.ci.awssandbox.ssh")
	if err != nil {
		return nil, err
	}

	if ec2TestEnv.ddAPIKey, err = credentialsManager.GetCredential(credentials.AWSSSMStore, "agent.ci.dev.apikey"); err != nil {
		return nil, err
	}

	ddAPPKey, err := credentialsManager.GetCredential(credentials.AWSSSMStore, "agent.ci.dev.appkey")
	if err != nil {
		return nil, err
	}

	stackOutput, err := infra.GetStackManager().GetStack(context.Background(), "aws/sandbox", fmt.Sprintf("process-agent-%s", name), nil, func(ctx *pulumi.Context) error {
		instance, err := ec2.CreateEC2Instance(ctx, fmt.Sprintf("process-agent-%s", name), "", ec2.AMD64Arch, "t3.large",
			"agent-ci-sandbox", initScript)
		if err != nil {
			return err
		}

		ctx.Export("private-ip", instance.PrivateIp)
		return nil
	})
	if err != nil {
		return nil, err
	}

	ec2TestEnv.stackOutput = stackOutput
	output, found := stackOutput.Outputs["private-ip"]
	if !found {
		return nil, errors.New("unable to find the ec2 host ip")
	}
	ec2TestEnv.instanceIP = output.Value.(string)

	if ec2TestEnv.sshClient, _, err = clients.GetSSHClient("ubuntu", fmt.Sprintf("%s:%d", ec2TestEnv.instanceIP, 22),
		sshKey, 2*time.Second, 30); err != nil {
		return nil, err
	}

	if err = os.Setenv("DD_API_KEY", ec2TestEnv.ddAPIKey); err != nil {
		return nil, err
	}
	if err = os.Setenv("DD_APP_KEY", ddAPPKey); err != nil {
		return nil, err
	}

	return ec2TestEnv, nil
}

// Close performs cleanup and destroys the ec2 stack
func (e *EC2TestEnv) Close() {
	_ = os.Unsetenv("DD_API_KEY")
	_ = os.Unsetenv("DD_APP_KEY")
	_ = e.sshClient.Close()
}

func createHostName(testName string) string {
	sl := strings.Split(testName, "/")
	hostName := fmt.Sprintf("%s-%d", sl[len(sl)-1], time.Now().UnixMilli())
	return hostName
}
