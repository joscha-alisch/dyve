package provider

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/database"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const Collection = "providers"

type Type string

const (
	TypeApps      Type = "apps"
	TypeGroups    Type = "groups"
	TypePipelines Type = "pipelines"
)

type Service interface {
	recon.JobProvider

	AddAppProvider(id string, p sdk.AppProvider) error
	GetAppProvider(id string) (sdk.AppProvider, error)
	DeleteAppProvider(id string) error

	AddPipelineProvider(id string, p sdk.PipelineProvider) error
	GetPipelineProvider(id string) (sdk.PipelineProvider, error)
	DeletePipelineProvider(id string) error

	AddGroupProvider(id string, p sdk.GroupProvider) error
	GetGroupProvider(id string) (sdk.GroupProvider, error)
	DeleteGroupProvider(id string) error
}

func NewService(db database.Database) Service {
	return &service{
		db:        db,
		providers: make(map[Type]map[string]interface{}),
	}
}

type service struct {
	db        database.Database
	providers map[Type]map[string]interface{}
}

func (s *service) AddAppProvider(id string, p sdk.AppProvider) error {
	return s.add(id, TypeApps, p)
}

func (s *service) GetAppProvider(id string) (sdk.AppProvider, error) {
	p, err := s.get(id, TypeApps)
	if err != nil {
		return nil, err
	}
	return p.(sdk.AppProvider), nil
}

func (s *service) DeleteAppProvider(id string) error {
	return s.delete(id, TypeApps)
}

func (s *service) AddPipelineProvider(id string, p sdk.PipelineProvider) error {
	return s.add(id, TypePipelines, p)
}

func (s *service) GetPipelineProvider(id string) (sdk.PipelineProvider, error) {
	p, err := s.get(id, TypePipelines)
	if err != nil {
		return nil, err
	}
	return p.(sdk.PipelineProvider), nil
}

func (s *service) DeletePipelineProvider(id string) error {
	return s.delete(id, TypePipelines)
}

func (s *service) AddGroupProvider(id string, p sdk.GroupProvider) error {
	return s.add(id, TypeGroups, p)
}

func (s *service) GetGroupProvider(id string) (sdk.GroupProvider, error) {
	p, err := s.get(id, TypeGroups)
	if err != nil {
		return nil, err
	}
	return p.(sdk.GroupProvider), nil
}

func (s *service) DeleteGroupProvider(id string) error {
	return s.delete(id, TypeGroups)
}

func (s *service) add(id string, providerType Type, p interface{}) error {
	if s.providers[providerType] == nil {
		s.providers[providerType] = make(map[string]interface{})
	}

	if s.providers[providerType][id] != nil {
		return ErrExists
	}
	err := s.db.EnsureCreated(Collection, bson.M{"id": id, "type": string(providerType)}, nil)
	if err != nil {
		return err
	}

	s.providers[providerType][id] = p
	return nil
}

func (s *service) get(id string, providerType Type) (interface{}, error) {
	if s.providers[providerType] == nil {
		return nil, ErrNotFound
	}

	if s.providers[providerType][id] == nil {
		return nil, ErrNotFound
	}

	return s.providers[providerType][id], nil
}

func (s *service) delete(id string, providerType Type) error {
	if s.providers[providerType] == nil {
		return ErrNotFound
	}

	err := s.db.DeleteOne(Collection, bson.M{"id": id, "type": providerType})
	if err != nil {
		return err
	}

	delete(s.providers[providerType], id)
	return nil
}

type provider struct {
	Id           string
	ProviderType string    `bson:"type"`
	LastUpdated  time.Time `bson:"lastUpdated"`
}

var currentTime = time.Now

func (s *service) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	t := currentTime()
	p := provider{}

	filter := bson.M{
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
	}
	update := bson.M{
		"$set": bson.M{
			"lastUpdated": t,
		},
	}
	err := s.db.UpdateOne(Collection, filter, false, update, &p)

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
