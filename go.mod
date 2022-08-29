module github.com/joscha-alisch/dyve

go 1.16

require (
	github.com/Pallinder/go-randomdata v1.2.0
	github.com/ajstarks/svgo v0.0.0-20210406150507-75cfd577ce75
	github.com/approvals/go-approval-tests v0.0.0-20210628084631-e4f9005a9e2e
	github.com/bradleyfalzon/ghinstallation v1.1.1
	github.com/cloudfoundry-community/go-cfclient v0.0.0-20210621174645-7773f7e22665
	github.com/fatih/structs v1.1.0
	github.com/go-pkgz/auth v1.18.0
	github.com/google/go-cmp v0.5.6
	github.com/google/go-github/v39 v39.2.0
	github.com/google/uuid v1.3.0
	github.com/gopherjs/gopherjs v0.0.0-20200217142428-fce0ec30dd00 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/gorilla/websocket v1.4.2
	github.com/jeremywohl/flatten v1.0.1
	github.com/kr/text v0.2.0 // indirect
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/zerolog v1.28.0
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/spf13/viper v1.10.1
	github.com/tryvium-travels/memongo v0.3.2
	go.mongodb.org/mongo-driver v1.8.1
	gonum.org/v1/gonum v0.9.3
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/h2non/gock.v1 v1.1.2
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

replace github.com/go-pkgz/auth => github.com/joscha-alisch/auth v1.18.1-0.20211006101921-5702e9e067f9
