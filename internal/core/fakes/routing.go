package fakes

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type MappingRoutesService struct {
	Routes map[string]sdk.AppRouting
}

func (m *MappingRoutesService) GetRoutes(app string) (sdk.AppRouting, error) {
	return m.Routes[app], nil
}

func (m *MappingRoutesService) UpdateRoutes(app string, routes sdk.AppRouting) error {
	m.Routes[app] = routes
	return nil
}
