package database

import (
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Database interface {
	FindOne(coll Collection, filter bson.M, res interface{}) error
	FindOneById(coll Collection, id string, res interface{}) error
	FindOneSorted(coll Collection, filter bson.M, sort bson.M, res interface{}) error
	FindMany(coll Collection, filter bson.M, each func(c *mongo.Cursor) error) error
	FindManyWithOptions(coll Collection, filter bson.M, each func(c *mongo.Cursor) error, sort bson.M, limit int) error
	ListPaginated(coll Collection, perPage int, page int, p *sdk.Pagination, each func(c *mongo.Cursor) error) error

	UpdateProvided(coll Collection, provider string, updates map[string]interface{}) error
	UpdateMany(coll Collection, filters map[string]interface{}, updates map[string]interface{}) error
	UpdateOne(coll Collection, filter bson.M, createIfMissing bool, update interface{}, res interface{}) error
	UpdateOneById(coll Collection, id string, createIfMissing bool, update interface{}, res interface{}) error

	DeleteOne(coll Collection, filter bson.M) error
	DeleteOneById(coll Collection, id string) error

	EnsureIndex(coll Collection, model mongo.IndexModel) error
}

const (
	ReconcileAppProvider      recon.Type = "apps"
	ReconcilePipelineProvider recon.Type = "pipelines"
	ReconcileGroupProvider    recon.Type = "groups"
)

type Collection string
