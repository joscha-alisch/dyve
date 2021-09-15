package concourse

import (
	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/go-concourse/concourse"
)

type API interface {
}

type ConcourseCli interface {
	Team(name string) ConcourseTeamCli
}

type ConcourseTeamCli interface {
	ListPipelines() ([]atc.Pipeline, error)
}

func main() {
	p, _, _ := concourse.NewClient().Team("bla").Pipeline("asd")

}

func getCli() ConcourseCli {
	return concourse.NewClient()
}
