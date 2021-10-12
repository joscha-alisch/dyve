package database

import (
	"context"
	"errors"
	prov "github.com/joscha-alisch/dyve/internal/core/provider"
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
		cli:              c,
		db:               db,
		providers:        db.Collection("providers"),
		apps:             db.Collection("apps"),
		pipelines:        db.Collection("pipelines"),
		pipelineRuns:     db.Collection("pipeline_runs"),
		pipelineVersions: db.Collection("pipeline_versions"),
		groups:           db.Collection("groups"),
	}

	m.pipelineRuns.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.M{
			"providerId": 1,
			"pipelineId": 1,
			"started":    1,
		},
		Options: options.Index().SetUnique(true),
	})

	return m, nil
}

type mongoDb struct {
	cli              *mongo.Client
	db               *mongo.Database
	providers        *mongo.Collection
	apps             *mongo.Collection
	pipelines        *mongo.Collection
	pipelineRuns     *mongo.Collection
	pipelineVersions *mongo.Collection
	groups           *mongo.Collection
}

func (m *mongoDb) AddProvider(providerId string, providerType prov.Type) error {
	return m.addProvider(providerId, string(providerType))
}

func (m *mongoDb) DeleteProvider(providerId string, providerType prov.Type) error {
	switch providerType {
	case prov.TypeApps:
		return m.deleteProvider(providerId, m.apps)
	case prov.TypePipelines:
		return m.deleteProvider(providerId, m.pipelines)
	case prov.TypeGroups:
		return m.deleteProvider(providerId, m.groups)
	default:
		return errors.New("unknown provider type " + string(providerType))
	}
}

type provider struct {
	Id           string
	ProviderType string    `bson:"type"`
	LastUpdated  time.Time `bson:"lastUpdated"`
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

	return recon.Job{
		Type:        recon.Type(p.ProviderType),
		Guid:        p.Id,
		LastUpdated: p.LastUpdated,
	}, true
}

func (m *mongoDb) updateCollection(c *mongo.Collection, providerId string, items map[string]interface{}) error {
	ids := make([]string, len(items))
	models := make([]mongo.WriteModel, len(items))
	i := 0
	for k, v := range items {
		ids[i] = k
		model := mongo.NewUpdateOneModel()
		model.SetUpsert(true).SetFilter(bson.M{
			"provider": bson.M{
				"$eq": providerId,
			},
			"id": bson.M{
				"$eq": k,
			},
		}).SetUpdate(bson.D{
			bson.E{Key: "$set", Value: bson.M{"provider": providerId}},
			bson.E{Key: "$set", Value: v},
		})
		models[i] = model
		i++
	}
	ctx := context.Background()

	_, err := c.BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	_, err = c.DeleteMany(ctx, bson.M{
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

func (m *mongoDb) addProvider(providerId, providerType string) error {
	_, err := m.providers.UpdateOne(context.Background(), bson.M{
		"id": providerId,
	}, bson.M{
		"$set": bson.M{
			"id":   providerId,
			"type": providerType,
		},
	}, options.Update().SetUpsert(true))

	return err
}

func (m *mongoDb) deleteProvider(providerId string, dependents ...*mongo.Collection) error {
	_, err := m.providers.DeleteOne(context.Background(), bson.M{
		"id": bson.M{
			"$eq": providerId,
		},
	})
	if err != nil {
		return err
	}

	for _, dependent := range dependents {
		_, err = dependent.DeleteMany(context.Background(), bson.M{
			"provider": bson.M{
				"$eq": providerId,
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mongoDb) listPaginated(coll *mongo.Collection, perPage, page int) (sdk.Pagination, *mongo.Cursor, error) {
	c, err := coll.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return sdk.Pagination{}, nil, err
	}
	pages := int(math.Ceil(float64(c) / float64(perPage)))

	res, err := coll.Find(context.Background(), bson.M{}, options.Find().
		SetSkip(int64(page*perPage)).
		SetLimit(int64(perPage)),
	)
	if err != nil {
		return sdk.Pagination{}, nil, err
	}

	return sdk.Pagination{
		TotalResults: int(c),
		TotalPages:   pages,
		PerPage:      perPage,
		Page:         page,
	}, res, nil
}
