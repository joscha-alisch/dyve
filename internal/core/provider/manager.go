package provider

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type Type string

const (
	TypeApps      Type = "apps"
	TypePipelines Type = "pipelines"
	TypeGroups    Type = "groups"
)

type Manager interface {
	AddAppProvider(id string, p sdk.AppProvider) error
	GetAppProvider(id string) (sdk.AppProvider, error)
	AddPipelineProvider(id string, p sdk.PipelineProvider) error
	GetPipelineProvider(id string) (sdk.PipelineProvider, error)
	AddGroupProvider(id string, p sdk.GroupProvider) error
	GetGroupProvider(id string) (sdk.GroupProvider, error)
}

type Store interface {
	AddProvider(id string, providerType Type) error
	DeleteProvider(id string, providerType Type) error
}

func NewManager(db Store) Manager {
	return &manager{
		db: db,
	}
}

type manager struct {
	appProviders      map[string]sdk.AppProvider
	pipelineProviders map[string]sdk.PipelineProvider
	groupProviders    map[string]sdk.GroupProvider
	db                Store
}

func (m *manager) AddGroupProvider(id string, p sdk.GroupProvider) error {
	if p == nil {
		return ErrNil
	}

	if m.groupProviders == nil {
		m.groupProviders = make(map[string]sdk.GroupProvider)
	}

	if m.groupProviders[id] != nil {
		return ErrExists
	}

	m.groupProviders[id] = p
	return m.db.AddProvider(id, TypeGroups)
}

func (m *manager) GetGroupProvider(id string) (sdk.GroupProvider, error) {
	if m.groupProviders[id] == nil {
		return nil, ErrNotFound
	}
	return m.groupProviders[id], nil
}

func (m *manager) AddPipelineProvider(id string, p sdk.PipelineProvider) error {
	if p == nil {
		return ErrNil
	}

	if m.pipelineProviders == nil {
		m.pipelineProviders = make(map[string]sdk.PipelineProvider)
	}

	if m.pipelineProviders[id] != nil {
		return ErrExists
	}

	m.pipelineProviders[id] = p
	return m.db.AddProvider(id, TypePipelines)
}

func (m *manager) GetPipelineProvider(id string) (sdk.PipelineProvider, error) {
	if m.pipelineProviders[id] == nil {
		return nil, ErrNotFound
	}
	return m.pipelineProviders[id], nil
}

func (m *manager) AddAppProvider(id string, p sdk.AppProvider) error {
	if p == nil {
		return ErrNil
	}

	if m.appProviders == nil {
		m.appProviders = make(map[string]sdk.AppProvider)
	}

	if m.appProviders[id] != nil {
		return ErrExists
	}

	m.appProviders[id] = p
	return m.db.AddProvider(id, TypeApps)
}

func (m *manager) GetAppProvider(id string) (sdk.AppProvider, error) {
	if m.appProviders[id] == nil {
		return nil, ErrNotFound
	}
	return m.appProviders[id], nil
}
