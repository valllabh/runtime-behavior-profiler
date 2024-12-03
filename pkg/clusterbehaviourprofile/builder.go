package clusterbehaviourprofile

import (
	clusterbehaviourprofiletype "runtime-behavior-profiler/pkg/clusterbehaviourprofile/type"
)

// profiles is a map of cluster names to ClusterBehaviourProfile objects.
// It is used to store the profiles of multiple clusters.
// The key is the cluster name and the value is the ClusterBehaviourProfile object.
var profiles = make(map[string]clusterbehaviourprofiletype.ClusterBehaviourProfile)

// GetClusterBehaviourProfile returns the ClusterBehaviourProfile object for the given cluster name.
func GetClusterBehaviourProfile(clusterName string) clusterbehaviourprofiletype.ClusterBehaviourProfile {

	// Check if the cluster profile already exists
	_, ok := profiles[clusterName]

	// If the cluster profile does not exist, create a new one
	if !ok {
		profiles[clusterName] = clusterbehaviourprofiletype.ClusterBehaviourProfile{
			Cluster:    clusterName,
			Namespaces: map[string]clusterbehaviourprofiletype.Namespace{},
		}
	}

	// Return the cluster profile
	return profiles[clusterName]
}
