package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/benweissmann/memongo"
	"github.com/google/go-cmp/cmp"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"github.com/rs/zerolog/log"
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

var baseAppState = map[string]interface{}{
	"providers": []bson.M{
		{"id": "provider-a", "type": "apps", "lastUpdated": "2006-01-01T14:58:00Z"},
		{"id": "provider-b", "type": "apps", "lastUpdated": "2006-01-01T14:58:30Z"},
	},
	"apps": []bson.M{
		{"provider": "provider-a", "name": "app-a", "id": "app-a"},
		{"provider": "provider-a", "name": "app-b", "id": "app-b"},
		{"provider": "provider-b", "name": "app-c", "id": "app-c"},
		{"provider": "provider-b", "name": "app-d", "id": "app-d"},
	},
}

var basePipelinesState = map[string]interface{}{
	"providers": []bson.M{
		{"id": "provider-a", "type": "pipelines", "lastUpdated": "2006-01-01T14:58:00Z"},
		{"id": "provider-b", "type": "pipelines", "lastUpdated": "2006-01-01T14:58:30Z"},
	},
	"pipelines": []bson.M{
		{"provider": "provider-a", "name": "pipeline-a", "id": "pipeline-a"},
		{"provider": "provider-a", "name": "pipeline-b", "id": "pipeline-b"},
		{"provider": "provider-b", "name": "pipeline-c", "id": "pipeline-c"},
		{"provider": "provider-b", "name": "pipeline-d", "id": "pipeline-d"},
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
		/**
		Providers
		*/
		{desc: "adds app provider", state: baseAppState, f: func(db Database, tt *testing.T) error {
			return db.AddAppProvider("provider-c")
		}},
		{desc: "adds pipeline provider", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			return db.AddPipelineProvider("provider-c")
		}},
		{desc: "deletes provider and apps", state: baseAppState, f: func(db Database, tt *testing.T) error {
			return db.DeleteAppProvider("provider-a")
		}},
		{desc: "deletes provider and pipelines", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			return db.DeletePipelineProvider("provider-a")
		}},
		/**
		Apps
		*/
		{desc: "updates existing apps", state: baseAppState, f: func(db Database, tt *testing.T) error {
			return db.UpdateApps("provider-a", []sdk.App{
				{Id: "app-a", Name: "new-app-a"},
				{Id: "app-b", Name: "new-app-b"},
			})
		}},
		{desc: "gets app", state: baseAppState, f: func(db Database, tt *testing.T) error {
			app, err := db.GetApp("app-a")
			expected := sdk.App{
				Id:   "app-a",
				Name: "app-a",
			}
			if !cmp.Equal(expected, app) {
				tt.Errorf("wrong app returned:\n%s\n", cmp.Diff(expected, app))
			}

			return err
		}},
		{desc: "adds new apps", state: baseAppState, f: func(db Database, tt *testing.T) error {
			return db.UpdateApps("provider-a", []sdk.App{
				{Id: "app-a", Name: "app-a"},
				{Id: "app-b", Name: "app-b"},
				{Id: "app-c", Name: "app-c"},
			})
		}},
		{desc: "removes old apps", state: baseAppState, f: func(db Database, tt *testing.T) error {
			return db.UpdateApps("provider-a", []sdk.App{
				{Id: "app-a", Name: "app-a"},
			})
		}},
		{desc: "lists apps paginated", state: baseAppState, f: func(db Database, tt *testing.T) error {
			apps, err := db.ListAppsPaginated(2, 1)
			if err != nil {
				return err
			}
			expected := sdk.AppPage{
				Pagination: sdk.Pagination{
					TotalResults: 4,
					TotalPages:   2,
					PerPage:      2,
					Page:         1,
				},
				Apps: []sdk.App{
					{Id: "app-c", Name: "app-c"},
					{Id: "app-d", Name: "app-d"},
				},
			}

			if !cmp.Equal(expected, apps) {
				tt.Errorf("wrong apps returned:\n%s\n", cmp.Diff(expected, apps))
			}
			return nil
		}},
		/**
		Pipelines
		*/
		{desc: "updates existing pipelines", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			return db.UpdatePipelines("provider-a", []sdk.Pipeline{
				{Id: "pipeline-a", Name: "new-pipeline-a"},
				{Id: "pipeline-b", Name: "new-pipeline-b"},
			})
		}},
		{desc: "gets pipeline", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			app, err := db.GetPipeline("pipeline-a")
			expected := sdk.Pipeline{
				Id:   "pipeline-a",
				Name: "pipeline-a",
			}
			if !cmp.Equal(expected, app) {
				tt.Errorf("wrong pipeline returned:\n%s\n", cmp.Diff(expected, app))
			}

			return err
		}},
		{desc: "adds new pipelines", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			return db.UpdatePipelines("provider-a", []sdk.Pipeline{
				{Id: "pipeline-a", Name: "pipeline-a"},
				{Id: "pipeline-b", Name: "pipeline-b"},
				{Id: "pipeline-c", Name: "pipeline-c"},
			})
		}},
		{desc: "removes old pipelines", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			return db.UpdatePipelines("provider-a", []sdk.Pipeline{
				{Id: "pipeline-a", Name: "pipeline-a"},
			})
		}},
		{desc: "lists pipelines paginated", state: basePipelinesState, f: func(db Database, tt *testing.T) error {
			apps, err := db.ListPipelinesPaginated(2, 1)
			if err != nil {
				return err
			}
			expected := sdk.PipelinePage{
				Pagination: sdk.Pagination{
					TotalResults: 4,
					TotalPages:   2,
					PerPage:      2,
					Page:         1,
				},
				Pipelines: []sdk.Pipeline{
					{Id: "pipeline-c", Name: "pipeline-c"},
					{Id: "pipeline-d", Name: "pipeline-d"},
				},
			}

			if !cmp.Equal(expected, apps) {
				tt.Errorf("wrong pipelines returned:\n%s\n", cmp.Diff(expected, apps))
			}
			return nil
		}},

		/**
		Fetch Jobs
		*/
		{desc: "fetch app provider job", state: bson.M{
			"providers": []bson.M{
				{"id": "provider-a", "type": "apps", "lastUpdated": someTime.Add(-30 * time.Second)},
				{"id": "provider-b", "type": "apps", "lastUpdated": someTime.Add(-90 * time.Second)},
				{"id": "provider-c", "type": "apps", "lastUpdated": someTime.Add(-60 * time.Second)},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcileAppProvider, Guid: "provider-b"}
			j, ok := db.AcceptReconcileJob(1 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch pipeline provider job", state: bson.M{
			"providers": []bson.M{
				{"id": "provider-a", "type": "pipelines", "lastUpdated": someTime.Add(-30 * time.Second)},
				{"id": "provider-b", "type": "pipelines", "lastUpdated": someTime.Add(-90 * time.Second)},
				{"id": "provider-c", "type": "apps", "lastUpdated": someTime.Add(-60 * time.Second)},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcilePipelineProvider, Guid: "provider-b"}
			j, ok := db.AcceptReconcileJob(1 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch no job", state: bson.M{
			"providers": []bson.M{
				{"id": "provider-a", "type": "apps", "lastUpdated": someTime.Add(-20 * time.Second)},
				{"id": "provider-b", "type": "apps", "lastUpdated": someTime.Add(-30 * time.Second)},
				{"id": "provider-c", "type": "apps", "lastUpdated": someTime.Add(-40 * time.Second)},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: "", Guid: ""}
			j, ok := db.AcceptReconcileJob(1 * time.Minute)
			if ok {
				tt.Errorf("expected no work to be done")
			}
			if !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch app job never updated", state: bson.M{
			"providers": []bson.M{
				{"id": "provider-a", "type": "apps"},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcileAppProvider, Guid: "provider-a"}
			j, ok := db.AcceptReconcileJob(1 * time.Minute)
			if !ok || !cmp.Equal(expected, j) {
				tt.Errorf("wrong job returned:\n%s\n", cmp.Diff(expected, j))
			}
			return nil
		}},
		{desc: "fetch pipeline job never updated", state: bson.M{
			"providers": []bson.M{
				{"id": "provider-a", "type": "pipelines"},
			},
		}, f: func(db Database, tt *testing.T) error {
			expected := recon.Job{Type: ReconcilePipelineProvider, Guid: "provider-a"}
			j, ok := db.AcceptReconcileJob(1 * time.Minute)
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
			m[k] = time.Unix(int64(t)/1000, 0).Format(time.RFC3339)
		}

		if t, ok := m[k].(primitive.A); ok {
			m[k] = ([]interface{})(t)
		}
	})

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
