package sdk

type AppProvider interface {
	ListApps() ([]App, error)
	GetApp(id string) (App, error)
	Search(term string, limit int) ([]AppSearchResult, error)
}
