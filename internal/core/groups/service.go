package groups

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
)

const Collection = "groups"

type Service interface {
	ListGroupsByProvider() (GroupByProviderMap, error)
	ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error)
	GetGroup(id string) (sdk.Group, error)
	DeleteGroup(id string) error
	UpdateGroups(guid string, groups []sdk.Group) error
}

func NewService(db database.Database, providers provider.Service) Service {
	return &service{
		db:        db,
		providers: providers,
	}
}

type service struct {
	db        database.Database
	providers provider.Service
}

func (s *service) ListGroupsByProvider() (GroupByProviderMap, error) {
	m := make(GroupByProviderMap)

	groupProviders, err := s.providers.ListGroupProviders()
	if err != nil {
		return nil, err
	}
	for _, groupProvider := range groupProviders {
		m[groupProvider.Id] = ProviderWithGroups{
			Provider: groupProvider.Id,
			Name:     groupProvider.Name,
		}
	}

	each := func(c database.Decodable) error {
		group := GroupWithProvider{}
		err := c.Decode(&group)
		if err != nil {
			return err
		}

		prov := m[group.Provider]
		prov.Groups = append(prov.Groups, group.Group)
		m[group.Provider] = prov

		return nil
	}

	err = s.db.FindMany(Collection, bson.M{}, each)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (s *service) ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error) {
	var res sdk.GroupPage
	err := s.db.ListPaginated(Collection, perPage, page, &res.Pagination, func(c database.Decodable) error {
		group := sdk.Group{}
		err := c.Decode(&group)
		if err != nil {
			return err
		}
		res.Groups = append(res.Groups, group)
		return nil
	})
	return res, err
}

func (s *service) GetGroup(id string) (sdk.Group, error) {
	t := sdk.Group{}
	return t, s.db.FindOneById(Collection, id, &t)
}

func (s *service) DeleteGroup(id string) error {
	return s.db.DeleteOneById(Collection, id)
}

func (s *service) UpdateGroups(providerId string, groups []sdk.Group) error {
	groupMap := make(map[string]interface{}, len(groups))
	for _, group := range groups {
		groupMap[group.Id] = group
	}
	return s.db.UpdateProvided(Collection, providerId, groupMap)
}
