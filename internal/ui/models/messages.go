package models

import "ki/internal/cmd"

// Message types
type (
	ClustersMsg      []cmd.Cluster
	NodesMsg         []cmd.Node
	ClusterDetailMsg cmd.Cluster
	MessageMsg       struct {
		Text    string
		MsgType string
	}
)