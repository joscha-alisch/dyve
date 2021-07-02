package cloudfoundry

import "time"

type ReconcileType int

const (
	ReconcileOrg ReconcileType = iota
	ReconcileSpace
	ReconcileApp
)

type ReconcileJob struct {
	Type        ReconcileType
	Guid        string
	LastUpdated time.Time
}

type Org struct {
	Guid string
	Name string
}

type Space struct {
	Guid string
	Name string
}

type App struct {
	Guid string
	Name string
}
