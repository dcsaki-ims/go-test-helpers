package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/davecgh/go-spew/spew"
)

var profile = "default"
var region = "ca-central-1"

func TestCreateAwsSessionWithDefaults(t *testing.T) {
	session := CreateAwsSessionWithDefaults(t)
	t.Logf("Created AWS session: %s", spew.Sdump(session))
}

func TestCreateAwsSessionWithOptions(t *testing.T) {
	var options = session.Options{
		Config: aws.Config{
			Region: aws.String(region),
			//LogLevel: aws.LogLevel(aws.LogDebug),
		},
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	}
	session := CreateAwsSessionWithOptions(t, &options)
	t.Logf("Created AWS session: %s", spew.Sdump(session))
}

func TestCreateAwsSessionWithParameters(t *testing.T) {
	session := CreateAwsSessionWithParameters(t, profile, region, true)
	t.Logf("Created AWS session: %s", spew.Sdump(session))
}
