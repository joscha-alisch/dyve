package cloudfoundry

import (
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

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

func (a App) toSdkApp() sdk.App {
	app := sdk.App{
		Id:   a.Guid,
		Name: a.Name,
		Meta: map[string]interface{}{},
	}

	if a.Space.Org.Name != "" {
		app.Meta["org"] = a.Space.Org.Name
	}

	if a.Space.Name != "" {
		app.Meta["space"] = a.Space.Name
	}

	if len(app.Meta) == 0 {
		app.Meta = nil
	}

	return app
}
