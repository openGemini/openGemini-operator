package utils

import (
	"os"
	"testing"
)

func TestGetClusterDomain(t *testing.T) {
	// Backup current environment variable
	old := os.Getenv("KUBERNETES_CLUSTER_DOMAIN")

	// Unset it for the first part of our test
	os.Unsetenv("KUBERNETES_CLUSTER_DOMAIN")
	domain := GetClusterDomain()
	if domain != "cluster.local" {
		t.Errorf("Expected default domain 'cluster.local', got '%s'", domain)
	}

	// Set it to a known value for the second part of our test
	os.Setenv("KUBERNETES_CLUSTER_DOMAIN", "mycluster.local")
	domain = GetClusterDomain()
	if domain != "mycluster.local" {
		t.Errorf("Expected domain 'mycluster.local', got '%s'", domain)
	}

	// Restore environment variable
	os.Setenv("KUBERNETES_CLUSTER_DOMAIN", old)
}
