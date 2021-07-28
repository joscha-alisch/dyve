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
	"sort"
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
		cli:          c,
		db:           db,
		providers:    db.Collection("providers"),
		apps:         db.Collection("apps"),
		pipelines:    db.Collection("pipelines"),
		pipelineRuns: db.Collection("pipeline_runs"),
	}

	return m, nil
}

type mongoDb struct {
	cli          *mongo.Client
	db           *mongo.Database
	providers    *mongo.Collection
	apps         *mongo.Collection
	pipelines    *mongo.Collection
	pipelineRuns *mongo.Collection
}

func (m *mongoDb) AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error {
	models := make([]mongo.WriteModel, len(runs))

	for i, run := range runs {
		model := mongo.NewUpdateOneModel()
		model.SetUpsert(true).SetFilter(bson.M{
			"provider": bson.M{
				"$eq": providerId,
			},
			"pipelineId": bson.M{
				"$eq": run.PipelineId,
			},
			"started": bson.M{
				"$eq": run.Started,
			},
		}).SetUpdate(bson.D{
			bson.E{Key: "$set", Value: bson.M{"provider": providerId}},
			bson.E{Key: "$set", Value: run},
		})
		models[i] = model
	}

	_, err := m.pipelineRuns.BulkWrite(context.Background(), models)
	return err
}

func (m *mongoDb) ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error) {
	res, err := m.pipelineRuns.Find(context.Background(), bson.M{
		"pipelineId": bson.M{
			"$eq": id,
		},
		"started": bson.M{
			"$lt":  toExcl,
			"$gte": fromIncl,
		},
	})
	if err != nil {
		return nil, err
	}

	var runs sdk.PipelineStatusList

	for res.Next(context.Background()) {
		run := sdk.PipelineStatus{}
		err = res.Decode(&run)
		if err != nil {
			return nil, err
		}

		runs = append(runs, run)
	}

	sort.Sort(runs)
	return runs, nil
}

func (m *mongoDb) ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error) {
	p, cursor, err := m.listPaginated(m.pipelines, perPage, page)
	if err != nil {
		return sdk.PipelinePage{}, err
	}

	var pipelines []sdk.Pipeline
	for cursor.Next(context.Background()) {
		pipeline := sdk.Pipeline{}
		err = cursor.Decode(&pipeline)
		if err != nil {
			return sdk.PipelinePage{}, err
		}
		pipelines = append(pipelines, pipeline)
	}

	return sdk.PipelinePage{
		Pagination: p,
		Pipelines:  pipelines,
	}, nil
}

func (m *mongoDb) GetPipeline(id string) (sdk.Pipeline, error) {
	res := m.pipelines.FindOne(context.Background(), bson.M{
		"id": id,
	})

	p := sdk.Pipeline{}
	err := res.Decode(&p)
	if err != nil {
		return sdk.Pipeline{}, err
	}

	return p, nil
}

func (m *mongoDb) AddPipelineProvider(providerId string) error {
	return m.addProvider(providerId, "pipelines")
}

func (m *mongoDb) DeletePipelineProvider(providerId string) error {
	return m.deleteProvider(providerId, m.pipelines)
}

func (m *mongoDb) UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error {
	pipelineMap := make(map[string]interface{}, len(pipelines))
	for _, pipeline := range pipelines {
		pipelineMap[pipeline.Id] = pipeline
	}

	return m.updateCollection(m.pipelines, providerId, pipelineMap)
}

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

func (m *mongoDb) AddAppProvider(providerId string) error {
	return m.addProvider(providerId, "apps")
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

	return recon.Job{
		Type: recon.Type(p.ProviderType),
		Guid: p.Id,
	}, true
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

func (m *mongoDb) DeleteAppProvider(providerId string) error {
	return m.deleteProvider(providerId, m.apps)
}

func (m *mongoDb) UpdateApps(providerId string, apps []sdk.App) error {
	appMap := make(map[string]interface{}, len(apps))
	for _, app := range apps {
		appMap[app.Id] = app
	}

	return m.updateCollection(m.apps, providerId, appMap)
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

	_, err := c.BulkWrite(context.Background(), models)
	if err != nil {
		return err
	}

	_, err = c.DeleteMany(context.Background(), bson.M{
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
