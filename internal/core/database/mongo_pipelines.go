package database

import (
	"context"
	"errors"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

func (m *mongoDb) AddPipelineVersions(providerId string, versions sdk.PipelineVersionList) error {
	models := make([]mongo.WriteModel, len(versions))

	for i, version := range versions {
		model := mongo.NewUpdateOneModel()
		model.SetUpsert(true).SetFilter(bson.M{
			"provider": bson.M{
				"$eq": providerId,
			},
			"pipelineId": bson.M{
				"$eq": version.PipelineId,
			},
			"created": bson.M{
				"$eq": version.Created,
			},
		}).SetUpdate(bson.D{
			bson.E{Key: "$set", Value: bson.M{"provider": providerId}},
			bson.E{Key: "$set", Value: version},
		})
		models[i] = model
	}

	_, err := m.pipelineVersions.BulkWrite(context.Background(), models)
	return err
}

func (m *mongoDb) ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error) {
	res, err := m.pipelineRuns.Find(context.Background(), bson.M{
		"pipelineId": bson.M{
			"$eq": id,
		},
		"started": bson.M{
			"$lt": toExcl,
		},
	}, options.Find().SetLimit(int64(limit)).SetSort(bson.M{"started": -1}))
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

func (m *mongoDb) ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error) {
	res, err := m.pipelineVersions.Find(context.Background(), bson.M{
		"pipelineId": bson.M{
			"$eq": id,
		},
		"created": bson.M{
			"$lt":  toExcl,
			"$gte": fromIncl,
		},
	})
	if err != nil {
		return nil, err
	}

	var versions sdk.PipelineVersionList

	for res.Next(context.Background()) {
		version := sdk.PipelineVersion{}
		err = res.Decode(&version)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	sort.Sort(versions)

	if len(versions) == 0 || versions[0].Created != fromIncl {
		res := m.pipelineVersions.FindOne(context.Background(), bson.M{
			"pipelineId": bson.M{
				"$eq": id,
			},
			"created": bson.M{
				"$lt": fromIncl,
			},
		}, options.FindOne().SetSort(bson.M{"created": -1}))
		version := sdk.PipelineVersion{}
		err = res.Decode(&version)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return versions, nil
		}

		versions = append(sdk.PipelineVersionList{version}, versions...)
	}

	return versions, nil
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

	_, err := m.pipelineRuns.BulkWrite(context.Background(), models, options.BulkWrite().SetOrdered(false))
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

func (m *mongoDb) UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error {
	pipelineMap := make(map[string]interface{}, len(pipelines))
	for _, pipeline := range pipelines {
		pipelineMap[pipeline.Id] = pipeline
	}

	return m.updateCollection(m.pipelines, providerId, pipelineMap)
}
