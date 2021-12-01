package reconciler

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/internal/core/service"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"time"
)

func NewReconciler(core service.Core, olderThan time.Duration) recon.Reconciler {
	if olderThan == 0 {
		olderThan = time.Minute
	}

	r := &reconciler{
		Reconciler: recon.NewReconciler(core.Providers, olderThan),
		core:       core,
	}

	r.Handler(provider.ReconcileAppProvider, r.reconcileAppProvider)
	r.Handler(provider.ReconcileRoutingProviders, r.reconcileAppRouting)
	r.Handler(provider.ReconcileInstancesProviders, r.reconcileAppInstances)

	r.Handler(provider.ReconcilePipelineProvider, r.reconcilePipelineProvider)
	r.Handler(provider.ReconcileGroupProvider, r.reconcileGroupProvider)

	return r
}

type reconciler struct {
	recon.Reconciler
	core service.Core
}

func (r *reconciler) reconcileAppProvider(j recon.Job) error {
	p, err := r.core.Providers.GetAppProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		return r.core.Providers.DeleteAppProvider(j.Guid)
	}
	if err != nil {
		return err
	}

	apps, err := p.ListApps()
	if err != nil {
		return err
	}

	return r.core.Apps.UpdateApps(j.Guid, apps)
}

func (r *reconciler) reconcileAppRouting(j recon.Job) error {
	p, err := r.core.Providers.GetRoutingProviders()
	if err != nil {
		return err
	}

	routing := sdk.AppRouting{}
	for _, routingProvider := range p {
		result, err := routingProvider.GetAppRouting(j.Guid)
		if err != nil {
			return err
		}
		routing.Routes = append(routing.Routes, result.Routes...)
	}

	err = r.core.Routing.UpdateRoutes(j.Guid, routing)
	if err != nil {
		return err
	}

	return nil
}

func (r *reconciler) reconcileAppInstances(j recon.Job) error {
	p, err := r.core.Providers.GetInstancesProviders()
	if err != nil {
		return err
	}

	instances := sdk.AppInstances{}
	for _, routingProvider := range p {
		result, err := routingProvider.GetAppInstances(j.Guid)
		if err != nil {
			return err
		}
		instances = append(instances, result...)
	}

	err = r.core.Instances.UpdateInstances(j.Guid, instances)
	if err != nil {
		return err
	}

	return nil
}

func (r *reconciler) reconcileGroupProvider(j recon.Job) error {
	p, err := r.core.Providers.GetGroupProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		return r.core.Providers.DeleteGroupProvider(j.Guid)
	}
	if err != nil {
		return err
	}

	groups, err := p.ListGroups()
	if err != nil {
		return err
	}

	return r.core.Groups.UpdateGroups(j.Guid, groups)
}

func (r *reconciler) reconcilePipelineProvider(j recon.Job) error {
	p, err := r.core.Providers.GetPipelineProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		return r.core.Providers.DeletePipelineProvider(j.Guid)
	}
	if err != nil {
		return err
	}

	pipelines, err := p.ListPipelines()
	if err != nil {
		return err
	}

	updates, err := p.ListUpdates(j.LastUpdated)

	err = r.core.Pipelines.UpdatePipelines(j.Guid, pipelines)
	if err != nil {
		return err
	}

	err = r.core.Pipelines.AddPipelineVersions(j.Guid, updates.Versions)
	if err != nil {
		return err
	}
	err = r.core.Pipelines.AddPipelineRuns(j.Guid, updates.Runs)
	if err != nil {
		return err
	}

	return nil
}
