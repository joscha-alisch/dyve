package db

import (
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseRecorder struct {
	Records []DatabaseRecord
}

type DatabaseRecord struct {
	Collection      database.Collection
	Filter          interface{}
	Id              string
	Sort            interface{}
	Limit           int
	PerPage         int
	Page            int
	Provider        string
	Updates         map[string]interface{}
	Filters         map[string]interface{}
	CreateIfMissing bool
	Update          interface{}
	Data            interface{}
}

func (r *DatabaseRecorder) Record(record DatabaseRecord) {
	r.Records = append(r.Records, record)
}

type Database struct {
	Recorder         *DatabaseRecorder
	Return           func(target interface{})
	ReturnEach       func(each func(decodable database.Decodable) error)
	ReturnPagination func(pagination *sdk.Pagination)
	Err              error
}

func (d *Database) FindOne(coll database.Collection, filter interface{}, res interface{}) error {
	if d.Err != nil {
		return d.Err
	}

	if d.Return != nil {
		d.Return(res)
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filter:     filter,
	})
	return nil
}

func (d *Database) FindOneById(coll database.Collection, id string, res interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	if d.Return != nil {
		d.Return(res)
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Id:         id,
	})
	return nil
}

func (d *Database) FindOneSorted(coll database.Collection, filter bson.M, sort bson.M, res interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	if d.Return != nil {
		d.Return(res)
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filter:     filter,
		Sort:       sort,
	})
	return nil
}

func (d *Database) FindMany(coll database.Collection, filter bson.M, each func(dec database.Decodable) error) error {
	if d.Err != nil {
		return d.Err
	}
	d.ReturnEach(each)
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filter:     filter,
	})
	return nil
}

func (d *Database) FindManyWithOptions(coll database.Collection, filter bson.M, each func(dec database.Decodable) error, sort bson.M, limit int) error {
	if d.Err != nil {
		return d.Err
	}
	d.ReturnEach(each)
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filter:     filter,
		Sort:       sort,
		Limit:      limit,
	})
	return nil
}

func (d *Database) ListPaginated(coll database.Collection, perPage int, page int, p *sdk.Pagination, each func(dec database.Decodable) error) error {
	if d.Err != nil {
		return d.Err
	}
	d.ReturnEach(each)
	d.ReturnPagination(p)
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		PerPage:    perPage,
		Page:       page,
	})
	return nil
}

func (d *Database) UpdateProvided(coll database.Collection, provider string, updates map[string]interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Provider:   provider,
		Updates:    updates,
	})
	return nil
}

func (d *Database) UpdateMany(coll database.Collection, filters map[string]interface{}, updates map[string]interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filters:    filters,
		Updates:    updates,
	})
	return nil
}

func (d *Database) UpdateOne(coll database.Collection, filter bson.M, createIfMissing bool, update interface{}, res interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	if d.Return != nil {
		d.Return(res)
	}
	d.Recorder.Record(DatabaseRecord{
		Collection:      coll,
		CreateIfMissing: createIfMissing,
		Filter:          filter,
		Update:          update,
	})
	return nil
}

func (d *Database) UpdateOneById(coll database.Collection, id string, createIfMissing bool, update interface{}, res interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	if d.Return != nil {
		d.Return(res)
	}
	d.Recorder.Record(DatabaseRecord{
		Collection:      coll,
		CreateIfMissing: createIfMissing,
		Id:              id,
		Update:          update,
	})
	return nil
}

func (d *Database) InsertOne(coll database.Collection, existsFilter interface{}, data interface{}) error {
	if d.Err != nil {
		return d.Err
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filter:     existsFilter,
		Data:       data,
	})
	return nil
}

func (d *Database) DeleteOne(coll database.Collection, filter bson.M) error {
	if d.Err != nil {
		return d.Err
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Filter:     filter,
	})
	return nil
}

func (d *Database) DeleteOneById(coll database.Collection, id string) error {
	if d.Err != nil {
		return d.Err
	}
	d.Recorder.Record(DatabaseRecord{
		Collection: coll,
		Id:         id,
	})
	return nil
}

func (d *Database) EnsureIndex(coll database.Collection, model mongo.IndexModel) error {
	if d.Err != nil {
		return d.Err
	}
	return nil
}
