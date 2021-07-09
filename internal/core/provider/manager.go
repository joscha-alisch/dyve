package provider

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type Manager interface {
	AddAppProvider(id string, p sdk.AppProvider) error
	GetAppProvider(id string) (sdk.AppProvider, error)
}

func NewManager(db database.Database) Manager {
	return &manager{
		db: db,
	}
}

type manager struct {
	appProviders map[string]sdk.AppProvider
	db           database.Database
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
