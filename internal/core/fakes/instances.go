package fakes

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type MappingInstancesService struct {
	Instances map[string]sdk.AppInstances
}

func (m *MappingInstancesService) GetInstances(app string) (sdk.AppInstances, error) {
	return m.Instances[app], nil
}

func (m *MappingInstancesService) UpdateInstances(app string, routes sdk.AppInstances) error {
	m.Instances[app] = routes
	return nil
}
