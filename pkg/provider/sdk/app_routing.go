package sdk

type RoutingProvider interface {
	GetAppRouting(id string) (AppRouting, error)
}

type AppRouting struct {
	Routes AppRoutes `json:"routes"`
}

type AppRoutes []AppRoute

type AppRoute struct {
	Host    string `json:"host"`
	Path    string `json:"path"`
	AppPort int    `json:"appPort"`
}
