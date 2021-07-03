package cloudfoundry

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
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

func (d *mongoDatabase) DeleteApp(guid string) {
	d.deleteByGuid(d.apps, guid)
}

func (d *mongoDatabase) DeleteSpace(guid string) {
	d.deleteByGuid(d.spaces, guid)
	d.deleteBySpace(d.apps, guid)
}

func (d *mongoDatabase) DeleteOrg(guid string) {
	d.deleteByGuid(d.orgs, guid)
	d.deleteByOrg(d.spaces, guid)
	d.deleteByOrg(d.apps, guid)
}

func (d *mongoDatabase) deleteByGuid(coll *mongo.Collection, guid string) (bool, error) {
	return d.deleteBy(coll, bson.M{"guid": bson.M{"$eq": guid}})
}

func (d *mongoDatabase) deleteByOrg(coll *mongo.Collection, guid string) (bool, error) {
	return d.deleteBy(coll, bson.M{"org": bson.M{"$eq": guid}})
}

func (d *mongoDatabase) deleteBySpace(coll *mongo.Collection, guid string) (bool, error) {
	return d.deleteBy(coll, bson.M{"space": bson.M{"$eq": guid}})
}

func (d *mongoDatabase) deleteBy(coll *mongo.Collection, filter bson.M) (bool, error) {
	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (d *mongoDatabase) AcceptReconcileJob(olderThan time.Time, againAt time.Time) (ReconcileJob, bool) {
	j, ok := d.acceptCollectionReconcileJob(d.orgs, olderThan, againAt)
	if ok {
		j.Type = ReconcileOrg
		return j, true
	}

	j, ok = d.acceptCollectionReconcileJob(d.spaces, olderThan, againAt)
	if ok {
		j.Type = ReconcileSpace
		return j, true
	}

	return ReconcileJob{}, false
}

func (d *mongoDatabase) acceptCollectionReconcileJob(coll *mongo.Collection, olderThan time.Time, againAt time.Time) (ReconcileJob, bool) {
	res := coll.FindOneAndUpdate(context.Background(), bson.M{
		"lastUpdated": bson.M{
			"$lte": olderThan,
		},
	}, bson.M{
		"$set": bson.M{
			"lastUpdated": againAt,
		},
	}, options.FindOneAndUpdate().SetSort(bson.D{{"lastUpdated", 1}}))

	j := ReconcileJob{}
	err := res.Decode(&j)
	if err != nil {
		return ReconcileJob{}, false
	}
	return j, true
}

func (d *mongoDatabase) UpsertOrg(o Org) error {
	err := d.upsertByGuid(d.orgs, o.Guid, o)
	if err != nil {
		return err
	}

	err = d.removeOutdatedSpaces(o)
	if err != nil {
		return err
	}

	err = d.removeOutdatedOrgApps(o)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDatabase) UpsertSpace(s Space) error {
	err := d.upsertByGuid(d.spaces, s.Guid, s)
	if err != nil {
		return err
	}

	err = d.removeOutdatedSpaceApps(s)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDatabase) UpsertApps(apps []App) error {
	for _, app := range apps {
		err := d.upsertByGuid(d.apps, app.Guid, app)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *mongoDatabase) removeOutdatedSpaces(org Org) error {
	if org.Spaces == nil {
		org.Spaces = []string{}
	}
	_, err := d.spaces.DeleteMany(context.Background(), bson.M{
		"org": bson.M{
			"$eq": org.Guid,
		},
		"guid": bson.M{
			"$nin": org.Spaces,
		},
	})
	return err
}

func (d *mongoDatabase) removeOutdatedOrgApps(org Org) error {
	if org.Spaces == nil {
		org.Spaces = []string{}
	}

	_, err := d.apps.DeleteMany(context.Background(), bson.M{
		"org": bson.M{
			"$eq": org.Guid,
		},
		"space": bson.M{
			"$nin": org.Spaces,
		},
	})
	return err
}

func (d *mongoDatabase) upsertByGuid(c *mongo.Collection, guid string, o interface{}) error {
	_, err := c.ReplaceOne(context.Background(), bson.M{
		"guid": guid,
	}, o, options.Replace().SetUpsert(true))

	return err
}

func (d *mongoDatabase) removeOutdatedSpaceApps(s Space) error {
	if s.Apps == nil {
		s.Apps = []string{}
	}
	_, err := d.apps.DeleteMany(context.Background(), bson.M{
		"space": bson.M{
			"$eq": s.Guid,
		},
		"guid": bson.M{
			"$nin": s.Apps,
		},
	})
	return err
}
