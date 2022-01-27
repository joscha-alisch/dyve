package cloudfoundry

import (
	"context"
	"errors"
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
	ctx := context.Background()
	c, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(l.Uri),
	)
	if err != nil {
		return nil, err
	}

	err = c.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	db := c.Database(l.DB)

	m := &mongoDatabase{
		ctx:     ctx,
		cli:     c,
		db:      db,
		cfInfos: db.Collection("cf_infos"),
		orgs:    db.Collection("orgs"),
		spaces:  db.Collection("spaces"),
		apps:    db.Collection("apps"),
		cache:   db.Collection("cache"),
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
	cache   *mongo.Collection
	ctx     context.Context
}

func (d *mongoDatabase) Cached(id string, duration time.Duration, cached interface{}, f func() (interface{}, error)) (interface{}, error) {
	cacheTime := currentTime().Add(-duration)
	cacheRes := d.cache.FindOne(d.ctx, bson.M{"id": id, "last": bson.M{"$gte": cacheTime}})

	err := cacheRes.Err()
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		return nil, err
	} else if err == nil {
		return nil, cacheRes.Decode(cached)
	}

	data, err := f()
	if err != nil {
		return nil, err
	}

	_, err = d.cache.UpdateOne(d.ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"id": id, "last": currentTime(), "data": data}}, options.Update().SetUpsert(true))
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (d *mongoDatabase) GetApp(id string) (App, error) {
	res := d.apps.FindOne(d.ctx, bson.M{
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
	c, err := d.apps.Find(d.ctx, filter, options)
	if err != nil {
		return nil, err
	}

	var apps []App
	for c.Next(d.ctx) {
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
	res, err := d.cfInfos.UpdateOne(d.ctx, bson.M{
		"guid": bson.M{
			"$eq": cfGuid,
		},
	}, bson.M{
		"$set": bson.M{
			"orgs":        orgGuids,
			"lastUpdated": currentTime(),
		},
	})

	if res.MatchedCount == 0 {
		return errNotFound
	}

	return err
}

func (d *mongoDatabase) updateOrg(orgGuid string, spaceGuids []string) error {
	_, err := d.orgs.UpdateOne(d.ctx, bson.M{
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
	_, err := d.spaces.UpdateOne(d.ctx, bson.M{
		"guid": spaceGuid,
	}, bson.M{
		"$set": bson.M{
			"apps":        appGuids,
			"lastUpdated": currentTime(),
		},
	})

	return err
}

func (d *mongoDatabase) DeleteApp(guid string) (bool, error) {
	return d.deleteByGuid(d.apps, guid, "")
}

func (d *mongoDatabase) DeleteSpace(guid string) (bool, error) {
	deletedSpace, err := d.deleteByGuid(d.spaces, guid, "")
	if err != nil {
		return deletedSpace, err
	}
	deletedApps, err := d.deleteByGuid(d.apps, guid, "space.")
	return deletedApps || deletedSpace, err
}

func (d *mongoDatabase) DeleteOrg(guid string) (bool, error) {
	deletedOrg, err := d.deleteByGuid(d.orgs, guid, "")
	if err != nil {
		return deletedOrg, err
	}
	deletedSpaces, err := d.deleteByGuid(d.spaces, guid, "org.")
	if err != nil {
		return deletedOrg || deletedSpaces, err
	}
	deletedApps, err := d.deleteByGuid(d.apps, guid, "space.org.")
	return deletedOrg || deletedSpaces || deletedApps, err
}

func (d *mongoDatabase) deleteByGuid(coll *mongo.Collection, guid string, nested string) (bool, error) {
	return d.deleteBy(coll, bson.M{nested + "guid": guid})
}

func (d *mongoDatabase) deleteBy(coll *mongo.Collection, filter bson.M) (bool, error) {
	res, err := coll.DeleteMany(d.ctx, filter)
	if err != nil {
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (d *mongoDatabase) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	t := currentTime()

	j, ok := d.acceptCollectionReconcileJob(ReconcileSpaces, d.orgs, t, olderThan)
	if ok {
		return j, true
	}

	j, ok = d.acceptCollectionReconcileJob(ReconcileApps, d.spaces, t, olderThan)
	if ok {
		return j, true
	}

	j, ok = d.acceptCollectionReconcileJob(ReconcileOrganizations, d.cfInfos, t, olderThan)
	if ok {
		return j, true
	}

	return recon.Job{}, false
}

func (d *mongoDatabase) acceptCollectionReconcileJob(typ recon.Type, coll *mongo.Collection, t time.Time, olderThan time.Duration) (recon.Job, bool) {
	lessThanTime := t.Add(-olderThan)
	res := coll.FindOneAndUpdate(d.ctx, bson.M{
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

	j.Type = typ

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
	_, err := c.ReplaceOne(d.ctx, bson.M{
		"guid": guid,
	}, o, options.Replace().SetUpsert(true))

	return err
}

func (d *mongoDatabase) upsertSpaces(spaces []Space) error {
	for _, s := range spaces {
		_, err := d.spaces.UpdateOne(d.ctx, bson.M{
			"guid": bson.M{
				"$eq": s.Guid,
			},
		}, bson.M{
			"$set": s.SpaceInfo,
		}, options.Update().SetUpsert(true))

		_, err = d.apps.UpdateMany(d.ctx, bson.M{
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
		_, err := d.orgs.UpdateOne(d.ctx, bson.M{
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
	_, err := d.cfInfos.UpdateOne(d.ctx, bson.M{
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
		_, err := d.apps.UpdateOne(d.ctx, bson.M{
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
	res := d.spaces.FindOne(d.ctx, bson.M{
		"guid": spaceGuid,
	})
	if res.Err() != nil {
		return Space{}, errNotFound
	}

	s := Space{}
	err := res.Decode(&s)
	if err != nil {
		return Space{}, errDecode
	}

	return s, nil
}

func (d *mongoDatabase) getOrg(orgGuid string) (Org, error) {
	res := d.orgs.FindOne(d.ctx, bson.M{
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
	res := d.cfInfos.FindOne(d.ctx, bson.M{
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

	_, err := c.DeleteMany(d.ctx, filter)
	return err
}
