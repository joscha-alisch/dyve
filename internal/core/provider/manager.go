package provider

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type Manager interface {
	AddAppProvider(id string, p sdk.AppProvider) error
	GetAppProvider(id string) (sdk.AppProvider, error)
	AddPipelineProvider(id string, p sdk.PipelineProvider) error
	GetPipelineProvider(id string) (sdk.PipelineProvider, error)
}

func NewManager(db database.Database) Manager {
	return &manager{
		db: db,
	}
}

type manager struct {
	appProviders      map[string]sdk.AppProvider
	pipelineProviders map[string]sdk.PipelineProvider
	db                database.Database
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
	return m.db.AddPipelineProvider(id)
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
	return m.db.AddAppProvider(id)
}

func (m *manager) GetAppProvider(id string) (sdk.AppProvider, error) {
	if m.appProviders[id] == nil {
		return nil, ErrNotFound
	}
	return m.appProviders[id], nil
}
