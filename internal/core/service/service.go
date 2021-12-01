package service

import (
	"github.com/joscha-alisch/dyve/internal/core/apps"
	"github.com/joscha-alisch/dyve/internal/core/groups"
	"github.com/joscha-alisch/dyve/internal/core/instances"
	"github.com/joscha-alisch/dyve/internal/core/pipelines"
	"github.com/joscha-alisch/dyve/internal/core/provider"
	"github.com/joscha-alisch/dyve/internal/core/routing"
	"github.com/joscha-alisch/dyve/internal/core/teams"
)

type Core struct {
	Teams     teams.Service
	Apps      apps.Service
	Groups    groups.Service
	Providers provider.Service
	Pipelines pipelines.Service
	Routing   routing.Service
	Instances instances.Service
}
