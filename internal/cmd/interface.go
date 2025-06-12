package cmd

// CommandInterface defines the interface for KIND operations
// This allows for easy testing by providing mock implementations
type CommandInterface interface {
	GetClusters() ([]Cluster, error)
	GetClusterNodes(clusterName string) ([]Node, error)
	GetClusterDetail(clusterName string) (Cluster, error)
	CreateCluster(name string) error
	DeleteCluster(name string) error
	LoadDockerImage(imageName, clusterName string) error
	BuildNodeImage(sourcePath string) error
	ExportLogs(clusterName, outputPath string) error
}

// DefaultCommands implements CommandInterface using the actual KIND commands
type DefaultCommands struct{}

func (d DefaultCommands) GetClusters() ([]Cluster, error) {
	return GetClusters()
}

func (d DefaultCommands) GetClusterNodes(clusterName string) ([]Node, error) {
	return GetClusterNodes(clusterName)
}

func (d DefaultCommands) GetClusterDetail(clusterName string) (Cluster, error) {
	return GetClusterDetail(clusterName)
}

func (d DefaultCommands) CreateCluster(name string) error {
	return CreateCluster(name)
}

func (d DefaultCommands) DeleteCluster(name string) error {
	return DeleteCluster(name)
}

func (d DefaultCommands) LoadDockerImage(imageName, clusterName string) error {
	return LoadDockerImage(imageName, clusterName)
}

func (d DefaultCommands) BuildNodeImage(sourcePath string) error {
	return BuildNodeImage(sourcePath)
}

func (d DefaultCommands) ExportLogs(clusterName, outputPath string) error {
	return ExportLogs(clusterName, outputPath)
}

// Global instance that can be replaced for testing
var Commands CommandInterface = DefaultCommands{}