package cloudfoundry

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

func NewProvider(db Database, cf API) *Provider {
	return &Provider{
		db: db,
		cf: cf,
	}
}

type Provider struct {
	db Database
	cf API
}

func (p *Provider) ListApps() ([]sdk.App, error) {
	cfApps, err := p.db.ListApps()
	if err != nil {
		return nil, err
	}

	var res []sdk.App
	for _, app := range cfApps {
		res = append(res, app.toSdkApp())
	}
	return res, nil
}

func (p *Provider) GetApp(id string) (sdk.App, error) {
	app, err := p.db.GetApp(id)
	if err != nil {
		return sdk.App{}, err
	}

	return app.toSdkApp(), nil
}

func (p *Provider) GetAppRouting(id string) (sdk.AppRouting, error) {
	cached := sdk.AppRouting{}

	res, err := p.db.Cached(id+"/routing", 5*time.Second, &cached, func() (interface{}, error) {
		routes, err := p.cf.GetRoutes(id)
		if err != nil {
			return nil, err
		}
		appRouting := sdk.AppRouting{}
		for _, route := range routes {
			appRouting.Routes = append(appRouting.Routes, sdk.AppRoute{
				Host:    route.Host,
				Path:    route.Path,
				AppPort: route.Port,
			})
		}
		return appRouting, nil
	})
	if err != nil {
		return sdk.AppRouting{}, err
	}
	if res != nil {
		return res.(sdk.AppRouting), nil
	}

	return cached, nil
}

func (p *Provider) GetAppInstances(id string) (sdk.AppInstances, error) {
	cached := sdk.AppInstances{}
	res, err := p.db.Cached(id+"/instances", 5*time.Second, &cached, func() (interface{}, error) {
		instances, err := p.cf.GetInstances(id)
		if err != nil {
			return nil, err
		}
		appInstances := sdk.AppInstances{}
		for _, instance := range instances {
			appInstances = append(appInstances, sdk.AppInstance{
				State: cfStateToSdkState(instance.State),
				Since: instance.Since,
			})
		}
		return appInstances, nil
	})
	if err != nil {
		return nil, err
	}
	if res != nil {
		return res.(sdk.AppInstances), nil
	}
	return cached, nil
}

func cfStateToSdkState(state string) sdk.AppState {
	switch state {
	case "STOPPED":
		return sdk.AppStateStopped
	case "STARTED":
		return sdk.AppStateRunning
	default:
		return sdk.AppStateUnknown
	}
}
