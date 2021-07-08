package sdk

type AppProvider interface {
	ListApps() ([]App, error)
	ListAppsPaged(perPage int, page int) (AppPage, error)
	GetApp(id string) (App, error)
	Search(term string, limit int) ([]AppSearchResult, error)
}
