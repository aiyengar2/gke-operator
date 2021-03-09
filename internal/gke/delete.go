package gke

import (
	"context"
	"strings"
	"time"

	gkev1 "github.com/rancher/gke-operator/pkg/apis/gke.cattle.io/v1"
	gkeapi "google.golang.org/api/container/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	waitSec      = 30
	backoffSteps = 12
)

// RemoveCluster attempts to delete a cluster and retries the delete request if the cluster is busy.
func RemoveCluster(ctx context.Context, client *gkeapi.Service, config *gkev1.GKEClusterConfig) error {
	backoff := wait.Backoff{
		Duration: waitSec * time.Second,
		Steps:    backoffSteps,
	}
	return wait.ExponentialBackoff(backoff, func() (bool, error) {
		_, err := client.Projects.
			Locations.
			Clusters.
			Delete(ClusterRRN(config.Spec.ProjectID, config.Spec.Region, config.Spec.ClusterName)).
			Context(ctx).
			Do()
		if err != nil && strings.Contains(err.Error(), "Please wait and try again once it is done") {
			return false, nil
		}
		if err != nil && strings.Contains(err.Error(), "notFound") {
			return true, nil
		}
		if err != nil {
			return false, err
		}
		return true, nil
	})
}
