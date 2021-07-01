package cloudfoundry

type Database interface {
	FetchReconcileJob() *ReconcileJob
	UpdateOrg(guid string, o Org) error
	UpdateSpace(guid string, s Space) error
	UpdateApp(guid string, a App) error
}

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
