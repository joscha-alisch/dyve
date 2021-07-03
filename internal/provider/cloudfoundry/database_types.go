package cloudfoundry

import "time"

type ReconcileType int

const (
	ReconcileOrg ReconcileType = iota
	ReconcileSpace
)

type ReconcileJob struct {
	Type ReconcileType
	Guid string
}

type Org struct {
	Guid        string
	Name        string
	Spaces      []string
	LastUpdated time.Time `bson:"lastUpdated"`
}

type Space struct {
	Guid        string
	Name        string
	Apps        []string
	LastUpdated time.Time `bson:"lastUpdated"`
}

type App struct {
	Guid  string
	Name  string
	Org   string
	Space string
}
