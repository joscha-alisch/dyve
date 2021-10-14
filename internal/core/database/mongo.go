package database

import (
	"context"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"time"
)

var currentTime = time.Now

type MongoLogin struct {
	Uri string
	DB  string
}

func NewMongoDB(l MongoLogin) (Database, error) {
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

	m := &mongoDb{
		ctx:         context.Background(),
		cli:         c,
		db:          db,
		collections: make(map[Collection]*mongo.Collection),
	}

	return m, nil
}

type mongoDb struct {
	ctx         context.Context
	cli         *mongo.Client
	db          *mongo.Database
	collections map[Collection]*mongo.Collection
}

func (m *mongoDb) EnsureIndex(coll Collection, model mongo.IndexModel) error {
	c := m.collection(coll)
	_, err := c.Indexes().CreateOne(m.ctx, model)
	return err
}

func (m *mongoDb) FindManyWithOptions(coll Collection, filter bson.M, each func(c *mongo.Cursor) error, sort bson.M, limit int) error {
	c := m.collection(coll)
	o := options.FindOptions{}
	if sort != nil {
		o.SetSort(sort)
	}

	if limit > 0 {
		o.SetLimit(int64(limit))
	}

	res, err := c.Find(m.ctx, filter, &o)
	if err != nil {
		return err
	}
	for res.Next(m.ctx) {
		err := each(res)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *mongoDb) FindOneSorted(coll Collection, filter bson.M, sort bson.M, res interface{}) error {
	c := m.collection(coll)
	findResult := c.FindOne(m.ctx, filter, options.FindOne().SetSort(sort))
	return findResult.Decode(&res)
}

func (m *mongoDb) FindMany(coll Collection, filter bson.M, each func(c *mongo.Cursor) error) error {
	return m.FindManyWithOptions(coll, filter, each, nil, 0)
}

func Set(d interface{}) bson.M {
	return bson.M{"$set": d}
}

func SetOrdered(data ...interface{}) bson.D {
	res := bson.D{}
	for _, datum := range data {
		res = append(res, bson.E{Key: "$set", Value: datum})
	}
	return res
}

func (m *mongoDb) UpdateMany(coll Collection, filters map[string]interface{}, updates map[string]interface{}) error {
	c := m.collection(coll)

	models := make([]mongo.WriteModel, len(updates))
	i := 0
	for k, v := range updates {
		filter := filters[k]
		model := mongo.NewUpdateOneModel()
		model.SetUpsert(true).SetFilter(filter).SetUpdate(SetOrdered(filter, v))
		models[i] = model
		i++
	}
	ctx := context.Background()

	_, err := c.BulkWrite(ctx, models)
	if err != nil {
		return err
	}

	return nil
}

func (m *mongoDb) EnsureCreated(coll Collection, data bson.M, res interface{}) error {
	return m.UpdateOne(coll, data, true, Set(data), res)
}

func (m *mongoDb) UpdateOne(coll Collection, filter bson.M, createIfMissing bool, update bson.M, res interface{}) error {
	c := m.collection(coll)

	if res == nil {
		o := options.UpdateOptions{}
		o.SetUpsert(createIfMissing)
		_, err := c.UpdateOne(m.ctx, filter, update, &o)
		return err
	}

	o := options.FindOneAndUpdateOptions{}
	o.SetUpsert(createIfMissing)
	findResult := c.FindOneAndUpdate(m.ctx, filter, update, &o)
	return findResult.Decode(&res)
}

func (m *mongoDb) UpdateOneById(coll Collection, id string, createIfMissing bool, update bson.M, res interface{}) error {
	return m.UpdateOne(coll, bson.M{"id": id}, createIfMissing, update, res)
}

func (m *mongoDb) DeleteOne(coll Collection, filter bson.M) error {
	c := m.collection(coll)
	_, err := c.DeleteOne(context.Background(), filter)
	return err
}

func (m *mongoDb) DeleteOneById(coll Collection, id string) error {
	return m.DeleteOne(coll, bson.M{"id": id})
}

func (m *mongoDb) FindOne(coll Collection, filter bson.M, res interface{}) error {
	c := m.collection(coll)
	findResult := c.FindOne(context.Background(), filter)
	return findResult.Decode(&res)
}

func (m *mongoDb) FindOneById(coll Collection, id string, res interface{}) error {
	return m.FindOne(coll, bson.M{"id": id}, res)
}

func (m *mongoDb) ListPaginated(collName Collection, perPage int, page int, p *sdk.Pagination, each func(c *mongo.Cursor) error) error {
	coll := m.collection(collName)

	c, err := coll.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		return err
	}
	pages := int(math.Ceil(float64(c) / float64(perPage)))

	cursor, err := coll.Find(context.Background(), bson.M{}, options.Find().
		SetSkip(int64(page*perPage)).
		SetLimit(int64(perPage)),
	)
	if err != nil {
		return err
	}

	*p = sdk.Pagination{
		TotalResults: int(c),
		TotalPages:   pages,
		PerPage:      perPage,
		Page:         page,
	}

	for cursor.Next(context.Background()) {
		err = each(cursor)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *mongoDb) collection(c Collection) *mongo.Collection {
	if m.collections[c] == nil {
		m.collections[c] = m.db.Collection(string(c))
	}
	return m.collections[c]
}

func (m *mongoDb) UpdateProvided(collName Collection, provider string, updates map[string]interface{}) error {
	c := m.collection(collName)

	ids := make([]string, len(updates))
	filterMap := make(map[string]interface{})
	for id, _ := range updates {
		filterMap[id] = bson.M{
			"provider": provider,
			"id":       id,
		}
	}

	err := m.UpdateMany(collName, filterMap, updates)
	if err != nil {
		return err
	}

	_, err = c.DeleteMany(m.ctx, bson.M{
		"provider": provider,
		"id": bson.M{
			"$nin": ids,
		},
	})
	if err != nil {
		return err
	}

	return nil
}
