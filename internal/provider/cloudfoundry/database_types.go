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

type Routes []Route
type Route struct {
	Host string `bson:"host"`
	Path string `bson:"path"`
	Port int    `bson:"port"`
}

type Instances []Instance
type Instance struct {
	State string
	Since time.Time
}

func (a App) toSdkApp() sdk.App {
	app := sdk.App{
		Id:       a.Guid,
		Name:     a.Name,
		Labels:   sdk.AppLabels{},
		Position: sdk.AppPosition{},
	}

	if a.Space.Org.Name != "" {
		app.Labels["org"] = a.Space.Org.Name
		app.Position = append(app.Position, a.Space.Org.Name)
	}

	if a.Space.Name != "" {
		app.Labels["space"] = a.Space.Name
		app.Position = append(app.Position, a.Space.Name)
	}

	if len(app.Labels) == 0 {
		app.Labels = nil
	}

	if len(app.Position) == 0 {
		app.Position = nil
	}

	return app
}
