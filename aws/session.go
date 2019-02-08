package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/stretchr/testify/require"
)

// CreateAwsSessionWithDefaults creates an AWS that only uses defaults
func CreateAwsSessionWithDefaults(t *testing.T) *session.Session {
	sess, err := session.NewSession()
	require.NoErrorf(t, err, "Creating AWS Session with defaults failed: %s", err)
	return sess
}

// CreateAwsSessionWithOptions creates a session using provided options
func CreateAwsSessionWithOptions(t *testing.T, options *session.Options) *session.Session {
	sess, err := session.NewSessionWithOptions(*options)
	require.NoErrorf(t, err, "Creating AWS Session with options failed: %s", err)
	return sess
}

// CreateAwsSessionWithParameters creates a session with shared config state enabled
// and based on session info settings and/or defaults
func CreateAwsSessionWithParameters(t *testing.T, profile string, region string, sharedConfigStateEnabled bool) *session.Session {
	var options *session.Options
	var sharedConfigState session.SharedConfigState
	if sharedConfigStateEnabled {
		sharedConfigState = session.SharedConfigEnable
	} else {
		sharedConfigState = session.SharedConfigDisable
	}
	if len(region) > 0 {
		options = &session.Options{
			Config: aws.Config{
				Region: aws.String(region),
				//LogLevel: aws.LogLevel(aws.LogDebug),
			},
			SharedConfigState: sharedConfigState,
		}
	} else {
		options = &session.Options{
			SharedConfigState: sharedConfigState,
		}
	}

	if len(profile) > 0 {
		options.Profile = profile
	}

	return CreateAwsSessionWithOptions(t, options)
}
