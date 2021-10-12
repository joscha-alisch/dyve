package reconciler

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

func NewReconciler(db database.Database, m provider.Manager, olderThan time.Duration) recon.Reconciler {
	if olderThan == 0 {
		olderThan = time.Minute
	}

	r := &reconciler{
		Reconciler: recon.NewReconciler(db, olderThan),
		db:         db,
		m:          m,
	}
	r.Handler(database.ReconcileAppProvider, r.reconcileAppProvider)
	r.Handler(database.ReconcilePipelineProvider, r.reconcilePipelineProvider)
	r.Handler(database.ReconcileGroupProvider, r.reconcileGroupProvider)

	return r
}

type reconciler struct {
	recon.Reconciler
	db database.Database
	m  provider.Manager
}

func (r *reconciler) reconcileAppProvider(j recon.Job) error {
	p, err := r.m.GetAppProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		return r.db.DeleteProvider(j.Guid, provider.TypeApps)
	}
	if err != nil {
		return err
	}

	apps, err := p.ListApps()
	if err != nil {
		return err
	}

	return r.db.UpdateApps(j.Guid, apps)
}

func (r *reconciler) reconcileGroupProvider(j recon.Job) error {
	p, err := r.m.GetGroupProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		return r.db.DeleteProvider(j.Guid, provider.TypeGroups)
	}
	if err != nil {
		return err
	}

	groups, err := p.ListGroups()
	if err != nil {
		return err
	}

	return r.db.UpdateGroups(j.Guid, groups)
}

func (r *reconciler) reconcilePipelineProvider(j recon.Job) error {
	p, err := r.m.GetPipelineProvider(j.Guid)
	if errors.Is(err, provider.ErrNotFound) {
		return r.db.DeleteProvider(j.Guid, provider.TypePipelines)
	}
	if err != nil {
		return err
	}

	pipelines, err := p.ListPipelines()
	if err != nil {
		return err
	}

	updates, err := p.ListUpdates(j.LastUpdated)

	err = r.db.UpdatePipelines(j.Guid, pipelines)
	if err != nil {
		return err
	}

	err = r.db.AddPipelineVersions(j.Guid, updates.Versions)
	if err != nil {
		return err
	}
	err = r.db.AddPipelineRuns(j.Guid, updates.Runs)
	if err != nil {
		return err
	}

	return nil
}
