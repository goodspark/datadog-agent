// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build secrets
// +build secrets

package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/go-secure-stdlib/awsutil"
	"github.com/pkg/errors"
)

// AwsCloudAuth is fetching a Datadog API Key for the requested org using an AWS based proof
func AwsCloudAuth() (value string, err error) {
	credsConfig := awsutil.CredentialsConfig{}
	creds, err := credsConfig.GenerateCredentialChain()
	if err != nil {
		return "", errors.Wrap(err, "unable to generate the credential chain")
	}
	data, err := generateAwsAuthData(creds, "us1.prod.dog", "us-east-1")

	if err != nil {
		return "", errors.Wrap(err, "unable to generate the AWS Auth Data")
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", errors.Wrap(err, "unable to marshal the data")
	}
	resp, err := http.Post("https://keyless-authentication-validation.us1.staging.dog/validate-proof", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.Wrap(err, "error while contacting the keyless authentication service")
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error when reading the response")
	}
	defer resp.Body.Close()

	returnedData := &ValidatorResponse{}
	err = json.Unmarshal(d, returnedData)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("unable to unmarshal the Keyless Authentication data from Datadog: %s", string(d)))
	}
	return returnedData.Key, nil
}
