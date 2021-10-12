package database

import (
	"context"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *mongoDb) GetApp(id string) (sdk.App, error) {
	res := m.apps.FindOne(context.Background(), bson.M{
		"id": id,
	})

	a := sdk.App{}
	err := res.Decode(&a)
	if err != nil {
		return sdk.App{}, err
	}

	return a, nil
}

func (m *mongoDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	p, cursor, err := m.listPaginated(m.apps, perPage, page)
	if err != nil {
		return sdk.AppPage{}, err
	}

	var apps []sdk.App
	for cursor.Next(context.Background()) {
		app := sdk.App{}
		err = cursor.Decode(&app)
		if err != nil {
			return sdk.AppPage{}, err
		}
		apps = append(apps, app)
	}

	return sdk.AppPage{
		Pagination: p,
		Apps:       apps,
	}, nil
}

func (m *mongoDb) UpdateApps(providerId string, apps []sdk.App) error {
	appMap := make(map[string]interface{}, len(apps))
	for _, app := range apps {
		appMap[app.Id] = app
	}

	return m.updateCollection(m.apps, providerId, appMap)
}
