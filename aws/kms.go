package aws

import (
	"os/exec"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	kms_client "github.com/aws/aws-sdk-go/service/kms"
	"github.com/stretchr/testify/require"
)

// KmsDecryptFile decrypt file
// - this function is using aws-encryption-cli which must already be installed
// - eg:
//       aws-encryption-cli --decrypt --input /tmp/dcsaki-secrets.encrypted --master-keys profile=cloudinfra-dev --metadata-output ~/metadata --output /tmp/dcsaki-secrets.decrypted
func KmsDecryptFile(t *testing.T, filename string, profile string) string {

	_, err := exec.LookPath("aws-encryption-cli")
	if err != nil {
		t.Fatalf("aws-encryption-cli exec not found: %s", err)
	}

	cmd := exec.Command("aws-encryption-cli", "--decrypt", "--input", filename, "--master-keys",
		"profile="+profile, "--metadata-output", "/tmp/metadata", "--output", "-")
	out, err := cmd.CombinedOutput()
	require.NoErrorf(t, err, "Error executing aws-encryption-cli command: %s", out)
	return string(out)
}

// KmsDecryptBlob decrypt a byte array
func KmsDecryptBlob(t *testing.T, sess *session.Session, data []byte) string {
	// Create KMS service client
	svc := kms_client.New(sess)

	/// create EncryptionContext
	var purpose = "dcsaki"
	var encryptionContext map[string]*string
	encryptionContext = make(map[string]*string)
	encryptionContext["purpose"] = &purpose

	// Decrypt the data
	result, err := svc.Decrypt(&kms.DecryptInput{
		CiphertextBlob:    data,
		EncryptionContext: encryptionContext,
	})

	require.NoErrorf(t, err, "Error decrypting data: %s", err)
	return string(result.Plaintext)
}
