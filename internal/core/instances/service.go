package instances

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetInstances(app string) (sdk.AppInstances, error)
	UpdateInstances(app string, instances sdk.AppInstances) error
}

const Collection = "instances"

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

func (s *service) GetInstances(app string) (sdk.AppInstances, error) {
	res := instancesData{}

	err := s.db.FindOneById(Collection, app, &res)
	return res.InstancesData, err
}

func (s *service) UpdateInstances(app string, instances sdk.AppInstances) error {
	return s.db.UpdateOne(Collection, bson.M{"id": app}, true, instancesData{
		Id:            app,
		InstancesData: instances,
	}, nil)
}

type instancesData struct {
	Id            string           `bson:"id"`
	InstancesData sdk.AppInstances `bson:"instancesData"`
}
