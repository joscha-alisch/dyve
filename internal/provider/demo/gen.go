package demo

import (
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/google/uuid"
	"strings"
)

func id() string {
	return uuid.New().String()
}

func appName() string {
	return randomdata.Adjective() + "-" + randomdata.Adjective() + "-" + appWords[randomdata.Number(0, len(appWords))]
}

func pipelineName() string {
	return pipelineWords[randomdata.Number(0, len(pipelineWords))] + "-" + randomdata.Adjective() + "-" + randomdata.Noun()
}

func version() string {
	return fmt.Sprintf("%d.%d.%d", randomdata.Number(0, 10), randomdata.Number(0, 10), randomdata.Number(0, 10))
}

func namespace() string {
	return randomdata.City()
}

func team() string {
	return strings.ToLower(randomdata.FirstName(randomdata.RandomGender)) + "s-" + randomdata.Adjective() + "-" + strings.ToLower(randomdata.Noun()) + "s"
}

func pipelineStep() string {
	return pipelineStepNames[randomdata.Number(0, len(pipelineStepNames)-1)]
}

var appWords = []string{
	"generator",
	"service",
	"api",
	"feeder",
	"db",
	"cms",
	"wordpress",
	"queue",
	"kafka-ingester",
	"consumer",
	"provider",
	"filter",
	"creator",
	"deployer",
	"collector",
	"retriever",
	"destroyer",
	"editor",
	"fetcher",
	"cluster",
	"pod",
	"infinidash",
	"server",
	"actuator",
	"updater",
	"search",
	"finder",
	"remapper",
	"locator",
	"definer",
	"worker",
	"threader",
	"builder",
	"distributor",
	"runner",
	"grapher",
	"drawer",
	"sleeper",
	"preview",
	"interface",
	"adapter",
	"executor",
	"lister",
	"stopper",
	"app",
	"frontend",
	"backend",
	"querier",
	"cache",
	"store",
	"connector",
	"client",
	"faker",
}

var pipelineWords = []string{
	"run",
	"deploy",
	"build",
	"test",
	"clean",
	"verify",
	"notify",
}

var pipelineStepNames = []string{
	"build docker",
	"build image",
	"remove image",
	"build app",
	"test",
	"cdc",
	"run-cdc",
	"verify cdc",
	"clean repo",
	"notify slack",
	"slack",
	"tool",
	"build",
	"verify",
}
