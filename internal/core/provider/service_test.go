package provider

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/fakes/db"
	"github.com/joscha-alisch/dyve/internal/core/fakes/fakeProvider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

func TestService_AppProvider(t *testing.T) {
	rec := &db.DatabaseRecorder{}
	d := &db.RecordingDatabase{Recorder: rec}
	s := NewService(d)

	p, err := s.GetAppProvider("fakeProvider")
	assertNil(t, "GetAppProvider should return nil in the beginning", p)
	assertErr(t, err, ErrNotFound)

	origProv := fakeProvider.AppProvider(nil)

	err = s.AddAppProvider("fakeProvider", "name", origProv)
	assertNil(t, "there should be no error", err)

	err = s.AddAppProvider("fakeProvider", "name", origProv)
	assertErr(t, err, ErrExists)

	p, err = s.GetAppProvider("fakeProvider")
	assertNil(t, "there should be no error", err)
	assertSame(t, p, origProv)

	err = s.DeleteAppProvider("fakeProvider")
	assertNil(t, "there should be no error", err)

	p, err = s.GetAppProvider("fakeProvider")
	assertNil(t, "GetAppProvider should return nil after deletion", p)
	assertErr(t, err, ErrNotFound)
}

func TestService_PipelineProvider(t *testing.T) {
	rec := &db.DatabaseRecorder{}
	d := &db.RecordingDatabase{Recorder: rec}
	s := NewService(d)

	p, err := s.GetPipelineProvider("fakeProvider")
	assertNil(t, "getProvider should return nil in the beginning", p)
	assertErr(t, err, ErrNotFound)

	origProv := fakeProvider.PipelineProvider(nil, sdk.PipelineUpdates{})

	err = s.AddPipelineProvider("fakeProvider", "name", origProv)
	assertNil(t, "there should be no error", err)

	err = s.AddPipelineProvider("fakeProvider", "name", origProv)
	assertErr(t, err, ErrExists)

	p, err = s.GetPipelineProvider("fakeProvider")
	assertNil(t, "there should be no error", err)
	assertSame(t, p, origProv)

	err = s.DeletePipelineProvider("fakeProvider")
	assertNil(t, "there should be no error", err)

	p, err = s.GetPipelineProvider("fakeProvider")
	assertNil(t, "getProvider should return nil after deletion", p)
	assertErr(t, err, ErrNotFound)
}

func TestService_GroupProvider(t *testing.T) {
	rec := &db.DatabaseRecorder{}
	returnFunc := func(each func(decodable database.Decodable) error) {
		_ = each(database.DecodableFunc(func(target interface{}) error {
			*target.(*Data) = Data{
				Id:   "fakeProvider",
				Name: "name",
			}
			return nil
		}))
	}
	d := &db.RecordingDatabase{
		Recorder:   rec,
		ReturnEach: func(each func(decodable database.Decodable) error) {},
	}
	s := NewService(d)

	list, err := s.ListGroupProviders()
	if list != nil {
		assertNil(t, "list should return nil in the beginning", list)
	}
	assertNil(t, "there should be no error", err)

	p, err := s.GetGroupProvider("fakeProvider")
	assertNil(t, "getProvider should return nil in the beginning", p)
	assertErr(t, err, ErrNotFound)

	origProv := fakeProvider.GroupProvider()

	err = s.AddGroupProvider("fakeProvider", "name", origProv)
	assertNil(t, "there should be no error", err)

	err = s.AddGroupProvider("fakeProvider", "name", origProv)
	assertErr(t, err, ErrExists)

	d.ReturnEach = returnFunc
	list, err = s.ListGroupProviders()
	assertEqual(t, list, []Data{{
		Id:   "fakeProvider",
		Name: "name",
	}})
	assertNil(t, "there should be no error", err)

	p, err = s.GetGroupProvider("fakeProvider")
	assertNil(t, "there should be no error", err)
	assertSame(t, p, origProv)

	err = s.DeleteGroupProvider("fakeProvider")
	assertNil(t, "there should be no error", err)

	p, err = s.GetGroupProvider("fakeProvider")
	assertNil(t, "getProvider should return nil after deletion", p)
	assertErr(t, err, ErrNotFound)
}

func TestService_RoutingProvider(t *testing.T) {
	rec := &db.DatabaseRecorder{}
	d := &db.RecordingDatabase{Recorder: rec}
	s := NewService(d)

	p, err := s.GetRoutingProviders()
	if p != nil {
		assertNil(t, "getProvider should return nil in the beginning", p)
	}
	assertErr(t, err, ErrNotFound)

	origProv := fakeProvider.RoutesProvider(nil)

	err = s.AddRoutingProvider("fakeProvider", "name", origProv)
	assertNil(t, "there should be no error", err)

	err = s.AddRoutingProvider("fakeProvider", "name", origProv)
	assertErr(t, err, ErrExists)

	p, err = s.GetRoutingProviders()
	assertNil(t, "there should be no error", err)
	assertSame(t, p[0], origProv)

	err = s.DeleteRoutingProvider("fakeProvider")
	assertNil(t, "there should be no error", err)

	p, err = s.GetRoutingProviders()
	if p != nil {
		assertNil(t, "getProvider should return nil after deletion", p)
	}
	assertErr(t, err, ErrNotFound)
}

func TestService_InstancesProvider(t *testing.T) {
	rec := &db.DatabaseRecorder{}
	d := &db.RecordingDatabase{Recorder: rec}
	s := NewService(d)

	p, err := s.GetInstancesProviders()
	if p != nil {
		assertNil(t, "getProvider should return nil in the beginning", p)
	}
	assertErr(t, err, ErrNotFound)

	origProv := fakeProvider.InstancesProvider(nil)

	err = s.AddInstancesProvider("fakeProvider", "name", origProv)
	assertNil(t, "there should be no error", err)

	err = s.AddInstancesProvider("fakeProvider", "name", origProv)
	assertErr(t, err, ErrExists)

	p, err = s.GetInstancesProviders()
	assertNil(t, "there should be no error", err)
	assertSame(t, p[0], origProv)

	err = s.DeleteInstancesProvider("fakeProvider")
	assertNil(t, "there should be no error", err)

	p, err = s.GetInstancesProviders()
	if p != nil {
		assertNil(t, "getProvider should return nil after deletion", p)
	}
	assertErr(t, err, ErrNotFound)
}

func TestReconcile(t *testing.T) {
	rec := &db.DatabaseRecorder{}
	d := &db.RecordingDatabase{Recorder: rec}
	s := NewService(d)

	d.Return = func(target interface{}) {
		*target.(*Provider) = Provider{
			ProviderType: "something",
			Data:         Data{Id: "id"},
			LastUpdated:  someTime,
		}
	}

	j, ok := s.AcceptReconcileJob(2 * time.Minute)
	assertEqual(t, j, recon.Job{
		Type:        "something",
		Guid:        "id",
		LastUpdated: someTime,
	})
	assertEqual(t, ok, true)

	err := s.RequestAppUpdate("app-id")
	assertNil(t, "no error", err)

	j, ok = s.AcceptReconcileJob(2 * time.Minute)
	assertEqual(t, j, recon.Job{
		Type: "routing",
		Guid: "app-id",
	})
	assertEqual(t, ok, true)

	j, ok = s.AcceptReconcileJob(2 * time.Minute)
	assertEqual(t, j, recon.Job{
		Type: "instances",
		Guid: "app-id",
	})
	assertEqual(t, ok, true)

	j, ok = s.AcceptReconcileJob(2 * time.Minute)
	assertEqual(t, j, recon.Job{
		Type:        "something",
		Guid:        "id",
		LastUpdated: someTime,
	})
	assertEqual(t, ok, true)
}

func assertNil(t *testing.T, desc string, a interface{}) {
	if a != nil {
		t.Fatal(desc)
	}
}

func assertErr(t *testing.T, err, should error) {
	if !errors.Is(err, should) {
		t.Fatalf("error mismatch: %s\n", cmp.Diff(should, err))
	}
}

func assertSame(t *testing.T, a, b interface{}) {
	if a != b {
		t.Fatal("the two objects are not the same")
	}
}

func assertEqual(t *testing.T, a, b interface{}) {
	if !cmp.Equal(a, b) {
		t.Fatalf("the two objects mismatch: %s\n", cmp.Diff(a, b))
	}
}
