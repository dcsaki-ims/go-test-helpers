package k8s

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"

	tf_k8s "github.com/gruntwork-io/terratest/modules/k8s"

	corev1 "k8s.io/api/core/v1"
)

// CheckNodesConditions checks the conditions of all nodes
// the k8s client must already be initialized for this to work
func CheckNodesConditions(t *testing.T) {

	// check all nodes in k8s cluster are 'ready'
	require.Truef(t, tf_k8s.AreAllNodesReady(t), "K8s nodes not ready")

	// get a list of all the nodes
	nodes := tf_k8s.GetNodes(t)
	require.Truef(t, len(nodes) > 0, "No nodes found in K8s cluster")

	// for each node, check its conditions
	for _, node := range nodes {
		t.Logf("Found node %s with status %s", node.Name, spew.Sdump(node.Status.Conditions)) // spew.Sdump(node))
		for _, condition := range node.Status.Conditions {
			switch condition.Type {
			case corev1.NodeReady:
				require.Truef(t, condition.Status == corev1.ConditionTrue, "K8s Node %s is not ready", node.Name)
			default:
				require.Truef(t, condition.Status == corev1.ConditionFalse, "K8s node %s has condition %s = True", node.Name, condition.Type)
			}
		}
	}

}

// GetOptions contructs an options instance from config file
func GetOptions(configPath string) *tf_k8s.KubectlOptions {
	options := tf_k8s.NewKubectlOptions("", configPath)
	return options
}
