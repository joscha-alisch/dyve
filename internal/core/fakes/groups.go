package fakes

import (
	"github.com/joscha-alisch/dyve/internal/core/groups"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type RecordingGroupsService struct {
	Err        error
	ByProvider groups.GroupByProviderMap
	Record     GroupsRecorder
}

func (r *RecordingGroupsService) ListGroupsByProvider() (groups.GroupByProviderMap, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	return r.ByProvider, nil
}

func (r *RecordingGroupsService) ListGroupsPaginated(perPage int, page int) (sdk.GroupPage, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordingGroupsService) GetGroup(id string) (sdk.Group, error) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordingGroupsService) DeleteGroup(id string) error {
	//TODO implement me
	panic("implement me")
}

func (r *RecordingGroupsService) UpdateGroups(guid string, groups []sdk.Group) error {
	//TODO implement me
	panic("implement me")
}

type GroupsRecorder struct {
}
