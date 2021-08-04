package reconciliation

import (
	"github.com/rs/zerolog/log"
	"time"
)

type JobProvider interface {
	AcceptReconcileJob(olderThan time.Duration) (Job, bool)
}

type Type string
type ReconcileHandler func(j Job) error
type Job struct {
	Type        Type
	Guid        string
	LastUpdated time.Time
}

type Reconciler interface {
	Run() (bool, error)
	Handler(t Type, f ReconcileHandler)
}

func NewReconciler(p JobProvider, olderThan time.Duration) Reconciler {
	return &reconciler{
		p:         p,
		mapping:   map[Type]ReconcileHandler{},
		olderThan: olderThan,
	}
}

type reconciler struct {
	p         JobProvider
	mapping   map[Type]ReconcileHandler
	olderThan time.Duration
}

func (r *reconciler) Handler(t Type, f ReconcileHandler) {
	r.mapping[t] = f
}

func (r *reconciler) Run() (bool, error) {
	j, ok := r.p.AcceptReconcileJob(r.olderThan)
	if !ok {
		return false, nil
	}

	log.Info().Interface("job", j).Msg("reconciling")

	f := r.mapping[j.Type]
	if f == nil {
		return true, nil
	}

	return true, f(j)
}
