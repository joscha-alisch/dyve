package fakes

import (
	"github.com/joscha-alisch/dyve/internal/core/apps"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

type RecordingAppsService struct {
	Err    error
	App    apps.App
	Page   sdk.AppPage
	Apps   []sdk.App
	Record AppsRecorder
}

type AppsRecorder struct {
	App        sdk.App
	PerPage    int
	Page       int
	AppId      string
	ProviderId string
	Apps       []sdk.App
}

func (a *RecordingAppsService) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	a.Record.PerPage = perPage
	a.Record.Page = page

	if a.Err != nil {
		return sdk.AppPage{}, a.Err
	}
	return a.Page, nil
}

func (a *RecordingAppsService) GetApp(id string) (apps.App, error) {
	a.Record.AppId = id

	if a.Err != nil {
		return apps.App{}, a.Err
	}
	return a.App, nil
}

func (a *RecordingAppsService) UpdateApps(providerId string, apps []sdk.App) error {
	a.Record.ProviderId = providerId
	a.Record.Apps = apps

	if a.Err != nil {
		return a.Err
	}
	return nil
}

func (a *RecordingAppsService) UpdateApp(app sdk.App) error {
	a.Record.App = app
	if a.Err != nil {
		return a.Err
	}
	return nil
}

type MappingAppsService struct {
	Apps map[string]apps.App
}

func (m *MappingAppsService) ListAppsPaginated(perPage int, page int) (sdk.AppPage, error) {
	panic("implement me")
}

func (m *MappingAppsService) GetApp(id string) (apps.App, error) {
	return m.Apps[id], nil
}

func (m *MappingAppsService) UpdateApps(providerId string, appList []sdk.App) error {
	for id, app := range m.Apps {
		if app.ProviderId == providerId {
			delete(m.Apps, id)
		}
	}

	for _, app := range appList {
		m.Apps[app.Id] = apps.App{
			App:        app,
			ProviderId: providerId,
		}
	}
	return nil
}

func (m *MappingAppsService) UpdateApp(app sdk.App) error {
	m.Apps[app.Id] = apps.App{
		App:        app,
		ProviderId: m.Apps[app.Id].ProviderId,
	}
	return nil
}
