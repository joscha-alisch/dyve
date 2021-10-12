package sdk

type GroupProvider interface {
	ListGroups() ([]Group, error)
	GetGroup(id string) (Group, error)
}

type GroupPage struct {
	Pagination
	Groups []Group `json:"groups"`
}

type Group struct {
	Id      string `json:"id" bson:"id"`
	Name    string `json:"name" bson:"name"`
	Members []Member
}

type Member struct {
	Id   string
	Name string
}
