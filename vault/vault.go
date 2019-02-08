package vaulthelpers

import (
	"bufio"
	"strings"
	"testing"

	"go-test-helpers/aws"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/require"
)

// InitializeVaultClient creates a new AWS session to get information from bucket
// in order to be able to create a vault 'test' client
func InitializeVaultClient(t *testing.T, s *session.Session, tmpDir, profile, region, bucket, username, address string) *api.Client {
	config := api.Config{Address: address}
	client, err := api.NewClient(&config)
	require.NoErrorf(t, err, "Error getting Vault client: %s", err)
	return client
}

// GetInitialRootToken get the Initial Root Token from the S3 bucket & kmsid
func GetInitialRootToken(t *testing.T, s *session.Session, tmpDir, profile, region, bucket, username string) string {

	secretFilename := username + "-secrets.encrypted"

	//FIXME the DecryptBlob function returns InvalidCiphertextException - not sure why
	// // download the secret file into a byte array
	// encryptedSecretsr := awsbucket.DownloadItem(sess, bucket, secretFilename)
	// if err != nil {
	// 	t.Fatalf("Error downloading secrets: %s", err)
	// }
	// fmt.Printf("encryptedSecrets len %d\n", len(encryptedSecrets)) //REMOVE ME !!!!!!!!!!!!!!

	// // decrypt the secrets from byte array
	// decryptedSecrets := awsbucket.DecryptBlob(sess, encryptedSecrets)

	// download file for debugging - REMOVE THIS !!!!!!!!!!!!!
	var tmpFilename = tmpDir + "/" + secretFilename
	aws.S3DownloadFile(t, s, bucket, secretFilename, tmpFilename)

	// decrypt the file
	//aws-encryption-cli --decrypt --input /tmp/dcsaki-secrets.encrypted --master-keys profile=cloudinfra-dev --metadata-output ~/metadata --output /tmp/dcsaki-secrets.decrypted
	decryptedSecrets := aws.KmsDecryptFile(t, tmpFilename, profile)

	// parse decrypted secrets and extract initial root token
	//   $ cat /tmp/secrets
	//   Recovery Key 1: EW0ImBFnWSu0Lk3qvB0X2aDjkuZxRoMHAkKy9qTCkWg=
	//   Initial Root Token: s.b3exM9kZ1vwq5QPUeKcqdFiQ
	//   Success! Vault is initialized
	//   Recovery key initialized with 1 key shares and a key threshold of 1. Please
	//   securely distribute the key shares printed above.
	reader := strings.NewReader(decryptedSecrets)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Initial Root Token:") {
			tokens := strings.Split(line, " ")
			return tokens[len(tokens)-1]
		}
	}
	t.Fatalf("Initial Root Token not found")
	return ""
}

// GetKubeConfig retrieves the kubeconfig from the vault s3 bucket
// returns the temporary file where it is stored
func GetKubeConfig(t *testing.T, s *session.Session, tmpDir, profile, region, bucket string) string {
	var tmpFilename = tmpDir + "/kubeconfig"
	aws.S3DownloadFile(t, s, bucket, "kubeconfig", tmpFilename)
	return tmpFilename
}
