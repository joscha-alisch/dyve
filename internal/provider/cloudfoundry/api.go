package cloudfoundry

type Login struct {
	Api  string
	User string
	Pass string
}

/**
A cloudfoundry.API is a wrapper around the official client.
*/
type API interface {
	GetOrg(guid string) (Org, error)
	GetSpace(guid string) (Space, error)
	GetApp(guid string) (App, error)
}

func NewApi() API {
	return &api{}
}

type api struct {
}

func (a *api) GetApp(guid string) (App, error) {
	return App{}, nil
}

func (a *api) GetSpace(guid string) (Space, error) {
	return Space{}, nil
}

func (a *api) GetOrg(guid string) (Org, error) {
	return Org{}, nil
}
