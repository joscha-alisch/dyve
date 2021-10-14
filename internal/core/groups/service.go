package groups

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/mongo"
)

const Collection = "groups"

type Service interface {
	ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error)
	GetGroup(id string) (sdk.Group, error)
	UpdateGroups(guid string, groups []sdk.Group) error
}

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

func (s *service) ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error) {
	var res sdk.GroupPage
	err := s.db.ListPaginated(Collection, perPage, page, &res.Pagination, func(c *mongo.Cursor) error {
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
