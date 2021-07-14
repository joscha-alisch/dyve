package sdk

type AppProvider interface {
	ListApps() ([]App, error)
	GetApp(id string) (App, error)
}

type AppPage struct {
	TotalResults int   `json:"totalResults"`
	TotalPages   int   `json:"totalPages"`
	PerPage      int   `json:"perPage"`
	Page         int   `json:"page"`
	Apps         []App `json:"apps"`
}

type App struct {
	Id   string                 `json:"id"`
	Name string                 `json:"name"`
	Meta map[string]interface{} `json:"meta,omitempty" bson:"meta,omitempty"`
}

type AppSearchResult struct {
	App   App     `json:"app"`
	Score float64 `json:"score"`
}
