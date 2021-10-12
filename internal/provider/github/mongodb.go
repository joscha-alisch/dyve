package github

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

func NewMongoDatabase(l MongoLogin, org string) (Database, error) {
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
		cli:   c,
		db:    db,
		orgs:  db.Collection("orgs"),
		teams: db.Collection("teams"),
	}

	err = m.setupBaseJob(org)

	return m, err
}

type mongoDatabase struct {
	cli   *mongo.Client
	db    *mongo.Database
	teams *mongo.Collection
	orgs  *mongo.Collection
}

func (d *mongoDatabase) UpdateTeamMembers(guid string, members []Member) error {
	res, err := d.teams.UpdateOne(context.Background(), bson.M{
		"guid": guid,
	}, bson.M{
		"$set": bson.M{
			"members": members,
		},
	})

	if res != nil && res.MatchedCount != 1 {
		return errNotFound
	}

	return err
}

func (d *mongoDatabase) GetTeam(guid string) (Team, error) {
	res := d.teams.FindOne(context.Background(), bson.M{
		"guid": guid,
	})
	if res.Err() != nil {
		return Team{}, errNotFound
	}

	team := Team{}
	err := res.Decode(&team)
	if err != nil {
		return Team{}, err
	}

	return team, nil
}

func (d *mongoDatabase) ListTeams() ([]Team, error) {
	return d.getTeams(bson.M{}, options.Find().
		SetSort(bson.M{"guid": 1}))
}

func (d *mongoDatabase) getTeams(filter bson.M, options *options.FindOptions) ([]Team, error) {
	c, err := d.teams.Find(context.Background(), filter, options)
	if err != nil {
		return nil, err
	}

	var teams []Team
	for c.Next(context.Background()) {
		team := Team{}
		err = c.Decode(&team)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}
	return teams, nil
}

var currentTime = time.Now

func (d *mongoDatabase) deleteBy(coll *mongo.Collection, filter bson.M) (bool, error) {
	res, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}

	return res.DeletedCount > 0, nil
}

func (d *mongoDatabase) AcceptReconcileJob(olderThan time.Duration) (recon.Job, bool) {
	t := currentTime()

	j, ok := d.acceptCollectionReconcileJob(ReconcileTeams, d.orgs, t, olderThan)
	if ok {
		return j, true
	}

	j, ok = d.acceptCollectionReconcileJob(ReconcileMembers, d.teams, t, olderThan)
	if ok {
		return j, true
	}

	return recon.Job{}, false
}

func (d *mongoDatabase) acceptCollectionReconcileJob(typ recon.Type, coll *mongo.Collection, t time.Time, olderThan time.Duration) (recon.Job, bool) {
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

	j.Type = typ

	return j, true
}

func (d *mongoDatabase) UpsertOrgTeams(orgGuid string, teams []Team) error {
	org, err := d.getOrg(orgGuid)
	if err != nil {
		return err
	}

	var teamIds []string
	for i, team := range teams {
		teamIds = append(teamIds, team.Guid)
		teams[i].Org = org.OrgInfo
	}

	err = d.updateOrg(orgGuid, teamIds)
	if err != nil {
		return err
	}

	err = d.upsertTeams(teams)
	if err != nil {
		return err
	}

	err = d.removeOutdatedIn(d.teams, "org.guid", orgGuid, "guid", teamIds)
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

func (d *mongoDatabase) upsertTeams(teams []Team) error {
	for _, s := range teams {
		_, err := d.teams.UpdateOne(context.Background(), bson.M{
			"guid": bson.M{
				"$eq": s.Guid,
			},
		}, bson.M{
			"$set": s.TeamInfo,
		}, options.Update().SetUpsert(true))

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *mongoDatabase) setupBaseJob(org string) error {
	_, err := d.orgs.UpdateOne(context.Background(), bson.M{
		"guid": bson.M{
			"$eq": org,
		},
	}, bson.M{
		"$set": bson.M{
			"guid": org,
		},
	}, options.Update().SetUpsert(true))

	return err
}

func (d *mongoDatabase) updateOrg(orgGuid string, teamGuids []string) error {
	_, err := d.orgs.UpdateOne(context.Background(), bson.M{
		"guid": orgGuid,
	}, bson.M{
		"$set": bson.M{
			"teams":       teamGuids,
			"lastUpdated": currentTime(),
		},
	})

	return err
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
