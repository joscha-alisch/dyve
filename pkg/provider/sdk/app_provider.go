package sdk

type AppProvider interface {
	ListApps() ([]App, error)
	GetApp(id string) (App, error)
}
