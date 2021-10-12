package github

import "time"

type Org struct {
	OrgInfo     `bson:",inline"`
	LastUpdated time.Time `bson:"lastUpdated"`
	Teams       []string
}

type OrgInfo struct {
	Guid string
}

type Team struct {
	TeamInfo    `bson:",inline"`
	LastUpdated time.Time `bson:"lastUpdated"`
	Members     []Member
}

type TeamInfo struct {
	Org  OrgInfo
	Guid string
	Name string
	Slug string
}

type Member struct {
	Guid string `bson:"guid"`
	Name string `bson:"name"`
}
