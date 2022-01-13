package reconciler

import (
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/joscha-alisch/dyve/internal/core/apps"
	"github.com/joscha-alisch/dyve/internal/core/fakes"
	"github.com/joscha-alisch/dyve/internal/core/pipelines"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/internal/core/service"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"testing"
	"time"
)

var someTime, _ = time.Parse(time.RFC3339, "2006-01-01T15:00:00Z")

func TestName(t *testing.T) {
	tests := []struct {
		desc             string
		job              recon.Job
		providerId       string
		appProvider      *fakes.Provider
		pipelineProvider *fakes.Provider
		routingProvider  *fakes.Provider
		appsBefore       *fakes.MappingAppsService
		appsAfter        *fakes.MappingAppsService
		pipelinesBefore  *fakes.MappingPipelinesService
		pipelinesAfter   *fakes.MappingPipelinesService
		routesBefore     *fakes.MappingRoutesService
		routesAfter      *fakes.MappingRoutesService
		expectedErr      error
		expectedWorked   bool
		recordedTime     time.Time
	}{
		{
			desc: "adds apps", job: recon.Job{
				Type: provider.ReconcileAppProvider,
				Guid: "app-provider",
			}, providerId: "app-provider", appProvider: fakes.AppProvider([]sdk.App{
				{Id: "app-a", Name: "app-a"},
				{Id: "app-b", Name: "app-b"},
			}), appsAfter: &fakes.MappingAppsService{Apps: map[string]apps.App{
				"app-a": {ProviderId: "app-provider", App: sdk.App{Id: "app-a", Name: "app-a"}},
				"app-b": {ProviderId: "app-provider", App: sdk.App{Id: "app-b", Name: "app-b"}},
			}}, expectedWorked: true,
		},
		{
			desc: "adds pipelines, runs and versions", job: recon.Job{
				Type:        provider.ReconcilePipelineProvider,
				Guid:        "pipeline-provider",
				LastUpdated: someTime,
			}, providerId: "pipeline-provider", pipelineProvider: fakes.PipelineProvider([]sdk.Pipeline{
				{Id: "pipeline-a", Name: "pipeline-a"},
				{Id: "pipeline-b", Name: "pipeline-b"},
			}, sdk.PipelineUpdates{
				Versions: sdk.PipelineVersionList{{PipelineId: "pipeline-a", Created: someTime}},
				Runs:     sdk.PipelineStatusList{{PipelineId: "pipeline-a", Started: someTime}},
			}), pipelinesAfter: &fakes.MappingPipelinesService{Pipelines: map[string]pipelines.Pipeline{
				"pipeline-a": {ProviderId: "pipeline-provider", Pipeline: sdk.Pipeline{Id: "pipeline-a", Name: "pipeline-a"}},
				"pipeline-b": {ProviderId: "pipeline-provider", Pipeline: sdk.Pipeline{Id: "pipeline-b", Name: "pipeline-b"}},
			}, Versions: map[string]sdk.PipelineVersionList{
				"pipeline-a": {{PipelineId: "pipeline-a", Created: someTime}},
			}, Runs: map[string]sdk.PipelineStatusList{
				"pipeline-a": {{PipelineId: "pipeline-a", Started: someTime}},
			},
			}, expectedWorked: true, recordedTime: someTime,
		},
		{
			desc: "adds routes", job: recon.Job{
				Type:        provider.ReconcileRoutingProviders,
				Guid:        "app-a",
				LastUpdated: someTime,
			}, providerId: "routes-provider", routingProvider: fakes.RoutesProvider(map[string]sdk.AppRouting{
				"app-a": {Routes: sdk.AppRoutes{{
					Host:    "host",
					Path:    "path",
					AppPort: 233,
				}}},
			}), routesAfter: &fakes.MappingRoutesService{Routes: map[string]sdk.AppRouting{
				"app-a": {Routes: sdk.AppRoutes{{
					Host:    "host",
					Path:    "path",
					AppPort: 233,
				}}},
			}}, expectedWorked: true,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(tt *testing.T) {
			providers := &fakes.ProviderService{
				Job:               &test.job,
				AppProviders:      map[string]sdk.AppProvider{},
				PipelineProviders: map[string]sdk.PipelineProvider{},
				RoutingProviders:  map[string]sdk.RoutingProvider{},
			}
			if test.appProvider != nil {
				providers.AppProviders[test.providerId] = test.appProvider
			}
			if test.pipelineProvider != nil {
				providers.PipelineProviders[test.providerId] = test.pipelineProvider
			}

			if test.routingProvider != nil {
				providers.RoutingProviders[test.providerId] = test.routingProvider
			}

			if test.appsBefore == nil {
				test.appsBefore = &fakes.MappingAppsService{Apps: map[string]apps.App{}}
			}
			if test.pipelinesBefore == nil {
				test.pipelinesBefore = &fakes.MappingPipelinesService{Pipelines: map[string]pipelines.Pipeline{}}
			}
			if test.routesBefore == nil {
				test.routesBefore = &fakes.MappingRoutesService{Routes: map[string]sdk.AppRouting{}}
			}

			r := NewReconciler(service.Core{
				Apps:      test.appsBefore,
				Providers: providers,
				Routing:   test.routesBefore,
				Pipelines: test.pipelinesBefore,
			}, 1*time.Minute)
			worked, err := r.Run()
			if !errors.Is(err, test.expectedErr) {
				tt.Errorf("\nwanted err %v\n   got err %v", test.expectedErr, err)
			}

			if worked != test.expectedWorked {
				tt.Errorf("\nwanted worked: %v\n   got worked: %v", test.expectedWorked, worked)
			}

			if test.appProvider != nil && test.appProvider.RecordedTime != test.recordedTime {
				tt.Errorf("\nwanted time: %v\n   got time: %v", test.recordedTime, test.appProvider.RecordedTime)
			} else if test.pipelineProvider != nil && test.pipelineProvider.RecordedTime != test.recordedTime {
				tt.Errorf("\nwanted time: %v\n   got time: %v", test.recordedTime, test.pipelineProvider.RecordedTime)
			} else if test.routingProvider != nil && test.routingProvider.RecordedTime != test.recordedTime {
				tt.Errorf("\nwanted time: %v\n   got time: %v", test.recordedTime, test.routingProvider.RecordedTime)
			}

			if test.appsAfter != nil && !cmp.Equal(test.appsAfter.Apps, test.appsBefore.Apps) {
				tt.Errorf("\napp service states don't match: \n%s\n", cmp.Diff(test.appsAfter.Apps, test.appsBefore.Apps))
			} else if test.pipelinesAfter != nil && !cmp.Equal(test.pipelinesAfter.Pipelines, test.pipelinesBefore.Pipelines) {
				tt.Errorf("\npipeline service states don't match: \n%s\n", cmp.Diff(test.pipelinesAfter.Pipelines, test.pipelinesBefore.Pipelines))
			} else if test.routesAfter != nil && !cmp.Equal(test.routesAfter.Routes, test.routesBefore.Routes) {
				tt.Errorf("\nrouting service states don't match: \n%s\n", cmp.Diff(test.routesAfter.Routes, test.routesBefore.Routes))
			}
		})
	}

}
