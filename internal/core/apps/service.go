package apps

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

const Collection database.Collection = "apps"

type Service interface {
	ListAppsPaginated(perPage int, page int) (sdk.AppPage, error)
	GetApp(id string) (App, error)
	UpdateApps(providerId string, apps []sdk.App) error
	UpdateApp(app sdk.App) error
}

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

func (m *service) GetApp(id string) (App, error) {
	a := App{}
	return a, m.db.FindOneById(Collection, id, &a)
}

func (m *service) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	var res sdk.AppPage
	err := m.db.ListPaginated(Collection, perPage, page, &res.Pagination, func(c database.Decodable) error {
		app := sdk.App{}
		err := c.Decode(&app)
		if err != nil {
			return err
		}
		res.Apps = append(res.Apps, app)
		return nil
	})
	return res, err
}

func (m *service) UpdateApps(providerId string, apps []sdk.App) error {
	appMap := make(map[string]interface{}, len(apps))
	for _, app := range apps {
		appMap[app.Id] = app
	}
	return m.db.UpdateProvided(Collection, providerId, appMap)
}

func (m *service) UpdateApp(app sdk.App) error {
	return m.db.UpdateOneById(Collection, app.Id, false, app, nil)
}
