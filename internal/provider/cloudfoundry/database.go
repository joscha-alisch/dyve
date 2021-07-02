package cloudfoundry

type Database interface {
	FetchReconcileJob() (ReconcileJob, bool)
	UpsertOrg(guid string, o Org) error
	UpsertSpace(guid string, s Space) error
	UpsertApp(guid string, a App) error
}
