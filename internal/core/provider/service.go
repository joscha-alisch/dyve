package provider

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/queue"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"strings"
	"time"
)

const Collection = "providers"

type Type string

const (
	ReconcileAppProvider        recon.Type = "apps"
	ReconcileRoutingProviders   recon.Type = "routing"
	ReconcilePipelineProvider   recon.Type = "pipelines"
	ReconcileGroupProvider      recon.Type = "groups"
	ReconcileInstancesProviders recon.Type = "instances"
)

const (
	TypeApps      Type = "apps"
	TypeGroups    Type = "groups"
	TypePipelines Type = "pipelines"
	TypeRouting   Type = "routing"
	TypeInstances Type = "instances"
)

type Service interface {
	recon.JobProvider

	AddAppProvider(id string, name string, p sdk.AppProvider) error
	GetAppProvider(id string) (sdk.AppProvider, error)
	DeleteAppProvider(id string) error
	RequestAppUpdate(id string) error

	AddRoutingProvider(id string, name string, p sdk.RoutingProvider) error
	GetRoutingProviders() ([]sdk.RoutingProvider, error)
	DeleteRoutingProvider(id string) error

	AddInstancesProvider(id string, name string, p sdk.InstancesProvider) error
	GetInstancesProviders() ([]sdk.InstancesProvider, error)
	DeleteInstancesProvider(id string) error

	AddPipelineProvider(id string, name string, p sdk.PipelineProvider) error
	GetPipelineProvider(id string) (sdk.PipelineProvider, error)
	DeletePipelineProvider(id string) error

	ListGroupProviders() ([]Data, error)
	AddGroupProvider(id string, name string, p sdk.GroupProvider) error
	GetGroupProvider(id string) (sdk.GroupProvider, error)
	DeleteGroupProvider(id string) error
}

func NewService(db database.Database) Service {
	return &service{
		db:                db,
		providers:         make(map[Type]map[string]interface{}),
		appUpdateRequests: queue.NewStringQueue(1000),
	}
}

type service struct {
	db                database.Database
	providers         map[Type]map[string]interface{}
	appUpdateRequests *queue.StringQueue
}

func (s *service) AddInstancesProvider(id string, name string, p sdk.InstancesProvider) error {
	return s.add(id, name, TypeInstances, p)
}

func (s *service) GetInstancesProviders() ([]sdk.InstancesProvider, error) {
	providers, err := s.getAll(TypeInstances)
	if err != nil {
		return nil, err
	}

	var res []sdk.InstancesProvider
	for _, provider := range providers.([]interface{}) {
		res = append(res, provider.(sdk.InstancesProvider))
	}

	return res, nil
}

func (s *service) DeleteInstancesProvider(id string) error {
	return s.delete(id, TypeInstances)
}

func (s *service) AddRoutingProvider(id string, name string, p sdk.RoutingProvider) error {
	return s.add(id, name, TypeRouting, p)
}

func (s *service) GetRoutingProviders() ([]sdk.RoutingProvider, error) {
	providers, err := s.getAll(TypeRouting)
	if err != nil {
		return nil, err
	}

	var res []sdk.RoutingProvider
	for _, provider := range providers.([]interface{}) {
		res = append(res, provider.(sdk.RoutingProvider))
	}

	return res, nil
}

func (s *service) DeleteRoutingProvider(id string) error {
	return s.delete(id, TypeRouting)
}

func (s *service) RequestAppUpdate(id string) error {
	err := s.appUpdateRequests.Push(string(TypeRouting) + "/" + id)
	if err != nil {
		return err
	}
	err = s.appUpdateRequests.Push(string(TypeInstances) + "/" + id)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) AddAppProvider(id string, name string, p sdk.AppProvider) error {
	return s.add(id, name, TypeApps, p)
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

func (s *service) AddPipelineProvider(id string, name string, p sdk.PipelineProvider) error {
	return s.add(id, name, TypePipelines, p)
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

func (s *service) ListGroupProviders() ([]Data, error) {
	return s.list(TypeGroups)
}

func (s *service) AddGroupProvider(id string, name string, p sdk.GroupProvider) error {
	return s.add(id, name, TypeGroups, p)
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

func (s *service) add(id string, name string, providerType Type, p interface{}) error {
	if s.providers[providerType] == nil {
		s.providers[providerType] = make(map[string]interface{})
	}

	if s.providers[providerType][id] != nil {
		return ErrExists
	}

	providerData := bson.M{"id": id, "name": name, "type": string(providerType)}

	err := s.db.UpdateOne(Collection, providerData, true, providerData, nil)
	if err != nil {
		return err
	}

	s.providers[providerType][id] = p
	return nil
}

func (s *service) list(providerType Type) ([]Data, error) {
	var res []Data
	err := s.db.FindMany(Collection, bson.M{"type": providerType}, func(c database.Decodable) error {
		data := Data{}
		err := c.Decode(&data)
		if err != nil {
			return err
		}

		res = append(res, data)
		return nil
	})
	return res, err
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

func (s *service) getAll(providerType Type) (interface{}, error) {
	if s.providers[providerType] == nil {
		return nil, ErrNotFound
	}

	var res []interface{}
	for _, provider := range s.providers[providerType] {
		res = append(res, provider)
	}
	return res, nil
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
	if len(s.providers[providerType]) == 0 {
		s.providers[providerType] = nil
	}
	return nil
}

var currentTime = time.Now

func (s *service) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	t := currentTime()
	p := Provider{}

	if request, ok := s.appUpdateRequests.Pop(); ok {
		parts := strings.Split(request, "/")
		t := Type(parts[0])
		appId := parts[1]

		switch t {
		case TypeRouting:
			return recon.Job{
				Type:        ReconcileRoutingProviders,
				Guid:        appId,
				LastUpdated: time.Time{},
			}, true
		case TypeInstances:
			return recon.Job{
				Type:        ReconcileInstancesProviders,
				Guid:        appId,
				LastUpdated: time.Time{},
			}, true
		}

	}

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
		"lastUpdated": t,
	}
	err := s.db.UpdateOne(Collection, filter, false, update, &p)
	if errors.Is(err, database.ErrNotFound) {
		return recon.Job{}, false
	}
	if err != nil {
		log.Error().Err(err).Msg("error when fetching job")
		return recon.Job{}, false
	}

	return recon.Job{
		Type:        recon.Type(p.ProviderType),
		Guid:        p.Id,
		LastUpdated: p.LastUpdated,
	}, true
}
