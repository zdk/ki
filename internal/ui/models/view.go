package models

// ViewMode represents different views in the application
type ViewMode int

const (
	MainMenuView ViewMode = iota
	ClusterListView
	ClusterDetailView
	NodeListView
	CreateClusterView
	LoadImageView
	BuildImageView
	ExportLogsView
	DeleteConfirmView
)