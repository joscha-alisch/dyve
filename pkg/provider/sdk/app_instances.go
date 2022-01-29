package sdk

import "time"

type InstancesProvider interface {
	GetAppInstances(id string) (AppInstances, error)
}

type AppInstances []AppInstance
type AppInstance struct {
	State AppState  `json:"state"`
	Since time.Time `json:"since"`
}

type AppState string

const (
	AppStateRunning  AppState = "running"
	AppStateStopped  AppState = "stopped"
	AppStateCrashed  AppState = "crashed"
	AppStateStarting AppState = "starting"
	AppStateUnknown  AppState = "unknown"
)
