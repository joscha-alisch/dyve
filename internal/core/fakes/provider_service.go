package fakes

import (
	"github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

type ProviderService struct {
	Job                *recon.Job
	AppProviders       map[string]sdk.AppProvider
	PipelineProviders  map[string]sdk.PipelineProvider
	RoutingProviders   map[string]sdk.RoutingProvider
	InstancesProviders map[string]sdk.InstancesProvider
	GroupProviders     map[string]sdk.GroupProvider
}

func (s *ProviderService) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	if s.Job == nil {
		return recon.Job{}, false
	}
	return *s.Job, true
}

func (s *ProviderService) AddAppProvider(id string, name string, p sdk.AppProvider) error {
	s.AppProviders[id] = p
	return nil
}

func (s *ProviderService) GetAppProvider(id string) (sdk.AppProvider, error) {
	if s.AppProviders[id] == nil {
		return nil, provider.ErrNotFound
	}

	return s.AppProviders[id], nil
}

func (s *ProviderService) DeleteAppProvider(id string) error {
	delete(s.AppProviders, id)
	return nil
}

func (s *ProviderService) RequestAppUpdate(id string) error {
	panic("implement me")
}

func (s *ProviderService) AddRoutingProvider(id string, name string, p sdk.RoutingProvider) error {
	s.RoutingProviders[id] = p
	return nil
}

func (s *ProviderService) GetRoutingProviders() ([]sdk.RoutingProvider, error) {
	var res []sdk.RoutingProvider
	for _, routingProvider := range s.RoutingProviders {
		res = append(res, routingProvider)
	}
	return res, nil
}

func (s *ProviderService) DeleteRoutingProvider(id string) error {
	delete(s.RoutingProviders, id)
	return nil
}

func (s *ProviderService) AddInstancesProvider(id string, name string, p sdk.InstancesProvider) error {
	s.InstancesProviders[id] = p
	return nil
}

func (s *ProviderService) GetInstancesProviders() ([]sdk.InstancesProvider, error) {
	var res []sdk.InstancesProvider
	for _, instancesProvider := range s.InstancesProviders {
		res = append(res, instancesProvider)
	}
	return res, nil
}

func (s *ProviderService) DeleteInstancesProvider(id string) error {
	panic("implement me")
}

func (s *ProviderService) AddPipelineProvider(id string, name string, p sdk.PipelineProvider) error {
	s.PipelineProviders[id] = p
	return nil
}

func (s *ProviderService) GetPipelineProvider(id string) (sdk.PipelineProvider, error) {
	if s.PipelineProviders[id] == nil {
		return nil, provider.ErrNotFound
	}
	return s.PipelineProviders[id], nil
}

func (s *ProviderService) DeletePipelineProvider(id string) error {
	delete(s.PipelineProviders, id)
	return nil
}

func (s *ProviderService) ListGroupProviders() ([]provider.Data, error) {
	var res []provider.Data
	for id, _ := range s.GroupProviders {
		res = append(res, provider.Data{
			Id:   id,
			Name: id,
		})
	}
	return res, nil
}

func (s *ProviderService) AddGroupProvider(id string, name string, p sdk.GroupProvider) error {
	s.GroupProviders[id] = p
	return nil
}

func (s *ProviderService) GetGroupProvider(id string) (sdk.GroupProvider, error) {
	if s.GroupProviders[id] == nil {
		return nil, provider.ErrNotFound
	}
	return s.GroupProviders[id], nil
}

func (s *ProviderService) DeleteGroupProvider(id string) error {
	delete(s.PipelineProviders, id)
	return nil
}
