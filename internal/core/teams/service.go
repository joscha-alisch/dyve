package teams

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const Collection = "teams"

type Service interface {
	ListTeamsPaginated(perPage int, page int) (TeamPage, error)
	GetTeam(id string) (Team, error)
	DeleteTeam(id string) error
	UpdateTeam(id string) error
}

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

func (s *service) ListTeamsPaginated(perPage int, page int) (TeamPage, error) {
	var res TeamPage
	err := s.db.ListPaginated(Collection, perPage, page, &res.Pagination, func(c *mongo.Cursor) error {
		team := Team{}
		err := c.Decode(&team)
		if err != nil {
			return err
		}
		res.Teams = append(res.Teams, team)
		return nil
	})
	return res, err
}

func (s *service) GetTeam(id string) (Team, error) {
	t := Team{}
	return t, s.db.FindOneById(Collection, id, &t)
}

func (s *service) DeleteTeam(id string) error {
	return s.db.DeleteOneById(Collection, id)
}

func (s *service) UpdateTeam(id string) error {
	return s.db.UpdateOneById(Collection, id, true, bson.M{}, nil)
}
