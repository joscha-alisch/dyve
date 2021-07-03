package cloudfoundry

import cf "github.com/cloudfoundry-community/go-cfclient"

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
	GetSpace(guid string) (Space, []App, error)
	GetApp(guid string) (App, error)
}

type CfCli interface {
	ListOrgs() ([]cf.Org, error)
}

func main() {
	cli, _ := cf.NewClient(&cf.Config{})
	NewApi(cli)
}

func NewApi(cli CfCli) API {

	return &api{}
}

type api struct {
}

func (a *api) GetApp(guid string) (App, error) {
	return App{}, nil
}

func (a *api) GetSpace(guid string) (Space, []App, error) {
	return Space{}, nil, nil
}

func (a *api) GetOrg(guid string) (Org, error) {
	return Org{}, nil
}
