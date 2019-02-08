package aws

import (
	"bytes"
	"net/http"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/require"
)

// S3Exists checks if a bucket exists
func S3Exists(t *testing.T, sess *session.Session, bucket string) bool {
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	_, err := svc.ListObjects(params)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				t.Logf("No bucket named %s exists", bucket)
				return false
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			// fmt.Println(err.Error())
		}
		t.Fatalf("Error checking if %s bucket exists: %s", bucket, err)
	}
	return true
}

// S3CreateWithInput a bucket using CreateBucketInput
func S3CreateWithInput(t *testing.T, sess *session.Session, input *s3.CreateBucketInput, ignoreErrIfExists bool) *s3.CreateBucketOutput {
	svc := s3.New(sess)

	result, err := svc.CreateBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeBucketAlreadyExists, s3.ErrCodeBucketAlreadyOwnedByYou:
				if ignoreErrIfExists {
					t.Logf("%s bucket already exists", *input.Bucket)
					return result
				}
			}
			t.Fatalf("Error creating %s bucket: %s", *input.Bucket, aerr.Error())
		} else {
			t.Fatalf("Error creating %s bucket: %s", *input.Bucket, err)
		}
	}
	t.Logf("Created %s bucket", *input.Bucket)
	return result
}

// S3Create a bucket using only bucket name as input
func S3Create(t *testing.T, sess *session.Session, bucket string, ignoreErrIfExists bool) *s3.CreateBucketOutput {
	input := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}
	return S3CreateWithInput(t, sess, input, ignoreErrIfExists)
}

// S3Delete a bucket using bucket name
func S3Delete(t *testing.T, sess *session.Session, bucket string, ignoreNoSuchBucket bool) bool {
	svc := s3.New(sess)
	input := &s3.DeleteBucketInput{
		Bucket: aws.String(bucket),
	}
	_, err := svc.DeleteBucket(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == s3.ErrCodeNoSuchBucket && ignoreNoSuchBucket {
				return true
			}
		}
		return false
	}
	t.Logf("Deleted %s bucket", bucket)
	return true
}

// S3ListItems return a list of a given bucket's items
func S3ListItems(t *testing.T, sess *session.Session, bucket string) *s3.ListObjectsOutput {
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	result, err := svc.ListObjects(params)
	require.NoErrorf(t, err, "Error listing %s bucket items: %s", bucket, err)
	return result
}

// S3ItemExists check if item exists in bucket
func S3ItemExists(t *testing.T, sess *session.Session, bucket string, item string) bool {
	svc := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(item),
	}

	listObjectsOutput, err := svc.ListObjects(params)
	require.NoErrorf(t, err, "Error checking %s bucket item %s exists: %s", bucket, item, err)

	// search for item
	for _, object := range listObjectsOutput.Contents {
		if *object.Key == item {
			return true
		}
	}
	return false
}

// S3DownloadItem downloads an item from an S3 Bucket in the region configured in the shared config
// or AWS_REGION environment variable.
func S3DownloadItem(t *testing.T, sess *session.Session, bucket, item string) []byte {

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a buffer to hold the item
	buff := &aws.WriteAtBuffer{}

	// Write the contents of S3 Object to the file
	_, err := downloader.Download(buff, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	require.NoErrorf(t, err, "Failed to download item %s from %s bucket: %v", bucket, item, err)

	t.Logf("Successfully downloaded item %s from %s", item, bucket)
	return buff.Bytes()
}

// S3DownloadFile download bucket item into file
func S3DownloadFile(t *testing.T, sess *session.Session, bucket, item string, filename string) *os.File {
	t.Helper()

	// Create a downloader with the session and default options
	downloader := s3manager.NewDownloader(sess)

	// Create a file to write the S3 Object contents to.
	f, err := os.Create(filename)
	require.NoErrorf(t, err, "Failed to create local file %s: %s", filename, err)
	defer f.Close()

	// Write the contents of S3 Object to the file
	_, err = downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(item),
	})
	require.NoErrorf(t, err, "Failed to download file %s: %s", item, err)

	t.Logf("Successfully downloaded file %s to %s", item, filename)
	return f
}

// S3UploadItem uploads an item (byte array) to a bucket
func S3UploadItem(t *testing.T, sess *session.Session, bucket, item string, value []byte) {

	contentType := getBufferContentType(value)

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(item),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: bytes.NewReader(value),

		ContentType: aws.String(contentType),
	})
	require.NoErrorf(t, err, "Unable to upload %s to %s: %s", item, bucket, err)

	t.Logf("Successfully uploaded %s to %s", item, bucket)
}

// S3UploadFile uploads a file to a bucket item
func S3UploadFile(t *testing.T, sess *session.Session, bucket, item string, filename string) {

	f, err := os.Open(filename)
	require.NoErrorf(t, err, "Failed to open file %s: %s", filename, err)

	contentType, err := getContentType(f)
	require.NoErrorf(t, err, "Could get the content-type of %s: %s", item, err)

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	uploader := s3manager.NewUploader(sess)

	// Upload the file's body to S3 bucket as an object with the key being the
	// same as the filename.
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),

		// Can also use the `filepath` standard library package to modify the
		// filename as need for an S3 object key. Such as turning absolute path
		// to a relative path.
		Key: aws.String(item),

		// The file to be uploaded. io.ReadSeeker is preferred as the Uploader
		// will be able to optimize memory when uploading large content. io.Reader
		// is supported, but will require buffering of the reader's bytes for
		// each part.
		Body: f,

		ContentType: aws.String(contentType),
	})
	require.NoErrorf(t, err, "Unable to upload %s to %s: %s", item, bucket, err)

	t.Logf("Successfully uploaded %s to %s", item, bucket)
}

func getContentType(input *os.File) (string, error) {
	defer input.Seek(0, 0)

	buffer := make([]byte, 512)
	_, err := input.Read(buffer)
	if err != nil {
		return "", err
	}

	return getBufferContentType(buffer), nil
}

func getBufferContentType(buffer []byte) string {
	return http.DetectContentType(buffer)
}
