package aws

import (
	"testing"

	aws_api "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/stretchr/testify/require"
)

// GetEksClusterInformation get info that describes the EKS cluster
func GetEksClusterInformation(t *testing.T, s *session.Session, clusterName string) *eks.DescribeClusterOutput {
	svc := eks.New(s)
	input := &eks.DescribeClusterInput{
		Name: aws_api.String(clusterName),
	}
	result, err := svc.DescribeCluster(input)

	require.NoErrorf(t, err, "Error getting info for cluster %s: %s", clusterName, err)
	return result
}

// IsEksClusterActive returns true is EKS Cluster status of active
func IsEksClusterActive(t *testing.T, s *session.Session, clusterName string) bool {
	info := GetEksClusterInformation(t, s, clusterName)
	return *info.Cluster.Status == eks.ClusterStatusActive
}
