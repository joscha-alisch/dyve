package routing

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
)

type Service interface {
	GetRoutes(app string) (sdk.AppRouting, error)
	UpdateRoutes(app string, routes sdk.AppRouting) error
}

const Collection = "routing"

func NewService(db database.Database) Service {
	return &service{
		db: db,
	}
}

type service struct {
	db database.Database
}

type routeData struct {
	Id        string         `bson:"id"`
	RouteData sdk.AppRouting `json:"routeData"`
}

func (s *service) GetRoutes(app string) (sdk.AppRouting, error) {
	res := routeData{}

	err := s.db.FindOneById(Collection, app, &res)
	return res.RouteData, err
}

func (s *service) UpdateRoutes(app string, routes sdk.AppRouting) error {
	return s.db.UpdateOne(Collection, bson.M{"id": app}, true, routeData{
		Id:        app,
		RouteData: routes,
	}, nil)
}
