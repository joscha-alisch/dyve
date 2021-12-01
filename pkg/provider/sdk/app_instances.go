package sdk

import "time"

type InstancesProvider interface {
	GetAppInstances(id string) (AppInstances, error)
}

type AppInstances []AppInstance
type AppInstance struct {
	State AppState
	Since time.Time
}

type AppState string

const (
	AppStateRunning  AppState = "running"
	AppStateStopped  AppState = "stopped"
	AppStateCrashed  AppState = "crashed"
	AppStateStarting AppState = "starting"
	AppStateUnknown  AppState = "unknown"
)
