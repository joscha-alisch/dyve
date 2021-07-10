package database

import (
	"context"
	"errors"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

var currentTime = time.Now

type MongoLogin struct {
	Uri string
	DB  string
}

func NewMongoDB(l MongoLogin) (Database, error) {
	c, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(l.Uri),
	)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	db := c.Database(l.DB)

	m := &mongoDb{
		cli:       c,
		db:        db,
		providers: db.Collection("providers"),
		apps:      db.Collection("apps"),
	}

	return m, nil
}

type mongoDb struct {
	cli       *mongo.Client
	db        *mongo.Database
	providers *mongo.Collection
	apps      *mongo.Collection
}

func (m *mongoDb) GetApp(id string) (sdk.App, error) {
	panic("implement me")
}

func (m *mongoDb) AddAppProvider(providerId string) error {
	_, err := m.providers.UpdateOne(context.Background(), bson.M{
		"id": providerId,
	}, bson.M{
		"$set": bson.M{
			"id":   providerId,
			"type": "apps",
		},
	}, options.Update().SetUpsert(true))

	return err
}

type provider struct {
	Id           string
	ProviderType string `bson:"type"`
}

func (m *mongoDb) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	t := currentTime()
	res := m.providers.FindOneAndUpdate(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{
				"lastUpdated": bson.M{
					"$lte": t.Add(-olderThan),
				},
			},
			bson.M{
				"lastUpdated": nil,
			},
		},
	}, bson.M{
		"$set": bson.M{
			"lastUpdated": t,
		},
	})

	p := provider{}
	err := res.Decode(&p)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return recon.Job{}, false
	}
	if err != nil {
		panic(err)
	}

	switch p.ProviderType {
	case "apps":
		return recon.Job{
			Type: ReconcileAppProvider,
			Guid: p.Id,
		}, true
	}

	return recon.Job{}, false
}

func (m *mongoDb) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	c, err := m.apps.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return sdk.AppPage{}, err
	}
	pages := int(math.Ceil(float64(c) / float64(perPage)))

	res, err := m.apps.Find(context.Background(), bson.M{}, options.Find().
		SetSkip(int64(page*perPage)).
		SetLimit(int64(perPage)),
	)
	if err != nil {
		return sdk.AppPage{}, err
	}

	var apps []sdk.App
	for res.Next(context.Background()) {
		app := sdk.App{}
		err = res.Decode(&app)
		if err != nil {
			return sdk.AppPage{}, err
		}
		apps = append(apps, app)
	}

	return sdk.AppPage{
		TotalResults: int(c),
		TotalPages:   pages,
		PerPage:      perPage,
		Page:         page,
		Apps:         apps,
	}, nil
}

func (m *mongoDb) DeleteAppProvider(providerId string) error {
	_, err := m.providers.DeleteOne(context.Background(), bson.M{
		"id": bson.M{
			"$eq": providerId,
		},
	})
	if err != nil {
		return err
	}
	_, err = m.apps.DeleteMany(context.Background(), bson.M{
		"provider": bson.M{
			"$eq": providerId,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDb) UpdateApps(providerId string, apps []sdk.App) error {
	var ids []string
	for _, app := range apps {
		ids = append(ids, app.Id)
		_, err := m.apps.UpdateOne(context.Background(), bson.M{
			"provider": bson.M{
				"$eq": providerId,
			},
			"id": bson.M{
				"$eq": app.Id,
			},
		}, bson.D{
			bson.E{Key: "$set", Value: bson.M{"provider": providerId}},
			bson.E{Key: "$set", Value: app},
		}, options.Update().SetUpsert(true))
		if err != nil {
			return err
		}
	}

	_, err := m.apps.DeleteMany(context.Background(), bson.M{
		"provider": bson.M{
			"$eq": providerId,
		},
		"id": bson.M{
			"$nin": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
