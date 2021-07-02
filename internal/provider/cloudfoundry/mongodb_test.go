package cloudfoundry

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/benweissmann/memongo"
	"github.com/google/go-cmp/cmp"
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

func TestMongoIntegration(t *testing.T) {
	tests := []struct {
		desc  string
		f     func(db Database, tt *testing.T) error
		err   error
		state bson.M
	}{
		{desc: "create app", f: func(db Database, tt *testing.T) error {
			return db.UpsertApp(App{Name: "my-app", Guid: "abc"})
		}},
		{desc: "update app", state: bson.M{
			"apps": []bson.M{
				{"name": "old-name", "guid": "abc"},
			},
		}, f: func(db Database, tt *testing.T) error {
			return db.UpsertApp(App{Name: "my-app", Guid: "abc"})
		}},
		{desc: "create space", f: func(db Database, tt *testing.T) error {
			return db.UpsertSpace(Space{Name: "my-space", Guid: "abc"})
		}},
		{desc: "update space", state: bson.M{
			"spaces": []bson.M{
				{"name": "old-name", "guid": "abc"},
			},
		}, f: func(db Database, tt *testing.T) error {
			return db.UpsertSpace(Space{Name: "my-space", Guid: "abc"})
		}},
		{desc: "create org", f: func(db Database, tt *testing.T) error {
			return db.UpsertOrg(Org{Name: "my-org", Guid: "abc"})
		}},
		{desc: "update org", state: bson.M{
			"orgs": []bson.M{
				{"name": "old-name", "guid": "abc"},
			},
		}, f: func(db Database, tt *testing.T) error {
			return db.UpsertOrg(Org{Name: "my-org", Guid: "abc"})
		}},
		{desc: "fetch job", state: bson.M{
			"jobs": []bson.M{
				{"type": float32(ReconcileOrg), "guid": "def", "lastUpdated": someTime.Add(1 * time.Second)},
				{"type": float32(ReconcileOrg), "guid": "abc", "lastUpdated": someTime},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := ReconcileJob{Type: ReconcileOrg, Guid: "abc"}
			j, ok := db.AcceptReconcileJob(someTime.Add(10*time.Second), someTime.Add(5*time.Minute))
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
	}

	mongo, err := memongo.Start("3.6.23")
	if err != nil {
		t.Fatal(err)
	}
	defer mongo.Stop()

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			fileName := strings.ReplaceAll(test.desc, " ", "_")
			acceptanceTesting(fileName, test.state, test.f, mongo, tt)
		})
	}
}

func acceptanceTesting(
	name string,
	state map[string]interface{},
	f func(db Database, tt *testing.T) error,
	s *memongo.Server,
	tt *testing.T,
) {
	dbName := memongo.RandomDatabase()
	db, err := NewMongoDatabase(s.URI(), dbName)
	if err != nil {
		tt.Fatal(err)
	}

	if state != nil {
		err = setState(state, s, dbName)
		if err != nil {
			tt.Fatal(err)
		}
	}

	err = f(db, tt)
	if err != nil {
		tt.Fatal(err)
	}

	contents, err := dumpContents(s, dbName)
	if err != nil {
		tt.Fatal(err)
	}

	walk(contents, func(m map[string]interface{}, k string) {
		if t, ok := m[k].(primitive.DateTime); ok {
			m[k] = time.Unix(int64(t)/1000, 0).Format(time.RFC3339)
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
