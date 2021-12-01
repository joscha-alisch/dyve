package live

import (
	"github.com/joscha-alisch/dyve/internal/core/service"
	"github.com/joscha-alisch/dyve/internal/core/ws"
	"github.com/joscha-alisch/dyve/pkg/provider/sdk"
	"github.com/rs/zerolog/log"
	"math/rand"
	"sync"
	"time"
)

type connections map[int]*ws.Connection

type AppUpdate struct {
	Routing   sdk.AppRouting   `json:"routing"`
	Instances sdk.AppInstances `json:"instances"`
}

func NewAppViewer(core service.Core) *AppViewer {
	return &AppViewer{
		core:    core,
		appSubs: make(map[string]connections),
		mu:      &sync.Mutex{},
	}
}

type AppViewer struct {
	core    service.Core
	appSubs map[string]connections
	mu      *sync.Mutex
}

func (v *AppViewer) Run() {
	go v.updateWorker()
}

func (v *AppViewer) AddWs(id string, c *ws.Connection) chan error {
	wsId := rand.Int()

	v.subscribe(id, wsId, c)

	errChan := make(chan error)
	c.On("update", func() {
		err := v.core.Providers.RequestAppUpdate(id)
		if err != nil {
			errChan <- err
		}
	})

	go v.wsWorker(wsId, id, c, errChan)

	return errChan
}

func (v *AppViewer) wsWorker(wsId int, appId string, c *ws.Connection, errChan chan error) {
	err := c.Run()
	if err != nil {
		errChan <- err
	}
	delete(v.appSubs[appId], wsId)
	close(errChan)
}

func (v *AppViewer) updateWorker() {
	t := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-t.C:
			v.mu.Lock()
			for appId, connections := range v.appSubs {
				routing, err := v.core.Routing.GetRoutes(appId)
				if err != nil {
					log.Error().Err(err).Str("app", appId).Msg("error getting app routing")
				}

				instances, err := v.core.Instances.GetInstances(appId)
				if err != nil {
					log.Error().Err(err).Str("app", appId).Msg("error getting app instances")
				}

				for wsId, connection := range connections {
					err := connection.Send(AppUpdate{
						Routing:   routing,
						Instances: instances,
					})
					if err != nil {
						log.Error().Err(err).Str("app", appId).Int("ws", wsId).Msg("error updating websocket")
						continue
					}
				}
			}
			v.mu.Unlock()
		}
	}
}

func (v *AppViewer) subscribe(app string, wsId int, conn *ws.Connection) {
	v.mu.Lock()
	if v.appSubs[app] == nil {
		v.appSubs[app] = make(connections)
	}

	v.appSubs[app][wsId] = conn
	v.mu.Unlock()
}
