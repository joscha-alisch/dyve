module github.com/joscha-alisch/dyve

go 1.16

require (
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75
	github.com/approvals/go-approval-tests v0.0.0-20210628084631-e4f9005a9e2e
	github.com/benweissmann/memongo v0.1.1
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20210621174645-7773f7e22665
	github.com/concourse/concourse v1.6.1-0.20210729233333-b047d257f253
	github.com/go-pkgz/auth v1.18.0
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v39 v39.1.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/rs/zerolog v1.25.0
	github.com/spf13/viper v1.7.0
	go.mongodb.org/mongo-driver v1.7.1
	gonum.org/v1/gonum v0.9.3
)

replace github.com/go-pkgz/auth => github.com/joscha-alisch/auth v1.18.1-0.20211006101921-5702e9e067f9
