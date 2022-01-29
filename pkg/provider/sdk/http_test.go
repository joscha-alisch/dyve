package sdk

import "testing"

func TestHttp(t *testing.T) {
	go func() {
		_ = ListenAndServe(":46283", ProviderConfig{
			Apps:      &fakeAppProvider{},
			Pipelines: &fakePipelineProvider{},
			Groups:    &fakeGroupProvider{},
			Routing:   &fakeRoutingProvider{},
			Instances: &fakeInstancesProvider{},
		})
	}()
}
