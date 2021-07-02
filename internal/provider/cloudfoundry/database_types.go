package cloudfoundry

type ReconcileType int

const (
	ReconcileOrg ReconcileType = iota
	ReconcileSpace
	ReconcileApp
)

type ReconcileJob struct {
	Type ReconcileType
	Guid string
}

type Org struct {
	Name string
}

type Space struct {
	Name string
}

type App struct {
	Name string
}
