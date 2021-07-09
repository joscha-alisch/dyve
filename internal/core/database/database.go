package database

import "github.com/joscha-alisch/dyve/pkg/provider/sdk"

type Database interface {
	ListAppsPaginated(perPage int, page int) (sdk.AppPage, error)
}
