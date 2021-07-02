package cloudfoundry

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(uri string, dbName string) (Database, error) {
	c, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(uri),
	)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	db := c.Database(dbName)

	return &mongoDatabase{
		cli:    c,
		db:     db,
		orgs:   db.Collection("orgs"),
		spaces: db.Collection("spaces"),
		apps:   db.Collection("apps"),
	}, nil
}

type mongoDatabase struct {
	cli    *mongo.Client
	db     *mongo.Database
	orgs   *mongo.Collection
	spaces *mongo.Collection
	apps   *mongo.Collection
}

func (d *mongoDatabase) FetchReconcileJob() (ReconcileJob, bool) {
	return ReconcileJob{}, false
}

func (d *mongoDatabase) UpsertOrg(o Org) error {
	_, err := d.orgs.ReplaceOne(context.Background(), bson.M{
		"guid": o.Guid,
	}, o, options.Replace().SetUpsert(true))

	return err
}

func (d *mongoDatabase) UpsertSpace(s Space) error {
	_, err := d.spaces.ReplaceOne(context.Background(), bson.M{
		"guid": s.Guid,
	}, s, options.Replace().SetUpsert(true))

	return err
}

func (d *mongoDatabase) UpsertApp(a App) error {
	_, err := d.apps.ReplaceOne(context.Background(), bson.M{
		"guid": a.Guid,
	}, a, options.Replace().SetUpsert(true))

	return err
}
