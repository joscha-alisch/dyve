package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
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

const (
	Subjects Collection = "subjects"
	Empty    Collection = "empty"
	Unsorted Collection = "unsorted"
	Provided Collection = "provided"
)

type testSubject struct {
	Id       string `bson:"id"`
	Property string `bson:"property"`
	Provider string `bson:"provider"`
}

var (
	subjectA = testSubject{
		Id:       "subject-a",
		Property: "a",
	}
	subjectB = testSubject{
		Id:       "subject-b",
		Property: "b",
	}
	subjectC = testSubject{
		Id:       "subject-c",
		Property: "c",
	}
	subjectNew = testSubject{
		Id:       "newItem",
		Property: "new",
	}
	providedA = testSubject{
		Id:       "provided-a",
		Provider: "provider-1",
	}
	providedB = testSubject{
		Id:       "provided-b",
		Provider: "provider-1",
	}
	providedC = testSubject{
		Id:       "provided-c",
		Provider: "provider-2",
	}
)

var baseState = map[string]interface{}{
	string(Subjects): toCollection(
		subjectA,
		subjectB,
		subjectC,
	),
	string(Empty): toCollection(),
	string(Unsorted): toCollection(
		subjectC,
		subjectB,
		subjectA,
	),
	string(Provided): toCollection(
		providedA,
		providedB,
		providedC,
	),
}

func TestMongoIntegration(t *testing.T) {
	currentTime = func() time.Time {
		return someTime
	}

	tests := []struct {
		desc            string
		f               func(db Database, res *testSubject, resList *[]testSubject, tt *testing.T) error
		expectedErr     error
		expectsOne      *testSubject
		expectsMultiple *[]testSubject
	}{
		/*
			Queries
		*/
		{desc: "finds a", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindOne(Subjects, bson.M{"id": subjectA.Id}, a)
		}, expectsOne: &subjectA},
		{desc: "returns not found", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindOne(Subjects, bson.M{"id": "not-existent"}, a)
		}, expectedErr: ErrNotFound},
		{desc: "finds a by id", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindOneById(Subjects, subjectA.Id, a)
		}, expectsOne: &subjectA},
		{desc: "returns not found by id", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindOneById(Subjects, "non-existent", a)
		}, expectedErr: ErrNotFound},
		{desc: "finds first match", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindOne(Unsorted, bson.M{}, a)
		}, expectsOne: &subjectC},
		{desc: "finds a when sorting", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindOneSorted(Unsorted, bson.M{}, bson.M{"id": 1}, a)
		}, expectsOne: &subjectA},
		{desc: "finds many", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindMany(Subjects, bson.M{}, decodeEach(resList))
		}, expectsMultiple: &[]testSubject{subjectA, subjectB, subjectC}},
		{desc: "finds many with limits and sort", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.FindManyWithOptions(Subjects, bson.M{}, decodeEach(resList), bson.M{"id": -1}, 2)
		}, expectsMultiple: &[]testSubject{subjectC, subjectB}},
		{desc: "lists paginated", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			page := sdk.Pagination{}
			err := db.ListPaginated(Subjects, 1, 1, &page, decodeEach(resList))
			requireEqual(page, sdk.Pagination{
				TotalResults: 3,
				TotalPages:   3,
				PerPage:      1,
				Page:         1,
			}, tt)
			return err
		}, expectsMultiple: &[]testSubject{subjectB}},
		/**
		Updates
		*/
		{desc: "adds new item", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			newProvided := map[string]interface{}{
				providedA.Id: providedA,
				providedB.Id: providedB,
				"newItem": testSubject{
					Id:       "newItem",
					Property: "some new value",
					Provider: providedA.Provider,
				},
			}
			return db.UpdateProvided(Provided, providedA.Provider, newProvided)
		}},
		{desc: "removes existing item", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			newProvided := map[string]interface{}{
				providedA.Id: providedA,
			}
			return db.UpdateProvided(Provided, providedA.Provider, newProvided)
		}},
		{desc: "updates existing item", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			newProvided := map[string]interface{}{
				providedA.Id: providedA,
				providedB.Id: testSubject{
					Id:       providedB.Id,
					Property: "changed-property",
					Provider: providedB.Provider,
				},
			}
			return db.UpdateProvided(Provided, providedA.Provider, newProvided)
		}},
		{desc: "updates multiple properties", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			filters := map[string]interface{}{
				subjectA.Id: bson.M{"id": subjectA.Id},
				subjectC.Id: bson.M{"id": subjectC.Id},
			}
			updates := map[string]interface{}{
				subjectA.Id: bson.M{"property": "changed-a"},
				subjectC.Id: bson.M{"property": "changed-c"},
			}
			return db.UpdateMany(Subjects, filters, updates)
		}},
		{desc: "updates single item", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.UpdateOne(Subjects, bson.M{"id": subjectA.Id}, false, bson.M{"property": "changed-a"}, a)
		}, expectsOne: subjectA.withProperty("changed-a")},
		{desc: "updates single item by id", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.UpdateOneById(Subjects, subjectA.Id, false, bson.M{"property": "changed-a"}, a)
		}, expectsOne: subjectA.withProperty("changed-a")},
		{desc: "creates item via update", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.UpdateOne(Subjects, bson.M{"id": subjectNew.Id}, true, subjectNew, a)
		}, expectsOne: &subjectNew},
		{desc: "creates item via update by id", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.UpdateOneById(Subjects, subjectNew.Id, true, subjectNew, a)
		}, expectsOne: &subjectNew},
		{desc: "update returns not found without createIfMissing", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.UpdateOne(Subjects, bson.M{"id": subjectNew.Id}, false, subjectNew, a)
		}, expectedErr: ErrNotFound},

		/**
		Delete
		DeleteOne(coll Collection, filter bson.M) error
		DeleteOneById(coll Collection, id string) error
		*/
		{desc: "deletes a", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.DeleteOne(Subjects, bson.M{"id": subjectA.Id})
		}},
		{desc: "deletes a by id", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.DeleteOneById(Subjects, subjectA.Id)
		}},
		{desc: "delete returns not found", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.DeleteOne(Subjects, bson.M{"id": "not-existent"})
		}, expectedErr: ErrNotFound},
		{desc: "delete by id returns not found", f: func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.DeleteOneById(Subjects, "not-existent")
		}, expectedErr: ErrNotFound},
		{desc: "inserts one", f: func(db Database, res *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.InsertOne(Subjects, bson.M{"id": "inserted"}, testSubject{
				Id:       "inserted",
				Property: "a",
			})
		}},
		{desc: "insert fails if exists", f: func(db Database, res *testSubject, resList *[]testSubject, tt *testing.T) error {
			return db.InsertOne(Subjects, bson.M{"id": "subject-a"}, testSubject{
				Id:       "inserted",
				Property: "a",
			})
		}, expectedErr: ErrExists},
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
			acceptanceTesting(fileName, baseState, test.f, test.expectsOne, test.expectsMultiple, test.expectedErr, mongodb, tt)
		})
	}
}

func toCollection(s ...testSubject) []bson.M {
	var res []bson.M
	for _, subject := range s {
		res = append(res, toBson(subject))
	}
	return res
}

func toBson(s testSubject) bson.M {
	return bson.M{"id": s.Id, "property": s.Property, "provider": s.Provider}
}

func (s testSubject) withProperty(change string) *testSubject {
	s.Property = change
	return &s
}

func decodeEach(list *[]testSubject) func(c *mongo.Cursor) error {
	return func(c *mongo.Cursor) error {
		res := testSubject{}
		err := c.Decode(&res)
		if err != nil {
			return err
		}
		*list = append(*list, res)
		return nil
	}
}

func requireEqual(a, b interface{}, tt *testing.T) {
	if !cmp.Equal(a, b) {
		tt.Errorf("expected the two objects to be equal: %s", cmp.Diff(a, b))
	}
}

func acceptanceTesting(
	name string,
	state map[string]interface{},
	f func(db Database, a *testSubject, resList *[]testSubject, tt *testing.T) error,
	expected *testSubject,
	expectedList *[]testSubject,
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

	db, err := NewMongoDB(MongoLogin{
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

	var res *testSubject
	if expected != nil {
		res = &testSubject{}
	}

	var resList *[]testSubject
	if expectedList != nil {
		resList = &[]testSubject{}
	}

	err = f(db, res, resList, tt)
	if err == nil && expectedErr != nil {
		tt.Errorf("expected an error but did not get one")
	}
	if err != nil && expectedErr == nil {
		tt.Errorf("expected no error but got one: %v", err)
	}
	if expectedErr != nil && err != nil && !errors.Is(err, expectedErr) {
		tt.Errorf("expected a different error: %v", cmp.Diff(expectedErr.Error(), err.Error()))
	}

	if expected != nil && !cmp.Equal(res, expected) {
		tt.Errorf("returned testSubject not correct. Diff:\n%+v", cmp.Diff(expected, res))
	}

	if expectedList != nil && !cmp.Equal(resList, expectedList) {
		tt.Errorf("returned list of test subjects not correct. Diff:\n%+v", cmp.Diff(expectedList, resList))
	}

	contents, err := dumpContents(s, dbName)
	if err != nil {
		tt.Fatal(err)
	}

	walk(contents, func(m map[string]interface{}, k string) {
		if t, ok := m[k].(primitive.DateTime); ok {
			m[k] = time.Unix(int64(t)/1000, 0).Format(time.RFC3339)
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

	acceptedBytes, _ := json.Marshal(acceptedContents)
	actualBytes, _ := json.Marshal(contents)

	if !cmp.Equal(acceptedBytes, actualBytes) {
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
