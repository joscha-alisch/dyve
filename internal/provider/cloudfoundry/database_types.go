package cloudfoundry

import "time"

type CFInfo struct {
	Guid string
}

type CF struct {
	CFInfo      `bson:",inline"`
	Orgs        []string
	LastUpdated time.Time `bson:"lastUpdated"`
}

type OrgInfo struct {
	Guid string
	Name string
	Cf   CFInfo
}

type Org struct {
	OrgInfo     `bson:",inline"`
	Spaces      []string
	LastUpdated time.Time `bson:"lastUpdated"`
}

type SpaceInfo struct {
	Guid string
	Name string
	Org  OrgInfo
}

type Space struct {
	SpaceInfo   `bson:",inline"`
	Apps        []string
	LastUpdated time.Time `bson:"lastUpdated"`
}

type AppInfo struct {
	Guid  string
	Name  string
	Space SpaceInfo
}

type App struct {
	AppInfo `bson:",inline"`
}
