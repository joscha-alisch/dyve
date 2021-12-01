package apps

import "github.com/joscha-alisch/dyve/pkg/provider/sdk"

type App struct {
	sdk.App    `json:",inline" bson:",inline"`
	ProviderId string `json:"providerId" bson:"providerId"`
}
