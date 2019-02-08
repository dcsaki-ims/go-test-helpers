package aws

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
)

func initSession(t *testing.T) *session.Session {
	var profile = "default"
	var region = "ca-central-1"
	return CreateAwsSessionWithParameters(t, profile, region, true)
}

func TestBucketExists(t *testing.T) {
	bucket := "test.dcsaki"
	sess := initSession(t)
	exists := S3Exists(t, sess, bucket)
	t.Log(exists)
}

func TestCreateBucket(t *testing.T) {
	bucket := "test2.dcsaki"
	sess := initSession(t)
	result := S3Create(t, sess, bucket, false)
	t.Log(result)
}

func TestDeleteBucket(t *testing.T) {
	bucket := "test2.dcsaki"
	sess := initSession(t)
	result := S3Delete(t, sess, bucket, false)
	t.Log(result)
}

func TestListBucketItems(t *testing.T) {
	bucket := "test.dcsaki"
	sess := initSession(t)
	items := S3ListItems(t, sess, bucket)
	t.Log(items)
}

func TestBucketItemExists(t *testing.T) {
	bucket := "test.dcsaki"
	item := "test_item_1xx"
	sess := initSession(t)
	exists := S3ItemExists(t, sess, bucket, item)
	t.Log(exists)
}

func TestUploadItem(t *testing.T) {
	bucket := "test.dcsaki"
	item := "test_item_1"
	s := "Hello"
	value := []byte(s)
	sess := initSession(t)
	S3UploadItem(t, sess, bucket, item, value)
}

func TestUploadFile(t *testing.T) {
	bucket := "test.dcsaki"
	item := "test_item_2"
	filename := "/tmp/test_item_1.txt"
	sess := initSession(t)
	S3UploadFile(t, sess, bucket, item, filename)
}

func TestDownloadItem(t *testing.T) {
	bucket := "test.dcsaki"
	item := "test_item_1"
	sess := initSession(t)
	S3DownloadItem(t, sess, bucket, item)
}

func TestDownloadFile(t *testing.T) {
	bucket := "test.dcsaki"
	item := "test_item_1"
	filename := "/tmp/test_item_1.txt"
	sess := initSession(t)
	file := S3DownloadFile(t, sess, bucket, item, filename)
	s := string(file.Name())
	t.Log(s)
}
