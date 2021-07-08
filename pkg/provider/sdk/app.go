package sdk

type AppPage struct {
	TotalResults int `json:"totalResults"`
	TotalPages int `json:"totalPages"`
	PerPage int `json:"perPage"`
	Page int `json:"page"`
	Apps []App `json:"apps"`
}

type App struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AppSearchResult struct {
	App   App     `json:"app"`
	Score float64 `json:"score"`
}
