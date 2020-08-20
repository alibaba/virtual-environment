# Quick start

Check the [examples folder](https://github.com/alibaba/virtual-environment/tree/master/examples) for the code structure description.

## Fetch demo code

Pull the git repository, and enter examples folder.

```bash
git clone https://github.com/alibaba/virtual-environment.git
cd virtual-environment/examples
```

## Verify virtual environment isolation

Without KtConnect toolkit, the isolation only works between pods in cluster. In order to observe the result, we could create a temporary pod in the cluster and send test request from there.

If there is already other suitable pod in the cluster, this step can be skip.

```bash
# Create a pod inside the cluster for sending request
kubectl create deployment sleep --image=virtualenvironment/sleep --dry-run -o yaml \
        | istioctl kube-inject -f - | kubectl apply -n default -f -
```

Use `app.sh` script to quickly create all required instances of VirtualEnvironment, Service and Deployment.

```bash
# Create demo resource instances
deploy/app.sh apply default
```

Use `kubectl get virtualenvironment`, `kubectl get service` and `kubectl get deployment` command to check resource creation progress, wait util all resources are ready to use.

Enter any pod in the same Namespace, e.g. the `sleep` pod created in previous step.

```bash
# Enter pod in cluster
kubectl exec -n default -it $(kubectl get -n default pod -l app=sleep -o jsonpath='{.items[0].metadata.name}') /bin/sh
```

Use `curl` tool call `app-js` service with different virtual-environment name in Header. Please notice that in this example, symbol `-` is used as the environment level splitter, and the virtual-environment Header key is configured as `ali-env-mark`.

All service instances would append formatted text as `[app-name @ virtual-environment-it-belongs-to] <- virtual-environment-header-on-request`. Observe the actual service instance response:

```bash
# Use dev.proj1 header
> curl -H 'ali-env-mark: dev.proj1' app-js:8080/demo
  [springboot @ dev.proj1] <-dev.proj1
  [go @ dev] <-dev.proj1
  [node @ dev.proj1] <-dev.proj1

# Use dev.proj1.feature1 header
> curl -H 'ali-env-mark: dev.proj1.feature1' app-js:8080/demo
  [springboot @ dev.proj1.feature1] <-dev.proj1.feature1
  [go @ dev] <-dev.proj1.feature1
  [node @ dev.proj1] <-dev.proj1.feature1

# Use dev.proj2 header
> curl -H 'ali-env-mark: dev.proj2' app-js:8080/demo
  [springboot @ dev] <-dev.proj2
  [go @ dev.proj2] <-dev.proj2
  [node @ dev] <-dev.proj2

# Without any header
# As AutoInject configure item is set, when request pass app-node service,
# the request is automatically headered with the virtual environment name where the Pod is located.
> curl app-js:8080/demo
  [springboot @ dev] <-dev
  [go @ dev] <-dev
  [node @ dev] <-empty
```

## Add local service instance into isolation

[KtConnect](https://github.com/alibaba/kt-connect) toolkit can setup a proxy between local and remote cluster, and add local service instance into any virtual environment in the cluster.

In order to let the shadow pod created via `ktctl` command follow the isolation rules, sidecar auto injector should be enabled in the target namespace:

```bash
kubectl label namespaces default istio-injection=enabled
```

- Use `ktctl connect` command allow local instance accessing remote cluster

```bash
# Notice: the label parameter specified the virtual environment name to join
sudo ktctl --label virtual-env=dev.proj2 --namespace default connect
```

Now, local shell can curl the `app-js` service inside the remote cluster directly.

```bash
# As envHeader.autoInject configure is enabled, request send from local is appended virtual environment header automatically
$ curl app-js:8080/demo
  [springboot @ dev] <-dev.proj2
  [go @ dev.proj2] <-dev.proj2
  [node @ dev] <-dev.proj2
```

- Use `ktctl mesh` command let remote service able to access local port.

```bash
# Start up a app-java service instance, listen on port 8080
# Notice: The envMark variable here will only show up in response output, and has nothing to do with actual routing control
cd examples/springboot
envMark=local mvn spring-boot:run

# the label parameter specified the virtual environment name to join
# app-java-dev is the name of app-java Deployment in shared environment
sudo ktctl --label virtual-env=dev.proj2 --namespace default mesh app-java-dev --expose 8080
```

Now there is a `app-java` service instance from local in the `dev.proj2` virtual environment, so the new route targets should be:


```
            +----------+
dev         |  app-js  |
            +----------+
                           +----------+   +-----------------+
dev.proj2                  |  app-go  |   | app-java(local) |
                           +----------+   +-----------------+
```

Make a request again from local shell, this time the request goes through the app-js and app-go service in the cluster, will route back to the local app-java instance.

```bash
$ curl app-js:8080/demo
  [springboot @ local] <-dev.proj2
  [go @ dev.proj2] <-dev.proj2
  [node @ dev] <-dev.proj2
```

## Clean up demo resources

```bash
# Delete related resources
deploy/app.sh delete default
# Delete the temporary pod
kubectl delete -n default deployment sleep
```
