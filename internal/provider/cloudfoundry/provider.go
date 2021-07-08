package cloudfoundry

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"math"
)

func NewAppProvider(db Database) sdk.AppProvider {
	return &provider{
		db: db,
	}
}

type provider struct {
	db Database
}

func (p *provider) ListAppsPaged(perPage int, page int) (sdk.AppPage, error) {
	c, apps, err := p.db.ListAppsPaged(page, perPage)
	if err != nil {
		return sdk.AppPage{}, err
	}

	totalPages := int(math.Ceil(float64(c) / float64(perPage)))
	cursor := perPage * page
	if cursor > c {
		return sdk.AppPage{}, sdk.ErrPageExceeded
	}

	var res []sdk.App
	for _, app := range apps {
		res = append(res, sdk.App{
			Id:   app.Guid,
			Name: app.Name,
		})
	}

	return sdk.AppPage{
		TotalResults: c,
		TotalPages: totalPages,
		PerPage: perPage,
		Page: page,
		Apps: res,
	}, nil
}

func (p *provider) ListApps() ([]sdk.App, error) {
	cfApps, err := p.db.ListApps()
	if err != nil {
		return nil, err
	}

	var res []sdk.App
	for _, app := range cfApps {
		res = append(res, sdk.App{
			Id:   app.Guid,
			Name: app.Name,
		})
	}
	return res, nil
}

func (p *provider) GetApp(id string) (sdk.App, error) {
	app, err := p.db.GetApp(id)
	if err != nil {
		return sdk.App{}, err
	}

	return sdk.App{
		Id:   app.Guid,
		Name: app.Name,
	}, nil
}

func (p *provider) Search(term string, limit int) ([]sdk.AppSearchResult, error) {
	panic("implement me")
}
