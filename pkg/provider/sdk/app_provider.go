package sdk

type AppProvider interface {
	Apps() []App
	App(id int) App
}
