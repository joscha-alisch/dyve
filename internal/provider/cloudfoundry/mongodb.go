package cloudfoundry

import (
	"context"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoLogin struct {
	Uri string
	DB  string
}

func NewMongoDatabase(l MongoLogin) (Database, error) {
	c, err := mongo.Connect(
		context.Background(),
		options.Client().ApplyURI(l.Uri),
	)
	if err != nil {
		return nil, err
	}

	err = c.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	db := c.Database(l.DB)

	m := &mongoDatabase{
		cli:     c,
		db:      db,
		cfInfos: db.Collection("cf_infos"),
		orgs:    db.Collection("orgs"),
		spaces:  db.Collection("spaces"),
		apps:    db.Collection("apps"),
	}

	err = m.setupBaseJob()
	return m, err
}

type mongoDatabase struct {
	cli     *mongo.Client
	db      *mongo.Database
	orgs    *mongo.Collection
	spaces  *mongo.Collection
	apps    *mongo.Collection
	cfInfos *mongo.Collection
}

func (d *mongoDatabase) GetApp(id string) (App, error) {
	res := d.apps.FindOne(context.Background(), bson.M{
		"guid": bson.M{
			"$eq": id,
		},
	})
	if res.Err() != nil {
		return App{}, res.Err()
	}

	a := App{}
	err := res.Decode(&a)
	if err != nil {
		return App{}, err
	}

	return a, nil
}

func (d *mongoDatabase) ListApps() ([]App, error) {
	return d.getApps(bson.M{}, options.Find().
		SetSort(bson.M{"guid": 1}))
}

func (d *mongoDatabase) getApps(filter bson.M, options *options.FindOptions) ([]App, error) {
	c, err := d.apps.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}

	var apps []App
	for c.Next(context.Background()) {
		app := App{}
		err = c.Decode(&app)
		if err != nil {
			return nil, err
		}
		apps = append(apps, app)
	}
	return apps, nil
}

var currentTime = time.Now

func (d *mongoDatabase) UpsertOrgs(cfGuid string, orgs []Org) error {
	var orgGuids []string
	for _, org := range orgs {
		orgGuids = append(orgGuids, org.Guid)
	}

	err := d.updateCfInfo(cfGuid, orgGuids)
	if err != nil {
		return err
	}

	cf, err := d.getCf(cfGuid)
	if err != nil {
		return err
	}

	for i := range orgs {
		orgs[i].Cf = cf.CFInfo
	}

	err = d.upsertOrgs(orgs)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.orgs, "cf.guid", cfGuid, "guid", orgGuids)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.spaces, "org.cf.guid", cfGuid, "org.guid", orgGuids)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.apps, "space.org.cf.guid", cfGuid, "space.org.guid", orgGuids)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDatabase) updateCfInfo(cfGuid string, orgGuids []string) error {
	_, err := d.cfInfos.UpdateOne(context.Background(), bson.M{
		"guid": bson.M{
			"$eq": cfGuid,
		},
	}, bson.M{
		"$set": bson.M{
			"orgs":        orgGuids,
			"lastUpdated": currentTime(),
		},
	})

	return err
}

func (d *mongoDatabase) updateOrg(orgGuid string, spaceGuids []string) error {
	_, err := d.orgs.UpdateOne(context.Background(), bson.M{
		"guid": orgGuid,
	}, bson.M{
		"$set": bson.M{
			"spaces":      spaceGuids,
			"lastUpdated": currentTime(),
		},
	})

	return err
}

func (d *mongoDatabase) updateSpace(spaceGuid string, appGuids []string) error {
	_, err := d.spaces.UpdateOne(context.Background(), bson.M{
		"guid": spaceGuid,
	}, bson.M{
		"$set": bson.M{
			"apps":        appGuids,
			"lastUpdated": currentTime(),
		},
	})

	return err
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

func (d *mongoDatabase) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	t := currentTime()

	j, ok := d.acceptCollectionReconcileJob(d.orgs, t, olderThan)
	if ok {
		j.Type = ReconcileSpaces
		return j, true
	}

	j, ok = d.acceptCollectionReconcileJob(d.spaces, t, olderThan)
	if ok {
		j.Type = ReconcileApps
		return j, true
	}

	j, ok = d.acceptCollectionReconcileJob(d.cfInfos, t, olderThan)
	if ok {
		j.Type = ReconcileOrganizations
		return j, true
	}

	return recon.Job{}, false
}

func (d *mongoDatabase) acceptCollectionReconcileJob(coll *mongo.Collection, t time.Time, olderThan time.Duration) (recon.Job, bool) {
	lessThanTime := t.Add(-olderThan)
	res := coll.FindOneAndUpdate(context.Background(), bson.M{
		"$or": bson.A{
			bson.M{
				"lastUpdated": bson.M{"$lte": lessThanTime},
			},
			bson.M{"lastUpdated": nil},
		},
	}, bson.M{
		"$set": bson.M{
			"lastUpdated": t,
		},
	}, options.FindOneAndUpdate().SetSort(bson.D{{"lastUpdated", 1}}))

	j := recon.Job{}
	err := res.Decode(&j)
	if err != nil {
		return recon.Job{}, false
	}
	return j, true
}

func (d *mongoDatabase) UpsertOrgSpaces(orgGuid string, spaces []Space) error {
	org, err := d.getOrg(orgGuid)
	if err != nil {
		return err
	}

	var spaceGuids []string
	for i, space := range spaces {
		spaceGuids = append(spaceGuids, space.Guid)
		spaces[i].Org = org.OrgInfo
	}

	err = d.updateOrg(orgGuid, spaceGuids)
	if err != nil {
		return err
	}

	err = d.upsertSpaces(spaces)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.spaces, "org.guid", orgGuid, "guid", spaceGuids)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.apps, "space.org.guid", orgGuid, "space.guid", spaceGuids)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDatabase) UpsertSpaceApps(spaceGuid string, apps []App) error {
	space, err := d.getSpace(spaceGuid)
	if err != nil {
		return err
	}

	var appGuids []string
	for i, app := range apps {
		appGuids = append(appGuids, app.Guid)
		apps[i].AppInfo.Space = space.SpaceInfo
	}

	err = d.updateSpace(space.Guid, appGuids)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.apps, "space.guid", space.Guid, "guid", appGuids)
	if err != nil {
		return err
	}

	err = d.upsertApps(apps)
	if err != nil {
		return err
	}

	return nil
}

func (d *mongoDatabase) upsertByGuid(c *mongo.Collection, guid string, o interface{}) error {
	_, err := c.ReplaceOne(context.Background(), bson.M{
		"guid": guid,
	}, o, options.Replace().SetUpsert(true))

	return err
}

func (d *mongoDatabase) upsertSpaces(spaces []Space) error {
	for _, s := range spaces {
		_, err := d.spaces.UpdateOne(context.Background(), bson.M{
			"guid": bson.M{
				"$eq": s.Guid,
			},
		}, bson.M{
			"$set": s.SpaceInfo,
		}, options.Update().SetUpsert(true))

		_, err = d.apps.UpdateMany(context.Background(), bson.M{
			"space.guid": bson.M{
				"$eq": s.Guid,
			},
		}, bson.M{
			"$set": bson.M{
				"space": s.SpaceInfo,
			},
		})

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *mongoDatabase) upsertOrgs(orgs []Org) error {
	for _, o := range orgs {
		_, err := d.orgs.UpdateOne(context.Background(), bson.M{
			"guid": o.Guid,
		}, bson.M{
			"$set": o.OrgInfo,
		}, options.Update().SetUpsert(true))

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *mongoDatabase) setupBaseJob() error {
	_, err := d.cfInfos.UpdateOne(context.Background(), bson.M{
		"guid": bson.M{
			"$eq": CFGuid,
		},
	}, bson.M{
		"$set": bson.M{
			"guid": CFGuid,
		},
	}, options.Update().SetUpsert(true))

	return err
}

func (d *mongoDatabase) upsertApps(apps []App) error {
	for _, app := range apps {
		_, err := d.apps.UpdateOne(context.Background(), bson.M{
			"guid": bson.M{
				"$eq": app.Guid,
			},
		}, bson.M{
			"$set": app.AppInfo,
		}, options.Update().SetUpsert(true))

		if err != nil {
			return err
		}
	}

	return nil
}

func (d *mongoDatabase) getSpace(spaceGuid string) (Space, error) {
	res := d.spaces.FindOne(context.Background(), bson.M{
		"guid": spaceGuid,
	})
	if res.Err() != nil {
		return Space{}, res.Err()
	}

	s := Space{}
	err := res.Decode(&s)
	if err != nil {
		return Space{}, err
	}

	return s, nil
}

func (d *mongoDatabase) getOrg(orgGuid string) (Org, error) {
	res := d.orgs.FindOne(context.Background(), bson.M{
		"guid": orgGuid,
	})
	if res.Err() != nil {
		return Org{}, errNotFound
	}

	o := Org{}
	err := res.Decode(&o)
	if err != nil {
		return Org{}, err
	}

	return o, nil
}

func (d *mongoDatabase) getCf(cfGuid string) (CF, error) {
	res := d.cfInfos.FindOne(context.Background(), bson.M{
		"guid": bson.M{
			"$eq": cfGuid,
		},
	})
	if res.Err() != nil {
		return CF{}, errNotFound
	}

	cf := CF{}
	err := res.Decode(&cf)
	if err != nil {
		return CF{}, err
	}

	return cf, nil
}

func (d *mongoDatabase) removeOutdatedIn(c *mongo.Collection, where, equals, and string, notIn []string) error {
	if notIn == nil {
		notIn = []string{}
	}

	filter := bson.M{
		where: bson.M{
			"$eq": equals,
		},
		and: bson.M{
			"$nin": notIn,
		},
	}

	_, err := c.DeleteMany(context.Background(), filter)
	return err
}
