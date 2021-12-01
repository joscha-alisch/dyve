package groups

import "github.com/joscha-alisch/dyve/pkg/provider/sdk"

type GroupByProviderMap map[string]ProviderWithGroups

type ProviderWithGroups struct {
	Provider string      `json:"provider"`
	Name     string      `json:"name"`
	Groups   []sdk.Group `json:"groups"`
}

type GroupWithProvider struct {
	Provider string    `bson:"provider"`
	Group    sdk.Group `bson:",inline"`
}
