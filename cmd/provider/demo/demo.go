package main

import (
	"github.com/joscha-alisch/dyve/internal/provider/demo"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
)

func main() {
	p := demo.NewProvider()

	err := sdk.ListenAndServe(":9003", sdk.ProviderConfig{
		Apps:      p,
		Pipelines: p,
	})
	if err != nil {
		panic(err)
	}
}
