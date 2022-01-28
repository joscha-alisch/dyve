package teams

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const Collection = "teams"

type Service interface {
	EnsureIndices() error

	ListTeamsPaginated(perPage int, page int) (TeamPage, error)
	GetTeam(id string) (Team, error)
	DeleteTeam(id string) error
	CreateTeam(id string, data TeamSettings) error
	UpdateTeam(id string, data TeamSettings) error
	TeamsForGroups(groups []string) (ByAccess, error)
}

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

func (s *service) TeamsForGroups(groups []string) (ByAccess, error) {
	res := ByAccess{}

	err := s.db.FindMany(Collection, bson.M{
		"$or": []bson.M{
			{"access.admin": bson.M{"$in": groups}},
			{"access.member": bson.M{"$in": groups}},
			{"access.viewer": bson.M{"$in": groups}},
		},
	}, func(c database.Decodable) error {
		t := Team{}
		err := c.Decode(&t)
		if err != nil {
			return err
		}

		if containsAny(t.Access.Admin, groups) {
			res.Admin = append(res.Admin, t)
		} else if containsAny(t.Access.Member, groups) {
			res.Member = append(res.Member, t)
		} else if containsAny(t.Access.Viewer, groups) {
			res.Viewer = append(res.Viewer, t)
		}
		return nil
	})

	return res, err
}

func containsAny(arr1, arr2 []string) bool {
	for _, value := range arr2 {
		if contains(arr1, value) {
			return true
		}
	}
	return false
}

func contains(arr []string, value string) bool {
	for _, s := range arr {
		if s == value {
			return true
		}
	}
	return false
}

func (s *service) EnsureIndices() error {
	return s.db.EnsureIndex(Collection, mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
}

func (s *service) ListTeamsPaginated(perPage int, page int) (TeamPage, error) {
	var res TeamPage
	err := s.db.ListPaginated(Collection, perPage, page, &res.Pagination, func(c database.Decodable) error {
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

func (s *service) CreateTeam(id string, data TeamSettings) error {
	return s.db.InsertOne(Collection, bson.M{"id": id}, Team{
		Id:           id,
		TeamSettings: data,
	})
}

func (s *service) UpdateTeam(id string, data TeamSettings) error {
	return s.db.UpdateOneById(Collection, id, false, Team{
		Id:           id,
		TeamSettings: data,
	}, nil)
}
