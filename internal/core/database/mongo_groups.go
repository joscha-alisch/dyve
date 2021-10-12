package database

import (
	"context"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *mongoDb) GetGroup(id string) (sdk.Group, error) {
	res := m.groups.FindOne(context.Background(), bson.M{
		"id": id,
	})

	g := sdk.Group{}
	err := res.Decode(&g)
	if err != nil {
		return sdk.Group{}, err
	}

	return g, nil
}

func (m *mongoDb) ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error) {
	p, cursor, err := m.listPaginated(m.apps, perPage, page)
	if err != nil {
		return sdk.GroupPage{}, err
	}

	var groups []sdk.Group
	for cursor.Next(context.Background()) {
		group := sdk.Group{}
		err = cursor.Decode(&group)
		if err != nil {
			return sdk.GroupPage{}, err
		}
		groups = append(groups, group)
	}

	return sdk.GroupPage{
		Pagination: p,
		Groups:     groups,
	}, nil
}

func (m *mongoDb) UpdateGroups(providerId string, groups []sdk.Group) error {
	groupMap := make(map[string]interface{}, len(groups))
	for _, group := range groups {
		groupMap[group.Id] = group
	}

	return m.updateCollection(m.groups, providerId, groupMap)
}
