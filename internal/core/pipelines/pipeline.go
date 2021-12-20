package pipelines

import "github.com/joscha-alisch/dyve/pkg/provider/sdk"

type Pipeline struct {
	sdk.Pipeline `json:",inline" bson:",inline"`
	ProviderId   string `json:"providerId" bson:"providerId"`
}
