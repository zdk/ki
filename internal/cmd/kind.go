package cmd

import (
	"fmt"
	"os/exec"
	"strings"
)

// Cluster represents a KIND cluster
type Cluster struct {
	Name        string
	Status      string
	Nodes       []Node
	KubeVersion string
}

// Node represents a Kubernetes node
type Node struct {
	Name       string
	Role       string
	Status     string
	Age        string
	Version    string
	InternalIP string
}

// GetClusters retrieves all KIND clusters
func GetClusters() ([]Cluster, error) {
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get clusters: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	clusters := make([]Cluster, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cluster := Cluster{
				Name:   line,
				Status: "running",
			}
			cluster = EnrichClusterInfo(cluster)
			clusters = append(clusters, cluster)
		}
	}

	return clusters, nil
}

// EnrichClusterInfo adds additional information to a cluster
func EnrichClusterInfo(c Cluster) Cluster {
	contextName := "kind-" + c.Name
	if c.Name == "kind" {
		contextName = "kind"
	}

	cmd := exec.Command("kubectl", "get", "nodes", "--context", contextName, "-o", "wide", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return c
	}

	nodes := ParseNodes(string(output))
	c.Nodes = nodes

	if len(nodes) > 0 {
		c.KubeVersion = nodes[0].Version
	}

	return c
}

// ParseNodes parses kubectl node output
func ParseNodes(output string) []Node {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	nodes := make([]Node, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) >= 6 {
			role := "worker"
			if strings.Contains(line, "control-plane") || strings.Contains(line, "master") {
				role = "control-plane"
			}

			node := Node{
				Name:       fields[0],
				Status:     fields[1],
				Role:       role,
				Age:        fields[3],
				Version:    fields[4],
				InternalIP: fields[5],
			}
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// GetClusterNodes retrieves nodes for a specific cluster
func GetClusterNodes(clusterName string) ([]Node, error) {
	contextName := "kind-" + clusterName
	if clusterName == "kind" {
		contextName = "kind"
	}

	cmd := exec.Command("kubectl", "get", "nodes", "--context", contextName, "-o", "wide", "--no-headers")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	return ParseNodes(string(output)), nil
}

// GetClusterDetail retrieves detailed information about a cluster
func GetClusterDetail(clusterName string) (Cluster, error) {
	cluster := Cluster{Name: clusterName, Status: "running"}
	cluster = EnrichClusterInfo(cluster)
	return cluster, nil
}

// CreateCluster creates a new KIND cluster
func CreateCluster(name string) error {
	args := []string{"create", "cluster"}
	if name != "" {
		args = append(args, "--name", name)
	}

	cmd := exec.Command("kind", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create cluster: %w\n%s", err, string(output))
	}

	return nil
}

// DeleteCluster deletes a KIND cluster
func DeleteCluster(name string) error {
	cmd := exec.Command("kind", "delete", "cluster", "--name", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete cluster: %w\n%s", err, string(output))
	}

	return nil
}

// LoadDockerImage loads a Docker image into a KIND cluster
func LoadDockerImage(imageName, clusterName string) error {
	args := []string{"load", "docker-image", imageName}
	if clusterName != "" {
		args = append(args, "--name", clusterName)
	}

	cmd := exec.Command("kind", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to load image: %w\n%s", err, string(output))
	}

	return nil
}

// BuildNodeImage builds a KIND node image from source
func BuildNodeImage(sourcePath string) error {
	args := []string{"build", "node-image"}
	if sourcePath != "" {
		args = append(args, sourcePath)
	}

	cmd := exec.Command("kind", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build node image: %w\n%s", err, string(output))
	}

	return nil
}

// ExportLogs exports cluster logs
func ExportLogs(clusterName, outputPath string) error {
	args := []string{"export", "logs"}
	if clusterName != "" {
		args = append(args, "--name", clusterName)
	}
	if outputPath != "" {
		args = append(args, outputPath)
	}

	cmd := exec.Command("kind", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to export logs: %w\n%s", err, string(output))
	}

	return nil
}