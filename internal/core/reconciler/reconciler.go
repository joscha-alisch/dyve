package reconciler

import (
	"errors"
	"github.com/joscha-alisch/dyve/internal/core/database"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	recon "github.com/joscha-alisch/dyve/internal/reconciliation"
	"time"
)

func NewReconciler(db database.Database, m provider.Manager, olderThan time.Duration) recon.Reconciler {
	r := &reconciler{
		Reconciler: recon.NewReconciler(db, olderThan),
		db:         db,
		m:          m,
	}
	r.Handler(database.ReconcileAppProvider, r.reconcileAppProvider)

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
		r.db.DeleteAppProvider(j.Guid)
		return nil
	}
	if err != nil {
		return err
	}

	apps, err := p.ListApps()
	if err != nil {
		return err
	}

	r.db.UpdateApps(j.Guid, apps)

	return nil
}
