package pylon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPylonComponentArgsValidate(t *testing.T) {
	// Test with valid args
	validArgs := PylonComponentArgs{
		Name:                "test-pylon",
		Namespace:           "default",
		ExecutionJwt:        "test-jwt",
		PylonImage:          "test-image:latest",
		PylonBlobBucketName: "test-bucket",
		Env: PylonEnv{
			PylonStartBlock:            "1000",
			PylonS3Url:                 "https://s3.example.com",
			PylonS3Region:              "us-west-2",
			PylonSenderAddress:         "0x1234567890123456789012345678901234567890",
			PylonNetworkSlotDuration:   "12",
			PylonNetworkSlotOffset:     "0",
			PylonRequestsPerSecond:     "100",
			PylonPort:                  "8080",
			AwsAccessKeyId:             "test-access-key",
			AwsSecretAccessKey:         "test-secret-key",
			AwsRegion:                  "us-west-2",
			PylonDbUrl:                 "postgresql://test:test@localhost:5432/test",
			PylonConsensusClientUrl:    "http://consensus:5052",
			PylonBlobscanBaseUrl:       "http://blobscan:3000",
			PylonNetworkStartTimestamp: "1234567890",
		},
	}

	err := validArgs.Validate()
	assert.NoError(t, err)

	// Test with missing name
	invalidArgs1 := PylonComponentArgs{
		Namespace:           "default",
		ExecutionJwt:        "test-jwt",
		PylonImage:          "test-image:latest",
		PylonBlobBucketName: "test-bucket",
		Env:                 validArgs.Env,
	}

	err = invalidArgs1.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")

	// Test with missing namespace
	invalidArgs2 := PylonComponentArgs{
		Name:                "test-pylon",
		ExecutionJwt:        "test-jwt",
		PylonImage:          "test-image:latest",
		PylonBlobBucketName: "test-bucket",
		Env:                 validArgs.Env,
	}

	err = invalidArgs2.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "namespace is required")

	// Test with missing executionJwt
	invalidArgs3 := PylonComponentArgs{
		Name:                "test-pylon",
		Namespace:           "default",
		PylonImage:          "test-image:latest",
		PylonBlobBucketName: "test-bucket",
		Env:                 validArgs.Env,
	}

	err = invalidArgs3.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "executionJwt is required")

	// Test with missing pylonImage
	invalidArgs4 := PylonComponentArgs{
		Name:                "test-pylon",
		Namespace:           "default",
		ExecutionJwt:        "test-jwt",
		PylonBlobBucketName: "test-bucket",
		Env:                 validArgs.Env,
	}

	err = invalidArgs4.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pylonImage is required")

	// Test with missing pylonBlobBucketName
	invalidArgs5 := PylonComponentArgs{
		Name:         "test-pylon",
		Namespace:    "default",
		ExecutionJwt: "test-jwt",
		PylonImage:   "test-image:latest",
		Env:          validArgs.Env,
	}

	err = invalidArgs5.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pylonBlobBucketName is required")
}
