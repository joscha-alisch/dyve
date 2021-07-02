package cloudfoundry

type Database interface {
	FetchReconcileJob() (ReconcileJob, bool)
	UpsertOrg(o Org) error
	UpsertSpace(s Space) error
	UpsertApp(a App) error
}
