package pipelines

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sort"
	"time"
)

const Collection database.Collection = "pipelines"
const CollectionRuns database.Collection = "pipeline_runs"
const CollectionVersions database.Collection = "pipeline_versions"

type Service interface {
	EnsureIndices() error

	ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error)
	GetPipeline(id string) (sdk.Pipeline, error)
	ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error)
	ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error)
	ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error)
	UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error
	AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error
	AddPipelineVersions(providerId string, versions sdk.PipelineVersionList) error
}

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

func (s *service) EnsureIndices() error {
	return s.db.EnsureIndex(CollectionRuns, mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "providerId", Value: 1},
			bson.E{Key: "pipelineId", Value: 1},
			bson.E{Key: "started", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
}

func (s *service) AddPipelineVersions(providerId string, versions sdk.PipelineVersionList) error {
	idMap := make(map[string]interface{}, len(versions))
	filterMap := make(map[string]interface{}, len(versions))
	for _, version := range versions {
		id := version.PipelineId + version.Created.Format(time.RFC3339)
		idMap[id] = version
		filterMap[id] = bson.M{
			"provider":   providerId,
			"pipelineId": version.PipelineId,
			"created":    version.Created,
		}
	}

	return s.db.UpdateMany(CollectionVersions, filterMap, idMap)
}

func (s *service) ListPipelineRunsLimit(id string, toExcl time.Time, limit int) (sdk.PipelineStatusList, error) {
	var runs sdk.PipelineStatusList

	filter := bson.M{
		"pipelineId": bson.M{
			"$eq": id,
		},
		"started": bson.M{
			"$lt": toExcl,
		},
	}
	each := func(c database.Decodable) error {
		run := sdk.PipelineStatus{}
		err := c.Decode(&run)
		if err != nil {
			return err
		}

		runs = append(runs, run)
		return nil
	}

	err := s.db.FindManyWithOptions(CollectionRuns, filter, each, bson.M{"started": -1}, limit)
	if err != nil {
		return nil, err
	}

	sort.Sort(runs)
	return runs, nil
}

func (s *service) ListPipelineVersions(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineVersionList, error) {
	var versions sdk.PipelineVersionList
	filter := bson.M{
		"pipelineId": id,
		"created": bson.M{
			"$lt":  toExcl,
			"$gte": fromIncl,
		},
	}
	each := func(c database.Decodable) error {
		version := sdk.PipelineVersion{}
		err := c.Decode(&version)
		if err != nil {
			return err
		}

		versions = append(versions, version)
		return nil
	}

	err := s.db.FindMany(CollectionVersions, filter, each)
	if err != nil {
		return nil, err
	}

	sort.Sort(versions)

	if len(versions) == 0 || versions[0].Created != fromIncl {
		version := sdk.PipelineVersion{}
		err := s.db.FindOneSorted(CollectionVersions, bson.M{
			"pipelineId": bson.M{
				"$eq": id,
			},
			"created": bson.M{
				"$lt": fromIncl,
			},
		}, bson.M{"created": -1}, &version)
		if err != nil {
			return nil, err
		}

		versions = append(sdk.PipelineVersionList{version}, versions...)
	}

	return versions, nil
}

func (s *service) AddPipelineRuns(providerId string, runs sdk.PipelineStatusList) error {
	idMap := make(map[string]interface{}, len(runs))
	filterMap := make(map[string]interface{}, len(runs))
	for _, run := range runs {
		filterMap[run.PipelineId] = bson.M{
			"provider":   providerId,
			"pipelineId": run.PipelineId,
			"started":    run.Started,
		}
		idMap[run.PipelineId] = run
	}

	return s.db.UpdateMany(CollectionRuns, filterMap, idMap)
}

func (s *service) ListPipelineRuns(id string, fromIncl time.Time, toExcl time.Time) (sdk.PipelineStatusList, error) {
	var runs sdk.PipelineStatusList

	filter := bson.M{
		"pipelineId": bson.M{
			"$eq": id,
		},
		"started": bson.M{
			"$lt":  toExcl,
			"$gte": fromIncl,
		},
	}
	each := func(c database.Decodable) error {
		run := sdk.PipelineStatus{}
		err := c.Decode(&run)
		if err != nil {
			return err
		}

		runs = append(runs, run)
		return err
	}

	err := s.db.FindMany(CollectionRuns, filter, each)
	if err != nil {
		return nil, err
	}

	sort.Sort(runs)
	return runs, nil
}

func (s *service) ListPipelinesPaginated(perPage int, page int) (sdk.PipelinePage, error) {
	var res sdk.PipelinePage
	err := s.db.ListPaginated(Collection, perPage, page, &res.Pagination, func(c database.Decodable) error {
		i := sdk.Pipeline{}
		err := c.Decode(&i)
		if err != nil {
			return err
		}
		res.Pipelines = append(res.Pipelines, i)
		return nil
	})
	return res, err
}

func (s *service) GetPipeline(id string) (sdk.Pipeline, error) {
	p := sdk.Pipeline{}
	return p, s.db.FindOneById(Collection, id, &p)
}

func (s *service) UpdatePipelines(providerId string, pipelines []sdk.Pipeline) error {
	pipelineMap := make(map[string]interface{}, len(pipelines))
	for _, pipeline := range pipelines {
		pipelineMap[pipeline.Id] = pipeline
	}
	return s.db.UpdateProvided(Collection, providerId, pipelineMap)
}
