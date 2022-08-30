// Copyright Hashicorp, Datadog, Inc.
// LICENSE: Mozilla Public License, version 2.0 https://github.com/hashicorp/go-secure-stdlib/blob/7849be51188ffe09900bf3232e94695389cfc8aa/awsutil/LICENSE
// CHANGES:
// 	- Update of the IamServerIDHeader
//	- Use of structures rather than a map[string]interface{} for handling the data
//  - Handle potential error returned by the signature procedure

//go:build secrets
// +build secrets

package providers

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/hashicorp/go-secure-stdlib/awsutil"
)

// SigningData is the data structure that represents the Data used to generate and AWS Proof
type SigningData struct {
	HeadersEncoded string `json:"iam_headers_encoded"`
	BodyEncoded    string `json:"iam_body_encoded"`
	URLEncoded     string `json:"iam_url_encoded"`
	Method         string `json:"iam_method"`
}

const (
	// IamServerIDHeader is the header we use to make sure the request was intended to the correct DC
	IamServerIDHeader = "X-DDOG-AWS-IAM-Server-ID"
	defaultRegion     = "us-east-1"
)

// GetCallerIdentityResponse is used to parse the response of the GetCallerIdentity API call by AWS
type GetCallerIdentityResponse struct {
	GetCallerIdentityResult []GetCallerIdentityResult `xml:"GetCallerIdentityResult"`
}

// GetCallerIdentityResult is used to parse the result of the GetCallerIdentity call by AWS
type GetCallerIdentityResult struct {
	Arn     string `xml:"Arn"`
	UserID  string `xml:"UserId"`
	Account string `xml:"Account"`
}

// OrgValidationConfig is used to store the configuration of a Datadog org about which proofs it can accept
type OrgValidationConfig struct {
	AuthorizedAccounts     []string
	AuthorizedArnsPatterns []string
}

// ValidatorResponse is the response structure used by the validator
type ValidatorResponse struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
}

func generateAwsAuthData(creds *credentials.Credentials, serverID, configuredRegion string) (*SigningData, error) {
	// This method follows the AWS Auth method as used and advertised by Vault
	// https://github.com/hashicorp/go-secure-stdlib/blob/bf6d78ef5b727b83b2f8f41b473eb9764b446df4/awsutil/generate_credentials.go
	// https://www.vaultproject.io/docs/auth/aws
	//a
	region, err := awsutil.GetRegion(configuredRegion)
	if err != nil {
		region = defaultRegion
	}
	stsSession, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials:      creds,
			Region:           &region,
			EndpointResolver: endpoints.ResolverFunc(stsSigningResolver),
		},
	})
	if err != nil {
		return nil, err
	}

	var params *sts.GetCallerIdentityInput
	svc := sts.New(stsSession)
	stsRequest, _ := svc.GetCallerIdentityRequest(params)

	// Inject the required auth header value, if supplied, and then sign the request including that header
	if serverID != "" {
		stsRequest.HTTPRequest.Header.Add(IamServerIDHeader, serverID)
	}
	err = stsRequest.Sign()
	if err != nil {
		return nil, err
	}

	headersJSON, err := json.Marshal(stsRequest.HTTPRequest.Header)
	if err != nil {
		return nil, err
	}
	requestBody, err := ioutil.ReadAll(stsRequest.HTTPRequest.Body)
	if err != nil {
		return nil, err
	}

	return &SigningData{
		HeadersEncoded: base64.StdEncoding.EncodeToString(headersJSON),
		BodyEncoded:    base64.StdEncoding.EncodeToString(requestBody),
		Method:         stsRequest.HTTPRequest.Method,
		URLEncoded:     base64.StdEncoding.EncodeToString([]byte(stsRequest.HTTPRequest.URL.String())),
	}, nil
}

// From https://github.com/hashicorp/go-secure-stdlib/blob/bf6d78ef5b727b83b2f8f41b473eb9764b446df4/awsutil/generate_credentials.go#L303
// STS is a really weird service that used to only have global endpoints but now has regional endpoints as well.
// For backwards compatibility, even if you request a region other than us-east-1, it'll still sign for us-east-1.
// See, e.g., https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_temp_enable-regions.html#id_credentials_temp_enable-regions_writing_code
// So we have to shim in this EndpointResolver to force it to sign for the right region
func stsSigningResolver(service, region string, optFns ...func(*endpoints.Options)) (endpoints.ResolvedEndpoint, error) {
	defaultEndpoint, err := endpoints.DefaultResolver().EndpointFor(service, region, optFns...)
	if err != nil {
		return defaultEndpoint, err
	}
	defaultEndpoint.SigningRegion = region
	return defaultEndpoint, nil
}
