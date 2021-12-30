# Helm

We provide a helm chart to deploy dyve into your kubernetes cluster.


???+ "Install Helm"

    In order to use the helm chart, you will need helm installed and configured to use your kubernetes cluster. You can find [the instructions here.](https://helm.sh/docs/intro/install/)

With helm installed, add the repository:

```bash
helm repo add dyve https://joscha-alisch.github.io/dyve
helm repo update
```

and install the chart:

```bash
helm upgrade --install -n dyve dyve dyve/dyve
```

## Configuration

In order to set configuration parameters, create a `values.yaml` and specify it when installing the chart:

```bash
helm upgrade --install -f values.yaml -n dyve dyve dyve/dyve 
```

The following shows the available options together with its default values:

```yaml title="values.yaml"
--8<-- "../dyve/infra/helm/dyve/values.yaml"
```
