package sdk

type App struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AppSearchResult struct {
	App   App     `json:"app"`
	Score float64 `json:"score"`
}
