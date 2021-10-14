package reconciler

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/internal/core/service"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
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
	r.Handler(database.ReconcileAppProvider, r.reconcileAppProvider)
	r.Handler(database.ReconcilePipelineProvider, r.reconcilePipelineProvider)
	r.Handler(database.ReconcileGroupProvider, r.reconcileGroupProvider)

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
