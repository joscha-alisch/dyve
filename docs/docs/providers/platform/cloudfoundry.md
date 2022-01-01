# CloudFoundry Provider

This provider retrieves apps from a CloudFoundry installation.

## Run
### With Helm

The CloudFoundry provider [is part of our helm chart.](../../install/helm.md)
Set `providers.cloudfoundry.enabled` to `true` in order to enable it. Also fill in `providers.cloudfoundry.config` to configure it with your CF Access. The minimal necessary config is as follows:
```yaml
providers:
  cloudfoundry:
    enabled: true
    config:
      cloudfoundry:
        api: "https://your_cloudfoundry_api"
        user: "username"
        password: "password"
      database:
        uri: mongodb://your_mongodb_server:27017
```

For a full list of configuration parameters, [see below.](#config)

### With Docker Image

Run the latest docker image with a mounted configuration file.

```bash
docker run -it --rm -v $(pwd)/my_config.yaml:/app/config.yaml ghcr.io/joscha-alisch/dyve-provider-cf:latest
```

For a full list of configuration parameters, [see below.](#config)

### As Binary

Download the latest binary for your OS [from the GitHub releases](https://github.com/joscha-alisch/dyve/releases).
Then run it, providing a configuration yaml file via `-c`:

```bash
dyve-provider-cf -c config.yaml
```

For a full list of configuration parameters, [see below.](#config)

## Config 

The CloudFoundry provider is configured via a yaml file with the following parameters and defaults:

```yaml
port: 9000      # The port to listen on
logLevel: info  # The log level (debug, info, warn, error)

cloudfoundry:
  api: ""       # The CloudFoundry API URL
  user: ""      # The User to authenticate with
  password: ""  # The password for the user

reconciliation:
  cacheSeconds: 20 # For how many to cache apps/spaces/orgs, before retrieving them again via the CF API

database:
  uri: mongodb://localhost:27017  # The MongoDB URL used for caching
  name: cf                        # The MongoDB database name
```