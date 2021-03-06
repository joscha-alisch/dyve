package cloudfoundry

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/rs/zerolog/log"
	"github.com/tryvium-travels/memongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

var baseState = map[string]interface{}{
	"cf_infos": []bson.M{
		{"guid": "main", "orgs": []string{"org-a-guid", "org-b-guid"}},
	},
	"orgs": []bson.M{
		{"name": "org-a-name", "guid": "org-a-guid", "cf": bson.M{"guid": "main"}, "spaces": []string{"space-a-guid", "space-b-guid"}},
		{"name": "org-b-name", "guid": "org-b-guid", "cf": bson.M{"guid": "main"}, "spaces": []string{"space-c-guid", "space-d-guid"}},
	},
	"spaces": []bson.M{
		{"name": "space-a-name", "guid": "space-a-guid", "org": bson.M{"guid": "org-a-guid", "name": "org-a-name", "cf": bson.M{"guid": "main"}}, "apps": []string{"app-a-guid", "app-b-guid"}},
		{"name": "space-b-name", "guid": "space-b-guid", "org": bson.M{"guid": "org-a-guid", "name": "org-a-name", "cf": bson.M{"guid": "main"}}, "apps": []string{"app-c-guid", "app-d-guid"}},
		{"name": "space-c-name", "guid": "space-c-guid", "org": bson.M{"guid": "org-b-guid", "name": "org-b-name", "cf": bson.M{"guid": "main"}}, "apps": []string{"app-e-guid", "app-f-guid"}},
		{"name": "space-d-name", "guid": "space-d-guid", "org": bson.M{"guid": "org-b-guid", "name": "org-b-name", "cf": bson.M{"guid": "main"}}, "apps": []string{"app-g-guid", "app-h-guid"}},
	},
	"apps": []bson.M{
		{"name": "app-a-name", "guid": "app-a-guid", "space": bson.M{"guid": "space-a-guid", "name": "space-a-name", "org": bson.M{"guid": "org-a-guid", "name": "org-a-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-b-name", "guid": "app-b-guid", "space": bson.M{"guid": "space-a-guid", "name": "space-a-name", "org": bson.M{"guid": "org-a-guid", "name": "org-a-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-c-name", "guid": "app-c-guid", "space": bson.M{"guid": "space-b-guid", "name": "space-b-name", "org": bson.M{"guid": "org-a-guid", "name": "org-a-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-d-name", "guid": "app-d-guid", "space": bson.M{"guid": "space-b-guid", "name": "space-b-name", "org": bson.M{"guid": "org-a-guid", "name": "org-a-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-e-name", "guid": "app-e-guid", "space": bson.M{"guid": "space-c-guid", "name": "space-c-name", "org": bson.M{"guid": "org-b-guid", "name": "org-b-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-f-name", "guid": "app-f-guid", "space": bson.M{"guid": "space-c-guid", "name": "space-c-name", "org": bson.M{"guid": "org-b-guid", "name": "org-b-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-g-name", "guid": "app-g-guid", "space": bson.M{"guid": "space-d-guid", "name": "space-d-name", "org": bson.M{"guid": "org-b-guid", "name": "org-b-name", "cf": bson.M{"guid": "main"}}}},
		{"name": "app-h-name", "guid": "app-h-guid", "space": bson.M{"guid": "space-d-guid", "name": "space-d-name", "org": bson.M{"guid": "org-b-guid", "name": "org-b-name", "cf": bson.M{"guid": "main"}}}},
	},
}

func TestMongoIntegration(t *testing.T) {
	currentTime = func() time.Time {
		return someTime
	}

	tests := []struct {
		desc  string
		f     func(db Database, tt *testing.T) error
		err   error
		state bson.M
	}{
		{desc: "updates space apps", state: baseState, f: func(db Database, tt *testing.T) error {
			return db.UpsertSpaceApps("space-a-guid", []App{
				{AppInfo: AppInfo{Name: "changed-name", Guid: "app-a-guid"}},
				{AppInfo: AppInfo{Name: "new-app", Guid: "new-app-guid"}},
			})
		}},
		{desc: "fails update apps for non-existent space", state: map[string]interface{}{}, f: func(db Database, tt *testing.T) error {
			return db.UpsertSpaceApps("space-a-guid", []App{
				{AppInfo: AppInfo{Name: "changed-name", Guid: "app-a-guid"}},
				{AppInfo: AppInfo{Name: "new-app", Guid: "new-app-guid"}},
			})
		}, err: errNotFound},
		{desc: "fails update apps for nonsense space data", state: map[string]interface{}{
			"spaces": []bson.M{
				{"name": 1, "guid": "space-a-guid", "org": []string{"hi"}},
			},
		}, f: func(db Database, tt *testing.T) error {
			return db.UpsertSpaceApps("space-a-guid", []App{
				{AppInfo: AppInfo{Name: "changed-name", Guid: "app-a-guid"}},
				{AppInfo: AppInfo{Name: "new-app", Guid: "new-app-guid"}},
			})
		}, err: errDecode},
		{desc: "updates org spaces", state: baseState, f: func(db Database, tt *testing.T) error {
			return db.UpsertOrgSpaces("org-a-guid", []Space{
				{SpaceInfo: SpaceInfo{Name: "changed-name", Guid: "space-a-guid"}},
				{SpaceInfo: SpaceInfo{Name: "new-space", Guid: "new-space-guid"}},
			})
		}},
		{desc: "updates cf orgs", state: baseState, f: func(db Database, tt *testing.T) error {
			return db.UpsertOrgs("main", []Org{
				{OrgInfo: OrgInfo{Name: "changed-name", Guid: "org-a-guid"}},
				{OrgInfo: OrgInfo{Name: "new-org", Guid: "new-org-guid"}},
			})
		}},
		{desc: "updates cf orgs for non-existent cf", state: baseState, f: func(db Database, tt *testing.T) error {
			return db.UpsertOrgs("doesnt-exist", []Org{
				{OrgInfo: OrgInfo{Name: "changed-name", Guid: "org-a-guid"}},
				{OrgInfo: OrgInfo{Name: "new-org", Guid: "new-org-guid"}},
			})
		}, err: errNotFound},
		{desc: "updates empty cf orgs list", state: baseState, f: func(db Database, tt *testing.T) error {
			return db.UpsertOrgs("main", nil)
		}},
		{desc: "updates cf orgs for missing collection", state: map[string]interface{}{}, f: func(db Database, tt *testing.T) error {
			return db.UpsertOrgs("doesnt-exist", []Org{
				{OrgInfo: OrgInfo{Name: "changed-name", Guid: "org-a-guid"}},
				{OrgInfo: OrgInfo{Name: "new-org", Guid: "new-org-guid"}},
			})
		}, err: errNotFound},
		{desc: "lists apps", state: map[string]interface{}{
			"apps": []bson.M{
				{"name": "app-a-name", "guid": "app-a-guid"},
				{"name": "app-b-name", "guid": "app-b-guid"},
			},
		}, f: func(db Database, tt *testing.T) error {
			apps, err := db.ListApps()
			if err != nil {
				return err
			}
			expected := []App{
				{AppInfo{Name: "app-a-name", Guid: "app-a-guid"}},
				{AppInfo{Name: "app-b-name", Guid: "app-b-guid"}},
			}

			if !cmp.Equal(expected, apps) {
				tt.Errorf("wrong apps returned:\n%s\n", cmp.Diff(expected, apps))
			}
			return nil
		}},
		{desc: "gets app", state: map[string]interface{}{
			"apps": []bson.M{
				{"name": "app-a-name", "guid": "app-a-guid"},
				{"name": "app-b-name", "guid": "app-b-guid"},
			},
		}, f: func(db Database, tt *testing.T) error {
			app, err := db.GetApp("app-a-guid")
			if err != nil {
				return err
			}
			expected := App{AppInfo{Name: "app-a-name", Guid: "app-a-guid"}}

			if !cmp.Equal(expected, app) {
				tt.Errorf("wrong app returned:\n%s\n", cmp.Diff(expected, app))
			}
			return nil
		}},
		{desc: "deletes app", state: baseState, f: func(db Database, tt *testing.T) error {
			_, err := db.DeleteApp("app-a-guid")
			return err
		}},
		{desc: "deletes space", state: baseState, f: func(db Database, tt *testing.T) error {
			_, err := db.DeleteSpace("space-a-guid")
			return err
		}},
		{desc: "deletes org", state: baseState, f: func(db Database, tt *testing.T) error {
			_, err := db.DeleteOrg("org-a-guid")
			return err
		}},
		{desc: "fetch org job", state: bson.M{
			"orgs": []bson.M{
				{"name": "b", "guid": "def", "lastUpdated": someTime.Add(-1 * time.Minute)},
				{"name": "a", "guid": "abc", "lastUpdated": someTime.Add(-3 * time.Minute)},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcileSpaces, Guid: "abc", LastUpdated: someTime.Add(-3 * time.Minute)}
			j, ok := db.AcceptReconcileJob(2 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch space job", state: bson.M{
			"spaces": []bson.M{
				{"name": "b", "guid": "def", "lastUpdated": someTime.Add(-1 * time.Minute)},
				{"name": "a", "guid": "abc", "lastUpdated": someTime.Add(-3 * time.Minute)},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcileApps, Guid: "abc", LastUpdated: someTime.Add(-3 * time.Minute)}
			j, ok := db.AcceptReconcileJob(2 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch cf info job", state: bson.M{
			"cf_infos": []bson.M{
				{"guid": "main", "lastUpdated": someTime.Add(-3 * time.Minute)},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcileOrganizations, Guid: "main", LastUpdated: someTime.Add(-3 * time.Minute)}
			j, ok := db.AcceptReconcileJob(2 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch org job never updated", state: bson.M{
			"orgs": []bson.M{
				{"name": "a", "guid": "abc"},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcileSpaces, Guid: "abc"}
			j, ok := db.AcceptReconcileJob(2 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "uses cache", state: bson.M{
			"cache": []bson.M{
				{"id": "a", "last": someTime.Add(-1 * time.Minute), "src": "cached"},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := cacheObj{
				Id:  "a",
				Src: "cached",
			}
			cached := cacheObj{}
			res, _ := db.Cached("a", 2*time.Minute, &cached, func() (interface{}, error) {
				return cacheObj{
					Id:  "a",
					Src: "func",
				}, nil
			})
			if !cmp.Equal(expected, cached) {
				tt.Errorf("wrong data returned:\n%s\n", cmp.Diff(expected, cached))
			}
			if res != nil {
				tt.Errorf("expected res to be nil")
			}
			return nil
		}},
		{desc: "doesnt use cache", state: bson.M{
			"cache": []bson.M{
				{"id": "a", "last": someTime.Add(-3 * time.Minute), "src": "cached"},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := cacheObj{
				Id:  "a",
				Src: "func",
			}
			cached := cacheObj{}
			res, _ := db.Cached("a", 2*time.Minute, &cached, func() (interface{}, error) {
				return cacheObj{
					Id:  "a",
					Src: "func",
				}, nil
			})
			if !cmp.Equal(expected, res) {
				tt.Errorf("wrong data returned:\n%s\n", cmp.Diff(expected, res))
			}
			return nil
		}},
	}

	opts := &memongo.Options{
		MongoVersion: "5.0.5",
	}
	if runtime.GOARCH == "arm64" {
		if runtime.GOOS == "darwin" {
			opts.DownloadURL = "https://fastdl.mongodb.org/osx/mongodb-macos-x86_64-5.0.5.tgz"
		}
	}

	mongodb, err := memongo.StartWithOptions(opts)
	if err != nil {
		t.Fatal(err)
	}
	defer mongodb.Stop()

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fileName := strings.ReplaceAll(test.desc, " ", "_")
			acceptanceTesting(fileName, test.state, test.f, test.err, mongodb, tt)
		})
	}
}

func acceptanceTesting(
	name string,
	state map[string]interface{},
	f func(db Database, tt *testing.T) error,
	expectedErr error,
	s *memongo.Server,
	tt *testing.T,
) {
	dbName := memongo.RandomDatabase()
	if state != nil {
		err := setState(state, s, dbName)
		if err != nil {
			tt.Fatal(err)
		}
	}

	db, err := NewMongoDatabase(MongoLogin{
		Uri: s.URI(),
		DB:  dbName,
	})
	if err != nil {
		tt.Fatal(err)
	}

	before, err := dumpContents(s, dbName)
	if err != nil {
		tt.Fatal(err)
	}
	walk(before, func(m map[string]interface{}, k string) {
		if t, ok := m[k].(primitive.DateTime); ok {
			m[k] = time.Unix(int64(t)/1000, 0).UTC().Format(time.RFC3339)
		}

		if t, ok := m[k].(primitive.A); ok {
			m[k] = ([]interface{})(t)
		}
	})

	err = f(db, tt)
	if !errors.Is(err, expectedErr) {
		tt.Fatalf("error mismatch: %s\n", cmp.Diff(expectedErr, err))
	}

	contents, err := dumpContents(s, dbName)
	if err != nil {
		tt.Fatal(err)
	}

	walk(contents, func(m map[string]interface{}, k string) {
		if t, ok := m[k].(primitive.DateTime); ok {
			m[k] = time.Unix(int64(t)/1000, 0).UTC().Format(time.RFC3339)
		}

		if t, ok := m[k].(primitive.A); ok {
			m[k] = ([]interface{})(t)
		}
	})

	_, testFilePath, _, _ := runtime.Caller(0)
	testDir := filepath.Dir(testFilePath)
	acceptedName := filepath.Join(testDir, "acceptance_tests", name+".accepted.json")
	actualName := filepath.Join(testDir, "acceptance_tests", name+".actual.json")

	acceptedContents := make(map[string]interface{})
	if _, err := os.Stat(acceptedName); !os.IsNotExist(err) {
		bytes, err := ioutil.ReadFile(acceptedName)
		if err != nil {
			tt.Fatal(err)
		}

		err = json.Unmarshal(bytes, &acceptedContents)
		if err != nil {
			tt.Fatal(err)
		}
	} else {
		log.Warn().Msg("first acceptance testing run. Diffing with 'before'-state")
		acceptedContents = before
	}

	if !cmp.Equal(acceptedContents, contents) {
		tt.Errorf("found diff between accepted and actual contents. Rename file to .accepted.json to accept changes:\n%s\n", cmp.Diff(acceptedContents, contents))

		bytes, err := json.MarshalIndent(contents, "", "    ")
		if err != nil {
			tt.Fatal("could not marshal actual into file")
		}

		_ = ioutil.WriteFile(actualName, bytes, 0666)
	} else {
		_ = os.Remove(actualName)
	}
}

func setState(data bson.M, s *memongo.Server, dbName string) error {
	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(s.URI()))
	if err != nil {
		return err
	}

	db := conn.Database(dbName)

	for coll, contents := range data {
		contentArr, ok := contents.([]bson.M)
		if !ok {
			return errors.New("data not in correct format")
		}

		for _, m := range contentArr {
			_, err := db.Collection(coll).InsertOne(context.Background(), m)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func dumpContents(s *memongo.Server, dbName string) (map[string]interface{}, error) {
	conn, err := mongo.Connect(context.Background(), options.Client().ApplyURI(s.URI()))
	if err != nil {
		return nil, err
	}

	db := conn.Database(dbName)
	colls, err := db.ListCollections(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})

	for colls.Next(context.Background()) {
		m := make(map[string]interface{})
		err = colls.Decode(&m)
		if err != nil {
			return nil, err
		}

		name := m["name"].(string)

		elem, err := db.Collection(name).Find(context.Background(), bson.D{})
		if err != nil {
			return nil, err
		}

		collContents := make([]interface{}, 0)
		for elem.Next(context.Background()) {
			doc := make(map[string]interface{})
			err = elem.Decode(&doc)
			if err != nil {
				return nil, err
			}
			delete(doc, "_id")
			collContents = append(collContents, doc)
		}
		res[name] = collContents
	}

	return res, nil
}

func walk(m bson.M, f func(map[string]interface{}, string)) {
	for k, v := range m {
		if sm, ok := v.(map[string]interface{}); ok {
			walk(sm, f)
		} else if ss, ok := v.([]interface{}); ok {
			for _, ms := range ss {
				if sm, ok := ms.(map[string]interface{}); ok {
					walk(sm, f)
				}
			}
		} else {
			f(m, k)
		}
	}
}

type cacheObj struct {
	Id  string
	Src string
}
